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
	"fmt"
	"github.com/chanjarster/gears/simplelog"
	"github.com/go-redis/redis/v7"
	"strconv"
)

// Config keys:
//  | Environment   |  Flag     |  Description             |
//  |---------------|-----------|--------------------------|
//  | HOST          | -host     |                          |
//  | PORT          | -port     |                          |
//  | PASSWORD      | -password |                          |
//  | POOL          | -pool     | Connection pool size     |
//  | MIN_IDLE      | -min-idle | Minimal idle connections |
//  | DB            | -db       | Database to be selected after connecting to the server. |
// Note: if RedisConf is nested in another struct, add corresponding prefix.
type RedisConf struct {
	Host     string // Redis host
	Port     int    // Redis port
	Password string // Redis password
	Pool     int    // Redis pool size
	MinIdle  int    // Redis min idle
	Db       int    // Database to be selected after connecting to the server. See https://redis.io/commands/select
}

type RedisOptionCustomizer func(ropt *redis.Options)

func NewRedisClient(rc *RedisConf, customizer RedisOptionCustomizer) *redis.Client {

	redisOpts := prepareRedisNativeConfig(rc, customizer)

	client := redis.NewClient(redisOpts)
	_, err := client.Ping().Result()
	if err != nil {
		simplelog.ErrLogger.Fatal("Redis connection error: ", err)
	}

	simplelog.StdLogger.Printf("Connected to Redis: %s:%d\n", rc.Host, rc.Port)
	return client

}

func (r *RedisConf) String() string {
	return fmt.Sprintf("{Host: %s, Port: %d, Password: ***, Pool: %d, MinIdle: %d}",
		r.Host, r.Port, r.Pool, r.MinIdle)
}

func prepareRedisNativeConfig(rc *RedisConf, customizer RedisOptionCustomizer) *redis.Options {
	redisOpts := &redis.Options{
		Addr:         rc.Host + ":" + strconv.Itoa(rc.Port),
		Password:     rc.Password,
		PoolSize:     rc.Pool,
		MinIdleConns: rc.MinIdle,
		DB:           rc.Db,
	}
	if customizer != nil {
		customizer(redisOpts)
	}
	return redisOpts
}
