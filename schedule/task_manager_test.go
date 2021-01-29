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
	"reflect"
	"testing"
	"time"
)

func Test_taskManager_Register(t *testing.T) {

	tm := NewTaskManager().(*taskManager)
	defer tm.UnregisterAll()

	tm.Register(NewFixedRateTask("foo", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	tm.Register(NewFixedDelayTask("bar", emptyFunc, time.Millisecond*20, time.Millisecond*10))

	got := len(tm.taskMap)
	want := 2
	if got != want {
		t.Errorf("len(taskMap) = %v, want %v", got, want)
	}
}

func Test_taskManager_Register_dup(t *testing.T) {

	tm := NewTaskManager().(*taskManager)
	defer tm.UnregisterAll()

	tm.Register(NewFixedRateTask("foo", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	got := tm.Register(NewFixedDelayTask("foo", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	want := errors.New("task[foo] duplicated")
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Register dup Task got err: %v, want %v", got, want)
	}
}

func Test_taskManager_Unregister(t *testing.T) {

	tm := NewTaskManager().(*taskManager)
	defer tm.UnregisterAll()

	tm.Register(NewFixedRateTask("foo", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	tm.Unregister("foo")
	tm.Unregister("foo")

	got := len(tm.taskMap)
	want := 0
	if got != want {
		t.Errorf("len(taskMap) = %v, want %v", got, want)
	}
}

func Test_taskManager_UnregisterAll(t *testing.T) {

	tm := NewTaskManager().(*taskManager)

	tm.Register(NewFixedRateTask("foo", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	tm.Register(NewFixedDelayTask("bar", emptyFunc, time.Millisecond*20, time.Millisecond*10))
	tm.Register(NewFixedRateTask("zoo", emptyFunc, time.Millisecond*20, time.Millisecond*10))

	tm.UnregisterAll()

	got := len(tm.taskMap)
	want := 0
	if got != want {
		t.Errorf("len(taskMap) = %v, want %v", got, want)
	}
}
