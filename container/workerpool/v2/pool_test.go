package workerpoolv2

import (
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestWorkerPool(t *testing.T) {
	wp := NewWorkerPool(100)
	t.Log(runtime.NumGoroutine())
	wg := &sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		if err := wp.Go(func() {
			time.Sleep(time.Second)
			wg.Done()
		}); err != nil {
			t.Error(err)
		}
	}
	t.Log(runtime.NumGoroutine())
	time.Sleep(time.Second * 5)
	t.Log(runtime.NumGoroutine())

	wg.Wait()
}

func TestBinarySearchNeedClean(t *testing.T) {
	dura := time.Second
	wp := NewWorkerPool(1, WithIdleDuration(dura)).(*workerPool)
	now := time.Now()
	line := now.Add(-dura)
	type fields struct {
		idleWorkers []*worker
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			fields: fields{
				idleWorkers: []*worker{
					{recycleTime: now.Add(-dura - 100)},
					{recycleTime: now.Add(-dura - 1000)},
					{recycleTime: now.Add(-dura)},
					{recycleTime: now.Add(-dura)},
					{recycleTime: now.Add(-dura + 100)},
					{recycleTime: now.Add(-dura + 1000)},
				},
			},
			want: 3,
		},
		{
			fields: fields{
				idleWorkers: []*worker{
					{recycleTime: now.Add(-dura - 100)},
					{recycleTime: now.Add(-dura - 1000)},
					{recycleTime: now.Add(-dura + 100)},
					{recycleTime: now.Add(-dura + 1000)},
				},
			},
			want: 1,
		},
		{
			fields: fields{
				idleWorkers: []*worker{
					{recycleTime: now.Add(-dura - 100)},
					{recycleTime: now.Add(-dura - 1000)},
					{recycleTime: now.Add(-dura - 2000)},
					{recycleTime: now.Add(-dura - 5000)},
				},
			},
			want: 3,
		},
		{
			fields: fields{
				idleWorkers: []*worker{
					{recycleTime: now.Add(-dura + 100)},
					{recycleTime: now.Add(-dura + 1000)},
					{recycleTime: now.Add(-dura + 2000)},
					{recycleTime: now.Add(-dura + 5000)},
				},
			},
			want: -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wp.idleWorkers = tt.fields.idleWorkers
			if got := wp.binarySearchNeedClean(line); got != tt.want {
				t.Errorf("binarySearchNeedClean() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWorkerPool_Go_NotEnoughWorker(t *testing.T) {
	poolSize := 1

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize)

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	taskNum := 100000
	var success int
	for i := 0; i < taskNum; i++ {
		ii := i
		if err := p.Go(func() {
			fmt.Printf("I(%d) am running\n", ii)
		}); err == nil {
			success++
		}
	}
	t.Logf("success run task num %+v", success)
	time.Sleep(time.Second * 2) // wait exec finished
	if success >= taskNum {
		t.Fatalf("success num is wrong: %+v", success)
	}

	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumActive() {
		t.Fatalf("goroutine num is wrong: %d != %d", n, startGoRoutines+1+p.NumActive())
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit

	if p.NumActive() != 0 || p.NumRunning() != 0 {
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumActive(), p.NumRunning())
	}
}

func TestWorkerPool_Go_EnoughWorker(t *testing.T) {
	taskNum := 100000
	poolSize := taskNum

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize)

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	var success int
	wg := &sync.WaitGroup{}
	for i := 0; i < taskNum; i++ {
		wg.Add(1)
		ii := i
		if err := p.Go(func() {
			fmt.Printf("I(%d) am running\n", ii)
			wg.Done()
		}); err == nil {
			success++
		}
	}
	t.Logf("success run task num %+v", success)
	wg.Wait()
	if success != taskNum {
		t.Fatalf("success num is wrong: %+v", success)
	}

	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumActive() {
		t.Fatalf("goroutine num is wrong: %d != %d", n, startGoRoutines+1+p.NumActive())
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit

	if p.NumActive() != 0 || p.NumRunning() != 0 {
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumActive(), p.NumRunning())
	}
}

func TestWorkerPool_Clean(t *testing.T) {
	taskNum := 1000
	poolSize := taskNum
	idle := time.Second * 3

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize, WithIdleDuration(idle))

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	for i := 0; i < taskNum; i++ {
		ii := i
		p.Go(func() {
			fmt.Printf("I(%d) am running\n", ii)
		})
	}
	time.Sleep(2 * idle) // ensure clean really happens
	t.Log(p.NumRunning(), p.NumActive())
	if x, y := p.NumRunning(), runtime.NumGoroutine(); x != 0 || y != startGoRoutines+1 {
		t.Fatalf("wrong worker or goroutine num: %+v, %+v", x, y)
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumActive() != 0 || p.NumRunning() != 0 {
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumActive(), p.NumRunning())
	}
}

func TestWorkerPool_CleanRunningWorker(t *testing.T) {
	taskNum := 1000
	poolSize := taskNum
	idle := time.Second * 2

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize, WithIdleDuration(idle))

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	for i := 0; i < taskNum; i++ {
		ii := i
		p.Go(func() {
			fmt.Printf("I(%d) am running\n", ii)
			<-make(chan struct{})
		})
	}
	time.Sleep(2 * idle) // ensure clean really happens
	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumRunning() {
		t.Fatalf("wrong goroutine num: %+v != %+v", n, startGoRoutines+1+p.NumRunning())
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumRunning() != taskNum {
		t.Fatalf("wrong goroutine num: %+v", p.NumRunning())
	}
}
