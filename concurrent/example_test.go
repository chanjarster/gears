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

package concurrent

func Example_lockSemaphore() {
	sema := NewLockSemaphore(100) // new a semaphore with 100 permits
	sema.Acquire()                // will be block if no permits available
	defer sema.Release()          // return permit

	if sema.TryAcquire() { // will return false if no permits available
		defer sema.Release() // return permit
	}
}

func Example_chanSemaphore() {
	sema := NewChanSemaphore(100) // new a semaphore with 100 permits
	sema.Acquire()                // will be block if no permits available
	defer sema.Release()          // return permit

	if sema.TryAcquire() { // will return false if no permits available
		defer sema.Release() // return permit
	}
}
