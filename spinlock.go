package main

import (
	"runtime"
	"sync/atomic"
)

type SpinLock uintptr

func (s *SpinLock) Lock() {
	for !atomic.CompareAndSwapUintptr((*uintptr)(s), 0, 1) {
		runtime.Gosched()
	}
}

func (s *SpinLock) Unlock() {
	atomic.StoreUintptr((*uintptr)(s), 0)
}
