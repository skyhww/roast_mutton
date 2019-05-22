package core

import (
	"sync"
)

type LinkedList struct {
	head *Node
	tail *Node
	Lock *sync.Mutex
}

type Node struct {
	Predecessor *Node
	Successor   *Node
	Data        interface{}
	parent      *LinkedList
}

func (list *LinkedList) GetHead() *Node {
	return list.head
}
func (list *LinkedList) GetTail() *Node {
	return list.tail
}

//O(1)
func (list *LinkedList) Add(data interface{}) *Node {
	list.Lock.Lock()
	node := &Node{parent: list, Data: data}
	if list.tail != nil {
		list.tail.Successor = node
		node.Successor = list.tail
	} else {
		list.head = node
		list.tail = node
	}
	list.Lock.Unlock()
	return node
}

//已增加parent的空间代价换取删除的效率O(1)
func (node *Node) Delete() *Node {
	node.parent.Lock.Lock()
	if node.Predecessor == nil {
		node.parent.head = nil
	} else {
		node.Predecessor.Successor = node.Successor
		if node.Successor != nil {
			node.Successor.Predecessor = node.Predecessor
		}
	}
	node.parent.Lock.Unlock()
	node.parent=nil
	return node
}

func (node *Node) Join(parent *LinkedList) *Node  {
	parent.Lock.Lock()
	if parent.tail != nil {
		parent.tail.Successor = node
		node.Successor = parent.tail
	} else {
		parent.head = node
		parent.tail = node
	}
	parent.Lock.Unlock()
	return node
}