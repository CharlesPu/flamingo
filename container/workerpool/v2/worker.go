package workerpoolv2

import (
	"time"

	"github.com/CharlesPu/flamingo/plog"
)

type (
	worker struct {
		pool *workerPool

		recycleTime time.Time
		taskCh      chan Task
	}
)

func (w *worker) run() {
	defer func() {
		if err := recover(); err != nil {
			plog.Errorf("[worker pool] get a panic: %+v", err)
		}
	}()

	for t := range w.taskCh {
		if t == nil { // quit
			plog.Infof("[worker pool] receive a sig to quit")
			return
		}
		t()

		w.pool.recycle(w)
	}
}

func (wp *workerPool) recycle(w *worker) {
	wp.mu.Lock()
	w.recycleTime = time.Now()
	wp.idleWorkers = append(wp.idleWorkers, w)
	wp.mu.Unlock()
}
