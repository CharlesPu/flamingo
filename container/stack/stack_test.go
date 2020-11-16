package stack

import "testing"

func TestStack(t *testing.T) {
	s := New(3)

	s.Push(11)
	s.Push("a")
	s.Push("3")

	t.Log(s)

	t.Log(s.Pop())
	t.Log(s.Pop())
	t.Log(s.Pop())
	t.Log(s.Pop())

	t.Log(s)

	s.Push("3")

	t.Log(s)
}

func TestStackInt(t *testing.T) {
	s := New(3)

	s.PushInt(11)
	s.PushInt(1)
	s.PushInt(3)

	t.Log(s)

	t.Log(s.PopInt())
	t.Log(s.PopInt())
	t.Log(s.PopInt())
	t.Log(s.PopInt())

	t.Log(s)

	s.PushInt(2)

	t.Log(s)
}
