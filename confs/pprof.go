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
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/pprofhandler"
	"net/http"
	"runtime"
	"strings"
)

type PprofConf struct {
	Username             string // /debug/pprof/* basic auth username
	Password             string // /debug/pprof/* basic auth password
	BlockProfileRate     int
	MutexProfileFraction int
}

func (p *PprofConf) String() string {
	return fmt.Sprintf("{Username: ***, Password: ***, BlockProfileRate: %d, MutexProfileFraction: %d}",
		p.BlockProfileRate, p.MutexProfileFraction)
}

// create a gin basic auth filter for pprof endpoint
func NewGinPprofBasicAuthFilter(conf *PprofConf) gin.HandlerFunc {

	// 如果没有设置，就不允许访问
	if conf.Username == "" || conf.Password == "" {
		return func(c *gin.Context) {
			// do nothing cut the chain
			c.String(http.StatusNotFound, "404 page not found")
			c.Abort()
			return
		}
	}

	// 设置Profile相关参数
	runtime.SetBlockProfileRate(conf.BlockProfileRate)
	runtime.SetMutexProfileFraction(conf.MutexProfileFraction)

	return gin.BasicAuth(gin.Accounts{
		conf.Username: conf.Password,
	})

}

// fasthttp

func emptyFasthttpRoutingHandler(ctx *routing.Context) error {
	return nil
}

// create a fasthttp routing handler for pprof handler
func NewFasthttpRoutingPprofHandler(conf *PprofConf) routing.Handler {

	// 如果没有设置，就不允许访问
	if conf.Username == "" || conf.Password == "" {
		return emptyFasthttpRoutingHandler
	}

	// 设置Profile相关参数
	runtime.SetBlockProfileRate(conf.BlockProfileRate)
	runtime.SetMutexProfileFraction(conf.MutexProfileFraction)

	return func(ctx *routing.Context) error {

		refuse := func(ctx *routing.Context) {
			ctx.Response.Header.Add("WWW-Authenticate", `Basic realm="User Visible Realm"`)
			ctx.SetStatusCode(fasthttp.StatusUnauthorized)
		}

		auth := strings.TrimPrefix(string(ctx.Request.Header.Peek("Authorization")), "Basic ")
		if auth == "" {
			refuse(ctx)
			return nil
		}

		b, err := base64.StdEncoding.DecodeString(auth)
		if err != nil {
			refuse(ctx)
			return nil
		}

		user_pass := strings.Split(string(b), ":")
		if len(user_pass) != 2 {
			refuse(ctx)
			return nil
		}

		if user_pass[0] != conf.Username || user_pass[1] != conf.Password {
			refuse(ctx)
			return nil
		}

		pprofhandler.PprofHandler(ctx.RequestCtx)
		return nil
	}

}
