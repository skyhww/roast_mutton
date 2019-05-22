package core

import (
	"time"
	"sync"
	"sync/atomic"
	"fmt"
)

//延迟队列,小顶堆
type DelayQueue interface {
	Add(compareAble CompareAble)
	//从堆顶获取一个元素，阻塞时间不超过duration
	Poll(duration *time.Duration) CompareAble
}

//最小堆实现，基于数组
type UnboundDelayQueue struct {
	elements []Expire
	Lock     *sync.Mutex
	size     int
	//避免数组的频繁删除
	//当数组中空闲空间大于threshold后，删除多余的节点
	Threshold int
	ChanList  *LinkedList
	ch        chan int
	writer    int64
	reader    int64
}

//获取左孩子的索引
func (queue *UnboundDelayQueue) getLeftChildIndex(index int) int {
	left := (index << 1) + 1
	if left >= queue.Size() {
		return -1
	}
	return left
}

func (queue *UnboundDelayQueue) compact() {
	//提升空间利用率
	tmp := len(queue.elements) - queue.size
	if tmp >= queue.size {
		queue.elements = queue.elements[0:queue.size]
	}
}

//获取右节点的索引
func (queue *UnboundDelayQueue) getRightChildIndex(index int) int {
	right := (index + 1) << 1
	if right >= queue.Size() {
		return -1
	}
	return right
}

//获取父节点的索引
func (queue *UnboundDelayQueue) getParentIndex(index int) int {
	return (index - 1) >> 1
}

//获取父节点的索引
func (queue *UnboundDelayQueue) first() CompareAble {
	if queue.size == 0 {
		return nil
	}
	return queue.elements[0]
}

//从index开始往上重构小顶堆
func (queue *UnboundDelayQueue) siftUp() {

	for coordinate := queue.getParentIndex(queue.size - 1); coordinate >= 0; coordinate-- {
		left := queue.getLeftChildIndex(coordinate)
		right := queue.getRightChildIndex(coordinate)
		less := left
		if right > 0 && queue.elements[left].Compare(queue.elements[right]) > 0 {
			less = right
		}
		if queue.elements[coordinate].Compare(queue.elements[less]) > 0 {
			queue.elements[coordinate], queue.elements[less] = queue.elements[less], queue.elements[coordinate]
		}
	}

}
func (queue *UnboundDelayQueue) peek() Expire {
	if queue.size == 0 {
		return nil
	}
	return queue.elements[queue.size-1]
}
func (queue *UnboundDelayQueue) siftDown() {
	current := 0
	for current >= 0 {
		left := queue.getLeftChildIndex(current)
		right := queue.getRightChildIndex(current)
		//叶子节点
		if left < 0 && right < 0 {
			return
		}
		less := left
		if left > 0 && right > 0 {
			if queue.elements[left].Compare(queue.elements[right]) > 0 {
				less = right
			}
		} else if right > 0 {
			less = right
		}
		if queue.elements[current].Compare(queue.elements[less]) > 0 {
			queue.elements[current], queue.elements[less] = queue.elements[less], queue.elements[current]
			current = less
		} else {
			//未发生交换，则子树性质保持
			return
		}
	}
}

//拿到第一个节点，并重构小顶堆
func (queue *UnboundDelayQueue) poll() Expire {
	if queue.size == 0 {
		return nil
	}
	e := queue.elements[0]
	queue.elements[0] = queue.elements[ queue.size-1]
	queue.elements[ queue.size-1] = nil
	queue.size--
	queue.compact()
	queue.siftDown()
	return e
}
func (queue *UnboundDelayQueue) readerIncrease() {
	atomic.AddInt64(&queue.writer, 1)
}
func (queue *UnboundDelayQueue) readerDecrease() {
	atomic.AddInt64(&queue.writer, -1)
}
func (queue *UnboundDelayQueue) notify() {
	if queue.writer > 0 && queue.reader == 0 {
		queue.Lock.Lock()
		queue.reader = 1
		if queue.writer > 0 && len(queue.ch) == 0 {
			queue.ch <- 1
		}
		if len(queue.ch) == 0 {
			queue.ch <- 1
		}
		queue.Lock.Unlock()
		queue.reader = 0
	}
}

