package confs

import (
	"fmt"
	"github.com/chanjarster/gears/simplelog"
	"github.com/elastic/go-elasticsearch/v7"
	"io/ioutil"
	"net/http"
	"strings"
)

type EsConf struct {
	Addrs               string
	Username            string
	Password            string
	MaxConnsPerHost     int
	MaxIdleConnsPerHost int
}

func (m *EsConf) String() string {
	return fmt.Sprintf("{Addrs: %s, Username: ***, Password: ***, MaxConnsPerHost: %v, MaxIdleConnsPerHost: %v}",
		m.Addrs, m.MaxConnsPerHost, m.MaxIdleConnsPerHost)
}

type EsConfigCustomizer func(mc *elasticsearch.Config)

func NewEsClient(conf *EsConf, customizer EsConfigCustomizer) *elasticsearch.Client {

	esConf := &elasticsearch.Config{
		Addresses: strings.Split(conf.Addrs, ","),
		Username:  conf.Username,
		Password:  conf.Password,
		Transport: &http.Transport{
			MaxConnsPerHost:     conf.MaxConnsPerHost,
			MaxIdleConnsPerHost: conf.MaxIdleConnsPerHost,
		},
	}

	if customizer != nil {
		customizer(esConf)
	}

	es, err := elasticsearch.NewClient(*esConf)
	if err != nil {
		panic(err)
	}

	pingRes, err := es.Ping()
	if err != nil {
		panic(err)
	}
	defer pingRes.Body.Close()

	if pingRes.IsError() {
		body, err := ioutil.ReadAll(pingRes.Body)
		if err != nil {
			panic(err)
		}
		panic(fmt.Sprintf("Connected to Elasticsearch error: %s", string(body)))
	}
	simplelog.StdLogger.Printf("Connected to Elasticsearch: %s", conf.Addrs)
	return es

}
