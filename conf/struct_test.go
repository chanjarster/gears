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

package autoconf

import "time"

//---
// struct with exported fields

type inner struct {
	I  int
	Ip *int
}

type outer struct {
	B  bool
	Bp *bool

	F  float64
	Fp *float64

	I  int
	Ip *int

	I64  int64
	I64p *int64

	S  string
	Sp *string

	Uint  uint
	Uintp *uint

	Inner  inner
	Innerp *inner

	Inner2 struct {
		I  int
		Ip *int
	}

	D  time.Duration
	Dp *time.Duration
}

//-----
// struct with exported and unexported fields

type bar struct {
	i  int
	ip *int

	I  int
	Ip *int
}

type foo struct {
	i    int
	ip   *int
	bar  bar
	barP *bar

	I    int
	Ip   *int
	Bar  bar
	BarP *bar

	D  time.Duration
	Dp *time.Duration

	baz struct {
		i int
	}

	Baz struct {
		i int
	}
}
