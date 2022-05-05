package confs

import (
	"testing"
	"time"
)

func testNewOracleDb(t *testing.T) {

	/*
		docker run --rm -p 1521:1521 \
		  wnameless/oracle-xe-11g-r2:latest
	*/
	oracleConf := &OracleConf{
		Host:            "localhost",
		Port:            1521,
		Username:        "system",
		Password:        "oracle",
		Service:         "",
		Sid:             "xe",
		Params:          "",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnectTimeout:  3,
		ConnMaxLifetime: time.Second * 5,
	}

	_ = NewOracleDb(oracleConf, nil)

}
