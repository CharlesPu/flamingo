package workerpool

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"
)

type (
	worker struct {
		name        string // only a id, not a uuid
		terminateCh chan sig
		pool        *workerPool

		lastActiveTime atomic.Value
		state          int32
		taskFunc       chan Task
	}
)

const (
	workerDown = iota
	workerIdle
	workerRunning
)

func (w *worker) work(selectCh chan<- *worker) {
	atomic.StoreInt32(&w.state, workerIdle)
	atomic.AddInt64(&w.pool.activeNum, 1)
	w.refreshActiveTime()

	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Printf("[worker pool] worker %v got a panic: %+v\n", w, err)
			}
			fmt.Printf("[worker pool] worker %+v receive terminate signal to quit\n", w)

			atomic.StoreInt32(&w.state, workerDown)
			w.pool.workerPool.Put(w)
			atomic.AddInt64(&w.pool.activeNum, -1)
		}()

		for {
			select {
			case <-w.terminateCh:
				return
			case selectCh <- w: // compete success, otherwise waiting task
				select {
				case <-w.terminateCh:
					return
				case t := <-w.taskFunc:
					fmt.Printf("[worker pool] worker %+v get a task to run\n", w)
					w.refreshActiveTime() // refresh state
					atomic.StoreInt32(&w.state, workerRunning)
					t()
				}
				atomic.StoreInt32(&w.state, workerIdle)
			}
		}
	}()
}

func (w *worker) refreshActiveTime() {
	w.lastActiveTime.Store(time.Now())
}

func (w *worker) getLastActiveTime() time.Time {
	return w.lastActiveTime.Load().(time.Time)
}

func (w *worker) isIdle() bool {
	return atomic.LoadInt32(&w.state) == workerIdle
}

func (w worker) String() string {
	var strBuf bytes.Buffer

	strBuf.WriteString("&{\"")
	strBuf.WriteString(w.name)
	strBuf.WriteString(fmt.Sprintf("\", %v", w.getLastActiveTime().Format("20060102-150405.999999999")))
	strBuf.WriteString("}")

	return strBuf.String()
}
