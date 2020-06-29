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

package confs

import (
	"github.com/go-sql-driver/mysql"
	"testing"
	"time"
)

func Test_prepareMySqlNativeConfig(t *testing.T) {

	mysqlConf := &MysqlConf{
		Host:            "localhost",
		Port:            1234,
		Username:        "foo",
		Password:        "bar",
		Database:        "good",
		MaxOpenConns:    10,
		MaxIdleConns:    11,
		ConnMaxLifetime: time.Second * 5,
		ReadTimeout:     time.Second * 1,
		WriteTimeout:    time.Second * 2,
		Timeout:         time.Second * 3,
		Loc:             "Asia/Shanghai",
		ParseTime:       true,
		Params:          "foo=1&bar=2&loo=%20",
	}
	customizer := func(mc *mysql.Config) {
		mc.Params["autocommit"] = "true"
		mc.Params["charset"] = "utf8"
	}

	mc := prepareMySqlNativeConfig(mysqlConf, customizer)

	if got, want := mc.Addr, "localhost:1234"; got != want {
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
	if got, want := mc.User, mysqlConf.Username; got != want {
		t.Errorf("mc.User = %v, want %v", got, want)
	}
	if got, want := mc.Passwd, mysqlConf.Password; got != want {
		t.Errorf("mc.Passwd = %v, want %v", got, want)
	}
	if got, want := mc.DBName, mysqlConf.Database; got != want {
		t.Errorf("mc.DBName = %v, want %v", got, want)
	}
	if got, want := mc.ReadTimeout, mysqlConf.ReadTimeout; got != want {
		t.Errorf("mc.ReadTimeout = %v, want %v", got, want)
	}
	if got, want := mc.WriteTimeout, mysqlConf.WriteTimeout; got != want {
		t.Errorf("mc.WriteTimeout = %v, want %v", got, want)
	}
	if got, want := mc.Timeout, mysqlConf.Timeout; got != want {
		t.Errorf("mc.Timeout = %v, want %v", got, want)
	}
	if got, want := mc.ParseTime, mysqlConf.ParseTime; got != want {
		t.Errorf("mc.ParseTime = %v, want %v", got, want)
	}
	if got, want := mc.Loc.String(), mysqlConf.Loc; got != want {
		t.Errorf("mc.Loc = %v, want %v", got, want)
	}
	if got, want := mc.Params["autocommit"], "true"; got != want {
		t.Errorf("mc.Params[\"autocommit\"] = %v, want %v", got, want)
	}
	if got, want := mc.Params["charset"], "utf8"; got != want {
		t.Errorf("mc.Params[\"charset\"] = %v, want %v", got, want)
	}
	if got, want := mc.Params["foo"], "1"; got != want {
		t.Errorf("mc.Params[\"foo\"] = %v, want %v", got, want)
	}
	if got, want := mc.Params["bar"], "2"; got != want {
		t.Errorf("mc.Params[\"bar\"] = %v, want %v", got, want)
	}
	if got, want := mc.Params["loo"], " "; got != want {
		t.Errorf("mc.Params[\"loo\"] = %v, want %v", got, want)
	}
}

func testNewMySqlDb(t *testing.T) {
	mysqlConf := &MysqlConf{
		Host:            "localhost",
		Port:            3306,
		Username:        "platform_openapi",
		Password:        "platform_openapi",
		Database:        "platform_openapi",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnMaxLifetime: time.Second * 5,
		ReadTimeout:     time.Second * 1,
		WriteTimeout:    time.Second * 2,
		Timeout:         time.Second * 3,
		Params:          "time_zone='Asia/Shanghai'",
	}

	_ = NewMySqlDb(mysqlConf, nil)

}
