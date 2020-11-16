package queue

import "testing"

func TestQueue(t *testing.T) {
	s := New(3)

	s.In(11)
	s.In("a")
	s.In("3")
	s.In("3")

	t.Log(s)

	t.Log(s.Out())
	t.Log(s.Out())
	t.Log(s.Out())
	t.Log(s.Out())

	t.Log(s)

	s.In("3")

	t.Log(s)
}
