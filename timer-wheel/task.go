package timer_wheel

import "time"

type Task interface {
	CallBack()
	//超时时间
	Expire() time.Time
	Owner() TaskList
}

type TaskList interface {
	AddObject(f func(), expire time.Time) Task
	Add(task Task) Task
	Iterator(func(task Task))
}

type simpleTask struct {
	previous  *simpleTask
	successor *simpleTask
	f         func()
	expire    time.Time
	owner     TaskList
}

func (task *simpleTask) CallBack() {
	task.f()
}
func (task *simpleTask) Expire() time.Time {
	return task.expire
}
func (task *simpleTask) Owner() TaskList {
	return task.owner
}

type SimpleList struct {
	header *simpleTask
}

func (list *SimpleList) Iterator(f func(task Task)) {
	tmp := list.header
	for tmp != nil {
		f(tmp)
		tmp = tmp.successor
	}

}
func (list *SimpleList) Delete(task Task) {
	v, ok := task.(*simpleTask)
	if ok {
		if v.successor != nil {
			v.successor.previous = v.previous
		}
		if v.previous != nil {
			v.previous.successor = v.successor
		}
		v.successor = nil
		v.previous = nil
	}
}
func (list *SimpleList) Add(task Task) Task {
	t := &simpleTask{}
	t.f = task.CallBack
	t.expire = task.Expire()
	t.owner = list
	if list.header == nil {
		list.header = t
	} else {
		list.header.previous = t
		t.successor = list.header
	}
	return task
}

func (list *SimpleList) AddObject(f func(), expire time.Time) Task {
	task := &simpleTask{}
	task.f = f
	task.expire = expire
	task.owner = list
	if list.header == nil {
		list.header = task
	} else {
		list.header.previous = task
		task.successor = list.header
	}
	return task
}
