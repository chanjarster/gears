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
	"fmt"
	"net/http"
	"os"
	"time"
)

func Example_circuitBreaker() {
	getGoogle := func(hc *http.Client, circuitBreaker Interface) {
		circuitBreaker.Do(
			func() error {
				resp, err := hc.Get("https://google.com")
				if err != nil {
					return err
				}
				resp.Write(os.Stdout)
				return nil
			},
			func(err error) {
				fmt.Println("google is not available. Error:", err)
			},
			func() {
				fmt.Println("google is not available. Circuit breaker is opened")
			},
		)
	}

	circuitBreaker := NewSyncCircuitBreaker(1, time.Second)
	client := &http.Client{Timeout: time.Nanosecond}

	// this call get will timeout
	getGoogle(client, circuitBreaker)
	// this call will be denied by circuit breaker because of it's opened
	getGoogle(client, circuitBreaker)

}
