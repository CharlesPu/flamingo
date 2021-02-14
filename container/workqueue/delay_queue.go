package workqueue

import (
	"container/heap"
	"sync"
	"time"
)

type (
	DelayQueue interface {
		Queue

		ShutDown()
		AddAfter(interface{}, time.Duration)
	}

	// note: concurrent safe
	delayQueue struct {
		Queue

		cond            *sync.Cond // protect Queue
		stopCh          chan struct{}
		waitingForAddCh chan *waiting
		waitingQueue    waitingQueue // unlimited size
	}
	waiting struct {
		data interface{}

		hopeAt time.Time
	}
	waitingQueue []*waiting
)

func NewDelayQueue(qSize int) DelayQueue {
	return newDelayQueue(NewRingQueue(qSize))
}

func NewDelayQueueWithCustomQueue(q Queue) DelayQueue {
	return newDelayQueue(q)
}

const (
	waitingChanSize = 1024 // todo
)

func newDelayQueue(q Queue) DelayQueue {
	dq := &delayQueue{
		Queue:           q,
		cond:            sync.NewCond(&sync.Mutex{}),
		stopCh:          make(chan struct{}),
		waitingForAddCh: make(chan *waiting, waitingChanSize),
	}
	heap.Init(&dq.waitingQueue)

	go dq.waitingLoop()

	return dq
}

func (dq *delayQueue) ShutDown() {
	close(dq.stopCh) // to stop waitingLoop
}

// may block due to wSize. todo add bool: isFull
func (dq *delayQueue) AddAfter(i interface{}, duration time.Duration) {
	if duration <= 0 {
		dq.Add(i)
		return
	}

	select {
	case <-dq.stopCh:
	case dq.waitingForAddCh <- &waiting{data: i, hopeAt: time.Now().Add(duration)}:
	}
}

func (dq *delayQueue) waitingLoop() {
	waitingItems := &dq.waitingQueue

	for {
		// first pick ready items
		for waitingItems.Len() > 0 {
			item := waitingItems.Peek().(*waiting)
			if item.hopeAt.After(time.Now()) { // means all behind now
				break
			}
			item = heap.Pop(waitingItems).(*waiting)
			dq.Add(item.data)
		}

		// then waiting
		select {
		case <-dq.stopCh:
			return
		case item := <-dq.waitingForAddCh: // get one
			if item.hopeAt.After(time.Now()) {
				heap.Push(waitingItems, item)
			} else {
				dq.Add(item.data)
			}
		default:
		}
	}
}

func (dq *delayQueue) Add(i interface{}) {
	dq.cond.L.Lock()
	dq.Queue.Add(i)
	dq.cond.Signal()
	dq.cond.L.Unlock()
}

// may block
func (dq *delayQueue) Get() interface{} {
	dq.cond.L.Lock()
	for dq.Queue.Len() == 0 {
		dq.cond.Wait()
	}
	r := dq.Queue.Get()
	dq.cond.L.Unlock()
	return r
}

func (dq *delayQueue) Len() int {
	dq.cond.L.Lock()
	r := dq.Queue.Len()
	dq.cond.L.Unlock()
	return r
}

func (wq waitingQueue) Len() int {
	return len(wq)
}

func (wq waitingQueue) Less(i, j int) bool {
	return wq[i].hopeAt.Before(wq[j].hopeAt)
}

func (wq waitingQueue) Swap(i, j int) {
	wq[i], wq[j] = wq[j], wq[i]
}

func (wq *waitingQueue) Push(x interface{}) {
	*wq = append(*wq, x.(*waiting))
}

func (wq *waitingQueue) Pop() interface{} {
	i := (*wq)[len(*wq)-1]
	*wq = (*wq)[:len(*wq)-1]

	return i
}

func (wq waitingQueue) Peek() interface{} {
	return wq[0]
}
