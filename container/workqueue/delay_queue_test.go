package workqueue

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestDelayQueue_AddAfter(t *testing.T) {
	dq := NewDelayQueue(10)

	tests := []struct {
		item  interface{}
		delay time.Duration
	}{
		{
			item:  1,
			delay: time.Second,
		},
		{
			item:  2,
			delay: time.Second * 2,
		},
		{
			item:  3,
			delay: time.Second * 3,
		},
	}

	for _, tc := range tests {
		t.Run("", func(t *testing.T) {
			start := time.Now()
			dq.AddAfter(tc.item, tc.delay)

			res := make(chan interface{})
			go func(cc chan interface{}) {
				cc <- dq.Get()
			}(res)
			tm := time.NewTimer(tc.delay)
			select {
			case <-tm.C:
				t.Fatalf("fail timeout, cost %+v", time.Now().Sub(start))
			case item := <-res:
				t.Logf("get item: %+v, cost: %+v", item, time.Now().Sub(start))
				if !reflect.DeepEqual(item, tc.item) {
					t.Fatalf("fail not qual")
				}
			}
		})
	}
}

func TestDelayQueue_Concurrent(t *testing.T) {
	dq := NewDelayQueue(200)

	for i := 0; i < 100; i++ {
		go func(item int) {
			dq.Add(item)
			dq.AddAfter(item+100, time.Second*time.Duration(item%3))
		}(i)
	}

	result := make(chan interface{}, 200)
	wg := &sync.WaitGroup{}
	for i := 0; i < 200; i++ {
		go func() {
			wg.Add(1)
			result <- dq.Get()
			wg.Done()
		}()
	}
	wg.Wait()

	if !reflect.DeepEqual(len(result), 200) || !reflect.DeepEqual(dq.Len(), 0) {
		t.Fatalf("wrong results len(%+v) or queue len(%+v)", len(result), dq.Len())
	}
	close(result)

	cmp := make(map[interface{}]struct{}, len(result))
	for v := range result {
		t.Log(v)
		if _, exist := cmp[v]; exist {
			t.Fatalf("wrong result: %+v", v)
		}
		cmp[v] = struct{}{}
	}
}
