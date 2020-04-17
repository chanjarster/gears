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
	"reflect"
	"testing"
	"time"
)

func TestNewSyncTokenBucket(t *testing.T) {
	var cap int = 10
	var irs int = 20

	bucket := NewSyncTokenBucket(cap, irs).(*SyncTokenBucket)

	if bucket.capacity != cap {
		t.Errorf("NewSyncTokenBucket().capacity = %v, want %v", bucket.capacity, cap)
	}
	if bucket.tokens != cap {
		t.Errorf("NewSyncTokenBucket().tokens = %v, want %v", bucket.tokens, cap)
	}
	if bucket.issueRatePerSecond != irs {
		t.Errorf("NewSyncTokenBucket().issueRatePerSecond = %v, want %v", bucket.issueRatePerSecond, irs)
	}
	if bucket.lastIssueTimestamp <= 0 {
		t.Errorf("NewSyncTokenBucket().lastIssueTimestamp = %v, want > 0", bucket.issueRatePerSecond)
	}

}

func TestSyncTokenBucket_Acquire(t *testing.T) {
	testTokenBucket_Acquire(t, func() (bucket TokenBucket, i int) {
		cap := 10
		return NewSyncTokenBucket(10, 10), cap
	})
}

func TestNewAtomicTokenBucket(t *testing.T) {
	var cap int = 10
	var irs int = 20 // issue rate

	bucket := NewAtomicTokenBucket(cap, irs).(*AtomicTokenBucket)

	if bucket.capacity != cap {
		t.Errorf("NewAtomicTokenBucket().capacity = %v, want %v", bucket.capacity, cap)
	}
	if bucket.tokens != int64(cap) {
		t.Errorf("NewAtomicTokenBucket().tokens = %v, want %v", bucket.tokens, cap)
	}
	if bucket.issueRatePerSecond != irs {
		t.Errorf("NewAtomicTokenBucket().issueRatePerSecond = %v, want %v", bucket.issueRatePerSecond, irs)
	}
	if bucket.lastIssueTimestamp <= 0 {
		t.Errorf("NewAtomicTokenBucket().lastIssueTimestamp = %v, want > 0", bucket.issueRatePerSecond)
	}

}

func TestAtomicTokenBucket_Acquire(t *testing.T) {

	testTokenBucket_Acquire(t, func() (bucket TokenBucket, i int) {
		cap := 10
		return NewAtomicTokenBucket(10, 10), cap
	})
}

func testTokenBucket_Acquire(t *testing.T, new func() (TokenBucket, int)) {

	t.Run("acquire 10 times", func(t *testing.T) {
		bucket, cap := new()
		for i := 0; i < cap; i++ {
			if got := bucket.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
	})

	t.Run("acquire 11 times", func(t *testing.T) {
		bucket, cap := new()
		for i := 0; i < cap; i++ {
			if got := bucket.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
		if got := bucket.Acquire(); got {
			t.Errorf("acquire() = %v, want %v", got, false)
		}
	})

	t.Run("acquire 10  times, wait 1 sec", func(t *testing.T) {
		bucket, cap := new()
		for i := 0; i < cap; i++ {
			if got := bucket.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
		time.Sleep(time.Millisecond * 1100)
		if got := bucket.Acquire(); !got {
			t.Errorf("acquire() = %v, want %v", got, true)
		}
	})

	t.Run("acquire for 2 seconds", func(t *testing.T) {
		// 连续不停的拿Token，肯定会经历从 可拿->不可拿->可拿->不可拿 的过程
		bucket, _ := new()

		prev := false
		history := make([]bool, 0, 4)
		timer := time.NewTimer(time.Millisecond * 1500)
		for {
			fin := false
			select {
			case <-timer.C:
				timer.Stop()
				fin = true
			default:
				acquired := bucket.Acquire()
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
