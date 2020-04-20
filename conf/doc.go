/*
A simple tool to load config from yaml, environment variables and flags by convention.

example.go:

  package main

  import (
    "fmt"
    "github.com/chanjarster/gears/conf"
  )

  type MyConf struct {
    Mysql *MysqlConf
    Redis *RedisConf
  }

  type MysqlConf struct {
    Host string
    Port int
    Database string
  }

  type RedisConf struct {
    Host string
    Port int
    Password string
  }

  func main() {
    myConf := &MyConf{}
    conf.Load(myConf, "c")
    fmt.Println(myConf)
  }

conf.yaml:

  mysql:
    host: localhost
    port: 3306
    database: test
  redis:
    host: localhost
    port: 6379
    password: test

Load from nothing:

  ./example

Load from yaml:

  ./example

Load from environment:

  ./example -c conf.yaml

Load from flag:

  ./example -redis-host=localhost

Load from env:

  MYSQL_HOST=localhost ./example

Mix them together (flag takes the highest precedence, then environment, then yaml file):

  MYSQL_PORT=3306 ./example -c conf.yaml -redis-host=localhost

Supported Field Types

Only the exported fields will be loaded. Support field types are:

  string
  bool
  int
  int64
  uint
  uint64
  float64
  string
  time.Duration, any string legal to time.ParseDuration(string)
  struct
  pointer to above types

Name Convention

This util load fields by convention. For example, field BazBoom path is Foo.Bar.BazBoom:

  type bar struct {
    BazBoom string
  }
  type foo struct {
    Bar *bar
  }
  type Conf struct {
    Foo *foo
  }

So the corresponding:

  flag name is `-foo-bar-baz-boom`
  env name is `FOO_BAR_BAZ_BOOM`

Since this util use gopkg.in/yaml.v3 to parse yaml, so the field name appear in yaml file should be lowercase,
unless you customized it the "yaml" name in the field tag:

  foo:
    bar:
      bazroom: ...

Source code and other details for the project are available at GitHub:

   https://github.com/chanjarster/gears

*/
package conf
