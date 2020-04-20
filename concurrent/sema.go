package concurrent

import "sync"

// Semaphore
type Semaphore interface {
	// Acquires a permit from this semaphore, blocking until one is available, or the thread is
	Acquire()
	// Releases a permit, returning it to the semaphore.
	// Over release will not make the permits less than zero.
	Release()
	// Acquires a permit from this semaphore, only if one is available at the time of invocation.
	TryAcquire() bool
}

func NewChanSemaphore(permits int64) *ChanSemaphore {
	return &ChanSemaphore{
		permits: make(chan int, permits),
	}
}

// A Semaphore implementation using channel
type ChanSemaphore struct {
	permits chan int
}

func (s *ChanSemaphore) Acquire() {
	// put into the channel, if the channel is full, will block
	s.permits <- 0
}

func (s *ChanSemaphore) Release() {
	select {
	case <-s.permits:
	default:
		// nothing in the channel
	}
}

func (s *ChanSemaphore) TryAcquire() bool {
	select {
	case s.permits <- 0:
		// channel not full
		return true
	default:
		// channel full
		return false
	}
}

func NewLockSemaphore(permits int64) *LockSemaphore {
	s := &LockSemaphore{
		permits: permits,
	}
	s.notFull = sync.NewCond(&s.lock)
	return s
}

// A Semaphore implementation using Lock
type LockSemaphore struct {
	notFull *sync.Cond
	lock    sync.Mutex
	permits int64
	count   int64
}

func (s *LockSemaphore) Acquire() {
	s.lock.Lock()
	for s.isFull() {
		s.notFull.Wait()
	}
	s.count++
	s.lock.Unlock()
}

func (s *LockSemaphore) Release() {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.count <= 0 {
		return
	}
	s.count--
	s.notFull.Broadcast()
}

func (s *LockSemaphore) TryAcquire() bool {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.isFull() {
		return false
	}
	s.count++
	return true
}

func (s *LockSemaphore) isFull() bool {
	return s.count >= s.permits
}
