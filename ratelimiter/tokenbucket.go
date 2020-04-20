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
	"github.com/chanjarster/gears"
	"sync"
	"sync/atomic"
	"time"
)

type TokenBucket interface {
	Interface
}

// 令牌桶
type SyncTokenBucket struct {
	lock               sync.Mutex
	capacity           int   // 桶尺寸
	tokens             int   // 令牌数量
	issueRatePerSecond int   // 每秒签发token数量
	lastIssueTimestamp int64 // 上一次签发的时间戳
	nowFn              gears.NowFunc
}

func NewSyncTokenBucket(capacity, issueRatePerSecond int) TokenBucket {
	return &SyncTokenBucket{
		capacity:           capacity,
		tokens:             capacity,
		issueRatePerSecond: issueRatePerSecond,
		lastIssueTimestamp: gears.SysNow(),
		nowFn:              gears.SysNow,
	}
}

func (t *SyncTokenBucket) Capacity() int {
	return t.capacity
}

func (t *SyncTokenBucket) Acquire() bool {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.issueIfNecessary()

	if t.tokens > 0 {
		t.tokens--
		return true
	}
	return false

}

func (t *SyncTokenBucket) issueIfNecessary() {
	if t.tokens >= t.capacity {
		return
	}
	now := t.nowFn()

	elapse := now - t.lastIssueTimestamp
	delta := elapse / int64(time.Second) * int64(t.issueRatePerSecond)

	if delta == 0 {
		return
	}

	t.tokens += int(delta)
	if t.tokens > t.capacity {
		t.tokens = t.capacity
	}

	t.lastIssueTimestamp = now
}

type AtomicTokenBucket struct {
	capacity           int   // 桶尺寸
	tokens             int64 // 令牌数量
	issueRatePerSecond int   // 每秒签发token数量
	lastIssueTimestamp int64 // 上一次签发的时间戳
	nowFn              gears.NowFunc
}

func NewAtomicTokenBucket(capacity, issueRatePerSecond int) TokenBucket {
	return &AtomicTokenBucket{
		capacity:           capacity,
		tokens:             int64(capacity),
		issueRatePerSecond: issueRatePerSecond,
		lastIssueTimestamp: gears.SysNow(),
		nowFn:              gears.SysNow,
	}
}

func (t *AtomicTokenBucket) Capacity() int {
	return t.capacity
}

func (t *AtomicTokenBucket) Acquire() bool {

	t.issueIfNecessary()

	for {
		if oldTokens := atomic.LoadInt64(&t.tokens); oldTokens > 0 {
			swapped := atomic.CompareAndSwapInt64(&t.tokens, oldTokens, oldTokens-1)
			if swapped {
				break
			}
		} else {
			return false
		}
	}
	return true

}

func (t *AtomicTokenBucket) issueIfNecessary() {

	for {
		oldTokens := atomic.LoadInt64(&t.tokens)
		if oldTokens >= int64(t.capacity) {
			return
		}

		oldLit := atomic.LoadInt64(&t.lastIssueTimestamp)
		now := t.nowFn()

		elapse := now - oldLit
		delta := elapse / int64(time.Second) * int64(t.issueRatePerSecond)

		if delta == 0 {
			return
		}

		if oldTokens+delta > int64(t.capacity) {
			delta = int64(t.capacity) - oldTokens
		}

		s1 := atomic.CompareAndSwapInt64(&t.tokens, oldTokens, oldTokens+delta)
		if !s1 {
			continue
		}

		s2 := atomic.CompareAndSwapInt64(&t.lastIssueTimestamp, oldLit, now)
		if !s2 {
			if s1 {
				atomic.AddInt64(&t.tokens, -delta)
			}
			continue
		}

		break
	}

}
