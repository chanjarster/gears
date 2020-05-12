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
	"log"
	"time"
)

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
}

func (m *MysqlConf) String() string {
	return fmt.Sprintf("{Host: %s, Port: %d, Username: ***, Password: ***, Database: %s}",
		m.Host, m.Port, m.Database)
}

type MysqlConfigCustomizer func(mc *mysql.Config)

func NewMySqlDb(conf *MysqlConf, customizer MysqlConfigCustomizer) *sql.DB {

	mc := prepareMySqlNativeConfig(conf, customizer)

	dsn := mc.FormatDSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		log.Fatal("MySQL connection error: ", err)
		panic(err)
	}

	log.Printf("Connected to MySQL: %s:%d", conf.Host, conf.Port)
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
	mc.Params = make(map[string]string)

	if customizer != nil {
		customizer(mc)
	}
	return mc
}
