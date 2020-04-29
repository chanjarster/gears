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
	"reflect"
	"testing"
	"time"
)


func TestNewSyncFixedWindow(t *testing.T) {
	cap := 500
	windowSize := 500 * time.Millisecond

	got := NewSyncFixedWindow(cap, windowSize)

	if got, want := got.count, 0; got != want {
		t.Errorf("NewSyncFixedWindow().count = %v, want %v", got, want)
	}
	if got, want := got.WindowSize(), windowSize; got != want {
		t.Errorf("NewSyncFixedWindow().WindowSize() = %v, want %v", got, want)
	}
	if got, want := got.Capacity(), cap; got != want {
		t.Errorf("NewSyncFixedWindow().Capacity() = %v, want %v", got, want)
	}
	if got, want := got.until, int64(0); got != want {
		t.Errorf("NewSyncFixedWindow().until = %v, want %v", got, want)
	}
	if got := got.nowFn; got == nil {
		t.Errorf("NewSyncFixedWindow().nowFn is nil, want not nil")
	}
}

func TestSyncFixedWindow_Acquire(t *testing.T) {
	new := func() *SyncFixedWindow {
		cap := 10
		windowSize := 50 * time.Millisecond
		return NewSyncFixedWindow(cap, windowSize)
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

	t.Run("acquire 10 times, wait 60ms", func(t *testing.T) {
		ratelimiter := new()
		for i := 0; i < ratelimiter.Capacity(); i++ {
			if got := ratelimiter.Acquire(); !got {
				t.Errorf("acquire() = %v, want %v", got, true)
			}
		}
		ratelimiter.nowFn = gtime.FixedNow(time.Now().Add(time.Millisecond * 60))

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

func TestSyncFixedWindow_contains(t *testing.T) {
	type fields struct {
		until int64
	}
	type args struct {
		ts int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			fields: fields{0},
			args:   args{0},
			want:   true,
		},
		{
			fields: fields{1},
			args:   args{0},
			want:   true,
		},
		{
			fields: fields{1},
			args:   args{2},
			want:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncFixedWindow{
				until: tt.fields.until,
			}
			if got := s.contains(tt.args.ts); got != tt.want {
				t.Errorf("contains() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSyncFixedWindow_reset(t *testing.T) {
	type fields struct {
		windowSize int64
		until      int64
	}
	type args struct {
		ts int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int64
	}{
		{
			fields: fields{10, 0},
			args:   args{1},
			want:   10,
		},
		{
			fields: fields{10, 1},
			args:   args{0},
			want:   1,
		},
		{
			fields: fields{10, 1},
			args:   args{9},
			want:   11,
		},
		{
			fields: fields{10, 1},
			args:   args{11},
			want:   11,
		},
		{
			fields: fields{10, 1},
			args:   args{12},
			want:   21,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SyncFixedWindow{
				windowSize: tt.fields.windowSize,
				until:      tt.fields.until,
			}
			s.reset(tt.args.ts)
			if got := s.until; got != tt.want {
				t.Errorf("s.until after reset() = %v, want %v", got, tt.want)
			}
			if got := s.count; got != 0 {
				t.Errorf("s.count after reset() = %v, want %v", got, 0)
			}
		})
	}
}
