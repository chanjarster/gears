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
	"math/rand"
	"testing"
	"time"
)

func BenchmarkSyncCircuitBreaker_Do(b *testing.B) {
	breaker := NewSyncCircuitBreaker(5, time.Millisecond*100)
	benchmarkCircuitBreaker_Do(b, breaker)

}

func benchmarkCircuitBreaker_Do(b *testing.B, breaker Interface) {

	err := errors.New("some error")
	task := func() error {
		// 1/10 的概率会fail
		if rand.Intn(10) == 0 {
			return err
		} else {
			return nil
		}
	}
	onError := func(err error) {
	}
	onOpen := func() {
	}
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			breaker.Do(task, onError, onOpen)
		}
	})

}
