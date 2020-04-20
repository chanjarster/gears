package gears

import "time"

type NowFunc func() int64

func SysNow() int64 {
	return time.Now().UnixNano()
}
