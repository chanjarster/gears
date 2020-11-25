package ratelimiter

import (
	"github.com/chanjarster/gears/simplelog"
	"github.com/go-redis/redis/v7"
	"time"
)

const (
	/*
			if exists(key:block)
			  return true, false

		  size = llen(key)
			if size < capacity
			  rpush key now
			  return false, false

		  if size - capacity > 0
		    lpop key (size - capacity) elements

			oldest = lrange key 0 0
			if now - oldest > window size
				rpush key now
		    lpop key
				return false, false

			SET key:block 1 EX timeout NX
			return true, true
	*/
	prefix = "_rl:"
	suffix = ":bl"

	// return block, triggered, ttl, msg
	script = `local key = KEYS[1]
local keyb = KEYS[2]
local cap = tonumber(ARGV[1])
local win = tonumber(ARGV[2])
local exp = tonumber(ARGV[3])
local now = tonumber(ARGV[4])
local msg = ARGV[5]

if cap <= 0 or win <= 0 or exp <= 0
then
  return {0, 0, 0, ''}
end

if redis.call('EXISTS', keyb) == 1
then
  local ttl = redis.call('TTL', keyb)
  local omsg = redis.call('GET', keyb)
  return {1, 0, ttl, omsg}
end

local size = redis.call('LLEN', key)
if size < cap
then
  redis.call('RPUSH', key, now)
  return {0, 0, 0, ''}
end

local diff = size - cap
if diff > 0
then
  for i = diff, 1, -1
  do
    redis.call('LPOP', key)
  end
end

local list = redis.call('LRANGE', key, 0, 0)

if table.getn(list) == 0
then
  redis.call('RPUSH', key, now)
  redis.call('EXPIRE', key, 3600)
  return {0, 0, 0, ''}
end

local oldest = tonumber(list[1])
if now - oldest > win
then
  redis.call('RPUSH', key, now)
  redis.call('LPOP', key)
  redis.call('EXPIRE', key, 3600)
  return {0, 0, 0, ''}
end

redis.call('SET', keyb, msg, 'EX', exp, 'NX')
return {1, 1, exp, msg}
`
)

var (
	scriptSha1 = ""
)

func NewRedisTtlRateLimiter(redisClient *redis.Client, params TtlRateLimiterParams) TtlRateLimiter {
	r := &redisTtlRateLimiter{
		params:      params,
		redisClient: redisClient,
	}
	return r
}

type redisTtlRateLimiter struct {
	params      TtlRateLimiterParams
	redisClient *redis.Client
}

func (r *redisTtlRateLimiter) ShouldBlock(key string, msg string) *Result {
	return r.ShouldBlock2(key, key, msg)
}

func (r *redisTtlRateLimiter) ShouldBlock2(key string, blockKey string, msg string) *Result {

	result := &Result{}

	if isParamsNotSet(r.params) {
		return result
	}

	key = prefix + key
	blockKey = blockKey + suffix

	//local key = KEYS[1]
	//local keyb = KEYS[2]
	//local cap = tonumber(ARGV[1])
	//local win = tonumber(ARGV[2])
	//local exp = tonumber(ARGV[3])
	//local now = tonumber(ARGV[4])
	//local msg = ARGV[5]

	now := time.Now().UnixNano() / int64(time.Second)

	raw, err := r.redisClient.EvalSha(
		scriptSha1,
		[]string{key, blockKey},
		r.params.GetCapacity(),
		r.params.GetWindowSizeSeconds(),
		r.params.GetTimeoutSeconds(),
		now,
		msg,
	).Result()

	if err != nil {
		simplelog.ErrLogger.Println("eval script ", scriptSha1, "error", err)
		return result
	}

	arr := raw.([]interface{})
	result.Block = arr[0].(int64) == 1
	result.Triggered = arr[1].(int64) == 1
	result.Ttl = int(arr[2].(int64))
	result.Msg = arr[3].(string)
	return result

}

func (r *redisTtlRateLimiter) GetWindowSizeSeconds() int {
	return r.params.GetWindowSizeSeconds()
}

func (r *redisTtlRateLimiter) GetTimeoutSeconds() int {
	return r.params.GetTimeoutSeconds()
}

func (r *redisTtlRateLimiter) GetCapacity() int {
	return r.params.GetCapacity()
}

func LoadScript(redisClient *redis.Client) {
	sha1, err := redisClient.ScriptLoad(script).Result()
	if err != nil {
		panic(err)
	}
	scriptSha1 = sha1
}
