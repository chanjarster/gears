package confs

import (
	"testing"
	"time"
)

func testNewPostgresqlDb(t *testing.T) {

	/*
	  docker run --rm -p 5432:5432 \
	    -e POSTGRES_USER=post \
	    -e POSTGRES_PASSWORD=test \
	    -e POSTGRES_DB=test2 \
	    postgres:alpine3.19
	*/
	oracleConf := &PostgresqlConf{
		Host:            "localhost",
		Port:            5432,
		Username:        "post",
		Password:        "test",
		Database:        "test2",
		Params:          "client_encoding=UTF8&sslmode=disable",
		MaxOpenConns:    1,
		MaxIdleConns:    1,
		ConnectTimeout:  3,
		ConnMaxLifetime: time.Second * 5,
	}

	_ = NewPostgresqlDb(oracleConf, nil)

}
