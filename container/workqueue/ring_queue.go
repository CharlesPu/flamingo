package workqueue

import "fmt"

type (
	// note: concurrent unsafe
	ringQueue struct {
		size int

		head int // 队头，读指针
		tail int // 队尾，写指针，指向下一个可写的位置
		len  int // len 的存在是为了区别head==tail时到底是队列空还是队列满，还可以采用1.维护bool量判断 2.少用一个cell
		q    []interface{}
	}
)

func NewRingQueue(s int) Queue {
	return &ringQueue{
		size: s,
		q:    make([]interface{}, s),
	}
}
func (r *ringQueue) Add(i interface{}) {
	if r.len >= r.size {
		fmt.Println("[ring queue] warning: ringQueue is full")
		return
	}
	r.q[r.tail] = i
	r.tail = r.nextIdx(r.tail)
	r.len++
}

func (r *ringQueue) Get() interface{} {
	if r.len == 0 {
		return nil
	}
	i := r.q[r.head]
	r.q[r.head] = nil // note: need drop to gc
	r.head = r.nextIdx(r.head)
	r.len--
	return i
}

func (r *ringQueue) Len() int {
	return r.len
}

func (r *ringQueue) nextIdx(now int) int {
	return (now + 1) % r.size
}
