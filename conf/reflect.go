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
)

// Initialize p's field. If field is exported and is a pointer,
// that field will be initialized to new(Type)
func initStruct(p interface{}) {
	initPtrValue(reflect.ValueOf(p), "")
}

// param path is just for debug
func initPtrValue(v reflect.Value, path string) {

	if v.Kind() == reflect.Ptr {
		if t := v.Type().Elem(); v.IsNil() && v.CanSet() {
			v.Set(reflect.New(t))
		}
		v = v.Elem()
	}

	if !v.CanSet() {
		return
	}

	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		f := v.Type().Field(i)

		if path == "" {
			initPtrValue(fv, f.Name)
		} else {
			initPtrValue(fv, path+"."+f.Name)
		}
	}

}

// Field visitor
//
// args:
//   path: dot separated field path
//   v: value of the field
type visitor func(path string, f reflect.StructField, v reflect.Value)

// Visit all export fields, support fields of type:
//   bool
//   time.Duration TODO
//   float64
//   int
//   int64
//   string
//   uint
//   uint64
//   pointer to above type
// If field is a struct or a pointer to a struct, it will be
// recursively visited
func visitExportedFields(p interface{}, fn visitor) {

	v := reflect.ValueOf(p)

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if !v.CanSet() {
		return
	}
	if v.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < v.NumField(); i++ {
		f := v.Type().Field(i)
		fv := v.Field(i)
		visitExportedFieldsPath(f, fv, f.Name, fn)
	}

}

func visitExportedFieldsPath(f reflect.StructField, v reflect.Value, path string, fn visitor) {

	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}

	if !v.CanSet() {
		return
	}

	switch k := v.Kind(); k {

	case reflect.Bool, reflect.Float64, reflect.Int,
		reflect.Int64, reflect.String, reflect.Uint, reflect.Uint64:
		fn(path, f, v)

	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			fv := v.Field(i)
			f := v.Type().Field(i)
			visitExportedFieldsPath(f, fv, path+"."+f.Name, fn)
		}

	default:

	}

}
