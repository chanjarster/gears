package ratelimiter

import (
	"github.com/chanjarster/gears/confs"
	"os"
	"testing"
	"time"
)

func Test_redisTtlRateLimiter_ShouldBlock(t *testing.T) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  1,
	}, nil)
	redisClient.FlushAll()
	defer redisClient.Close()

	LoadScript(redisClient)
	params := NewFixedParams(1, 1, 1)
	r := NewRedisTtlRateLimiter(redisClient, params)

	type args struct {
		key   string
		msg   string
		sleep time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		// 下面是相同key，相同block key，不同msg
		{
			args: args{"foo", "foo1", 0},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		{
			args: args{"foo", "foo2", 0},
			want: &Result{
				Block:     true,
				Triggered: true,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "foo2",
			},
		},
		{
			args: args{"foo", "foo3", 0},
			want: &Result{
				Block:     true,
				Triggered: false,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "foo2",
			},
		},
		{
			args: args{"foo", "foo4", time.Second * 2},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		// 下面3个是不同key，但是共享block key
		{
			args: args{"bar", "bar1", 0},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		{
			args: args{"bar", "bar2", 0},
			want: &Result{
				Block:     true,
				Triggered: true,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "bar2",
			},
		},
		{
			args: args{"zoo", "zoo", 0},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.sleep != 0 {
				time.Sleep(tt.args.sleep)
			}
			got := r.ShouldBlock(tt.args.key, tt.args.msg)
			if got.Block != tt.want.Block {
				t.Errorf("ShouldBlock() gotBlock = %v, want %v", got.Block, tt.want.Block)
			}
			if got.Triggered != tt.want.Triggered {
				t.Errorf("ShouldBlock() gotTriggered = %v, want %v", got.Triggered, tt.want.Triggered)
			}
			if got.Ttl != tt.want.Ttl {
				t.Errorf("ShouldBlock() gotTtl = %v, want %v", got.Ttl, tt.want.Ttl)
			}
			if got.Msg != tt.want.Msg {
				t.Errorf("ShouldBlock() gotMsg = %v, want %v", got.Msg, tt.want.Msg)
			}
		})
	}

}

func Test_redisTtlRateLimiter_ShouldBlock2(t *testing.T) {

	val, hit := os.LookupEnv("INTEGRATION_TEST")
	if !hit || val != "true" {
		t.Skip("skip integration test")
	}

	redisClient := confs.NewRedisClient(&confs.RedisConf{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		Pool:     10,
		MinIdle:  1,
	}, nil)
	redisClient.FlushAll()
	defer redisClient.Close()

	LoadScript(redisClient)
	params := NewFixedParams(1, 1, 1)
	r := NewRedisTtlRateLimiter(redisClient, params)

	type args struct {
		key      string
		blockKey string
		msg      string
		sleep    time.Duration
	}
	tests := []struct {
		name string
		args args
		want *Result
	}{
		// 下面是相同key，相同block key，不同msg
		{
			args: args{"foo", "foo", "foo1", 0},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		{
			args: args{"foo", "foo", "foo2", 0},
			want: &Result{
				Block:     true,
				Triggered: true,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "foo2",
			},
		},
		{
			args: args{"foo", "foo", "foo3", 0},
			want: &Result{
				Block:     true,
				Triggered: false,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "foo2",
			},
		},
		{
			args: args{"foo", "foo", "foo4", time.Second * 2},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		// 下面3个是不同key，但是共享block key
		{
			args: args{"bar", "bar", "bar1", 0},
			want: &Result{
				Block:     false,
				Triggered: false,
				Ttl:       0,
				Msg:       "",
			},
		},
		{
			args: args{"bar", "bar", "bar2", 0},
			want: &Result{
				Block:     true,
				Triggered: true,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "bar2",
			},
		},
		{
			args: args{"zoo", "bar", "zoo", 0},
			want: &Result{
				Block:     true,
				Triggered: false,
				Ttl:       params.GetTimeoutSeconds(),
				Msg:       "bar2",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.sleep != 0 {
				time.Sleep(tt.args.sleep)
			}
			got := r.ShouldBlock2(tt.args.key, tt.args.blockKey, tt.args.msg)
			if got.Block != tt.want.Block {
				t.Errorf("ShouldBlock2() gotBlock = %v, want %v", got.Block, tt.want.Block)
			}
			if got.Triggered != tt.want.Triggered {
				t.Errorf("ShouldBlock2() gotTriggered = %v, want %v", got.Triggered, tt.want.Triggered)
			}
			if got.Ttl != tt.want.Ttl {
				t.Errorf("ShouldBlock2() gotTtl = %v, want %v", got.Ttl, tt.want.Ttl)
			}
			if got.Msg != tt.want.Msg {
				t.Errorf("ShouldBlock2() gotMsg = %v, want %v", got.Msg, tt.want.Msg)
			}
		})
	}

}
