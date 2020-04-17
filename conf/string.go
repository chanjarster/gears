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

// Foo.BarBaz -> foo-bar-baz
func flagStyle(path string) string {
	bs := []byte(path)
	rs := make([]byte, 0, len(bs))

	prevDot := true
	for i := 0; i < len(bs); i++ {
		b := bs[i]
		if isUpperASCII(b) {
			if !prevDot {
				rs = append(rs, '-')
			}
			rs = append(rs, lowerASCII(b))
			prevDot = false
		} else if b == '.' {
			rs = append(rs, '-')
			prevDot = true
		} else {
			rs = append(rs, b)
			prevDot = false
		}
	}

	return string(rs)
}

// Foo.BarBaz -> FOO_BAR_BAZ
func envStyle(path string) string {

	bs := []byte(path)
	rs := make([]byte, 0, len(bs))

	prevDot := true
	for i := 0; i < len(bs); i++ {
		b := bs[i]
		if isUpperASCII(b) {
			if !prevDot {
				rs = append(rs, '_')
			}
			rs = append(rs, b)
			prevDot = false
		} else if b == '.' {
			rs = append(rs, '_')
			prevDot = true
		} else {
			rs = append(rs, upperASCII(b))
			prevDot = false
		}
	}

	return string(rs)

}

func lowerASCII(b byte) byte {
	if isUpperASCII(b) {
		return b + ('a' - 'A')
	}
	return b
}

func upperASCII(b byte) byte {
	if isLowerASCII(b) {
		return b - ('a' - 'A')
	}
	return b
}

func isLowerASCII(b byte) bool {
	return 'a' <= b && b <= 'z'
}

func isUpperASCII(b byte) bool {
	return 'A' <= b && b <= 'Z'
}
