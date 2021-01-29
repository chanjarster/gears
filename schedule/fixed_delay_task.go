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
	"fmt"
	"github.com/chanjarster/gears/simplelog"
	"strings"
	"sync"
	"time"
)

type FixedDelayTask struct {
	lock *sync.Mutex // guard below
	taskCommon
}

// Create a new FixedDelayTask.
//
// Run jobs after another when previous one finished and interval time elapsed.
func NewFixedDelayTask(name string, job JobFunc, interval time.Duration, timeout time.Duration) Task {

	if strings.TrimSpace(name) == "" {
		panic(fmt.Sprintf("task name must not be blank"))
	}

	if timeout <= 0 {
		panic(fmt.Sprintf("task[%s] timeout must > 0", name))
	}
	if interval <= 0 {
		panic(fmt.Sprintf("task[%s] interval must > 0", name))
	}
	if job == nil {
		panic(fmt.Sprintf("task[%s] job must be not nil", name))
	}
	t := &FixedDelayTask{
		taskCommon: taskCommon{
			name:     name,
			job:      job,
			stopSig:  make(chan struct{}),
			interval: interval,
			timeout:  timeout,
			status:   new,
			jobQ:     make(chan struct{}, 0),
		},
		lock: &sync.Mutex{},
	}
	return t
}

func (t *FixedDelayTask) MustStart() {
	err := t.Start()
	if err != nil {
		panic(err)
	}
}

func (t *FixedDelayTask) Start() error {

	t.lock.Lock()
	defer t.lock.Unlock()

	err := t.taskCommon.markStart()
	if err != nil {
		return err
	}

	go t.startConsuming()
	// trigger first schedule immediately
	if !t.appendJobQueue() {
		return nil
	}

	simplelog.StdLogger.Printf("task[%s] started", t.name)

	return nil
}

func (t *FixedDelayTask) startConsuming() {
	for {
		select {
		case <-t.jobQ:
			if t.getStatus() == stopped {
				return
			}
			t.runJob(t.nextNo())
			go t.appendJobQueueFuture()
		case <-t.stopSig:
			// task is stopped
			return
		}
	}
}

func (t *FixedDelayTask) appendJobQueueFuture() {
	time.AfterFunc(t.interval, func() {
		t.appendJobQueue()
	})
}

func (t *FixedDelayTask) getStatus() int {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.status
}

func (t *FixedDelayTask) MustStop() {
	err := t.Stop()
	if err != nil {
		panic(err)
	}
}

func (t *FixedDelayTask) Stop() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	err := t.taskCommon.markStopSendSignal()
	if err != nil {
		return err
	}
	simplelog.StdLogger.Printf("task[%s] stop signal sent", t.name)
	return nil
}

func (t *FixedDelayTask) nextNo() int64 {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.no++
	return t.no
}
