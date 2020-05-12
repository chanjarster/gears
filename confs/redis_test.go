package confs

import (
	"github.com/go-redis/redis/v7"
	"testing"
)

func Test_prepareRedisNativeConfig(t *testing.T) {

	redisConf := &RedisConf{
		Host:     "localhost",
		Port:     1234,
		Password: "bar",
		Pool:     1,
		MinIdle:  2,
	}
	customizer := func(ropt *redis.Options) {
		ropt.MaxRetries = 2
	}
	redisOpts := prepareRedisNativeConfig(redisConf, customizer)
	if got, want := redisOpts.Addr, "localhost:1234"; got != want {
		t.Errorf("redisOpts.Addr = %v, want %v", got, want)
	}
	if got, want := redisOpts.Password, redisConf.Password; got != want {
		t.Errorf("redisOpts.Password = %v, want %v", got, want)
	}
	if got, want := redisOpts.PoolSize, redisConf.Pool; got != want {
		t.Errorf("redisOpts.PoolSize = %v, want %v", got, want)
	}
	if got, want := redisOpts.MinIdleConns, redisConf.MinIdle; got != want {
		t.Errorf("redisOpts.MinIdleConns = %v, want %v", got, want)
	}
	if got, want := redisOpts.MaxRetries, 2; got != want {
		t.Errorf("redisOpts.MaxRetries = %v, want %v", got, want)
	}
}
