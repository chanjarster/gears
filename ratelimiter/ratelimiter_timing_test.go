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
	"sync/atomic"
	"testing"
	"time"
)

func BenchmarkSyncTokenBucket_Acquire(b *testing.B) {
	bucket := NewSyncTokenBucket(10, 10)
	benchmarkTokenBucket_Acquire(b, bucket)

}

func BenchmarkAtomicTokenBucket_Acquire(b *testing.B) {
	bucket := NewAtomicTokenBucket(10, 10)
	benchmarkTokenBucket_Acquire(b, bucket)
}

func BenchmarkSyncSlidingWindow_Acquire(b *testing.B) {
	bucket := NewSyncSlidingWindow(10, time.Second)
	benchmarkTokenBucket_Acquire(b, bucket)
}

func BenchmarkSyncFixedWindow_Acquire(b *testing.B) {
	bucket := NewSyncFixedWindow(10, time.Second)
	benchmarkTokenBucket_Acquire(b, bucket)
}

func benchmarkTokenBucket_Acquire(b *testing.B, bucket Interface) {
	capacity := bucket.Capacity()
	var acquiredCount int64

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			if acquired := bucket.Acquire(); acquired {
				atomic.AddInt64(&acquiredCount, 1)
			} else {
				atomic.StoreInt64(&acquiredCount, 0)
			}
			ac := atomic.LoadInt64(&acquiredCount)
			if ac > int64(capacity) {
				b.Fatalf("continous acquiredCount: %d. Expecting <= %v", acquiredCount, capacity)
			}
		}
	})

}
