package workerpoolv2

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CharlesPu/flamingo/plog"
)

type (
	Pool interface {
		Go(f Task) error

		NumRunning() int // nums of goroutine that currently run task. = (total active workers)-(idle workers)
		NumActive() int  // nums of goroutine that currently run = (total active workers)

		Shutdown()
	}

	Task func()

	workerPool struct {
		options      *workerPoolOption
		workerMaxNum int
		activeNum    int64

		mu           sync.Locker
		state        int32
		workersCache *sync.Pool
		idleWorkers  []*worker // recycled and can be reused // time old -> time new
	}
)

const (
	poolRunning = iota
	poolClosed
)

var (
	ErrPoolClose = fmt.Errorf("[worker pool err] worker pool is closed")
	ErrOverload  = fmt.Errorf("[worker pool err] not enough goroutines")
)

func NewWorkerPool(size int, opts ...OptionFunc) Pool {
	options := defaultOptions
	for _, o := range opts {
		o(&options)
	}
	p := &workerPool{
		options:      &options,
		workerMaxNum: size,
		mu:           new(spinLock),
		state:        poolRunning,
		idleWorkers:  make([]*worker, 0, size),
	}
	p.workersCache = &sync.Pool{
		New: func() interface{} {
			return &worker{
				pool:   p,
				taskCh: make(chan Task, taskChanCap),
			}
		},
	}

	go p.clean()

	return p
}

var (
	// inspired by fasthttp
	// https://github.com/valyala/fasthttp/blob/master/workerpool.go#L155
	taskChanCap = func() int {
		// Use blocking workerChan if GOMAXPROCS=1.
		// This immediately switches Serve to WorkerFunc, which results
		// in higher performance (under go1.5 at least).
		if runtime.GOMAXPROCS(0) == 1 {
			return 0
		}

		// Use non-blocking workerChan if GOMAXPROCS>1,
		// since otherwise the Serve caller (Acceptor) may lag accepting
		// new connections if WorkerFunc is CPU-bound.
		return 1
	}()
)

func (wp *workerPool) clean() {
	tk := time.NewTicker(wp.options.idleDuration)
	defer tk.Stop()

	for range tk.C {
		if wp.isClosed() {
			plog.Infof("[worker pool] quit clean")
			return
		}

		wp.mu.Lock()

		l := len(wp.idleWorkers)
		plog.Infof("[worker pool] start to recycle %d workers", l)
		if l == 0 {
			wp.mu.Unlock()
			continue
		}
		idx := wp.binarySearchNeedClean(time.Now().Add(-wp.options.idleDuration))
		if idx == -1 {
			wp.mu.Unlock()
			continue
		}
		var workersNeedClean []*worker
		workersNeedClean = append(workersNeedClean, wp.idleWorkers[:idx+1]...)
		m := copy(wp.idleWorkers, wp.idleWorkers[idx+1:]) // move
		for i := m; i < l; i++ {                          // avoid mem leak
			wp.idleWorkers[i] = nil
		}
		wp.idleWorkers = wp.idleWorkers[:m]

		wp.mu.Unlock()

		plog.Infof("[worker pool] try to release %d workers", len(workersNeedClean))
		for _, w := range workersNeedClean { // non-blocking because wp.idleWorkers has quit Task
			w.taskCh <- nil
		}
	}
}

func (wp *workerPool) Go(f Task) error {
	if wp.isClosed() {
		return ErrPoolClose
	}
	w := wp.getWorker()
	if w == nil {
		return ErrOverload
	}
	w.taskCh <- f
	return nil
}

func (wp *workerPool) NumRunning() int {
	wp.mu.Lock()
	r := int(atomic.LoadInt64(&wp.activeNum)) - len(wp.idleWorkers)
	wp.mu.Unlock()
	return r
}

func (wp *workerPool) NumActive() int {
	return int(atomic.LoadInt64(&wp.activeNum))
}

func (wp *workerPool) Shutdown() {
	plog.Infof("[worker pool] shutdown...")
	atomic.StoreInt32(&wp.state, poolClosed)
	wp.mu.Lock()
	for _, w := range wp.idleWorkers {
		w.taskCh <- nil
	}
	wp.idleWorkers = nil
	wp.mu.Unlock()
}

func (wp *workerPool) getWorker() (w *worker) {
	newWorker := func() {
		wp.incWorkerNum()
		w = wp.workersCache.Get().(*worker)
		go func() {
			w.run()
			wp.decWorkerNum()
			wp.workersCache.Put(w)
		}()
	}

	wp.mu.Lock()

	ready := wp.idleWorkers
	if len(ready) == 0 { // try new
		if int(atomic.LoadInt64(&wp.activeNum)) < wp.workerMaxNum { // can new
			plog.Infof("[worker pool] get a new worker")
			newWorker()
		}
	} else { // get one. latest can avoid slice copy
		l := len(wp.idleWorkers)
		w = wp.idleWorkers[l-1]
		wp.idleWorkers[l-1] = nil // avoid mem leak
		wp.idleWorkers = wp.idleWorkers[:l-1]
		plog.Infof("[worker pool] get a recycled worker, reuse it")
	}

	wp.mu.Unlock()
	return
}

func (wp *workerPool) isClosed() bool {
	return atomic.LoadInt32(&wp.state) == poolClosed
}

func (wp *workerPool) incWorkerNum() {
	atomic.AddInt64(&wp.activeNum, 1)
}

func (wp *workerPool) decWorkerNum() {
	atomic.AddInt64(&wp.activeNum, -1)
}

// right bound binary search
// [0, r] need clean
func (wp *workerPool) binarySearchNeedClean(line time.Time) int {
	l, r, mid := 0, len(wp.idleWorkers)-1, 0

	for l <= r {
		mid = l + (r-l)/2
		if line.Before(wp.idleWorkers[mid].recycleTime) {
			r = mid - 1
		} else { // lock right when equal
			l = mid + 1
		}
	}
	// 3 scenes: 1.r=-1 2.r=length-1 3.found
	return r
}
