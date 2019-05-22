package core

import (
	"testing"
	"sync"
	"time"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"fmt"
	"math/rand"
	"sync/atomic"
)

func TestUnboundDelayQueue(test *testing.T) {
	runtime.GOMAXPROCS(4)
	var count int32
	var expired int32
	var added int32
	queue := &UnboundDelayQueue{Lock: &sync.Mutex{}, Threshold: 30, ChanList: &LinkedList{Lock: &sync.Mutex{}}}
	go func() {
		log.Println(http.ListenAndServe("localhost:18080", nil))
	}()
	qps := 2000
	ch := make(chan int, qps)
	go func() {
		for {
		lab:
			i := 0
			select {
			case <-time.After(time.Second):
				for {
					select {
					case <-ch:
						i++
						if i == qps {
							goto lab
						}
					}
				}
			}
		}
	}()
	for j := 0; j < 1000; j++ {
		go func() {
			for true {
				ch <- 1
				ex := time.Now().Add(time.Duration(30) * time.Millisecond)
				queue.Add(&ExpireElement{ExpireAt: &ex})
				atomic.AddInt32(&added, 1)
			}
		}()
	}

	for j := 0; j < 10; j++ {
		go func() {
			for true {
				e := time.Duration(4000) * time.Millisecond
				q := queue.Poll(&e)
				if q != nil && !q.Expired() {
					atomic.AddInt32(&expired, 1)
				}
				if q!=nil{
					atomic.AddInt32(&count, 1)
				}
			}
		}()
	}
	time.Sleep(time.Duration(2) * time.Minute)
	fmt.Println(count)
	fmt.Println(expired)
	fmt.Println(added)
}

type Integer struct {
	Value int
}

func (v *Integer) GetExpireAt() *time.Time {
	return nil
}
func (v *Integer) Expired() bool {
	return true
}
func (v *Integer) Compare(c CompareAble) int {
	e, ok := c.(*Integer)
	if !ok {
		panic(ok)
	}
	return v.Value - e.Value
}
func TestUnboundDelayQueue_Poll(t *testing.T) {
	queue := &UnboundDelayQueue{Lock: &sync.Mutex{}, Threshold: 30, ChanList: &LinkedList{Lock: &sync.Mutex{}}}
	limit := 1000
	for i := limit; i > 0; i-- {
		queue.Add(&Integer{Value: rand.Int()})
	}
	for i := 0; i < limit; i++ {
		c := queue.poll()
		if c == nil {
			panic("")
			continue
		}
		e, ok := c.(*Integer)
		if !ok {
			panic(ok)
		}
		fmt.Println(e.Value)
	}
}

/*func TestUnboundDelayQueue_PollV2(t *testing.T) {
	runtime.GOMAXPROCS(2)
	a:=time.Duration(100) * time.Millisecond
	b := time.NewTimer(a)

	select {
	case <-b.C:
		fmt.Println("ok")
	}

	//var count int32
	var expired int32
	queue := &UnboundDelayQueue{Ch:make(chan int,1),Lock: &sync.Mutex{}, Threshold: 30, ChanList: &LinkedList{Lock: &sync.Mutex{}}}
	go func() {
		log.Println(http.ListenAndServe("localhost:18080", nil))
	}()
	for j := 0; j < 100; j++ {
		go func() {
			for true {
				ex := time.Now().Add(time.Duration(300) * time.Millisecond)
				queue.AddV2(&ExpireElement{ExpireAt: &ex})
				time.Sleep(time.Duration(100) * time.Millisecond)
			}
		}()
	}

	for j := 0; j < 10; j++ {
		go func() {
			for true {
				e := time.Duration(100) * time.Millisecond
				fmt.Println(queue.PollV2(&e))
			}
		}()
	}
	time.Sleep(time.Duration(2) * time.Minute)
	fmt.Println(expired)
}*/
