package timer_wheel

import (
	"testing"
	"time"
	"fmt"
)

func TestUnboundDelayQueue_Put(t *testing.T) {
	queue := NewUnboundDelayQueue(10)
	now := time.Now()
	for i := 1; i <= 1000; i++ {
		queue.Put(&Bucket{now.Add(time.Duration(i+1) * time.Second), nil})
	}
	for  i := 1; i <= 1000; i++ {
		fmt.Print(queue.Take(time.Second))
	}
}
func TestUnboundDelayQueue_Take(t *testing.T) {

}
