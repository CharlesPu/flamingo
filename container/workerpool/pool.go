package workerpool

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/CharlesPu/flamingo/plog"
)

type (
	Pool interface {
		Go(f Task) (doing bool)

		NumRunning() int

		Shutdown()
	}
	Task func()

	workerPool struct {
		options      *workerPoolOption
		workerMaxNum int64
		activeNum    int64
		mu           *sync.Mutex // protect workers list and sync new worker operate

		workers    *list.List
		workerPool *sync.Pool

		state       int32
		selectCh    chan *worker
		stopCleanCh chan sig
	}

	sig struct{}
)

const (
	poolRunning = iota
	poolClosed
)

func NewWorkerPool(size int, opts ...OptionFunc) Pool {
	opt := defaultOptions
	for _, o := range opts {
		o(opt)
	}

	r := &workerPool{
		options:      opt,
		workerMaxNum: int64(size),
		mu:           &sync.Mutex{},
		workers:      list.New(),
		state:        poolRunning,
		selectCh:     make(chan *worker),
		stopCleanCh:  make(chan sig),
	}
	r.workerPool = &sync.Pool{
		New: func() interface{} {
			return &worker{
				pool:        r,
				terminateCh: make(chan sig),
				taskFunc:    make(chan Task),
			}
		},
	}

	go r.clean()
	// todo metrics

	return r
}

// non-block
func (wp *workerPool) Go(f Task) (doing bool) {
	worker := wp.getWorker()
	if worker == nil {
		return false
	}
	worker.taskFunc <- f
	return true
}

func (wp *workerPool) NumRunning() int {
	return int(atomic.LoadInt64(&wp.activeNum))
}

func (wp *workerPool) Shutdown() {
	plog.Infof("[worker pool] shutdown...")
	atomic.StoreInt32(&wp.state, poolClosed)
	close(wp.stopCleanCh)
	wp.mu.Lock()
	for e := wp.workers.Front(); e != nil; e = e.Next() {
		close(e.Value.(*worker).terminateCh)
	}
	wp.workers = wp.workers.Init()
	wp.mu.Unlock()
}

// release least active worker
func (wp *workerPool) clean() {
	idleDuration := wp.options.cleanDuration
	if idleDuration == 0 {
		plog.Infof("[worker pool] disable clean idle workers")
		return
	}

	tk := time.NewTicker(idleDuration)
	defer tk.Stop()

	for {
		select {
		case <-wp.stopCleanCh:
			plog.Infof("[worker pool] quit clean")
			return
		case <-tk.C:
			plog.Infof("[worker pool] start to clean idle workers")
			wp.mu.Lock()
			var needDelete []*list.Element
			now := time.Now()
			for e := wp.workers.Front(); e != nil; e = e.Next() {
				w := e.Value.(*worker)
				if now.Sub(w.getLastActiveTime()) >= idleDuration && w.isIdle() {
					needDelete = append(needDelete, e)
				}
			}
			for _, e := range needDelete {
				w := e.Value.(*worker)
				w.terminateCh <- sig{}
				wp.workers.Remove(e)
				plog.Infof("[worker pool] release worker: %+v", w)
			}
			wp.mu.Unlock()
		}
	}
}

func (wp *workerPool) getWorker() *worker {
	if wp.isClosed() {
		return nil
	}

	select {
	case w := <-wp.selectCh: // get a active worker
		plog.Infof("[worker pool] get a ACTIVE worker: %+v", w)
		return w
	default: // try to create new worker
		wp.mu.Lock()
		if n := atomic.LoadInt64(&wp.activeNum); n >= wp.workerMaxNum {
			plog.Infof("[worker pool] worker pool is full(%+v), no active worker", n)
			wp.mu.Unlock()
			return nil
		}

		w := wp.workerPool.Get().(*worker)
		w.name = fmt.Sprintf("%d", time.Now().UnixNano())
		plog.Infof("[worker pool] create a NEW worker: %+v", w.name)
		w.work(wp.selectCh)
		wp.workers.PushBack(w)
		wp.mu.Unlock()

		return wp.getWorker() // after create, try to catch it, but may caught by others
	}
}

func (wp *workerPool) isClosed() bool {
	return atomic.LoadInt32(&wp.state) == poolClosed
}
