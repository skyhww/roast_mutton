package core

import (
	"testing"
	"sync"
	"time"
	"fmt"
)

func TestLinkedList_Add(t *testing.T) {
	list := &LinkedList{Lock: &sync.Mutex{}}
	go func() {

	}()
	for i := 0; i < 10; i++ {
		go func() {
			for {
				list.Add(nil)
				fmt.Println("插入成功！")
			}
		}()
	}
	time.Sleep(time.Hour)
}
