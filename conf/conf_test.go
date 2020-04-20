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
	"os"
	"reflect"
	"testing"
)

type conf struct {
	FieldA string
	FieldB string
	FieldC string
}

func TestLoad(t *testing.T) {

	backup := flag.NewFlagSet("backup", flag.CommandLine.ErrorHandling())
	flag.CommandLine.VisitAll(func(f *flag.Flag) {
		backup.Var(f.Value, f.Name, f.Usage)
	})

	t.Run("without yaml", func(t *testing.T) {
		defer func() {
			// restore flag.CommandLine after test
			flag.CommandLine = backup
		}()

		os.Setenv("FIELD_A", "env-a")
		os.Setenv("FIELD_B", "env-b")

		args := os.Args
		os.Args = append(os.Args, "-field-b=flag-b")

		defer func() {
			os.Args = args
			os.Unsetenv("FIELD_A")
			os.Unsetenv("FIELD_B")
		}()

		f := &conf{}
		Load(f, "")

		want := &conf{
			FieldA: "env-a",
			FieldB: "flag-b",
			FieldC: "",
		}

		if !reflect.DeepEqual(f, want) {
			t.Errorf("f = %v, want %v", f, want)
		}

	})

	t.Run("with yaml", func(t *testing.T) {
		defer func() {
			// restore flag.CommandLine after test
			flag.CommandLine = backup
		}()

		os.Setenv("FIELD_B", "env-b")
		os.Setenv("FIELD_C", "env-c")

		args := os.Args
		os.Args = append(os.Args, "-c=conf.yaml")
		os.Args = append(os.Args, "-field-c=flag-c")

		defer func() {
			os.Args = args
			os.Unsetenv("FIELD_B")
			os.Unsetenv("FIELD_C")
		}()

		cf := &conf{}
		Load(cf, "c")

		want := &conf{
			FieldA: "yaml-a",
			FieldB: "env-b",
			FieldC: "flag-c",
		}

		if !reflect.DeepEqual(cf, want) {
			t.Errorf("conf = %v, want %v", cf, want)
		}

	})

}
