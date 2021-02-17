package workerpool

import (
	"bytes"
	"fmt"
	"sync/atomic"
	"time"
)

type (
	worker struct {
		name        string // only a name, not a uuid
		terminateCh chan sig
		pool        *workerPool

		lastActiveTime atomic.Value
		taskFunc       chan Task
	}
)

func (w *worker) work(selectCh chan *worker) {
	atomic.AddInt64(&w.pool.activeNum, 1)
	w.refreshActiveTime()

	go func() {
		defer func() {
			w.pool.workerPool.Put(w)
			atomic.AddInt64(&w.pool.activeNum, -1)
			if err := recover(); err != nil {
				fmt.Printf("worker %v got a panic: %+v\n", w, err)
			}
			fmt.Printf("worker %+v receive terminate signal to quit\n", w)
		}()

		for {
			select {
			case <-w.terminateCh:
				return
			case selectCh <- w: // compete success, otherwise waiting task
				w.refreshActiveTime() // refresh state
				select {
				case <-w.terminateCh:
					return
				case t := <-w.taskFunc:
					fmt.Printf("worker %+v get a task to run\n", w)
					t()
				}
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

func (w worker) String() string {
	var strBuf bytes.Buffer

	strBuf.WriteString("&{\"")
	strBuf.WriteString(w.name)
	strBuf.WriteString(fmt.Sprintf("\", %v", w.getLastActiveTime().Format("20060102-150405.999999999")))
	strBuf.WriteString("}")

	return strBuf.String()
}
