package confs

import (
	"database/sql"
	"fmt"
	"github.com/SkyAPM/go2sky"
	sqlPlugin "github.com/SkyAPM/go2sky-plugins/sql"
	"github.com/chanjarster/gears/simplelog"
	go_ora "github.com/sijms/go-ora/v2"
	"net/url"
	"strconv"
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
//  | SERVICE           | -service           | Oracle service name                                        |
//  | SID               | -sid               | Oracle SID, you need to provide either service or sid option  |
//  | SERVERS           | -servers           | Add more servers (nodes) in case of RAC in form of "srv1:port,srv2:port  |
//  | CONNECT_TIMEOUT   | -connect-timeout   | Connection timeout seconds.                                |
//  | PARAMS            |                    | Connection parameters, eg, foo=1&bar=%20                   |
//  |                   |                    | 1. All string values should be quoted with '               |
//  |                   |                    | 2. All value in PARAMS should be url.QueryEscaped          |
//  |                   |                    | more details:                                              |
//  |                   |                    | https://github.com/sijms/go-ora                            |
//  | MAX_OPEN_CONNS    | -max-open-conns    | Maximum number of open connections to the database.        |
//  |                   |                    | If == 0 means unlimited                                    |
//  | MAX_IDLE_CONNS    | -max-idle-conns    | Maximum number of connections in the idle connection pool. |
//  |                   |                    | If == 0 no idle connections are retained                   |
//  | CONN_MAX_LIFETIME | -conn-max-lifetime | Maximum amount of time a connection may be reused.         |
//  |                   |                    | If == 0, connections are reused forever.                   |
//  |                   |                    | (when using parseTime=true)                                |
//  |                   |                    | "Local" sets the system's location.                        |
//  |                   |                    | See time.LoadLocation for details.                         |
//
// Note: if OracleConf is nested in another struct, add corresponding prefix.
//
// more details: https://github.com/sijms/go-ora
type OracleConf struct {
	Host           string // Oracle host
	Port           int    // Oracle port
	Username       string // Oracle username
	Password       string // Oracle password
	Service        string // Oracle service
	Sid            string // or Oracle SID
	Servers        string // Add more servers (nodes) in case of RAC in form of "srv1:port,srv2:port
	ConnectTimeout int    // Connection timeout seconds
	Params         string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c *OracleConf) String() string {
	return fmt.Sprintf("{Host: %s, Port: %d, Username: ***, Password: ***, Service: %s, Sid: %s, Servers: %s, MaxOpenConns: %v, MaxIdleConns: %v, ConnMaxLifetime: %v, ConnectTimeout: %v, Params: %v}",
		c.Host, c.Port, c.Service, c.Sid, c.Servers, c.MaxIdleConns, c.MaxIdleConns,
		c.ConnMaxLifetime, c.ConnectTimeout, c.Params)
}

type OracleUrlOptionsCustomizer func(options map[string]string)

func NewOracleDb(conf *OracleConf, customizer OracleUrlOptionsCustomizer) *sql.DB {

	options := prepareUrlOptions(conf, customizer)
	connectionString := go_ora.BuildUrl(conf.Host, conf.Port, conf.Service, conf.Username, conf.Password, options)
	db, err := sql.Open("oracle", connectionString)
	if err != nil {
		simplelog.ErrLogger.Fatal(err)
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		simplelog.ErrLogger.Fatal("Oracle connection error: ", err)
		panic(err)
	}

	simplelog.StdLogger.Printf("Connected to Oracle: %s:%d", conf.Host, conf.Port)
	return db
}

func prepareUrlOptions(conf *OracleConf, customizer OracleUrlOptionsCustomizer) map[string]string {
	options := map[string]string{}
	if conf.Sid == "" && conf.Service == "" {
		panic("You need to provide either service or sid option")
	}
	if conf.Sid != "" {
		options["SID"] = conf.Sid
	}
	if conf.Servers != "" {
		options["server"] = conf.Servers
	}
	if conf.ConnectTimeout != 0 {
		options["CONNECTION TIMEOUT"] = strconv.Itoa(conf.ConnectTimeout)
	}
	if customizer != nil {
		customizer(options)
	}

	if conf.Params != "" {
		strs := strings.Split(conf.Params, "&")
		for _, str := range strs {
			p := strings.Split(str, "=")
			if len(p) != 2 {
				continue
			}
			unescape, err := url.QueryUnescape(p[1])
			if err != nil {
				continue
			}
			options[p[0]] = unescape
		}
	}

	return options
}
func NewOracleDbWithTracer(conf *OracleConf, tracer *go2sky.Tracer, customizer OracleUrlOptionsCustomizer) *sql.DB {

	options := prepareUrlOptions(conf, customizer)
	connectionString := go_ora.BuildUrl(conf.Host, conf.Port, conf.Service, conf.Username, conf.Password, options)
	db, err := sqlPlugin.Open("oracle", connectionString, tracer,
		sqlPlugin.WithSQLDBType(sqlPlugin.IPV4),
		sqlPlugin.WithQueryReport(),
		sqlPlugin.WithParamReport(),
	)
	if err != nil {
		simplelog.ErrLogger.Fatal(err)
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		simplelog.ErrLogger.Fatal("Oracle connection error: ", err)
		panic(err)
	}

	simplelog.StdLogger.Printf("Connected to Oracle: %s:%d", conf.Host, conf.Port)
	return db.DB

}
