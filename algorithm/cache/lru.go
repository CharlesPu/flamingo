package cache

import (
	"container/list"
	"fmt"
)

// LRU
type LRU interface {
	Put(k interface{}, v interface{})
	Get(k interface{}) (v interface{}, exist bool)
	Len() int
	Reset()

	removeOldest()
}

type (
	lru struct {
		ll       *list.List
		size     int
		cacheMap map[interface{}]*list.Element
	}
	entry struct {
		key   interface{}
		value interface{}
	}
)

func NewLRU(size int) LRU {
	return &lru{
		size:     size,
		ll:       list.New(),
		cacheMap: make(map[interface{}]*list.Element, size),
	}
}

func (l *lru) Put(key interface{}, value interface{}) {
	if e, ok := l.cacheMap[key]; ok {
		l.ll.MoveToFront(e)
		e.Value.(*entry).value = value
		return
	}
	l.cacheMap[key] = l.ll.PushFront(&entry{key: key, value: value})
	if l.ll.Len() > l.size {
		l.removeOldest()
	}
}

func (l *lru) Get(key interface{}) (value interface{}, exist bool) {
	if e, hit := l.cacheMap[key]; hit {
		l.ll.MoveToFront(e)
		return e.Value.(*entry).value, hit
	}
	return
}

func (l *lru) Len() int {
	return l.ll.Len()
}

func (l *lru) Reset() {
	l.ll.Init()
	l.cacheMap = make(map[interface{}]*list.Element, l.size)
}

func (l *lru) removeOldest() {
	e := l.ll.Back()
	if e != nil {
		l.ll.Remove(e)
		delete(l.cacheMap, e.Value.(*entry).key)
	}
}

func (l lru) String() string {
	var res []interface{}

	p := l.ll.Front()
	for p != nil {
		res = append(res, p.Value.(*entry).value)
		p = p.Next()
	}

	return fmt.Sprintf("%+v", res)
}
