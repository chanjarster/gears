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
	gtime "github.com/chanjarster/gears/time"
	"math"
	"sync"
	"time"
)

// 固定时间窗口限流器
//
// 将时间按照 windowSize 分割一个一个interval，如果当前interval的请求次数超过了 capacity 那么就拒绝请求
//
// 举个具体的例子，当前interval为(当前时间, 当前时间+1分钟]，在这个范围内如果请求次数超过了 100次 那么就拒绝请求
//
// 固定时间窗口的问题在于，如果在前一个interval的结束和后一个interval的开始里请求数达到饱和，那么就会出现超量请求，看下图:
//  |        o | o        |
//  |<- 1min ->|<- 1min ->|
// 图中的两个o代表饱和的请求高峰，这两个峰都为capacity，时间距离很近，小于1分钟，
// 那么在这段时间里实际发生的请求数量超过了 capacity
//
// 如果要完美的控制流量，请使用SlidingWindow。
// 不过固定时间窗口也具有它的优势：节省内存，它只需要计数就行了，而不需要记录每次请求的时间戳。
type FixedWindow interface {
	Interface
	WindowSize() time.Duration
}

// New a SyncFixedWindow.
//  capacity: window capacity
//  windowSize: time interval for each window
func NewSyncFixedWindow(capacity int, windowSize time.Duration) *SyncFixedWindow {
	return &SyncFixedWindow{
		capacity:   capacity,
		windowSize: int64(windowSize),
		nowFn:      gtime.SysNow,
	}
}

type SyncFixedWindow struct {
	lock       sync.Mutex
	capacity   int
	windowSize int64
	until      int64 // 该时间窗口能够覆盖到的未来的某个时间
	count      int   // 时间窗口里的请求数量
	nowFn      gtime.NowFunc
}

func (s *SyncFixedWindow) Acquire() bool {
	now := s.nowFn().UnixNano()

	s.lock.Lock()
	defer s.lock.Unlock()

	if !s.contains(now) {
		s.reset(now)
	}

	if s.count < s.capacity {
		s.count++
		return true
	}
	return false

}

func (s *SyncFixedWindow) Capacity() int {
	return s.capacity
}

func (s *SyncFixedWindow) WindowSize() time.Duration {
	return time.Duration(s.windowSize)
}

// ts 是否 <= s.until
func (s *SyncFixedWindow) contains(ts int64) bool {
	return s.until >= ts
}

// 按照 s.windowSize 为单位，扩展 s.until 到能够包含 ts
func (s *SyncFixedWindow) reset(ts int64) {
	if s.contains(ts) {
		return
	}
	s.until += int64(math.Ceil(float64(ts-s.until)/float64(s.windowSize))) * s.windowSize
	s.count = 0
}
