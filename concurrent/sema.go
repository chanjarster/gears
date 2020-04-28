/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package concurrent

import "sync"

// Semaphore
type Semaphore interface {
	// Acquire a permit from this semaphore, blocking until one is available.
	Acquire()
	// Return a permit to the semaphore.
	// Over release will not make the available permits exceeds capacity.
	Release()
	// Acquire a permit from this semaphore, only if one is available at the time of invocation,
	// never blocking.
	TryAcquire() bool
}

// New a ChanSemaphore.
//
// permits: max permits this semaphore could be acquired.
func NewChanSemaphore(permits int64) *ChanSemaphore {
	return &ChanSemaphore{
		permits: make(chan int, permits),
	}
}

// New a LockSemaphore.
//
// permits: max permits this semaphore could be acquired.
func NewLockSemaphore(permits int64) *LockSemaphore {
	s := &LockSemaphore{
		permits: permits,
	}
	s.notFull = sync.NewCond(&s.lock)
	return s
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

// A Semaphore implementation using sync.Mutex
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
