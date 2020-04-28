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

package event

import (
	"fmt"
)

// Create a FanOutBus and two receivers.
func Example_fanOutBus() {

	bus := NewFanOutBus(1024)
	recvFoo := bus.NewRecv("foo", 1024)

	go func() {
		bus.C <- "hello"
	}()

	go func() {
		for v := range recvFoo.C {
			fmt.Println(v)
		}
	}()

	// make a new Receiver on the fly
	recvBar := bus.NewRecv("bar", 1024)

	go func() {
		for v := range recvBar.C {
			fmt.Println(v)
		}
	}()

	go func() {
		// close Receiver if you want
		recvBar.Close()
	}()

	go func() {
		// close the bus, also close all the receivers.
		bus.Close()
	}()

}
