package collection

type Collection interface {
	Size() int
	Get(i int) interface{}
	Delete(i int) interface{}
}
