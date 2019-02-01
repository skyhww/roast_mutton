package syn

import (
	"testing"
	"fmt"
)

func TestSimpleWaitNotify_Wait(t *testing.T) {
	c := make(chan struct{}, 2)
	fmt.Println(len(c))
	c <- struct{}{}
	fmt.Println(len(c))
	<-c
}
