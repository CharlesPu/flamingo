package workerpool

import (
	"container/list"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type (
	Pool interface {
		Go(f Task) (doing bool)

		NumRunning() int

		Shutdown()
	}
	Task func()

	workerPool struct {
		workerMaxNum int64
		activeNum    int64
		mu           *sync.Mutex // protect workers list

		workers    *list.List
		workerPool *sync.Pool

		state       int32
		selectCh    chan *worker
		stopCleanCh chan sig
	}

	sig struct{}
)

const (
	running = iota
	closed
)

func NewWorkerPool(size int) Pool {
	r := &workerPool{
		workerMaxNum: int64(size),
		mu:           &sync.Mutex{},
		workers:      list.New(),
		state:        running,
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
	atomic.StoreInt32(&wp.state, closed)
	close(wp.stopCleanCh)
	wp.mu.Lock()
	for e := wp.workers.Front(); e != nil; e = e.Next() {
		close(e.Value.(*worker).terminateCh)
	}
	wp.workers = wp.workers.Init()
	wp.mu.Unlock()
}

const (
	idleDuration = time.Second * 5
)

// release least active worker
func (wp *workerPool) clean() {
	tk := time.NewTicker(idleDuration)

	for {
		select {
		case <-wp.stopCleanCh:
			fmt.Printf("quit clean\n")
			return
		case <-tk.C:
			fmt.Printf("start to clean idle workers\n")
			wp.mu.Lock()
			var needDelete []*list.Element
			for e := wp.workers.Front(); e != nil; e = e.Next() {
				if time.Now().Sub(e.Value.(*worker).lastActiveTime) >= idleDuration {
					needDelete = append(needDelete, e)
				}
			}
			for _, e := range needDelete {
				w := e.Value.(*worker)
				w.terminateCh <- sig{}
				wp.workers.Remove(e)
				fmt.Printf("release worker: %+v\n", w)
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
		fmt.Printf("get a ACTIVE worker: %+v\n", w)
		return w
	default: // try to create new worker
		// create
		wp.mu.Lock()
		if n := atomic.LoadInt64(&wp.activeNum); n >= wp.workerMaxNum {
			fmt.Printf("worker pool is full(%+v), no active worker\n", n)
			wp.mu.Unlock()
			return nil
		}

		w := wp.workerPool.Get().(*worker)
		w.name = fmt.Sprintf("%d", time.Now().UnixNano())
		fmt.Printf("create a NEW worker: %+v\n", w.name)
		w.work(wp.selectCh)
		wp.workers.PushBack(w)
		wp.mu.Unlock()

		return wp.getWorker() // after create, try to catch it, but may caught by others
	}
}

func (wp *workerPool) isClosed() bool {
	return atomic.LoadInt32(&wp.state) == closed
}
