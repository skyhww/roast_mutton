package core

import (
	"time"
	"fmt"
)

type CompareAble interface {
	Compare(c CompareAble) int
}
type Expire interface {
	CompareAble
	GetExpireAt() *time.Time
	Expired() bool
}
type ExpireElement struct {
	ExpireAt *time.Time
}

func (element *ExpireElement) Compare(c CompareAble) int {
	e, ok := c.(*ExpireElement)
	if !ok {
		panic(ok)
	}
	return element.ExpireAt.Nanosecond() - e.ExpireAt.Nanosecond()
}
func (element *ExpireElement) GetExpireAt() *time.Time {
	return element.ExpireAt
}
func (element *ExpireElement) Expired() bool {
	return time.Now().Before(*element.ExpireAt)||time.Now().Equal(*element.ExpireAt)
}

func (element *ExpireElement) String() string {
	return fmt.Sprintf("%d", element.ExpireAt.UnixNano())
}
