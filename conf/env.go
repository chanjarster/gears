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
	"flag"
	"fmt"
	"os"
	"reflect"
	"time"
)

type envResolver struct {
	environ []string
	flagSet *flag.FlagSet
}

func (r *envResolver) init(p interface{}) {

	if r.environ == nil {
		r.environ = make([]string, 0)
	}

	r.flagSet = flag.NewFlagSet("environ", flag.ContinueOnError)

	flagSet := r.flagSet

	visitExportedFields(p, func(path string, f reflect.StructField, v reflect.Value) {

		envName := envStyle(path)
		r.environ = append(r.environ, envName)

		vi := v.Addr().Interface()
		vt := v.Type()
		switch vt {
		case reflect.TypeOf(time.Nanosecond):
			flagSet.DurationVar(vi.(*time.Duration), envName, time.Duration(v.Int()), path)
		default:
			switch k := vt.Kind(); k {
			case reflect.Bool:
				flagSet.BoolVar(vi.(*bool), envName, v.Bool(), path)
			case reflect.Float64:
				flagSet.Float64Var(vi.(*float64), envName, v.Float(), path)
			case reflect.Int:
				flagSet.IntVar(vi.(*int), envName, int(v.Int()), path)
			case reflect.Int64:
				flagSet.Int64Var(vi.(*int64), envName, v.Int(), path)
			case reflect.String:
				flagSet.StringVar(vi.(*string), envName, v.String(), path)
			case reflect.Uint:
				flagSet.UintVar(vi.(*uint), envName, uint(v.Uint()), path)
			case reflect.Uint64:
				flagSet.Uint64Var(vi.(*uint64), envName, v.Uint(), path)
			default:
			}
		}

	})

}

func (r *envResolver) Resolve(p interface{}) error {

	flagSet := r.flagSet
	for _, env := range r.environ {
		v, ok := os.LookupEnv(env)
		if !ok {
			continue
		}
		if err := flagSet.Set(env, v); err != nil {
			return r.failf("invalid value %q for env %s: %v", v, env, err)
		}
	}
	return nil
}

func (r *envResolver) failf(format string, a ...interface{}) error {
	err := fmt.Errorf(format, a...)
	fmt.Fprintln(r.flagSet.Output(), err)
	return err
}
