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
	"strings"
	"sync"
	"time"
)

type FixedRateTask struct {
	lock *sync.Mutex // guard below
	taskCommon
	ticker *time.Ticker
}

// Create a new FixedRateTask.
//
// Run jobs at a fixed rate. Be careful that if interval <= timeout, or
// job does not respond to cancel channel (see JobFunc), jobs may be overlapped.
func NewFixedRateTask(name string, job JobFunc, interval time.Duration, timeout time.Duration) Task {

	if strings.TrimSpace(name) == "" {
		panic(fmt.Sprintf("task name must not be blank"))
	}

	if timeout <= 0 {
		panic(fmt.Sprintf("task[%s] timeout must > 0", name))
	}
	if interval <= 0 {
		panic(fmt.Sprintf("task[%s] interval must > 0", name))
	}
	if interval <= timeout {
		stdLogger.Printf("task[%s] interval <= timeout, jobs may be overlapped", name)
	}
	if job == nil {
		panic(fmt.Sprintf("task[%s] job must be not nil", name))
	}
	t := &FixedRateTask{
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

func (t *FixedRateTask) MustStart() {
	err := t.Start()
	if err != nil {
		panic(err)
	}
}

func (t *FixedRateTask) Start() error {
	t.lock.Lock()
	defer t.lock.Unlock()

	err := t.taskCommon.markStart()
	if err != nil {
		return err
	}

	t.ticker = time.NewTicker(t.interval)

	go t.startConsuming()
	go t.startProviding()

	stdLogger.Printf("task[%s] started", t.name)

	return nil
}

func (t *FixedRateTask) MustStop() {
	err := t.Stop()
	if err != nil {
		panic(err)
	}
}

func (t *FixedRateTask) Stop() error {
	t.lock.Lock()
	defer t.lock.Unlock()
	err := t.taskCommon.markStopSendSignal()
	if err != nil {
		return err
	}

	t.ticker.Stop()
	stdLogger.Printf("task[%s] stop signal sent", t.name)
	return nil
}

func (t *FixedRateTask) getStatus() int {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.status
}

func (t *FixedRateTask) startProviding() {
	// trigger first schedule immediately
	if !t.appendJobQueue() {
		return
	}
	for {
		select {
		case <-t.ticker.C:
			if !t.appendJobQueue() {
				return
			}
		case <-t.stopSig:
			// task is stopped
			return
		}
	}
}

func (t *FixedRateTask) startConsuming() {
	for {
		select {
		case <-t.jobQ:
			if t.getStatus() == stopped {
				return
			}
			go t.runJob(t.nextNo())
		case <-t.stopSig:
			// task is stopped
			return
		}
	}
}

func (t *FixedRateTask) nextNo() int64 {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.no++
	return t.no
}
