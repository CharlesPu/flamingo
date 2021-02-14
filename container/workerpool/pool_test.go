package workerpool

import (
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPool_TryDo_NotEnoughWorker(t *testing.T) {
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
			if p.TryDo(func() interface{} {
				fmt.Printf("I(%d) am running\n", ii)
				return r
			}, res) {
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

	if p.(*workerPool).activeNum > poolSize {
		t.Fatalf("pool active worker num is wrong: %+v", p.(*workerPool).activeNum)
	}
	if n := runtime.NumGoroutine(); n != startGoRoutines+1+p.(*workerPool).activeNum {
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
	if p.NumRunning() != 0 || runtime.NumGoroutine() != startGoRoutines {
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumRunning(), runtime.NumGoroutine())
	}
}

func TestWorkerPool_TryDo_EnoughWorker(t *testing.T) {
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
			if p.TryDo(func() interface{} {
				fmt.Printf("I(%d) am running\n", ii)
				return r
			}, res) {
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
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumRunning(), runtime.NumGoroutine())
	}
}

func TestWorkerPool_Clean(t *testing.T) {
	taskNum := 1000
	poolSize := taskNum

	startGoRoutines := runtime.NumGoroutine()
	p := NewWorkerPool(poolSize)

	if n := runtime.NumGoroutine(); n != startGoRoutines+1 {
		t.Fatalf("goroutine num is wrong: %+v", n)
	}

	res := make(chan interface{}, taskNum)
	for i := 0; i < taskNum; i++ {
		f := func(ii int) {
			r := fmt.Errorf("error %d", ii)
			p.TryDo(func() interface{} {
				fmt.Printf("I(%d) am running\n", ii)
				return r
			}, res)
		}
		f(i) // hold i
	}
	time.Sleep(2 * idleDuration) // ensure clean really happens
	if x, y := p.NumRunning(), runtime.NumGoroutine(); x != 0 || y != startGoRoutines+1 {
		t.Fatalf("wrong worker or goroutine num: %+v, %+v", x, y)
	}

	p.Shutdown()
	time.Sleep(time.Second * 1) // wait all worker to quit
	t.Log(p.NumRunning(), runtime.NumGoroutine())
	if p.NumRunning() != 0 || runtime.NumGoroutine() != startGoRoutines {
		t.Fatalf("wrong goroutine num: %+v, %+v", p.NumRunning(), runtime.NumGoroutine())
	}
}

// todo
func TestWorkerPool(t *testing.T) {
	runtime.GOMAXPROCS(1)

	wg := &sync.WaitGroup{}
	taskNum := 1000000
	p := NewWorkerPool(taskNum)
	s := time.Now()
	res := make(chan interface{}, taskNum)
	var success int64
	for i := 0; i < taskNum; i++ {
		wg.Add(1)

		f := func(ii int) {
			if p.TryDo(func() interface{} {
				return nil
			}, res) {
				atomic.AddInt64(&success, 1)
			}
			wg.Done()
		}
		f(i)
	}

	wg.Wait()

	t.Logf("success run task num %+v", success)
	t.Logf("aaa %+v", time.Now().Sub(s))
}

func TestWorkerPool_Go(t *testing.T) {
	runtime.GOMAXPROCS(1)
	s := time.Now()
	wg := &sync.WaitGroup{}
	taskNum := 1000000

	res := make(chan interface{}, taskNum)
	var success int64
	for i := 0; i < taskNum; i++ {
		wg.Add(1)

		f := func(ii int) interface{} {
			atomic.AddInt64(&success, 1)
			return nil
		}
		go func(ii int) {
			res <- f(ii) // hold i
			wg.Done()
		}(i)
	}
	wg.Wait()

	t.Logf("success run task num %+v", success)
	t.Logf("aaa %+v", time.Now().Sub(s))
}
