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
	"testing"
)

func Test_envStyle(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{path: "Foo.BarBaz"}, want: "FOO_BAR_BAZ"},
		{args: args{path: "Foo.barBaz"}, want: "FOO_BAR_BAZ"},
		{args: args{path: "foo.barBaz"}, want: "FOO_BAR_BAZ"},
		{args: args{path: "foo.barbaz"}, want: "FOO_BARBAZ"},
		{args: args{path: "FOO.BARBAZ"}, want: "F_O_O_B_A_R_B_A_Z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := envStyle(tt.args.path); got != tt.want {
				t.Errorf("envStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_flagStyle(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{args: args{path: "Foo.BarBaz"}, want: "foo-bar-baz"},
		{args: args{path: "Foo.barBaz"}, want: "foo-bar-baz"},
		{args: args{path: "foo.barBaz"}, want: "foo-bar-baz"},
		{args: args{path: "foo.barbaz"}, want: "foo-barbaz"},
		{args: args{path: "FOO.BARBAZ"}, want: "f-o-o-b-a-r-b-a-z"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := flagStyle(tt.args.path); got != tt.want {
				t.Errorf("flagStyle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lowerASCII(t *testing.T) {
	type args struct {
		b byte
	}

	type test struct {
		name string
		args args
		want byte
	}

	tests := make([]test, 0)

	for i := 'A'; i <= 'Z'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i + ('a' - 'A'))})
	}
	for i := 'a'; i <= 'z'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i)})
	}
	for i := '0'; i <= '9'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i)})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lowerASCII(tt.args.b); got != tt.want {
				t.Errorf("lowerASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_upperASCII(t *testing.T) {
	type args struct {
		b byte
	}
	type test struct {
		name string
		args args
		want byte
	}
	tests := make([]test, 0)

	for i := 'a'; i <= 'z'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i - ('a' - 'A'))})
	}
	for i := 'A'; i <= 'Z'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i)})
	}
	for i := '0'; i <= '9'; i++ {
		tests = append(tests, test{args: args{byte(i)}, want: byte(i)})
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := upperASCII(tt.args.b); got != tt.want {
				t.Errorf("upperASCII() = %v, want %v", got, tt.want)
			}
		})
	}
}
