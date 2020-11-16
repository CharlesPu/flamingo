package pool

import (
	"context"
	"fmt"
)

type (
	PoolManager interface {
		Do(interface{})
		DoWait(interface{}) error

		Destroy()
	}

	poolManager struct {
		size        int
		competeChan chan *worker
		context     context.Context
		ctxCancel   func()
		workers     []*worker
	}
	HandleFunc func(interface{}) error
)

func NewPoolManager(n int, h HandleFunc) PoolManager {
	p := &poolManager{
		size:        n,
		competeChan: make(chan *worker),
	}

	p.context, p.ctxCancel = context.WithCancel(context.Background())

	for i := 0; i < n; i++ {
		p.workers = append(p.workers, newWorker(i, p.context, p.competeChan, h))
	}

	return p
}

func (p *poolManager) Do(params interface{}) {
	w := <-p.competeChan
	fmt.Println("start to process on worker", w.id)
	w.input <- inputParams{
		isSync: false,
		params: params,
	}
}

func (p *poolManager) DoWait(params interface{}) error {
	w := <-p.competeChan
	fmt.Println("start to process wait on worker", w.id)
	w.input <- inputParams{
		isSync: true,
		params: params,
	}

	return <-w.output
}

func (p *poolManager) Destroy() {
	fmt.Println("destroying all workers in pool...")
	p.ctxCancel()
}
