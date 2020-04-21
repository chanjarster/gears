package concurrent

func Example_lockSemaphore() {
	sema := NewLockSemaphore(100) // new a semaphore with 100 permits
	sema.Acquire()                // will be block if no permits available
	defer sema.Release()          // return permit

	if sema.TryAcquire() { // will return false if no permits available
		defer sema.Release() // return permit
	}
}

func Example_chanSemaphore() {
	sema := NewChanSemaphore(100) // new a semaphore with 100 permits
	sema.Acquire()                // will be block if no permits available
	defer sema.Release()          // return permit

	if sema.TryAcquire() { // will return false if no permits available
		defer sema.Release() // return permit
	}
}
