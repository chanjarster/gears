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

package schedule

import (
	"errors"
	"github.com/chanjarster/gears/testutil"
	"sync"
	"testing"
	"time"
)

func TestNewFixedRateTask(t *testing.T) {
	type args struct {
		name     string
		job      JobFunc
		interval time.Duration
		timeout  time.Duration
	}
	tests := []struct {
		name        string
		args        args
		wantPanic   bool
		expectPanic interface{}
	}{
		{
			name: "all good",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: time.Second * 2,
				timeout:  time.Second,
			},
			wantPanic: false,
		},
		{
			name: "job nil",
			args: args{
				name:     "abc",
				job:      nil,
				interval: time.Second * 2,
				timeout:  time.Second,
			},
			wantPanic:   true,
			expectPanic: "task[abc] job must be not nil",
		},
		{
			name: "name blank",
			args: args{
				name:     "  ",
				job:      emptyFunc,
				interval: time.Second * 2,
				timeout:  time.Second,
			},
			wantPanic:   true,
			expectPanic: "task name must not be blank",
		},
		{
			name: "interval zero",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: 0,
				timeout:  time.Second,
			},
			wantPanic:   true,
			expectPanic: "task[abc] interval must > 0",
		},
		{
			name: "interval < zero",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: -1,
				timeout:  time.Second,
			},
			wantPanic:   true,
			expectPanic: "task[abc] interval must > 0",
		},
		{
			name: "timeout zero",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: time.Second,
				timeout:  0,
			},
			wantPanic:   true,
			expectPanic: "task[abc] timeout must > 0",
		},
		{
			name: "timeout < zero",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: time.Second,
				timeout:  -1,
			},
			wantPanic:   true,
			expectPanic: "task[abc] timeout must > 0",
		},
		{
			name: "interval = timeout",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: time.Second,
				timeout:  time.Second,
			},
			wantPanic: false,
		},
		{
			name: "interval < timeout",
			args: args{
				name:     "abc",
				job:      emptyFunc,
				interval: 1,
				timeout:  2,
			},
			wantPanic: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				testutil.ExpectPanic(t, "NewFixedRateTask()", tt.expectPanic, func() {
					NewFixedRateTask(tt.args.name, tt.args.job, tt.args.interval, tt.args.timeout)
				})
			} else {
				testutil.ShouldNotPanic(t, "NewFixedRateTask()", func() {
					got := NewFixedRateTask(tt.args.name, tt.args.job, tt.args.interval, tt.args.timeout)
					if got == nil {
						t.Errorf("NewFixedRateTask() = nil")
					}
					gotTask := got.(*FixedRateTask)
					if gotTask.name != tt.args.name {
						t.Errorf("NewFixedRateTask().name = %v, want %v", gotTask.name, tt.args.name)
					}
					if gotTask.status != new {
						t.Errorf("NewFixedRateTask().status = %v, want %v", gotTask.status, new)
					}
					if gotTask.job == nil {
						t.Errorf("NewFixedRateTask().job = nil")
					}
					if gotTask.stopSig == nil {
						t.Errorf("NewFixedRateTask().stopSig = nil")
					}
					if gotTask.timeout != tt.args.timeout {
						t.Errorf("NewFixedRateTask().timeout = %v, want %v", gotTask.timeout, tt.args.timeout)
					}
					if gotTask.interval != tt.args.interval {
						t.Errorf("NewFixedRateTask().interval = %v, want %v", gotTask.interval, tt.args.interval)
					}
				})
			}
		})
	}
}

func TestFixedRateTask_Start(t *testing.T) {

	t.Run("good", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		task.Start()
		if task.status != started {
			t.Errorf("Start().status = %v, want %v", task.status, started)
		}
		task.Stop()
	})

	t.Run("start twice", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		task.Start()
		testutil.ExpectPanic(t, "Start()", errors.New("task[foo] start error, you cannot start a Task twice"), func() {
			task.MustStart()
		})
		task.Stop()
	})

	t.Run("start stopped", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		task.Start()
		task.Stop()
		testutil.ExpectPanic(t, "Start()", errors.New("task[foo] start error, you cannot start a stopped Task"), func() {
			task.MustStart()
		})
	})

}

func TestFixedRateTask_Stop(t *testing.T) {

	t.Run("good", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		task.Start()
		task.Stop()
		if task.status != stopped {
			t.Errorf("Start().status = %v, want %v", task.status, stopped)
		}
	})

	t.Run("stop new", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		testutil.ExpectPanic(t, "Stop()", errors.New("task[foo] stop error, Task is not started"), func() {
			task.MustStop()
		})
	})

	t.Run("stop twice", func(t *testing.T) {
		task := NewFixedRateTask("foo", emptyFunc, time.Second*2, time.Second).(*FixedRateTask)
		task.Start()
		task.Stop()
		testutil.ExpectPanic(t, "Stop()", errors.New("task[foo] stop error, task is already stopped"), func() {
			task.MustStop()
		})
	})

}

func TestFixedRateTask_Start_scheduled_several_times(t *testing.T) {

	execCount := 0
	execCountFunc := func(ch <-chan struct{}) error {
		execCount++
		return nil
	}
	task := NewFixedRateTask("foo", execCountFunc, time.Millisecond*2, time.Millisecond).(*FixedRateTask)
	task.Start()

	time.Sleep(time.Millisecond * 5)
	task.Stop()

	if execCount < 2 {
		t.Errorf("execCount = %v, want >= 2", execCount)
	}

}

func TestFixedRateTask_Start_cancel_timeout_job(t *testing.T) {

	r := &countRecord{
		lock: &sync.Mutex{},
	}

	sleepFunc := func(cancel <-chan struct{}) error {
		fin := make(chan struct{})
		go func() {
			time.Sleep(time.Millisecond * 100)
			close(fin)
		}()
		select {
		case <-fin:
		case <-cancel:
			r.AddCount()
		}
		return nil
	}
	task := NewFixedRateTask("foo", sleepFunc, time.Millisecond*20, time.Millisecond*10).(*FixedRateTask)
	task.Start()
	time.Sleep(time.Millisecond * 50)
	task.Stop()
	time.Sleep(time.Millisecond * 50)
	if r.GetCount() < 2 {
		t.Errorf("timeout job is not canceled %v times, want >= 2", r.GetCount())
	}
}

func TestFixedRateTask_Start_overlap(t *testing.T) {
	o := &overlapRecord{
		lock:         &sync.Mutex{},
		overlapped:   false,
		runningCount: 0,
	}

	sleepFunc := func(cancel <-chan struct{}) error {
		if o.GetRunningCount() != 0 {
			o.SetOverlapped(true)
		}
		o.AddRunningCount(1)
		defer func() {
			o.AddRunningCount(-1)
		}()

		fin := make(chan struct{})
		go func() {
			time.Sleep(time.Millisecond * 100)
			close(fin)
		}()
		select {
		case <-fin:
		case <-cancel:
		}
		return nil
	}
	task := NewFixedRateTask("foo", sleepFunc, time.Millisecond*10, time.Millisecond*20).(*FixedRateTask)
	task.Start()
	time.Sleep(time.Millisecond * 100)
	task.Stop()
	time.Sleep(time.Millisecond * 50)
	if !o.IsOverlapped() {
		t.Errorf("job overlapped = false, expecting true")
	}
}
