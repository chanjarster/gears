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
	"fmt"
	"reflect"
	"testing"
	"time"
)

func Test_initStruct(t *testing.T) {

	t.Run("nil", func(t *testing.T) {
		initStruct(nil)
	})

	t.Run("nil ptr", func(t *testing.T) {
		initStruct((*foo)(nil))
	})

	t.Run("*foo nil", func(t *testing.T) {
		var p *foo
		initStruct(p)
		if p != nil {
			t.Errorf("p = %v, want nil", p)
		}
	})

	t.Run("*int", func(t *testing.T) {
		i := 1
		initStruct(&i)
		if i != 1 {
			t.Errorf("i = %v, want 1", i)
		}
	})

	t.Run("non ptr", func(t *testing.T) {
		i := 1
		initStruct(i)
		if i != 1 {
			t.Errorf("i = %v, want 1", i)
		}
	})

	t.Run("ptr", func(t *testing.T) {
		p := &foo{}
		initStruct(p)

		want := &foo{
			i:    0,
			ip:   nil,
			bar:  bar{},
			barP: nil,
			I:    0,
			Ip:   new(int),
			D:    0,
			Dp:   new(time.Duration),
			Bar: bar{
				I:  0,
				Ip: new(int),
			},
			BarP: &bar{
				I:  0,
				Ip: new(int),
			},
			baz: struct {
				i int
			}{},
			Baz: struct {
				i int
			}{},
		}

		if !reflect.DeepEqual(p, want) {
			t.Errorf("p = %v, want %v", p, want)
		}
	})

}

type fakeVisitor []string

func (fv *fakeVisitor) visit(path string, f reflect.StructField, v reflect.Value) {
	fmt.Println(path)
	*fv = append(*fv, path)
}

func Test_visitExportedFields(t *testing.T) {

	t.Run("", func(t *testing.T) {

		f := &foo{}
		initStruct(f)

		vis := fakeVisitor(make([]string, 0))
		visP := &vis
		visitExportedFields(f, visP.visit)

		want := []string{"I", "Ip", "Bar.I", "Bar.Ip", "BarP.I", "BarP.Ip", "D", "Dp"}
		if !reflect.DeepEqual([]string(*visP), want) {
			t.Errorf("visP = %v, want %v", *visP, want)
		}
	})

}