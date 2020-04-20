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
	"errors"
	"testing"
	"time"
)

var (
	nowTs     = time.Now().UnixNano()
	mockNowFn = func() int64 {
		return nowTs
	}
)

func TestNewSyncCircuitBreaker(t *testing.T) {
	breaker := NewSyncCircuitBreaker(5, 10)

	if got, want := breaker.failureThreshold, 5; got != want {
		t.Errorf("failureThreshold = %v, want %v", got, want)
	}
	if got, want := breaker.resetTimeout, int64(10); got != want {
		t.Errorf("resetTimeout = %v, want %v", got, want)
	}
	if got, want := breaker.lastFailureTs, int64(0); got != want {
		t.Errorf("lastFailureTs = %v, want %v", got, want)
	}
	if got, want := breaker.failures, 0; got != want {
		t.Errorf("lastFailureTs = %v, want %v", got, want)
	}
	if got := breaker.nowFn; got == nil {
		t.Errorf("nowFn is nil, want not nil")
	}
}

func TestSyncCircuitBreaker_Do(t *testing.T) {

	t.Run("always success", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs:    0,
			failures:         0,
			failureThreshold: 1,
			resetTimeout:     1000,
			nowFn:            mockNowFn,
		}
		successTimes := 0
		onErrorTimes := 0
		onOpenTimes := 0
		task := func() error {
			successTimes++
			return nil
		}
		onError := func(err error) {
			onErrorTimes++
		}
		onOpen := func() {
			onOpenTimes++
		}
		for i := 0; i < 10; i++ {
			cb.Do(task, onError, onOpen)
		}
		if got, want := successTimes, 10; got != want {
			t.Errorf("successTimes = %v, want %v", got, want)
		}
		if got, want := onErrorTimes, 0; got != want {
			t.Errorf("onErrorTimes = %v, want %v", got, want)
		}
		if got, want := onOpenTimes, 0; got != want {
			t.Errorf("onOpenTimes = %v, want %v", got, want)
		}
	})

	t.Run("success(close), failure(open)", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs:    0,
			failures:         0,
			failureThreshold: 1,
			resetTimeout:     1000,
			nowFn:            mockNowFn,
		}
		successTimes := 0

		onErrorTimes := 0
		onOpenTimes := 0
		successTask := func() error {
			successTimes++
			return nil
		}
		failureTask := func() error {
			return errors.New("on purpose")
		}
		onError := func(err error) {
			onErrorTimes++
		}
		onOpen := func() {
			onOpenTimes++
		}
		cb.Do(successTask, onError, onOpen)
		cb.Do(failureTask, onError, onOpen)

		if got, want := successTimes, 1; got != want {
			t.Errorf("successTimes = %v, want %v", got, want)
		}
		if got, want := onErrorTimes, 1; got != want {
			t.Errorf("onErrorTimes = %v, want %v", got, want)
		}
		if got, want := onOpenTimes, 0; got != want {
			t.Errorf("onOpenTimes = %v, want %v", got, want)
		}

	})

	t.Run("failure(open), success(open)", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs:    0,
			failures:         0,
			failureThreshold: 1,
			resetTimeout:     1000,
			nowFn:            mockNowFn,
		}
		successTimes := 0

		onErrorTimes := 0
		onOpenTimes := 0
		successTask := func() error {
			successTimes++
			return nil
		}
		failureTask := func() error {
			return errors.New("on purpose")
		}
		onError := func(err error) {
			onErrorTimes++
		}
		onOpen := func() {
			onOpenTimes++
		}
		cb.Do(failureTask, onError, onOpen)
		cb.Do(successTask, onError, onOpen)

		if got, want := successTimes, 0; got != want {
			t.Errorf("successTimes = %v, want %v", got, want)
		}
		if got, want := onErrorTimes, 1; got != want {
			t.Errorf("onErrorTimes = %v, want %v", got, want)
		}
		if got, want := onOpenTimes, 1; got != want {
			t.Errorf("onOpenTimes = %v, want %v", got, want)
		}

	})

	t.Run("failure(open), reset timeout, success(close)", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs:    0,
			failures:         0,
			failureThreshold: 1,
			resetTimeout:     10,
			nowFn:            mockNowFn,
		}

		successTimes := 0
		onErrorTimes := 0
		onOpenTimes := 0
		successTask := func() error {
			successTimes++
			return nil
		}
		failureTask := func() error {
			return errors.New("on purpose")
		}
		onError := func(err error) {
			onErrorTimes++
		}
		onOpen := func() {
			onOpenTimes++
		}
		cb.Do(failureTask, onError, onOpen)
		nowTs += 100 // mock reset timeout
		cb.Do(successTask, onError, onOpen)

		if got, want := successTimes, 1; got != want {
			t.Errorf("successTimes = %v, want %v", got, want)
		}
		if got, want := onErrorTimes, 1; got != want {
			t.Errorf("onErrorTimes = %v, want %v", got, want)
		}
		if got, want := onOpenTimes, 0; got != want {
			t.Errorf("onOpenTimes = %v, want %v", got, want)
		}

	})

}

