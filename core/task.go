package core

import "time"

type TaskEntity struct {
	ExpireAt     *time.Time
	Task         Task
	ErrorHandler ErrorHandler
}
type Task interface {
	//执行任务
	Execute() error
	//终止任务
	Abort() error
}

type ErrorHandler interface {
	Handle()
}
