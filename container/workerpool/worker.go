package workerpool

import (
	"bytes"
	"fmt"
	"time"
)

type (
	worker struct {
		name        string // only a name, not a uuid
		terminateCh chan sig

		lastActiveTime time.Time
		token          chan *workerToken
	}
	workerToken struct {
		task Task
		res  chan<- interface{}
	}
)

func (w *worker) work(selectCh chan *worker) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("worker %v get a panic: %+v\n", w, err)
		}
		fmt.Printf("worker %+v receive terminate signal to quit\n", w)
	}()

	w.lastActiveTime = time.Now()
	for {
		select {
		case <-w.terminateCh:
			return
		case selectCh <- w: // compete success, otherwise waiting task
			w.lastActiveTime = time.Now() // refresh state
			select {
			case <-w.terminateCh:
				return
			case token := <-w.token:
				fmt.Printf("worker %+v get a task to run\n", w)
				res := token.task()
				if token.res != nil {
					token.res <- res
				}
			}
		}
	}
}

func (w worker) String() string {
	var strBuf bytes.Buffer

	strBuf.WriteString("&{\"")
	strBuf.WriteString(w.name)
	strBuf.WriteString(fmt.Sprintf("\", %v", w.lastActiveTime.Format("20060102-150405.999999999")))
	strBuf.WriteString("}")

	return strBuf.String()
}