func (queue *UnboundDelayQueue) Add(e Expire) {
	queue.Lock.Lock()
	if queue.ch == nil {
		queue.ch = make(chan int, 1)
	}
	queue.add(e)
	if queue.writer > 0 && len(queue.ch) == 0 {
		queue.ch <- 1
	}
	queue.Lock.Unlock()

}

func (queue *UnboundDelayQueue) Poll(duration *time.Duration) Expire {
	c := time.NewTimer(*duration)
	defer c.Stop()
	begin := time.Now()
	queue.Lock.Lock()
	fmt.Println(time.Now().Sub(begin).Seconds())
	if queue.ch == nil {
		queue.ch = make(chan int, 1)
	}
	q := queue.peek()
	if q != nil && q.Expired() {
		q = queue.poll()
		queue.Lock.Unlock()
		return q
	}
	queue.Lock.Unlock()
	queue.readerIncrease()
	defer queue.readerDecrease()
	for {
		select {
		case <-c.C:
			{
				return nil
			}
		case <-queue.ch:
			{
				queue.Lock.Lock()
				q := queue.peek()
				if q != nil && q.Expired() {
					e := queue.poll()
					if e != nil {
						queue.Lock.Unlock()
						return e
					}
				}
				queue.Lock.Unlock()
			}
		}
	}
	return nil
}

/*func (queue *UnboundDelayQueue) Poll(duration *time.Duration) Expire {
	c := time.NewTimer(*duration)
	defer c.Stop()
	queue.Lock.Lock()
	q := queue.peek()
	if q != nil && q.Expired() {
		q = queue.poll()
		queue.Lock.Unlock()
		return q
	}
	queue.Lock.Unlock()
	ch := make(chan int, 1)
	defer close(ch)
	s := &signal{ch: ch}
	node := queue.ChanList.Add(s)
	for {
		select {
		case <-c.C:
			{
				node.Delete()
				return nil
			}
		case <-s.Wait():
			{
				node.Delete()
				queue.Lock.Lock()
				q := queue.peek()
				if q != nil && q.Expired() {
					e := queue.poll()
					if e != nil {
						queue.Lock.Unlock()
						return e
					}
				}
				queue.Lock.Unlock()
				//插到队尾
				node.Join(queue.ChanList)
			}
		}
	}
	return nil
}*/
func (queue *UnboundDelayQueue) add(e Expire) {
	if queue.size == 0 {
		queue.elements = make([]Expire, 0)
	}
	if queue.size < len(queue.elements) {
		queue.elements[queue.size] = e
	} else {
		queue.elements = append(queue.elements, e)
	}
	queue.size++
	queue.siftUp()
}

/*func (queue *UnboundDelayQueue) Add(e Expire) {
	queue.Lock.Lock()
	queue.add(e)
	queue.Lock.Unlock()
	head := queue.ChanList.GetHead()
	for head != nil {
		c, _ := head.Data.(*signal)
		if c != nil && !c.Notify() {
			head = head.Successor
		} else {
			return
		}
	}
}*/
func (queue *UnboundDelayQueue) Cap() int {
	return cap(queue.elements)
}
func (queue *UnboundDelayQueue) Size() int {
	return queue.size
}

type signal struct {
	ch  chan int
	tmp int32
}

func (s *signal) Wait() chan int {
	return s.ch
}

func (s *signal) Release() {
	//可能重复被唤醒
	atomic.StoreInt32(&s.tmp, 1)
}

func (s *signal) Notify() bool {
	success := atomic.CompareAndSwapInt32(&s.tmp, 1, 2)
	if success {
		if len(s.ch) > 0 {
			return true
		} else {
			s.ch <- 1
			return true
		}
	}
	return false

	/*
		if s.wait == 0 {
			return false
		}
		b := atomic.CompareAndSwapInt32(&s.count, 1, 2)
		if b {
			s.ch <- 1
			s.wait = 0
			return true
		}
		return false*/

}
