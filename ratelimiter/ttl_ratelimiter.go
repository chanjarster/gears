package ratelimiter

// Legal results are:
//
// 1. block:false, triggered:false, ttl:0, msg:""
//
// 2. block:true, triggered:true, ttl:>0, msg:"some message"
//
// 3. block:true, triggered:false, ttl:>0, msg:"some message"
type Result struct {
	Block     bool   // true: blocked，false: passed
	Triggered bool   // first time blocking，otherwise false
	Ttl       int    // how many seconds blocking will last
	Msg       string // message recorded when first time blocking
}

// A rate limiter that will prevent further request for ttl seconds after first time request rate exceeds the limit.
type TtlRateLimiter interface {

	// same as ShouldBlock2(key, key, msg)
	ShouldBlock(key string, msg string) *Result

	// When the request rate of `key` exceeds the limit, blocking will be triggered(record on `blockKey`)
	// and last for `timeout` seconds(ttl).
	// After `timeout` seconds, `blockKey` will be released and request `key` can be passed again.
	//
	// `msg` is the message for first time blocking.
	//
	// Note: different `key` can share same `blockKey`, same `key` MUST NOT share different `blockKey`
	ShouldBlock2(key string, blockKey string, msg string) *Result

	// capacity: window capacity
	// time range the window look back
	GetWindowSizeSeconds() int
	// window capacity
	GetCapacity() int
	// how many seconds blocking will last after first time blocking happened
	GetTimeoutSeconds() int
}

// Interface for the need of runtime rate limit parameters
type TtlRateLimiterParams interface {
	GetWindowSizeSeconds() int
	GetTimeoutSeconds() int
	GetCapacity() int
}

func isParamsNotSet(params TtlRateLimiterParams) bool {
	return params.GetCapacity() <= 0 || params.GetTimeoutSeconds() <= 0 || params.GetWindowSizeSeconds() <= 0
}

func NewFixedParams(capacity, windowsSizeSec, timeoutSec int) TtlRateLimiterParams {
	return &fixedParams{
		capacity:      capacity,
		windowSizeSec: windowsSizeSec,
		timeoutSec:    timeoutSec,
	}
}

type fixedParams struct {
	timeoutSec    int
	windowSizeSec int
	capacity      int
}

func (f *fixedParams) GetWindowSizeSeconds() int {
	return f.windowSizeSec
}

func (f *fixedParams) GetTimeoutSeconds() int {
	return f.timeoutSec
}

func (f *fixedParams) GetCapacity() int {
	return f.capacity
}
