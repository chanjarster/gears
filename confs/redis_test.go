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

package confs

import (
	"github.com/go-redis/redis/v7"
	"testing"
)

func Test_prepareRedisNativeConfig(t *testing.T) {

	redisConf := &RedisConf{
		Host:     "localhost",
		Port:     1234,
		Password: "bar",
		Pool:     1,
		MinIdle:  2,
		Db:       3,
	}
	customizer := func(ropt *redis.Options) {
		ropt.MaxRetries = 2
	}
	redisOpts := prepareRedisNativeConfig(redisConf, customizer)
	if got, want := redisOpts.Addr, "localhost:1234"; got != want {
		t.Errorf("redisOpts.Addr = %v, want %v", got, want)
	}
	if got, want := redisOpts.Password, redisConf.Password; got != want {
		t.Errorf("redisOpts.Password = %v, want %v", got, want)
	}
	if got, want := redisOpts.PoolSize, redisConf.Pool; got != want {
		t.Errorf("redisOpts.PoolSize = %v, want %v", got, want)
	}
	if got, want := redisOpts.MinIdleConns, redisConf.MinIdle; got != want {
		t.Errorf("redisOpts.MinIdleConns = %v, want %v", got, want)
	}
	if got, want := redisOpts.MaxRetries, 2; got != want {
		t.Errorf("redisOpts.MaxRetries = %v, want %v", got, want)
	}
	if got, want := redisOpts.DB, 3; got != want {
		t.Errorf("redisOpts.DB = %v, want %v", got, want)
	}
}
