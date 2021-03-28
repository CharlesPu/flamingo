package batch

import (
	"testing"
	"time"
)

func TestBatchProcessor(t *testing.T) {
	threshold := 5
	prodNum := 100

	pc := &producerConsumer{t: t, threshold: threshold, result: make(map[int]struct{})}
	b := NewBatchProcessor(prodNum, pc, WithBatchThresholdStrategy(threshold), WithMaxWaitStrategy(time.Second))
	b.Run()
	defer b.Shutdown()
	pc.bp = b

	// produce
	go func() {
		for i := 0; i < prodNum; i++ {
			if ok := pc.bp.TryPut(i); !ok {
				t.Logf("put %+v not ok", i)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(time.Second)
	pc.check(prodNum)
}

func TestBatchProcessor_Wait(t *testing.T) {
	threshold := 5
	prodNum := 100

	pc := &producerConsumer{t: t, threshold: threshold, result: make(map[int]struct{})}
	b := NewBatchProcessor(prodNum, pc, WithBatchThresholdStrategy(threshold), WithMaxWaitStrategy(time.Second))
	b.Run()
	defer b.Shutdown()
	pc.bp = b

	// produce
	go func() {
		for i := 0; i < prodNum+1; i++ {
			if ok := pc.bp.TryPut(i); !ok {
				t.Logf("put %+v not ok", i)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(time.Second * 2)
	pc.check(prodNum + 1)
}

func TestBatchProcessor_Shutdown(t *testing.T) {
	threshold := 5
	prodNum := 100

	pc := &producerConsumer{t: t, threshold: threshold, result: make(map[int]struct{})}
	b := NewBatchProcessor(prodNum, pc, WithBatchThresholdStrategy(threshold), WithMaxWaitStrategy(time.Second))
	b.Run()
	pc.bp = b

	b.Shutdown()

	if ok := pc.bp.TryPut(0); !ok {
		t.Logf("put not ok")
	}

	time.Sleep(time.Second)
	pc.check(1)
}

func TestBatchProcessor_Weight(t *testing.T) {
	threshold := 4
	prodNum := 100

	pc := &producerConsumer{t: t, threshold: threshold, result: make(map[int]struct{})}
	b := NewBatchProcessor(prodNum, pc,
		WithBatchThresholdStrategy(threshold),
		WithItemsWeightStrategy(func(i interface{}) int {
			return i.(int)
		}, 10))
	b.Run()
	pc.bp = b
	defer b.Shutdown()

	// produce
	go func() {
		for i := 0; i < prodNum; i++ {
			if ok := pc.bp.TryPut(i); !ok {
				t.Logf("put %+v not ok", i)
			}
			time.Sleep(time.Millisecond)
		}
	}()

	time.Sleep(time.Second)
	pc.check(prodNum)
}

type (
	producerConsumer struct {
		bp BatchProcessor
		t  *testing.T

		threshold int
		result    map[int]struct{}
	}
)

func (pc *producerConsumer) Consume(values []interface{}) {
	pc.t.Logf("[%+v] receive len %d", time.Now(), len(values))
	if pc.threshold < len(values) {
		pc.t.FailNow()
	}

	for _, v := range values {
		pc.result[v.(int)] = struct{}{}
	}
}

func (pc *producerConsumer) check(sum int) {
	if len(pc.result) != sum {
		pc.t.FailNow()
	}

	for i := 0; i < sum; i++ {
		if _, exist := pc.result[i]; !exist {
			pc.t.Errorf("%v not exist", i)
		}
	}
}