func TestSyncCircuitBreaker_recordFailure(t *testing.T) {

	t.Run("", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs: 10,
			failures:      5,
			nowFn:         mockNowFn,
		}
		cb.recordFailure()
		if cb.failures != 6 {
			t.Errorf("failures = %v, want %v", cb.failures, 6)
		}
		if cb.lastFailureTs == 10 {
			t.Errorf("lastFailureTs = %v, want neq %v", cb.lastFailureTs, 10)
		}
	})

}

func TestSyncCircuitBreaker_reset(t *testing.T) {
	t.Run("", func(t *testing.T) {
		cb := &SyncCircuitBreaker{
			lastFailureTs: 10,
			failures:      5,
		}
		cb.reset()
		if cb.failures != 0 {
			t.Errorf("failures = %v, want %v", cb.failures, 0)
		}
		if cb.lastFailureTs != 0 {
			t.Errorf("lastFailureTs = %v, want %v", cb.lastFailureTs, 0)
		}
	})
}

func TestSyncCircuitBreaker_state(t *testing.T) {

	type fields struct {
		failureThreshold int
		lastFailureTs    int64
		failures         int
		resetTimeout     int64
		nowFn            func() int64
	}

	tests := []struct {
		name   string
		fields fields
		want   state
	}{
		{
			name: "failureThreshold:0, failures:0, not reach resetTimeout",
			fields: fields{
				failureThreshold: 0,
				lastFailureTs:    nowTs - 1,
				failures:         0,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: open,
		},
		{
			name: "failureThreshold:0, failures:0, reach resetTimeout",
			fields: fields{
				failureThreshold: 0,
				lastFailureTs:    nowTs - 11,
				failures:         0,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: halfOpen,
		},
		{
			name: "failureThreshold:1, failures:1, not reach resetTimeout",
			fields: fields{
				failureThreshold: 1,
				lastFailureTs:    nowTs - 1,
				failures:         1,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: open,
		},
		{
			name: "failureThreshold:1, failures:1, reach resetTimeout",
			fields: fields{
				failureThreshold: 1,
				lastFailureTs:    nowTs - 11,
				failures:         1,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: halfOpen,
		},
		{
			name: "failureThreshold:1, failures:2, not reach resetTimeout",
			fields: fields{
				failureThreshold: 1,
				lastFailureTs:    nowTs - 1,
				failures:         2,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: open,
		},
		{
			name: "failureThreshold:1, failures:2, reach resetTimeout",
			fields: fields{
				failureThreshold: 1,
				lastFailureTs:    nowTs - 11,
				failures:         2,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: halfOpen,
		},
		{
			name: "failureThreshold:1, failures:0",
			fields: fields{
				failureThreshold: 1,
				lastFailureTs:    0,
				failures:         0,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: closed,
		},
		{
			name: "failureThreshold:2, failures:1",
			fields: fields{
				failureThreshold: 2,
				lastFailureTs:    nowTs,
				failures:         1,
				resetTimeout:     10,
				nowFn:            mockNowFn,
			},
			want: closed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncCircuitBreaker{
				failureThreshold: tt.fields.failureThreshold,
				lastFailureTs:    tt.fields.lastFailureTs,
				failures:         tt.fields.failures,
				resetTimeout:     tt.fields.resetTimeout,
				nowFn:            tt.fields.nowFn,
			}
			if got := s.state(); got != tt.want {
				t.Errorf("state() = %v, want %v", got, tt.want)
			}
		})
	}
}
