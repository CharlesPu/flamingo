package list

import "testing"

func TestSinglyLinkedList(t *testing.T) {
	l := NewSList("2", "3", "4", 10, 11)

	t.Log(l)
	l.Reverse()
	t.Log(l)
}
