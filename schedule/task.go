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
	"context"
	"errors"
	"fmt"
	"time"
)

type Task interface {
	// Return the Task name
	Name() string
	// Start the Task, cannot start a Task twice or start a stopped Task
	Start() error
	// Start the Task, panic if error happened
	MustStart()
	// Stop the task, cannot stop a not started Task or stop a Task twice.
	// This function only send the stop signal to the Task, does not wait for JobFunc to finish.
	Stop() error
	// Stop the Task, panic if error happened
	MustStop()
}

const (
	new = iota
	started
	stopped
)

// User defined job.
// It's recommended that JobFunc respond to cancel channel.
type JobFunc func(cancel <-chan struct{}) error

type taskCommon struct {
	name     string        // name for the task
	job      JobFunc       // user provided JobFunc
	stopSig  chan struct{} // internal stop signal channel
	interval time.Duration // interval between two runs
	timeout  time.Duration // timeout for the job
	jobQ     chan struct{} //
	no       int64         // count how many time current task has been run
	status   int           // task status, new, started, stopped
}

func (t *taskCommon) Name() string {
	return t.name
}

// should be called guard by lock
func (t *taskCommon) markStart() error {
	if t.status == started {
		return errors.New(fmt.Sprintf("task[%s] start error, you cannot start a Task twice", t.name))
	}

	if t.status == stopped {
		return errors.New(fmt.Sprintf("task[%s] start error, you cannot start a stopped Task", t.name))
	}

	t.status = started
	return nil
}

// should be called guard by lock
func (t *taskCommon) markStopSendSignal() error {
	if t.status == new {
		return errors.New(fmt.Sprintf("task[%s] stop error, Task is not started", t.name))
	}
	if t.status == stopped {
		return errors.New(fmt.Sprintf("task[%s] stop error, task is already stopped", t.name))
	}

	t.status = stopped
	close(t.stopSig)
	// we do not close jobQ because there will be race condition on sending and closing a channel
	// https://golang.org/doc/articles/race_detector.html#Unsynchronized_send_and_close_operations
	// close(t.jobQ)
	return nil
}

func (t *taskCommon) runJob(nextNo int64) {
	stdLogger.Printf("task[%s] job[#%v] start\n", t.name, nextNo)
	timeoutCtx, cancel := context.WithTimeout(context.Background(), t.timeout)
	defer cancel()

	fin := make(chan struct{})

	go func() {
		startAt := time.Now()

		err := t.job(timeoutCtx.Done())

		finishAt := time.Now()
		elapse := finishAt.Sub(startAt)

		close(fin)

		if err != nil {
			errLogger.Printf("task[%s] job[#%v] error: %v\n", t.name, nextNo, err)
		}

		if ctxErr := timeoutCtx.Err(); ctxErr == context.Canceled {
			stdLogger.Printf("task[%s] job[#%v] canceled by stop signal, time elapse %v\n", t.name, nextNo, elapse)
		} else if ctxErr == context.DeadlineExceeded {
			stdLogger.Printf("task[%s] job[#%v] canceled due to timeout(%v), time elapse %v\n", t.name, nextNo, t.timeout, elapse)
		} else {
			stdLogger.Printf("task[%s] job[#%v] finished, time elapse %v\n", t.name, nextNo, elapse)
		}

	}()

	select {
	case <-timeoutCtx.Done():
		<-fin
	case <-t.stopSig:
		cancel()
		<-fin
	case <-fin:
		// job normally finished
	}

}

// Append a JobFunc to jobQ.
// Return false if Task is stopped and no job appended, otherwise return true.
func (t *taskCommon) appendJobQueue() bool {
	// because we do not close t.jobQ on stop, so we should check manually whether Task is stopped
	select {
	case t.jobQ <- struct{}{}:
		return true
	case <-t.stopSig:
		// task is stopped
		return false
	}
	return false
}
