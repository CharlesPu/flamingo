package workerpoolv2

import (
	"runtime"
	"sync/atomic"
)

type (
	spinLock uint32
)

// note: non reentrant
func (sl *spinLock) Lock() {
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}
