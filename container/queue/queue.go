package queue

import (
	"bytes"
	"container/list"
	"fmt"
)

type Queue struct {
	l    *list.List
	size int
}

func (q *Queue) Init(size uint32) *Queue {
	if q == nil {
		return q
	}

	q.l = list.New()
	q.size = int(size)
	return q
}

func New(size uint32) *Queue {
	return new(Queue).Init(size)
}

func (q *Queue) In(e interface{}) {
	if q.Full() {
		return
	}

	q.l.PushBack(e)
}

func (q *Queue) Out() interface{} {
	if q.Empty() {
		return nil
	}

	return q.l.Remove(q.l.Front())
}

func (q *Queue) Empty() bool {
	return q.l.Len() == 0
}

func (q *Queue) Full() bool {
	return q.l.Len() == q.size
}

func (q *Queue) Len() int {
	return q.l.Len()
}

func (q *Queue) Size() int {
	return q.size
}

func (q *Queue) InInt(e int) {
	if q.Full() {
		return
	}

	q.l.PushBack(e)
}

func (q *Queue) OutInt() int {
	if q.Empty() {
		return 0
	}

	return q.l.Remove(q.l.Front()).(int)
}

func (q *Queue) String() string {
	tmp := list.New()

	tmp.PushBackList(q.l)

	var strBuf bytes.Buffer
	strBuf.WriteString("Q:{[")

	for tmp.Len() != 0 {
		v := tmp.Front()
		strBuf.WriteString(fmt.Sprintf("%+v", v.Value))
		tmp.Remove(v)
		if tmp.Len() != 0 {
			strBuf.WriteString(", ")
		}
	}

	strBuf.WriteString(fmt.Sprintf("], size: %+v}", q.size))

	return strBuf.String()
}
