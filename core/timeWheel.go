package core

import "time"

type TimeWheel interface {
	//添加任务
	Add(entity *TaskEntity) bool
	//向前移动,并返回已过期的任务
	Advance(duration *time.Duration) []*TaskEntity
	//获取精度
	GetPrecision() *time.Duration
	//获取轮子的大小
	GetSize() int64
	//获取最小的过期时间
	GetMinTime() *time.Time
}

type TimeWheelFactory interface {
	CreateTimeWheel(duration *time.Duration, size int64) TimeWheel
}

type DelayQueueTimeWheelFactory struct {
}

func (factory *DelayQueueTimeWheelFactory) CreateTimeWheel(duration *time.Duration, size int64) TimeWheel {
	return nil
}

type SimpleTimeWheel struct {
	Size     int64
	Precious *time.Duration
	//当前指向
	Point  int64
	parent *SimpleTimeWheel
}

/*func (wheel *SimpleTimeWheel) Add(entity *TaskEntity) bool {

}
func (wheel *SimpleTimeWheel) Advance(duration *time.Duration) []*TaskEntity {

}
func (wheel *SimpleTimeWheel) GetPrecision() *time.Duration {
	return wheel.Precious
}
func (wheel *SimpleTimeWheel) GetSize() int64 {
	return wheel.Size
}
func (wheel *SimpleTimeWheel) GetMinTime() *time.Time {
	return nil
}
*/