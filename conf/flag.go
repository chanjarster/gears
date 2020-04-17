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
	"flag"
	"os"
	"reflect"
	"time"
)

type flagResolver struct {
	args    []string
	flagSet *flag.FlagSet
}

func (r *flagResolver) init(p interface{}) {

	if r.args == nil {
		r.args = os.Args[1:]
	}

	if r.flagSet == nil {
		r.flagSet = flag.NewFlagSet("flag", flag.ContinueOnError)
	}

	flagSet := r.flagSet

	visitExportedFields(p, func(path string, f reflect.StructField, v reflect.Value) {

		flagName := flagStyle(path)

		vi := v.Addr().Interface()

		vt := v.Type()
		switch vt {
		case reflect.TypeOf(time.Nanosecond):
			flagSet.DurationVar(vi.(*time.Duration), flagName, time.Duration(v.Int()), path)
		default:
			switch k := vt.Kind(); k {
			case reflect.Bool:
				flagSet.BoolVar(vi.(*bool), flagName, v.Bool(), path)
			case reflect.Float64:
				flagSet.Float64Var(vi.(*float64), flagName, v.Float(), path)
			case reflect.Int:
				flagSet.IntVar(vi.(*int), flagName, int(v.Int()), path)
			case reflect.Int64:
				flagSet.Int64Var(vi.(*int64), flagName, v.Int(), path)
			case reflect.String:
				flagSet.StringVar(vi.(*string), flagName, v.String(), path)
			case reflect.Uint:
				flagSet.UintVar(vi.(*uint), flagName, uint(v.Uint()), path)
			case reflect.Uint64:
				flagSet.Uint64Var(vi.(*uint64), flagName, v.Uint(), path)
			default:
			}
		}

	})

}
func (r *flagResolver) Resolve(p interface{}) error {

	flagSet := r.flagSet
	return flagSet.Parse(r.args)
}
