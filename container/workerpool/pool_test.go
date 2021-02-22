package workerpool

import (
	"fmt"
	"reflect"
	"runtime"
	"testing"
	"time"
)

func TestWorkerPool_Go_NotEnoughWorker(t *testing.T) {
	poolSize := 1

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize)

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	taskNum := 100000
	res := make(chan interface{}, taskNum)
	resExpect := make(map[string]struct{})
	var success int
	for i := 0; i < taskNum; i++ {
		f := func(ii int) {
			r := fmt.Errorf("error %d", ii)
			if p.Go(func() {
				fmt.Printf("I(%d) am running\n", ii)
				res <- r
			}) {
				success++
				resExpect[r.Error()] = struct{}{}
			}
		}
		f(i) // hold i
	}
	t.Logf("success run task num %+v", success)
	time.Sleep(time.Second * 1) // wait res chan has been written
	if success >= taskNum {
		t.Fatalf("success num is wrong: %+v", success)
	}

	if p.NumRunning() > poolSize {
		t.Fatalf("pool active worker num is wrong: %+v", p.NumRunning())
	}
	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumRunning() {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	if success != len(resExpect) || success != len(res) {
		t.Fatalf("wrong results num: %d, %d, %d", success, len(resExpect), len(res))
	}

	close(res)
	// compare results
	cmp := make(map[string]struct{}, len(res))
	for i := range res {
		cmp[i.(error).Error()] = struct{}{}
	}

	if !reflect.DeepEqual(cmp, resExpect) {
		t.Fatalf("res and resExpect is not equal, \n%+v,\n%+v", cmp, resExpect)
	}

	p.Shutdown()
	time.Sleep(time.Second * 2) // wait all worker to quit
	if p.NumRunning() != 0 || runtime.NumGoroutine() != startGoRoutines {
		t.Fatalf("wrong goroutine num: %+v, %+v, %+v", p.NumRunning(), runtime.NumGoroutine(), startGoRoutines)
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

	res := make(chan interface{}, taskNum)
	resExpect := make(map[string]struct{})
	var success int
	for i := 0; i < taskNum; i++ {
		f := func(ii int) {
			r := fmt.Errorf("error %d", ii)
			if p.Go(func() {
				fmt.Printf("I(%d) am running\n", ii)
				res <- r
			}) {
				success++
				resExpect[r.Error()] = struct{}{}
			}
		}
		f(i) // hold i
	}
	t.Logf("success run task num %+v", success)
	time.Sleep(time.Second * 1) // wait res chan has been written
	if success != taskNum {
		t.Fatalf("success num is wrong: %+v", success)
	}

	if p.NumRunning() > poolSize {
		t.Fatalf("pool active worker num is wrong: %+v", p.NumRunning())
	}
	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumRunning() {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	if success != len(resExpect) || success != len(res) {
		t.Fatalf("wrong results num: %d, %d, %d", success, len(resExpect), len(res))
	}

	close(res)
	// compare results
	cmp := make(map[string]struct{}, len(res))
	for i := range res {
		cmp[i.(error).Error()] = struct{}{}
	}

	if !reflect.DeepEqual(cmp, resExpect) {
		t.Fatalf("res and resExpect is not equal, \n%+v,\n%+v", cmp, resExpect)
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumRunning() != 0 || runtime.NumGoroutine() != startGoRoutines {
		t.Fatalf("wrong goroutine num: %+v, %+v, %+v", p.NumRunning(), runtime.NumGoroutine(), startGoRoutines)
	}
}

func TestWorkerPool_Clean(t *testing.T) {
	taskNum := 1000
	poolSize := taskNum
	idle := time.Second * 3

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize, WithCleanDuration(idle))

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	for i := 0; i < taskNum; i++ {
		f := func(ii int) {
			p.Go(func() {
				fmt.Printf("I(%d) am running\n", ii)
			})
		}
		f(i) // hold i
	}
	time.Sleep(2 * idle) // ensure clean really happens
	if x, y := p.NumRunning(), runtime.NumGoroutine(); x != 0 || y != startGoRoutines+1 {
		t.Fatalf("wrong worker or goroutine num: %+v, %+v", x, y)
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumRunning() != 0 || runtime.NumGoroutine() != startGoRoutines {
		t.Fatalf("wrong goroutine num: %+v, %+v, %+v", p.NumRunning(), runtime.NumGoroutine(), startGoRoutines)
	}
}

func TestWorkerPool_CleanRunningWorker(t *testing.T) {
	taskNum := 1000
	poolSize := taskNum
	idle := time.Second * 2

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize, WithCleanDuration(idle))

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	for i := 0; i < taskNum; i++ {
		f := func(ii int) {
			p.Go(func() {
				fmt.Printf("I(%d) am running\n", ii)
				<-make(chan struct{})
			})
		}
		f(i) // hold i
	}
	time.Sleep(2 * idle) // ensure clean really happens
	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.NumRunning() {
		t.Fatalf("wrong goroutine num: %+v != %+v", n, startGoRoutines+1+p.NumRunning())
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumRunning() == 0 {
		t.Fatalf("wrong goroutine num: %+v", p.NumRunning())
	}
}
