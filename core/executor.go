package core

type Executor interface {
	Execute(taskEntity *TaskEntity)
	Submit(tasks []*TaskEntity)
}