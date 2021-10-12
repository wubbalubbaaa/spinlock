package internal

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type spinLock uint32

func (sl *spinLock) Lock() {
	var iter byte
	// 加锁过程就是原子操作看是否能修改锁变量, 不能的话就用runtime.Gosched把当前goroutine从P上放下来,把时间片让给别人
	// ,goroutine还是在就绪队列,等下一次轮到自己
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		// 做一个简单的自旋
		if iter<2{iter++;continue}
		runtime.Gosched()
	}
}

func (sl *spinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

func NewSpinLock() sync.Locker {
	return new(spinLock)
}
