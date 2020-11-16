package queue

type (
	PriorityQueue interface {
		Insert(interface{})
		Peek() interface{}
		Pop() interface{}
	}
	LessFunc         func(a interface{}, b interface{}) bool // whether a < b
	lessFuncAdaptor  func(i, j int) bool                     // whether q[i] < q[j]
	priorityQueueMax struct {
		cap  int
		less lessFuncAdaptor
		q    []interface{}
	}
	priorityQueueMin struct {
	}
)

func exchange(data []interface{}, i, j int) {
	data[i], data[j] = data[j], data[i]
}

func parent(i int) int {
	return i / 2
}
func left(i int) int {
	return i * 2
}

func right(i int) int {
	return i*2 + 1
}

func NewMaxPriorityQueue(size int, cmp LessFunc) PriorityQueue {
	p := &priorityQueueMax{
		cap: size + 1,
		q:   make([]interface{}, 1, size+1),
	}
	p.less = func(i, j int) bool {
		return cmp(p.q[i], p.q[j])
	}

	return p
}

func (mp *priorityQueueMax) Insert(k interface{}) {
	if len(mp.q) >= mp.cap {
		return
	}
	mp.q = append(mp.q, k)
	mp.swim(len(mp.q) - 1)
}

func (mp *priorityQueueMax) Peek() interface{} {
	return mp.q[1]
}

func (mp *priorityQueueMax) Pop() interface{} {
	ret := mp.q[1]
	exchange(mp.q, 0, len(mp.q)-1)
	mp.q = mp.q[:len(mp.q)-1]
	mp.sink(1)
	return ret
}

// 上浮
func (mp *priorityQueueMax) swim(i int) {
	for i > 1 && mp.less(parent(i), i) {
		exchange(mp.q, i, parent(i))
		i = parent(i)
	}
}

// 下沉
func (mp *priorityQueueMax) sink(i int) {
	for left(i) < len(mp.q) {
		max := left(i)
		if mp.less(max, right(i)) {
			max = right(i)
		}
		if mp.less(max, i) {
			break
		}
		exchange(mp.q, i, max)
		i = max
	}
}
