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
	"github.com/valyala/fasthttp"
	"testing"
	"time"
)

func TestNewFastHttpClient(t *testing.T) {

	hcConf := &FastHttpClientConf{
		CertChain:           "",
		PrivateKey:          "",
		SslTrustMode:        "OS",
		SslTrustCerts:       "",
		ReadTimeout:         time.Second,
		WriteTimeout:        time.Second * 2,
		MaxConnDuration:     time.Second * 3,
		MaxConnsPerHost:     10,
		MaxIdleConnDuration: time.Second * 4,
		MaxConnWaitTimeout:  time.Second * 5,
	}
	customizer := func(hc *fasthttp.Client) {
		hc.ReadBufferSize = 10
	}

	hc := NewFastHttpClient(hcConf, customizer)

	if got, want := hc.ReadTimeout, hcConf.ReadTimeout; got != want {
		t.Errorf("hc.ReadTimeout = %v, want %v", got, want)
	}
	if got, want := hc.WriteTimeout, hcConf.WriteTimeout; got != want {
		t.Errorf("hc.WriteTimeout = %v, want %v", got, want)
	}
	if got, want := hc.MaxConnDuration, hcConf.MaxConnDuration; got != want {
		t.Errorf("hc.MaxConnDuration = %v, want %v", got, want)
	}
	if got, want := hc.MaxConnsPerHost, hcConf.MaxConnsPerHost; got != want {
		t.Errorf("hc.MaxConnsPerHost = %v, want %v", got, want)
	}
	if got, want := hc.MaxIdleConnDuration, hcConf.MaxIdleConnDuration; got != want {
		t.Errorf("hc.MaxIdleConnDuration = %v, want %v", got, want)
	}
	if got, want := hc.MaxConnWaitTimeout, hcConf.MaxConnWaitTimeout; got != want {
		t.Errorf("hc.MaxConnWaitTimeout = %v, want %v", got, want)
	}
	if got := hc.TLSConfig; got != nil {
		t.Errorf("hc.TLSConfig = %v, want nil", got)
	}
	if got, want := hc.ReadBufferSize, 10; got != want {
		t.Errorf("hc.ReadBufferSize = %v, want %v", got, want)
	}

}
