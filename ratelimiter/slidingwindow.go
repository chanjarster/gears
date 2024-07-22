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
	gtime "github.com/chanjarster/gears/util/time"
	"sync"
	"time"
)

// SlidingWindow 滑动时间窗口限流器
//
// 在当前时间往前的 windowSize 范围内，如果请求次数超过 capacity 那么就拒绝请求
//
// 举个具体的例子，当前时间的往前 1分钟 内，如果请求次数超过了 100次 那么就拒绝请求，这样就限定了请求速率恒定在 100次/分钟
type SlidingWindow interface {
	Interface
	ConfigurableWindow
	WindowSize() time.Duration // time range the window look back
}

// NewSyncSlidingWindow New a SyncFixedWindow.
//
//	capacity: window capacity
//	windowSize: time range the window look back
func NewSyncSlidingWindow(capacity int, windowSize time.Duration) *SyncSlidingWindow {
	if capacity < 0 {
		capacity = 0
	}
	return &SyncSlidingWindow{
		capacity:   capacity,
		windowSize: int64(windowSize),
		records:    list.New(),
		nowFn:      gtime.SysNow,
	}
}

type SyncSlidingWindow struct {
	lock       sync.RWMutex
	capacity   int
	windowSize int64
	records    *list.List // 请求记录，其实就是时间戳
	nowFn      gtime.NowFunc
}

func (s *SyncSlidingWindow) Acquire() bool {
	now := s.nowFn().UnixNano()

	s.lock.Lock()
	defer s.lock.Unlock()

	if s.capacity == 0 {
		return true
	}

	if s.records.Len() < s.capacity {
		// 未到容量，可以通过
		s.records.PushBack(now)
		return true
	}

	oldestV := s.records.Front()
	oldest := time.Duration(oldestV.Value.(int64))
	if now-int64(oldest) > s.windowSize {
		// 容量超过，且最远的一次请求超过了windowSize，可以通过
		s.records.Remove(oldestV)
		s.records.PushBack(now)
		return true
	}
	return false
}

func (s *SyncSlidingWindow) Capacity() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.capacity
}

func (s *SyncSlidingWindow) WindowSize() time.Duration {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return time.Duration(s.windowSize)
}

func (s *SyncSlidingWindow) UpdateConfig(capacity int, windowSize time.Duration) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if capacity < 0 {
		capacity = 0
	}

	for s.records.Len() > 0 && s.records.Len() > capacity {
		// 针对缩容的情况
		s.records.Remove(s.records.Front())
	}

	s.capacity = capacity
	s.windowSize = int64(windowSize)

}
