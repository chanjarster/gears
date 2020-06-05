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
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"net/url"
	"strings"
	"time"
)

// Config keys:
//  | Environment       |  Flag              |  Description                                               |
//  |-------------------|--------------------|------------------------------------------------------------|
//  | HOST              | -host              |                                                            |
//  | PORT              | -port              |                                                            |
//  | USERNAME          | -username          |                                                            |
//  | PASSWORD          | -password          |                                                            |
//  | DATABASE          | -database          |                                                            |
//  | MAX_OPEN_CONNS    | -max-open-conns    | Maximum number of open connections to the database.        |
//  |                   |                    | If == 0 means unlimited                                    |
//  | MAX_IDLE_CONNS    | -max-idle-conns    | Maximum number of connections in the idle connection pool. |
//  |                   |                    | If == 0 no idle connections are retained                   |
//  | CONN_MAX_LIFETIME | -conn-max-lifetime | Maximum amount of time a connection may be reused.         |
//  |                   |                    | If == 0, connections are reused forever.                   |
//  | READ_TIMEOUT      | -read-timeout      | I/O read timeout                                           |
//  | WRITE_TIMEOUT     | -write-timeout     | I/O write timeout                                          |
//  | TIMEOUT           | -timeout           | Timeout for establishing connections, aka dial timeout.    |
//  | PARAMS            |                    | Connection parameters, eg, foo=1&bar=%20           |
// Note: if MysqlConf is nested in another struct, add corresponding prefix.
// more details: https://github.com/go-sql-driver/mysql
type MysqlConf struct {
	Host     string // MySQL host
	Port     int    // MySQL port
	Username string // MySQL username
	Password string // MySQL password
	Database string // MySQL database

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration

	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Timeout      time.Duration

	Params       string
}

func (m *MysqlConf) String() string {
	return fmt.Sprintf("{Host: %s, Port: %d, Username: ***, Password: ***, Database: %s, MaxOpenConns: %v, MaxIdleConns: %v, ConnMaxLifetime: %v, ReadTimeout: %v, WriteTimeout: %v, Timeout: %v, Params: %v}",
		m.Host, m.Port, m.Database, m.MaxIdleConns, m.MaxIdleConns, m.ConnMaxLifetime, m.ReadTimeout, m.WriteTimeout, m.Timeout, m.Params)
}

type MysqlConfigCustomizer func(mc *mysql.Config)

func NewMySqlDb(conf *MysqlConf, customizer MysqlConfigCustomizer) *sql.DB {

	mc := prepareMySqlNativeConfig(conf, customizer)

	dsn := mc.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		errLogger.Fatal(err)
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		errLogger.Fatal("MySQL connection error: ", err)
		panic(err)
	}

	stdLogger.Printf("Connected to MySQL: %s:%d", conf.Host, conf.Port)
	return db
}

func prepareMySqlNativeConfig(conf *MysqlConf, customizer MysqlConfigCustomizer) *mysql.Config {
	mc := mysql.NewConfig()
	mc.Timeout = conf.Timeout
	mc.ReadTimeout = conf.ReadTimeout
	mc.WriteTimeout = conf.WriteTimeout
	mc.User = conf.Username
	mc.Passwd = conf.Password
	mc.Net = "tcp"
	mc.Addr = fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	mc.DBName = conf.Database
	mc.Params = make(map[string]string,0)

	if conf.Params != "" {
		strs := strings.Split(conf.Params, "&")
		for _, str := range strs{
			p := strings.Split(str, "=")
			if len(p) != 2 {
				continue
			}
			unescape, err := url.QueryUnescape(p[1])
			if err != nil {
				continue
			}
			mc.Params[p[0]] = unescape
		}
	}

	if customizer != nil {
		customizer(mc)
	}
	return mc
}
