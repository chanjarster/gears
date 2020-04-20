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

package ratelimiter

import (
	"container/list"
	"reflect"
	"testing"
	"time"
)

func TestNewSyncSlidingWindow(t *testing.T) {
	cap := 500
	windowSize := 500 * time.Millisecond

	got := NewSyncSlidingWindow(cap, windowSize)

	if got, want := got.WindowSize(), windowSize; got != want {
		t.Errorf("NewSyncSlidingWindow().WindowSize() = %v, want %v", got, want)
	}
	if got, want := got.Capacity(), cap; got != want {
		t.Errorf("NewSyncSlidingWindow().Capacity() = %v, want %v", got, want)
	}
	if got, want := got.records, list.New(); !reflect.DeepEqual(got, want) {
		t.Errorf("NewSyncSlidingWindow().records = %v, want %v", got, want)
	}
	if got := got.nowFn; got == nil {
		t.Errorf("NewSyncSlidingWindow().nowFn is nil, want not nil")
	}

}

func TestSyncSlidingWindow_Acquire(t *testing.T) {

	new := func() *SyncSlidingWindow {
		cap := 10
		windowSize := 50 * time.Millisecond
		return NewSyncSlidingWindow(cap, windowSize)
	}

	t.Run("acquire 10 times", func(t *testing.T) {
		ratelimiter := new()
		for i := 0; i < ratelimiter.Capacity(); i++ {
			if got := ratelimiter.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
	})

	t.Run("acquire 11 times", func(t *testing.T) {
		ratelimiter := new()
		for i := 0; i < ratelimiter.Capacity(); i++ {
			if got := ratelimiter.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
		if got := ratelimiter.Acquire(); got {
			t.Errorf("acquire() = %v, want %v", got, false)
		}
	})

	// 因为是sliding window，所以只需要等60ms，就可以有新的请求通过
	t.Run("acquire 10 times, wait 60ms", func(t *testing.T) {
		ratelimiter := new()
		for i := 0; i < ratelimiter.Capacity(); i++ {
			if got := ratelimiter.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
		ratelimiter.nowFn = now(time.Now().Add(60 * time.Millisecond).UnixNano())
		if got := ratelimiter.Acquire(); !got {
			t.Errorf("acquire() = %v, want %v", got, true)
		}
	})

	t.Run("acquire for 60ms", func(t *testing.T) {
		// 连续不停的Acquire，肯定会经历从 可拿->不可拿->可拿->不可拿 的过程
		// |<- 50ms ->|<- 50ms ->|
		ratelimiter := new()

		prev := false
		history := make([]bool, 0, 4)
		timer := time.NewTimer(time.Millisecond * 60)
		for {
			fin := false
			select {
			case <-timer.C:
				timer.Stop()
				fin = true
			default:
				acquired := ratelimiter.Acquire()
				if prev != acquired {
					history = append(history, acquired)
				}
				prev = acquired
			}
			if fin {
				break
			}
		}
		expected := []bool{true, false, true, false}
		if !reflect.DeepEqual(history[:4], expected) {
			t.Errorf("history[:4] = %v, want %v", history, expected)
		}
	})

}
