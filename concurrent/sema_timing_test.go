package concurrent

import "testing"

func BenchmarkChanSemaphore_Acquire(b *testing.B) {
	permits := 100000
	s := NewChanSemaphore(int64(permits))
	benchmarkSemaphoreAcquire(b, s)
}

func BenchmarkLockSemaphore_Acquire(b *testing.B) {
	permits := 100000
	s := NewLockSemaphore(int64(permits))
	benchmarkSemaphoreAcquire(b, s)
}

func BenchmarkChanSemaphore_TryAcquire(b *testing.B) {
	permits := 100000
	s := NewChanSemaphore(int64(permits))
	benchmarkSemaphoreTryAcquire(b, s)
}

func BenchmarkLockSemaphore_TryAcquire(b *testing.B) {
	permits := 100000
	s := NewLockSemaphore(int64(permits))
	benchmarkSemaphoreTryAcquire(b, s)
}


func benchmarkSemaphoreAcquire(b *testing.B, s Semaphore) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Acquire()
			s.Release()
		}
	})
}

func benchmarkSemaphoreTryAcquire(b *testing.B, s Semaphore) {
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.TryAcquire()
			s.Release()
		}
	})
}