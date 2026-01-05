package spinlock

import (
	"runtime"
	"sync/atomic"
)

type L uint32

var spins int = 200

func init() {
	spins = runtime.NumCPU() * 50
}
func (sl *L) Lock() {

	for range spins {
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
