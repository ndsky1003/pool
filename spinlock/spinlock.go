package spinlock

import (
	"runtime"
	"sync/atomic"
)

type L uint32

func (sl *L) Lock() {
	for range 200 {
		if atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
			return
		}
		runtime.Gosched()
	}
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		runtime.Gosched()
	}
}

func (sl *L) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}
