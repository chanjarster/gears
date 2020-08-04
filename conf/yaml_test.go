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

package conf

import (
	"reflect"
	"testing"
	"time"
)

func TestYamlFileResolver_Resolve(t *testing.T) {

	t.Run("exported fields", func(t *testing.T) {

		o := &outer{}
		initStruct(o)

		y := &yamlFileResolver{
			File: "outer.yaml",
		}

		y.Resolve(o)

		var e_b = true
		var e_f float64 = 3
		var e_i = 3
		var e_i64 int64 = 3
		var e_s = "foo"
		var e_uint uint = 3
		var e_d = 3 * time.Second

		want := &outer{
			B:     e_b,
			Bp:    &e_b,
			F:     e_f,
			Fp:    &e_f,
			I:     e_i,
			Ip:    &e_i,
			I64:   e_i64,
			I64p:  &e_i64,
			S:     e_s,
			Sp:    &e_s,
			Uint:  e_uint,
			Uintp: &e_uint,
			Inner: inner{
				I:  e_i,
				Ip: &e_i,
			},
			Innerp: &inner{
				I:  e_i,
				Ip: &e_i,
			},
			Inner2: struct {
				I  int
				Ip *int
			}{
				I:  e_i,
				Ip: &e_i,
			},
			D:  e_d,
			Dp: &e_d,
		}

		if !reflect.DeepEqual(o, want) {
			t.Errorf("o = %v, want %v", o, want)
		}

	})

	t.Run("unexported fields", func(t *testing.T) {

		f := &foo{}
		initStruct(f)

		y := &yamlFileResolver{
			File: "foo.yaml",
		}

		y.Resolve(f)

		e_i := 3
		e_d := 3 * time.Second

		want := &foo{
			i:  0,
			ip: nil,
			bar: bar{
				i:  0,
				ip: nil,
				I:  0,
				Ip: nil,
			},
			barP: nil,
			I:    3,
			Ip:   &e_i,
			Bar: bar{
				i:  0,
				ip: nil,
				I:  3,
				Ip: &e_i,
			},
			BarP: &bar{
				i:  0,
				ip: nil,
				I:  3,
				Ip: &e_i,
			},
			baz: struct {
				i int
			}{
				i: 0,
			},
			Baz: struct {
				i int
			}{
				i: 0,
			},
			D:  e_d,
			Dp: &e_d,
		}

		if !reflect.DeepEqual(f, want) {
			t.Errorf("f = %v, want %v", f, want)
		}

	})

	t.Run("embedding", func(t *testing.T) {

		f := &embedding{}
		initStruct(f)

		y := &yamlFileResolver{
			File: "embedded.yaml",
		}

		y.Resolve(f)

		e_ap := "xyz"

		want := &embedding{}
		want.A = "abc"
		want.Ap = &e_ap

		if !reflect.DeepEqual(f, want) {
			t.Errorf("o = %v, want %v", f, want)
		}

	})


}
