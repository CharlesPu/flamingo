package pool

import (
	"context"
	"fmt"

	_ "github.com/Jeffail/tunny"
)

type (
	worker struct {
		id      int
		handler HandleFunc

		ctx         context.Context
		competeChan chan *worker
		input       chan inputParams
		output      chan error
	}
	inputParams struct {
		isSync bool
		params interface{}
	}
)

func newWorker(id int, ctx context.Context, c chan *worker, h HandleFunc) *worker {
	w := &worker{
		id:      id,
		handler: h,

		ctx:         ctx,
		competeChan: c,
		input:       make(chan inputParams),
		output:      make(chan error),
	}

	go w.run()

	return w
}

func (w *worker) run() {
	defer func() {
		close(w.output)
		close(w.input)
		fmt.Println("exit worker", w.id)
	}()
	for {
		select {
		case <-w.ctx.Done():
			return
		case w.competeChan <- w: // try to compete previlege
			// compete successfully
			select {
			case in := <-w.input:
				o := w.handler(in.params)
				if in.isSync {
					w.output <- o
				}

				fmt.Println("done process", w.id)
			case <-w.ctx.Done():
				return
			}
		}
	}
}
