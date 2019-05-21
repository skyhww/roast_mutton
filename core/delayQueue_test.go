package core

import (
	"testing"
	"sync"
	"time"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync/atomic"
)

func TestUnboundDelayQueue(test *testing.T) {
	runtime.GOMAXPROCS(2)
	var count int32
	var expired int32
	queue := &UnboundDelayQueue{Lock: &sync.Mutex{}, Threshold: 30, ChanList: &LinkedList{Lock: &sync.Mutex{}}}
	go func() {
		log.Println(http.ListenAndServe("localhost:18080", nil))
	}()
	for j := 0; j < 1000; j++ {
		go func() {
			for true {
				ex := time.Now().Add(time.Duration(300) * time.Millisecond)
				queue.Add(&ExpireElement{ExpireAt: &ex})
				time.Sleep(time.Duration(200) * time.Millisecond)
			}
		}()
	}

	for j := 0; j < 100; j++ {
		go func() {
			for true {
				e := time.Duration(30) * time.Millisecond
				q := queue.Poll(&e)
				if q != nil && !q.Expired() {
					atomic.AddInt32(&expired,1)
				}
				atomic.AddInt32(&count,1)
			}
		}()
	}
	time.Sleep(time.Duration(2) * time.Minute)
}

type Integer struct {
	Value int
}

func (v *Integer) Compare(c CompareAble) int {
	e, ok := c.(*Integer)
	if !ok {
		panic(ok)
	}
	return v.Value - e.Value
}
func TestUnboundDelayQueue_Poll(t *testing.T) {
	/*queue := &UnboundDelayQueue{Lock: &sync.Mutex{}, Threshold: 30,ChanList:&LinkedList{Lock:&sync.Mutex{}}}
	limit:=1000
	for i:=limit;i>0;i--{
		queue.Add(&Integer{Value:rand.Int()})
	}
	for i:=0;i<limit;i++  {
		c:=queue.poll()
		if c==nil{
			panic("")
			continue
		}
		e, ok := c.(*Integer)
		if !ok {
			panic(ok)
		}
		 fmt.Println(e.Value)
	}*/
}
