package workerpool

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type (
	Pool interface {
		Do(f Task, res chan<- interface{})
		TryDo(f Task, res chan<- interface{}) (doing bool)
		DoWait(f Task) interface{}

		NumRunning() int

		Shutdown()
	}
	Task func() interface{}

	workerPool struct {
		workerMaxNum int
		activeNum    int
		mu           *sync.Mutex

		workers    *list.List
		workerPool *sync.Pool

		stopCh   chan sig
		selectCh chan *worker
	}

	sig struct{}
)

func NewWorkerPool(size int) Pool {
	r := &workerPool{
		workerMaxNum: size,
		mu:           &sync.Mutex{},
		workers:      list.New(),
		stopCh:       make(chan sig),
		selectCh:     make(chan *worker),
	}
	r.workerPool = &sync.Pool{
		New: func() interface{} {
			return &worker{
				terminateCh: make(chan sig),
				token:       make(chan *workerToken),
			}
		},
	}

	go r.clean()
	// todo metrics

	return r
}

func (wp *workerPool) Shutdown() {
	close(wp.stopCh)
	wp.mu.Lock()
	for e := wp.workers.Front(); e != nil; e = e.Next() {
		close(e.Value.(*worker).terminateCh)
		wp.workerPool.Put(e.Value.(*worker))
	}
	wp.workers = wp.workers.Init()
	wp.activeNum = 0
	wp.mu.Unlock()
}

// non-block
func (wp *workerPool) Do(f Task, res chan<- interface{}) {
	worker := wp.getWorker()
	if worker == nil {
		return
	}
	worker.token <- &workerToken{
		task: f,
		res:  res,
	}
}

// non-block
func (wp *workerPool) TryDo(f Task, res chan<- interface{}) (doing bool) {
	worker := wp.getWorker()
	if worker == nil {
		return false
	}
	worker.token <- &workerToken{
		task: f,
		res:  res,
	}
	return true
}

// may block
func (wp *workerPool) DoWait(f Task) interface{} {
	worker := wp.getWorker()
	if worker == nil {
		fmt.Println("no worker to run")
		return nil
	}
	res := make(chan interface{})
	worker.token <- &workerToken{
		task: f,
		res:  res,
	}
	return <-res
}

func (wp *workerPool) NumRunning() int {
	wp.mu.Lock()
	r := wp.activeNum
	wp.mu.Unlock()
	return r
}

const (
	idleDuration = time.Second * 5
)

// release least active worker
func (wp *workerPool) clean() {
	defer func() {
		fmt.Println("quit clean")
	}()
	tk := time.NewTicker(idleDuration)
	for {
		select {
		case <-wp.stopCh:
			return
		case <-tk.C:
			fmt.Println("start to clean idle workers")
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
				wp.workerPool.Put(w)
				wp.activeNum--
				fmt.Printf("release worker: %+v\n", w)
			}
			wp.mu.Unlock()
		}
	}
}

func (wp *workerPool) getWorker() *worker {
	select {
	case <-wp.stopCh:
		return nil
	case w := <-wp.selectCh: // get a active worker
		fmt.Printf("get a ACTIVE worker: %+v\n", w)
		return w
	default: // try to create new worker
		wp.mu.Lock()
		if wp.activeNum >= wp.workerMaxNum {
			fmt.Printf("worker pool is full(%+v), no active worker\n", wp.activeNum)
			wp.mu.Unlock()
			return nil
		}
		// create
		w := wp.workerPool.Get().(*worker)
		//w := &worker{
		//	terminateCh: make(chan sig),
		//	token:       make(chan *workerToken),
		//}
		w.name = fmt.Sprintf("%d", time.Now().UnixNano())
		go w.work(wp.selectCh)
		fmt.Printf("create a NEW worker: %+v\n", w.name)
		wp.workers.PushBack(w)
		wp.activeNum++
		wp.mu.Unlock()
		return wp.getWorker() // after create, try to catch it, but may caught by others
	}
}
