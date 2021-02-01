package confs

import (
	"os"
	"testing"
)

func TestNewEsClient(t *testing.T) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	esConf := &EsConf{
		Addrs:               "http://127.0.0.1:9200",
		MaxConnsPerHost:     10,
		MaxIdleConnsPerHost: 10,
	}
	esClient := NewEsClient(esConf, nil)
	if esClient == nil {
		t.Errorf("NewEsClient() got nil")
	}

}
