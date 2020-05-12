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
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
	if got, want := mc.WriteTimeout, mysqlConf.WriteTimeout; got != want {
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
	if got, want := mc.Timeout, mysqlConf.Timeout; got != want {
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
	if got, want := mc.Params["autocommit"], "true"; got != want {
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
	if got, want := mc.Params["charset"], "utf8"; got != want {
		t.Errorf("mc.Addr = %v, want %v", got, want)
	}
}
