package core

import "time"

type Timer struct {
	Wheel    TimeWheel
	Executor Executor
}

func (timer *Timer) Start() {
	for true {
		//时间轮不是一格一格的去推进，而是推进一个时间段
		m:=timer.Wheel.GetMinTime()
		//任务已经过期了
		if m!=nil&&time.Now().Sub(*m) < 0 {
			timer.Executor.Submit(timer.Wheel.Advance(timer.Wheel.GetPrecision()))
		} else {
			interval:=time.Now().Sub(*m)
			timer.Executor.Submit(timer.Wheel.Advance(&interval))
		}
	}
}

func (timer *Timer) AddTask(entity *TaskEntity) {
	if !timer.Wheel.Add(entity) {
		timer.Executor.Execute(entity)
	}
}
