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
	"github.com/chanjarster/gears/simplelog"
	"github.com/go-redis/redis/v7"
	"strings"
	"time"
)

const (
	/*
			if exists(key:block)
			  return true, false

		  size = llen(key)
			if size < capacity
			  rpush key now
			  return false, false

		  if size - capacity > 0
		    lpop key (size - capacity) elements

			oldest = lrange key 0 0
			if now - oldest > window size
				rpush key now
		    lpop key
				return false, false

			SET key:block 1 EX timeout NX
			return true, true
	*/
	prefix = "_rl:"
	suffix = ":bl"

	// return block, triggered, ttl, msg
	script1 = `local key = KEYS[1]
local keyb = KEYS[2]
local cap = tonumber(ARGV[1])
local win = tonumber(ARGV[2])
local exp = tonumber(ARGV[3])
local now = tonumber(ARGV[4])
local msg = ARGV[5]

if cap <= 0 or win <= 0 or exp <= 0
then
  return {0, 0, 0, ''}
end

if redis.call('EXISTS', keyb) == 1
then
  local ttl = redis.call('TTL', keyb)
  local omsg = redis.call('GET', keyb)
  return {1, 0, ttl, omsg}
end

local size = redis.call('LLEN', key)
if size < cap
then
  redis.call('RPUSH', key, now)
  redis.call('EXPIRE', key, exp * 2)
  return {0, 0, 0, ''}
end

local diff = size - cap
if diff > 0
then
  for i = diff, 1, -1
  do
    redis.call('LPOP', key)
  end
end

local list = redis.call('LRANGE', key, 0, 0)

if table.getn(list) == 0
then
  redis.call('RPUSH', key, now)
  redis.call('EXPIRE', key, exp * 2)
  return {0, 0, 0, ''}
end

local oldest = tonumber(list[1])
if now - oldest > win
then
  redis.call('RPUSH', key, now)
  redis.call('LPOP', key)
  redis.call('EXPIRE', key, exp * 2)
  return {0, 0, 0, ''}
end

redis.call('SET', keyb, msg, 'EX', exp, 'NX')
return {1, 1, exp, msg}
`

	// return block, triggered, ttl, msg
	script2 = `local keyb = KEYS[1]

if redis.call('EXISTS', keyb) == 1
then
  local ttl = redis.call('TTL', keyb)
  local omsg = redis.call('GET', keyb)
  return {1, 0, ttl, omsg}
end

return {0, 0, 0, ''}
`
)

var (
	scriptSha1 = ""
	scriptSha2 = ""
)

func NewRedisTtlRateLimiter(redisClient *redis.Client, params TtlRateLimiterParams) TtlRateLimiter {
	r := &redisTtlRateLimiter{
		params:      params,
		redisClient: redisClient,
	}
	return r
}

// NewRedisTtlRateLimiterCluster create a redis rate limiter for Redis Cluster environment.
//  hashTag: redis hash tag value, helps to ensure all keys be in the same slot.
// see: https://redis.io/topics/cluster-tutorial#redis-cluster-data-sharding
func NewRedisTtlRateLimiterCluster(redisClient *redis.Client, params TtlRateLimiterParams, hashTag string) TtlRateLimiter {
	r := &redisTtlRateLimiter{
		params:      params,
		redisClient: redisClient,
		hashTag:     "{" + strings.Trim(hashTag, "{}") + "}",
	}
	return r
}

type redisTtlRateLimiter struct {
	params      TtlRateLimiterParams
	redisClient *redis.Client
	hashTag     string
}

func (r *redisTtlRateLimiter) ShouldBlock(key string, msg string) *Result {
	return r.ShouldBlock2(key, key, msg)
}

func (r *redisTtlRateLimiter) ShouldBlock2(key string, blockKey string, msg string) *Result {
	result := &Result{}

	if isParamsNotSet(r.params) {
		return result
	}

	key = prefix + key
	blockKey = blockKey + suffix

	if r.hashTag != "" {
		key = key + r.hashTag
		blockKey = blockKey + r.hashTag
	}

	//local key = KEYS[1]
	//local keyb = KEYS[2]
	//local cap = tonumber(ARGV[1])
	//local win = tonumber(ARGV[2])
	//local exp = tonumber(ARGV[3])
	//local now = tonumber(ARGV[4])
	//local msg = ARGV[5]

	now := time.Now().UnixNano() / int64(time.Second)

	raw, err := r.redisClient.EvalSha(
		scriptSha1,
		[]string{key, blockKey},
		r.params.GetCapacity(),
		r.params.GetWindowSizeSeconds(),
		r.params.GetTimeoutSeconds(),
		now,
		msg,
	).Result()

	if err != nil {
		simplelog.ErrLogger.Println("eval script1 ", scriptSha1, "error", err)
		return result
	}

	arr := raw.([]interface{})
	result.Block = arr[0].(int64) == 1
	result.Triggered = arr[1].(int64) == 1
	result.Ttl = int(arr[2].(int64))
	result.Msg = arr[3].(string)
	return result

}

func (r *redisTtlRateLimiter) IsBlocked(blockKey string) *Result {
	result := &Result{}
	blockKey = blockKey + suffix

	//local keyb = KEYS[1]
	raw, err := r.redisClient.EvalSha(
		scriptSha2,
		[]string{blockKey},
	).Result()

	if err != nil {
		simplelog.ErrLogger.Println("eval script2 ", scriptSha2, "error", err)
		return result
	}

	arr := raw.([]interface{})
	result.Block = arr[0].(int64) == 1
	result.Triggered = arr[1].(int64) == 1
	result.Ttl = int(arr[2].(int64))
	result.Msg = arr[3].(string)
	return result
}

func (r *redisTtlRateLimiter) GetWindowSizeSeconds() int {
	return r.params.GetWindowSizeSeconds()
}

func (r *redisTtlRateLimiter) GetTimeoutSeconds() int {
	return r.params.GetTimeoutSeconds()
}

func (r *redisTtlRateLimiter) GetCapacity() int {
	return r.params.GetCapacity()
}

func LoadScript(redisClient *redis.Client) {
	loadScript(redisClient, script1, func(scriptSha string) {
		scriptSha1 = scriptSha
	})
	loadScript(redisClient, script2, func(scriptSha string) {
		scriptSha2 = scriptSha
	})
}

func loadScript(redisClient *redis.Client, script string, callback func(scriptSha string)) {
	sha, err := redisClient.ScriptLoad(script).Result()
	if err != nil {
		panic(err)
	}
	callback(sha)
}
