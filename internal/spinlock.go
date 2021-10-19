// Copyright 2019 Andy Pan & Dietoad. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package internal

import (
	"runtime"
	"sync"
	"sync/atomic"
)

type SpinLock uint32

const maxBackoff = 64

func (sl *SpinLock) Lock() {
	backoff,cnt := 1,1
	lockedOSThread:=false
	for !atomic.CompareAndSwapUint32((*uint32)(sl), 0, 1) {
		// Leverage the exponential backoff algorithm, see https://en.wikipedia.org/wiki/Exponential_backoff.
		if cnt==backoff{lockedOSThread=true;runtime.LockOSThread()}
		if cnt < backoff {
			cnt++
			continue
		}
		for i := 0; i < backoff; i++ {
			runtime.Gosched()
		}
		if backoff < maxBackoff {
			backoff <<= 1
		}
	}
	if lockedOSThread{
		runtime.UnlockOSThread()
	}
}

func (sl *SpinLock) Unlock() {
	atomic.StoreUint32((*uint32)(sl), 0)
}

// NewSpinLock instantiates a spin-lock.
func NewSpinLock() sync.Locker {
	return new(SpinLock)
}
