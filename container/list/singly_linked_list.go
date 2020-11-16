package list

import "fmt"

type (
	SList struct {
		root SNode
		len  int
	}
	SNode struct {
		value interface{}
		next  *SNode
	}
)

func NewSList(values ...interface{}) *SList {
	l := &SList{}
	p := &l.root
	for _, v := range values {
		p.next = &SNode{
			value: v,
		}
		p = p.next
		l.len++
	}

	return l
}

func (l *SList) Reverse() {
	if l.root.next == nil {
		return
	}
	p, pNext := l.root.next, l.root.next.next

	for pNext != nil {
		p.next = pNext.next
		pNext.next = l.root.next
		l.root.next = pNext
		pNext = p.next
	}
}

func (l SList) String() string {
	var res []interface{}

	p := l.root.next
	for p != nil {
		res = append(res, p.value)
		p = p.next
	}

	return fmt.Sprintf("{len:%+v, list:%+v}", l.len, res)
}
