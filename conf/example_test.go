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
	"os"
)

type MyConf struct {
	Mysql *MysqlConf
	Redis *RedisConf
}

func (c *MyConf) String() string {
	return fmt.Sprint("{mysql: ", c.Mysql, ", redis: ", c.Redis, "}")
}

type MysqlConf struct {
	Host     string
	Port     int
	Database string
}

func (c *MysqlConf) String() string {
	return fmt.Sprint("{host: ", c.Host, ", port: ", c.Port, ", database: ", c.Database, "}")
}

type RedisConf struct {
	Host     string
	Port     int
	Password string
}

func (c *RedisConf) String() string {
	return fmt.Sprint("{host: ", c.Host, ", port: ", c.Port, ", password: ", c.Password, "}")
}

func ExampleLoad() {
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", "3360")
	os.Setenv("MYSQL_DATABASE", "test")
	os.Setenv("REDIS_HOST", "localhost")
	os.Setenv("REDIS_PORT", "6379")
	os.Setenv("REDIS_PASSWORD", "foobar")
	myConf := &MyConf{}
	Load(myConf, "conf")
	fmt.Println(myConf)
	// Output: {mysql: {host: localhost, port: 3360, database: test}, redis: {host: localhost, port: 6379, password: foobar}}
}
