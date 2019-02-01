package timer_wheel

import (
	"time"
	"sync"
)

type Wheel struct {
	//父级
	parent *Wheel
	//时间间隔
	duration time.Duration
	//轮子的开始时间
	currentTime time.Time
	//刻度槽
	bucket []DelayElement
	//槽位数
	slots int
	//时间间隔
	interval int64
	mutex    sync.Mutex
	queue    DelayQueue
}

//事实上,将指针移到t
func (wheel *Wheel) Advance(t *time.Time) {
	if t.Nanosecond() >= wheel.currentTime.Add(wheel.duration).Nanosecond() {
		//降噪
		wheel.currentTime = t.Add(-1 * time.Duration(t.UnixNano()%wheel.duration.Nanoseconds()))

		if wheel.parent != nil {
			wheel.parent.Advance(&wheel.currentTime)
		}
	}

}

//增加一个元素
func (wheel *Wheel) Add(task Task) bool {
	expireNans := task.Expire().Nanosecond()
	if expireNans >= wheel.currentTime.Add(wheel.duration).Nanosecond() {
		return false
	}
	endTime := wheel.currentTime.Add(time.Duration(wheel.interval)).Nanosecond()
	if endTime-expireNans >= 0 {
		if wheel.parent == nil {
			wheel.parent = &Wheel{nil, time.Duration(wheel.interval), wheel.currentTime, make([]DelayElement, wheel.slots), wheel.slots, 0, sync.Mutex{}, wheel.queue}
			wheel.parent.interval = int64(wheel.parent.slots) * wheel.parent.duration.Nanoseconds()
		}
		return wheel.parent.Add(task)
	}
	index := task.Expire().Sub(wheel.currentTime) % wheel.duration
	bucket := wheel.bucket[index]
	if bucket == nil {
		wheel.mutex.Lock()
		bucket = wheel.bucket[index]
		//double check
		if bucket == nil {
			wheel.bucket[index] = &Bucket{task.Expire(), &SimpleList{}}
			bucket = wheel.bucket[index]
		}
		wheel.queue.Put(bucket)
		wheel.mutex.Unlock()
	}
	bucket.Add(task)
	return true
}

type Bucket struct {
	time time.Time
	list TaskList
}

func (bucket *Bucket) Compare(other Comparable) int8 {
	v, ok := other.(*Bucket)
	if ok {
		a := bucket.time.Nanosecond()
		b := v.time.Nanosecond()
		if a > b {
			return 1
		}
		if a == b {
			return 0
		}
	}
	return -1
}

func (bucket *Bucket) Expire() time.Duration {
	return time.Now().Sub(bucket.time)
}
func (bucket *Bucket) SetExpire(time time.Time) {
	bucket.time = time
}

func (bucket *Bucket) Add(task Task) {
	bucket.list.Add(task)
}
