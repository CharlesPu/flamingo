package batch

import (
	"sync/atomic"
	"time"

	"github.com/CharlesPu/flamingo/plog"
)

type (
	batchProcessor struct {
		stopChan   chan struct{}
		buffer     chan interface{}
		bufferSize int

		options *options
		state   int32

		timer    *time.Timer
		consumer Consumer
	}

	timerSig struct{}
)

const (
	batchProcRunning = iota
	batchProcClosed
)

func NewBatchProcessor(s int, c Consumer, opts ...OptionFunc) BatchProcessor {
	opt := defaultOptions

	for _, o := range opts {
		o(&opt)
	}
	if opt.batchThreshold == 0 || opt.batchThreshold > s {
		opt.batchThreshold = s
	}

	r := &batchProcessor{
		stopChan:   make(chan struct{}, 1),
		buffer:     make(chan interface{}, s),
		bufferSize: s,
		options:    &opt,
		state:      batchProcRunning,
		timer:      time.NewTimer(opt.interval),

		consumer: c,
	}

	return r
}

func (b *batchProcessor) Put(i interface{}) {
	if b.isClosed() {
		plog.Warnf("[batchProc] shutdown but got a item and force flush")
		b.flush([]interface{}{i})
		return
	}

	b.buffer <- i
}

func (b *batchProcessor) TryPut(i interface{}) bool {
	if b.isClosed() {
		plog.Warnf("[batchProc] shutdown but got a item and force flush")
		b.flush([]interface{}{i})
		return true
	}

	select {
	case b.buffer <- i:
	default:
		plog.Warnf("[batchProc] full buffer!")
		return false
	}
	return true
}

func (b *batchProcessor) Run() {
	go b.notify()
	go b.timerSignal()
}

func (b *batchProcessor) Shutdown() {
	atomic.StoreInt32(&b.state, batchProcClosed)
	b.timer.Stop()
	b.buffer <- nil
	b.stopChan <- struct{}{}
}

func (b *batchProcessor) Num() int {
	return len(b.buffer)
}

func (b *batchProcessor) flush(batch []interface{}) {
	b.timer.Reset(b.options.interval)
	if len(batch) == 0 {
		return
	}
	b.consumer.Consume(batch)
}

func (b *batchProcessor) notify() {
	batch := make([]interface{}, 0, b.options.batchThreshold)
	var itemsWeight int
	flushFunc := func() {
		b.flush(batch)
		for i := range batch {
			batch[i] = nil
		}
		batch = batch[:0]
		itemsWeight = 0
	}

	for obj := range b.buffer {
		if obj == nil {
			plog.Warnf("[batchProc] buffer receive sig to quit and force flush(size %d)", len(batch))
			flushFunc()
			return
		}
		switch obj.(type) {
		case *timerSig:
			plog.Debugf("[batchProc] timeout and force flush(size %d)", len(batch))
			flushFunc()
		default:
			batch = append(batch, obj)
			// cal weight
			if b.options.calculator != nil {
				itemsWeight += b.options.calculator(obj)
				if itemsWeight >= b.options.weightMax {
					plog.Debugf("[batchProc] over weight and flush(size %d, weight %d)", len(batch), itemsWeight)
					flushFunc()
					continue
				}
			}
			if len(batch) >= b.options.batchThreshold {
				plog.Debugf("[batchProc] size out and flush(size %d)", len(batch))
				flushFunc()
			}
		}
	}
}

func (b *batchProcessor) timerSignal() {
	for {
		select {
		case <-b.timer.C:
			select {
			case b.buffer <- new(timerSig):
			case <-b.stopChan:
				return
			}
		case <-b.stopChan:
			return
		}
	}
}

func (b *batchProcessor) isClosed() bool {
	return atomic.LoadInt32(&b.state) == batchProcClosed
}
