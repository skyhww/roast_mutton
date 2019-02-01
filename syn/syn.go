package syn

import (
	"time"
	"sync/atomic"
)

type WaitNotify interface {
	Wait(duration time.Duration) (time.Duration, bool)
	Notify()
}

type WaitNotifyQueue struct {

}

type TimeLock struct {
	time time.Duration
	c    chan struct{}
}

func NewTimeLock(t time.Duration) *TimeLock {
	return &TimeLock{t, make(chan struct{}, 1)}
}
func (lock *TimeLock) Lock() {
	select {
	case <-lock.c:
	case <-time.After(lock.time):
	}
}
func (lock *TimeLock) Unlock() {
	lock.c <- struct{}{}
}

type SimpleWaitNotify struct {
	notifier int32
	interval time.Duration
}

func NewWaitNotify() WaitNotify {
	return &SimpleWaitNotify{0, time.Millisecond}
}

//检测notifier有没有变更过，若有，则表明被唤醒过
func (s *SimpleWaitNotify) Wait(duration time.Duration) (time.Duration, bool) {
	now := time.Now()
	notifier := atomic.LoadInt32(&s.notifier)
	interval := s.interval
	if interval > duration {
		interval = duration
	}
	for {
		select {
		case <-time.After(interval):
			t := time.Now().Sub(now)
			get := notifier != atomic.LoadInt32(&s.notifier)
			if get {
				return t, get
			}
			//时间不够
			if t >= duration {
				return t, get
			}
			if t+interval > duration {
				interval = duration - t
			}
		}
	}
}
func (s *SimpleWaitNotify) Notify() {
	atomic.AddInt32(&s.notifier, 1)
}
