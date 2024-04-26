package confs

import (
	"database/sql"
	"fmt"
	"github.com/SkyAPM/go2sky"
	sqlPlugin "github.com/SkyAPM/go2sky-plugins/sql"
	"github.com/chanjarster/gears/simplelog"
	_ "github.com/lib/pq"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// Config keys:
//
//	| Environment       |  Flag              |  Description                                               |
//	|-------------------|--------------------|------------------------------------------------------------|
//	| HOST              | -host              |                                                            |
//	| PORT              | -port              |                                                            |
//	| USERNAME          | -username          |                                                            |
//	| PASSWORD          | -password          |                                                            |
//	| CONNECT_TIMEOUT   | -connect-timeout   | Connection timeout seconds.                                |
//	| PARAMS            |                    | Connection parameters, eg, foo=1&bar=%20                   |
//	|                   |                    | 1. All string values should be quoted with '               |
//	|                   |                    | 2. All value in PARAMS should be url.QueryEscaped          |
//	|                   |                    | more details:                                              |
//	|                   |                    | https://github.com/sijms/go-ora                            |
//	| MAX_OPEN_CONNS    | -max-open-conns    | Maximum number of open connections to the database.        |
//	|                   |                    | If == 0 means unlimited                                    |
//	| MAX_IDLE_CONNS    | -max-idle-conns    | Maximum number of connections in the idle connection pool. |
//	|                   |                    | If == 0 no idle connections are retained                   |
//	| CONN_MAX_LIFETIME | -conn-max-lifetime | Maximum amount of time a connection may be reused.         |
//	|                   |                    | If == 0, connections are reused forever.                   |
//	|                   |                    | (when using parseTime=true)                                |
//	|                   |                    | "Local" sets the system's location.                        |
//	|                   |                    | See time.LoadLocation for details.                         |
//
// Note: if PostgresqlConf is nested in another struct, add corresponding prefix.
//
// more details: https://pkg.go.dev/github.com/lib/pq
//
// https://www.postgresql.org/docs/current/libpq-connect.html#LIBPQ-CONNSTRING-URIS
type PostgresqlConf struct {
	Host           string // Postgresql host
	Port           int    // Postgresql port
	Username       string // Postgresql username
	Password       string // Postgresql password
	Database       string // Postgresql database
	ConnectTimeout int    // Connection timeout seconds
	Params         string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

func (c *PostgresqlConf) String() string {
	return fmt.Sprintf("{Host: %s, Port: %d, Username: ***, Password: ***, Database: %s, MaxOpenConns: %v, MaxIdleConns: %v, ConnMaxLifetime: %v, ConnectTimeout: %v, Params: %v}",
		c.Host, c.Port, c.Database, c.MaxIdleConns, c.MaxIdleConns,
		c.ConnMaxLifetime, c.ConnectTimeout, c.Params)
}

type PostgresqlUrlOptionsCustomizer func(options map[string]string)

func NewPostgresqlDb(conf *PostgresqlConf, customizer PostgresqlUrlOptionsCustomizer) *sql.DB {

	options := preparePgUrlOptions(conf, customizer)
	dsn := buildPgDsn(conf, options)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		simplelog.ErrLogger.Fatal(err)
		panic(err)
	}
	db.SetMaxOpenConns(conf.MaxOpenConns)
	db.SetMaxIdleConns(conf.MaxIdleConns)
	db.SetConnMaxLifetime(conf.ConnMaxLifetime)

	err = db.Ping()
	if err != nil {
		simplelog.ErrLogger.Fatal("Postgresql connection error: ", err)
		panic(err)
	}

	simplelog.StdLogger.Printf("Connected to Postgresql: %s:%d", conf.Host, conf.Port)
	return db
}

func buildPgDsn(conf *PostgresqlConf, options map[string]string) string {
	ret := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?", conf.Username, conf.Password, conf.Host, conf.Port, conf.Database)
	for key, val := range options {
		val = strings.TrimSpace(val)
		for _, temp := range strings.Split(val, ",") {
			temp = strings.TrimSpace(temp)
			ret += fmt.Sprintf("%s=%s&", key, url.QueryEscape(temp))
		}
	}
	ret = strings.TrimRight(ret, "&")
	return ret
}

func preparePgUrlOptions(conf *PostgresqlConf, customizer PostgresqlUrlOptionsCustomizer) map[string]string {
	options := map[string]string{}

	if conf.ConnectTimeout != 0 {
		options["connect_timeout"] = strconv.Itoa(conf.ConnectTimeout)
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

func NewPostgresqlDbWithTracer(conf *PostgresqlConf, tracer *go2sky.Tracer, customizer PostgresqlUrlOptionsCustomizer) *sql.DB {

	options := preparePgUrlOptions(conf, customizer)
	db, err := sqlPlugin.Open("postgres", buildPgDsn(conf, options), tracer,
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
		simplelog.ErrLogger.Fatal("Postgresql connection error: ", err)
		panic(err)
	}

	simplelog.StdLogger.Printf("Connected to Postgresql: %s:%d", conf.Host, conf.Port)
	return db.DB

}
