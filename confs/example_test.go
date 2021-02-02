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
	"github.com/chanjarster/gears/conf"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/go-sql-driver/mysql"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

func ExampleNewRedisClient() {
	redisConf := &RedisConf{}
	conf.Load(redisConf, "")
	customizer := func(ropt *redis.Options) {
		ropt.MaxRetries = 2
	}
	redisClient := NewRedisClient(redisConf, customizer)
	redisClient.Close()
}

func ExampleNewMySqlDb() {
	mysqlConf := &MysqlConf{}
	conf.Load(mysqlConf, "")
	customizer := func(mc *mysql.Config) {
		mc.Params["autocommit"] = "true"
		mc.Params["charset"] = "utf8"
	}
	mySqlDb := NewMySqlDb(mysqlConf, customizer)
	mySqlDb.Close()
}

func ExampleNewFastHttpClient() {
	hcConf := &FastHttpClientConf{}
	conf.Load(hcConf, "")
	customizer := func(hc *fasthttp.Client) {
		// do something
	}
	hc := NewFastHttpClient(hcConf, customizer)
	hc.Get(make([]byte, 0), "https://baidu.com")
}

func ExampleNewGinPprofBasicAuthFilter() {
	pprofConf := &PprofConf{
		Username:             "abc",
		Password:             "xyz",
		BlockProfileRate:     1,
		MutexProfileFraction: 1,
	}
	r := gin.Default()
	pprofGroup := r.Group("debug/pprof", NewGinPprofBasicAuthFilter(pprofConf))
	pprof.RouteRegister(pprofGroup, "")
}

func ExampleNewFasthttpRoutingPprofHandler() {
	pprofConf := &PprofConf{
		Username:             "abc",
		Password:             "xyz",
		BlockProfileRate:     1,
		MutexProfileFraction: 1,
	}
	router := routing.New()
	router.Any("/debug/pprof/*", NewFasthttpRoutingPprofHandler(pprofConf))
}
