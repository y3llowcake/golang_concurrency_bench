package main

import (
	"runtime"
	"sync"
	"testing"
)

func BenchmarkMutex(b *testing.B) {
	var mu sync.Mutex
	for n := 0; n < b.N; n++ {
		mu.Lock()
		mu.Unlock()
	}
}

func BenchmarkRWMutex(b *testing.B) {
	var mu sync.RWMutex
	for n := 0; n < b.N; n++ {
		mu.Lock()
		mu.Unlock()
	}
}

func BenchmarkSpinlock(b *testing.B) {
	var mu SpinLock
	for n := 0; n < b.N; n++ {
		mu.Lock()
		mu.Unlock()
	}
}

/*func BenchmarkBufferedChannelWrite(b *testing.B) {
	for n := 0; n < b.N; n++ {
		c := make(chan struct{}, 1)
		c <- struct{}{}
	}
}

func BenchmarkUnbufferedChannelWrite(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		c := make(chan struct{})
		go func() {
			<-c
		}()
		runtime.Gosched()
		b.StartTimer()
		c <- struct{}{}
	}
}

func BenchmarkContendedMutex(b *testing.B) {
	var mu sync.Mutex
	for n := 0; n < b.N; n++ {
		RunConcurrently(b, func() {
			mu.Lock()
			mu.Unlock()
		})
	}
}

func BenchmarkContendedRWMutex_Lock(b *testing.B) {
	var mu sync.RWMutex
	for n := 0; n < b.N; n++ {
		RunConcurrently(b, func() {
			mu.Lock()
			mu.Unlock()
		})
	}
}

func BenchmarkContendedRWMutex_RLock(b *testing.B) {
	var mu sync.RWMutex
	for n := 0; n < b.N; n++ {
		RunConcurrently(b, func() {
			mu.RLock()
			mu.RUnlock()
		})
	}
}

func BenchmarkContendedSpinLock(b *testing.B) {
	var mu SpinLock
	for n := 0; n < b.N; n++ {
		RunConcurrently(b, func() {
			mu.Lock()
			mu.Unlock()
		})
	}
}*/

func BenchmarkContendedBufferedChannel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		b.StopTimer()
		done := make(chan struct{})
		c := make(chan int, CONCURRENCY)
		go func() {
			for <-c != 0 {
			}
			close(done)
		}()
		RunConcurrently(b, func() {
			c <- 1
		})
		close(c)
		<-done
	}
}

const CONCURRENCY = 1024

func RunConcurrently(b *testing.B, f func()) {
	b.StopTimer()
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < CONCURRENCY; i++ {
		wg.Add(1)
		go func() {
			<-start
			f()
			wg.Done()
		}()
	}
	runtime.Gosched()
	b.StartTimer()
	close(start)
	wg.Wait()
}
