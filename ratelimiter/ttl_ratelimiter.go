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

// Legal results are:
//
// 1. block:false, triggered:false, ttl:0, msg:""
//
// 2. block:true, triggered:true, ttl:>0, msg:"some message"
//
// 3. block:true, triggered:false, ttl:>0, msg:"some message"
type Result struct {
	Block     bool   // true: blocked，false: passed
	Triggered bool   // first time blocking，otherwise false
	Ttl       int    // how many seconds blocking will last
	Msg       string // message recorded when first time blocking
}

// A rate limiter that will prevent further request for ttl seconds after first time request rate exceeds the limit.
type TtlRateLimiter interface {

	// same as ShouldBlock2(key, key, msg)
	ShouldBlock(key string, msg string) *Result

	// When the request rate of `key` exceeds the limit, blocking will be triggered(record on `blockKey`)
	// and last for `timeout` seconds(ttl).
	// After `timeout` seconds, `blockKey` will be released and request `key` can be passed again.
	//
	// `msg` is the message for first time blocking.
	//
	// Note: different `key` can share same `blockKey`, same `key` MUST NOT share different `blockKey`
	ShouldBlock2(key string, blockKey string, msg string) *Result

	// capacity: window capacity
	// time range the window look back
	GetWindowSizeSeconds() int
	// window capacity
	GetCapacity() int
	// how many seconds blocking will last after first time blocking happened
	GetTimeoutSeconds() int
}

// Interface for the need of runtime rate limit parameters
type TtlRateLimiterParams interface {
	GetWindowSizeSeconds() int
	GetTimeoutSeconds() int
	GetCapacity() int
}

func isParamsNotSet(params TtlRateLimiterParams) bool {
	return params.GetCapacity() <= 0 || params.GetTimeoutSeconds() <= 0 || params.GetWindowSizeSeconds() <= 0
}

func NewFixedTtlRateLimiterParams(capacity, windowsSizeSec, timeoutSec int) TtlRateLimiterParams {
	return &fixedTtlRateLimiterParams{
		capacity:      capacity,
		windowSizeSec: windowsSizeSec,
		timeoutSec:    timeoutSec,
	}
}

type fixedTtlRateLimiterParams struct {
	timeoutSec    int
	windowSizeSec int
	capacity      int
}

func (f *fixedTtlRateLimiterParams) GetWindowSizeSeconds() int {
	return f.windowSizeSec
}

func (f *fixedTtlRateLimiterParams) GetTimeoutSeconds() int {
	return f.timeoutSec
}

func (f *fixedTtlRateLimiterParams) GetCapacity() int {
	return f.capacity
}
