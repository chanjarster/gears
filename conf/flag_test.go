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

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestFlagResolver_Resolve(t *testing.T) {

	t.Run("exported fields", func(t *testing.T) {

		o := &outer{}
		initStruct(o)

		r := &flagResolver{
			args: []string{
				"-b=true",
				"-bp=true",
				"-f=3",
				"-fp=3",
				"-i=3",
				"-ip=3",
				"-i64=3",
				"-i64p=3",
				"-s=foo",
				"-sp=foo",
				"-uint=3",
				"-uintp=3",
				"-inner-i=3",
				"-inner-ip=3",
				"-innerp-i=3",
				"-innerp-ip=3",
				"-inner2-i=3",
				"-inner2-ip=3",
				"-d=3s",
				"-dp=3s",
			},
		}

		r.init(o)
		r.Resolve(o)

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
		type foo struct {
			i int
		}

		f := &foo{}
		initStruct(f)

		r := &flagResolver{}
		r.args = []string{"-i=3"}

		r.init(f)
		err := r.Resolve(f)

		if r.flagSet.NFlag() > 0 {
			t.Errorf("NFlag = %v, want %v", r.flagSet.NFlag(), 0)
		}
		wantErr := "flag provided but not defined: -i"
		if err.Error() != wantErr {
			t.Errorf("err = %s, want %s", err, wantErr)
		}
	})

	t.Run("parse error", func(t *testing.T) {

		f := &outer{}
		initStruct(f)

		r := &flagResolver{}
		r.args = []string{"-i=abc"}

		r.init(f)
		err := r.Resolve(f)

		wantErr := `invalid value "abc" for flag -i:`
		if !strings.HasPrefix(err.Error(), wantErr) {
			t.Errorf("err prefix = %s, want %s", err, wantErr)
		}
	})

	t.Run("no flags", func(t *testing.T) {

		f := &outer{}
		initStruct(f)

		r := &flagResolver{}
		r.args = make([]string, 0)

		r.init(f)
		err := r.Resolve(f)

		if r.flagSet.NFlag() > 0 {
			t.Errorf("NFlag = %v, want %v", r.flagSet.NFlag(), 0)
		}

		if err != nil {
			t.Errorf("err = %s, want nil", err)
		}
	})

}
