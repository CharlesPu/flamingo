package workqueue

type (
	Queue interface {
		Add(interface{})
		Get() interface{}
		Len() int
	}
)
