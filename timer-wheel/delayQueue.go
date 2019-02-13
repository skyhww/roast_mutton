package timer_wheel

import (
	"sync"
	"time"
	"roast_mutton/syn"
)

//延迟队列
type DelayQueue interface {
	//阻塞一段时间，直到一个元素可用
	Take(duration time.Duration) DelayElement
	Put(e DelayElement)
	Size() int
	Cap() int
}

//最小堆实现，基于数组
type UnboundDelayQueue struct {
	elements []DelayElement
	lock     sync.Mutex
	size     int
	//避免数组的频繁删除
	//当数组中空闲空间大于threshold后，删除多余的节点
	threshold  int
	waitNotify syn.WaitNotify
}

func NewUnboundDelayQueue(threshold int) DelayQueue {
	queue := &UnboundDelayQueue{}
	queue.threshold = threshold
	queue.waitNotify = syn.NewWaitNotify()
	return queue
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
	if queue.threshold < len(queue.elements)-queue.size {
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
func (queue *UnboundDelayQueue) first() DelayElement {
	if queue.size == 0 {
		return nil
	}
	return queue.elements[0]
}

//从index开始往上重构小顶堆
func (queue *UnboundDelayQueue) siftUp(index int) {
	if index >= queue.size-1 {
		return
	}
	coordinate := index
	for coordinate > 0 {
		//父索引
		parent := queue.getParentIndex(coordinate)
		if queue.elements[parent].Compare(queue.elements[coordinate]) > 0 {
			queue.elements[parent], queue.elements[coordinate] = queue.elements[coordinate], queue.elements[parent]
		}
		//另一颗子树索引
		other := 0
		if index%2 == 0 {
			//获取左子树
			other = coordinate - 1
		} else {
			//获取右子树
			other = coordinate + 1
		}
		if queue.size-1 >= other {
			if queue.elements[parent].Compare(queue.elements[other]) > 0 {
				queue.elements[parent], queue.elements[other] = queue.elements[other], queue.elements[parent]
			}
		}
		coordinate = parent
	}

}

//拿到第一个节点，并重构小顶堆
func (queue *UnboundDelayQueue) poll() DelayElement {
	len := len(queue.elements)
	if len == 0 {
		return nil
	}
	e := queue.elements[0]
	queue.elements[0] = queue.elements[len-1]
	queue.elements[len-1] = nil
	queue.compact()
	queue.siftUp(queue.size - 1)
	return e
}

func (queue *UnboundDelayQueue) notify() {
	queue.waitNotify.Notify()
}

//至多等待duration，事实上，这个队列的竞争较小
func (queue *UnboundDelayQueue) Take(duration time.Duration) DelayElement {
	//duration包含了获取锁的时间，可能并不精确
	for duration > 0 {
		queue.lock.Lock()
		first := queue.first()
		if first == nil {
			duration, _ = queue.waitNotify.Wait(duration)
		} else {
			expire := first.Expire()
			if expire < 0 {
				queue.lock.Unlock()
				return queue.poll()
			}
			if expire > duration {
				duration, _ = queue.waitNotify.Wait(duration)
			} else {
				duration, _ = queue.waitNotify.Wait(expire)
			}
		}
	}
	return nil
}

func (queue *UnboundDelayQueue) Put(e DelayElement) {
	queue.lock.Lock()
	defer queue.lock.Unlock()
	if queue.elements == nil {
		queue.elements = make([]DelayElement, 10)
	}
	if len(queue.elements) > queue.size {
		queue.elements = append(queue.elements, e)
	} else {
		queue.elements[queue.Size()] = e
	}
	queue.size++
	queue.siftUp(queue.size - 1)
	queue.notify()
}
func (queue *UnboundDelayQueue) Cap() int {
	return cap(queue.elements)
}
func (queue *UnboundDelayQueue) Size() int {
	return queue.size
}

type Comparable interface {
	Compare(other Comparable) int8
}

type DelayElement interface {
	Compare(other Comparable) int8
	Expire() time.Duration
	Add(task Task)
	SetExpire(duration time.Time)
}
