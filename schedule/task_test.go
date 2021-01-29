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

import "sync"

var emptyFunc = func(cancel <-chan struct{}) error {
	return nil
}

type overlapRecord struct {
	lock         *sync.Mutex
	overlapped   bool
	runningCount int
}

func (o *overlapRecord) GetRunningCount() int {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.runningCount
}

func (o *overlapRecord) AddRunningCount(v int) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.runningCount = o.runningCount + v
}

func (o *overlapRecord) IsOverlapped() bool {
	o.lock.Lock()
	defer o.lock.Unlock()
	return o.overlapped
}

func (o *overlapRecord) SetOverlapped(v bool) {
	o.lock.Lock()
	defer o.lock.Unlock()
	o.overlapped = v
}


type countRecord struct {
	lock      *sync.Mutex
	execCount int
}

func (r *countRecord) GetCount() int {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.execCount
}

func (r *countRecord) AddCount() {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.execCount++
}
