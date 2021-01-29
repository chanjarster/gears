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
	"fmt"
	"sync"
)

// Task manager
type TaskManager interface {
	// Register a Task and start it.
	// Once the Task is managed by TaskManager, you should not call Task
	// lifecycle function(Task.Stop, Task.Start) any more.
	Register(task Task) error
	// Call Task.Stop and unregister it
	// This function only send the stop signal to Task(s), returning doesn't mean JobFunc is stopped
	Unregister(name string)
	// Stop all the Task and unregister them
	UnregisterAll()
}

type taskManager struct {
	lock    *sync.Mutex
	taskMap map[string]Task
}

func NewTaskManager() TaskManager {
	return &taskManager{
		lock:    &sync.Mutex{},
		taskMap: make(map[string]Task, 10),
	}
}

func (m *taskManager) Register(task Task) error {

	m.lock.Lock()
	defer m.lock.Unlock()

	if _, hit := m.taskMap[task.Name()]; hit {
		return errors.New(fmt.Sprintf("task[%s] duplicated", task.Name()))
	}

	m.taskMap[task.Name()] = task
	task.MustStart()
	return nil
}

func (m *taskManager) Unregister(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	task := m.taskMap[name]
	if task == nil {
		return
	}
	task.Stop()
	delete(m.taskMap, name)
}

func (m *taskManager) UnregisterAll() {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, task := range m.taskMap {
		task.Stop()
	}

	m.taskMap = make(map[string]Task, 10)
}
