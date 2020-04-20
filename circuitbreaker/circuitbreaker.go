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

package circuitbreaker

import (
	"github.com/chanjarster/gears"
	"sync"
	"time"
)

// task func
// return: error
type task func() error

// callback for error task returns
type onError func(error)

// callback for circuit break opened
type onOpen func()

// Following the following design:
// https://martinfowler.com/bliki/CircuitBreaker.html
type CircuitBreaker interface {
	// Do task
	//  task: task should be done
	//  onError: called when task returns error
	//  onOpen: called if state == open
	Do(task task, onError onError, onOpen onOpen)
}

//---------------------------
// 同步的CircuitBreaker
//---------------------------

type state int

const (
	open     state = iota // 断开
	closed   state = iota // 闭合
	halfOpen state = iota // 半开
)

//  failureThreshold: 出现几次错误就进入断开状态
//  resetTimeout: 断开状态持续多长时间，进入半开状态
func NewSyncCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *SyncCircuitBreaker {
	return &SyncCircuitBreaker{
		failureThreshold: failureThreshold,
		resetTimeout:     int64(resetTimeout),
		nowFn:            gears.SysNow,
	}
}

type SyncCircuitBreaker struct {
	lock             sync.RWMutex
	failureThreshold int           // 多少次失败之后就断开
	lastFailureTs    int64         // 上次失败的时间戳
	failures         int           // 失败次数
	resetTimeout     int64         // 当断开之后多久，把断路器重置到half-open状态
	nowFn            gears.NowFunc // 获得当前时间的函数
}

func (s *SyncCircuitBreaker) Do(task task, onError onError, onOpen onOpen) {
	switch s.state() {
	case open:
		onOpen()
	case closed, halfOpen:
		if err := task(); err != nil {
			s.recordFailure()
			onError(err)
		} else {
			s.reset()
		}
	default:
		panic("Unreachable code")
	}
}

func (s *SyncCircuitBreaker) recordFailure() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.failures++
	s.lastFailureTs = s.nowFn()
}

func (s *SyncCircuitBreaker) state() state {
	now := s.nowFn()
	s.lock.RLock()
	defer s.lock.RUnlock()

	if s.failures >= s.failureThreshold && (now-s.lastFailureTs) > s.resetTimeout {
		return halfOpen
	} else if s.failures >= s.failureThreshold {
		return open
	} else {
		return closed
	}
}

func (s *SyncCircuitBreaker) reset() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.failures = 0
	s.lastFailureTs = 0
}

//---------------------------
// 永远不会断开的CircuitBreaker
//---------------------------

// 永远不会断开的CircuitBreaker
var NeverOpen neverOpen

type neverOpen struct{}

func (n neverOpen) Do(task task, onError onError, onOpen onOpen) {
	if err := task(); err != nil {
		onError(err)
	}
}
