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

import (
	"testing"
)

func TestChanSemaphore_Acquire(t *testing.T) {
	s := NewChanSemaphore(1)
	s.Acquire()
	testSemaphoreIntegrity(t, s, 10)

}

func TestChanSemaphore_Release(t *testing.T) {
	s := NewLockSemaphore(1)
	n := 10
	// over release
	for i := 0; i < n; i++ {
		s.Release()
	}
	s.Acquire()
	testSemaphoreIntegrity(t, s, n)
}

func TestLockSemaphore_TryAcquire(t *testing.T) {
	s := NewLockSemaphore(1)
	if got, want := s.TryAcquire(), true; got != want {
		t.Errorf("s.TryAcquire() = %v, want: %v", got, want)
	}
	n := 10
	// all TryAcquire() should fail
	for i := 0; i < n; i++ {
		if got, want := s.TryAcquire(), false; got != want {
			t.Errorf("s.TryAcquire() = %v, want: %v", got, want)
		}
	}
	testSemaphoreIntegrity(t, s, n)
}

func TestLockSemaphore_Acquire(t *testing.T) {
	s := NewLockSemaphore(1)
	s.Acquire()
	testSemaphoreIntegrity(t, s, 10)

}

func TestLockSemaphore_Release(t *testing.T) {
	s := NewLockSemaphore(1)
	n := 10
	// over release
	for i := 0; i < n; i++ {
		s.Release()
	}
	s.Acquire()
	testSemaphoreIntegrity(t, s, n)
}

func TestChanSemaphore_TryAcquire(t *testing.T) {
	s := NewLockSemaphore(1)
	if got, want := s.TryAcquire(), true; got != want {
		t.Errorf("s.TryAcquire() = %v, want: %v", got, want)
	}
	n := 10
	// all TryAcquire() should fail
	for i := 0; i < n; i++ {
		if got, want := s.TryAcquire(), false; got != want {
			t.Fatalf("s.TryAcquire() = %v, want: %v", got, want)
		}
	}
	testSemaphoreIntegrity(t, s, n)
}

// Follow the test pattern introduced in cond_test.go TestCondSignal
func testSemaphoreIntegrity(t *testing.T, s Semaphore, n int) {

	running := make(chan bool, n)
	awake := make(chan bool, n)

	for i := 0; i < n; i++ {
		go func() {
			running <- true
			s.Acquire()
			awake <- true
		}()
	}

	for i := 0; i < n; i++ {
		<-running // Wait for everyone to run.
	}

	for n > 0 {
		select {
		case <-awake:
			t.Fatal("goroutine not blocked")
		default:
		}
		s.Release()
			<-awake // Will deadlock if no goroutine wakes up
		select {
		case <-awake:
			t.Fatal("too many goroutines acquire permits")
		default:
		}
		n--
	}

}
