package cache

import "testing"

func TestLRU(t *testing.T) {
	l := NewLRU(3)
	l.Put(1, 1)
	l.Put(2, 2)
	l.Put(3, 3)
	l.Put(4, 4)

	t.Log(l)
	t.Log(l.Get(2))
	t.Log(l)
}
