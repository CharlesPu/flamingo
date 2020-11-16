package queue

import (
	"testing"
)

func TestNewMaxPriorityQueue(t *testing.T) {
	q := NewMaxPriorityQueue(11, func(i interface{}, j interface{}) bool {
		return i.(int) < j.(int)
	})
	for i := 0; i < 10; i++ {
		if 10-i == 8 {
			continue
		}
		q.Insert(10 - i)
	}
	q.Insert(8)

	t.Log(q.(*priorityQueueMax).q)
}
