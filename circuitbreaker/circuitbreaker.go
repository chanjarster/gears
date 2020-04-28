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

// A circuit breaker has 3 states:
//
// open: no task will be executed
//
// closed: task will be executed
//
// halfOpen: task will be executed, if failed, back to state of open,
// if success, transfer to state of close
type Interface interface {
	// Do task
	//  task: task to be done
	//  onError: handle errors returned from task
	//  onOpen: be called if state == open
	Do(task func() error, onError func(error), onOpen func())
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

// failureThreshold: once the failures reach this threshold, the circuit breaker will be opened
//
// resetTimeout: if circuit breaker is in state of open,
// after this timeout, the circuit breaker will be in half-open state
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

func (s *SyncCircuitBreaker) Do(task func() error, onError func(error), onOpen func()) {
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

// A circuit breaker would never be opened
var NeverOpen neverOpen

type neverOpen struct{}

func (n neverOpen) Do(task func() error, onError func(error), onOpen func()) {
	if err := task(); err != nil {
		onError(err)
	}
}
