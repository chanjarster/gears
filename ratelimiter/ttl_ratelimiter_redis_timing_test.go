package ratelimiter

import (
	"github.com/chanjarster/gears/confs"
	"math/rand"
	"os"
	"testing"
)

var (
	_dict = []string{
		"foo", "bar", "zoo", "alice",
		"daniel", "fox", "dog", "fish",
		"cat", "mouse", "kangaroo",
		"dwarf", "insect", "zebra", "sheep",
		"pig", "bear", "duck", "lizard", "shepard",
	}
)

func randWord() string {
	return _dict[rand.Intn(len(_dict))]
}

func Benchmark_redisTtlRateLimiter_ShouldBlock(b *testing.B) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		b.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  5,
	}, nil)
	redisClient.FlushAll()
	defer redisClient.Close()

	LoadScript(redisClient)
	loginParams := NewFixedParams(10, 2, 1)
	r := NewRedisTtlRateLimiter(redisClient, loginParams)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			r.ShouldBlock(randWord(), randWord())
		}
	})

}
