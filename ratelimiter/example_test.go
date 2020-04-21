package ratelimiter

import (
	"fmt"
	"time"
)

func Example_tokenBucketAtomic() {
	bucket := NewAtomicTokenBucket(100, 10)
	if bucket.Acquire() {
		fmt.Println("Acquired")
	} else {
		fmt.Println("Rejected")
	}
}

func Example_tokenBucketSync() {
	bucket := NewSyncTokenBucket(100, 10)
	if bucket.Acquire() {
		fmt.Println("Acquired")
	} else {
		fmt.Println("Rejected")
	}
}

func Example_slidingWindow() {
	bucket := NewSyncSlidingWindow(100, time.Second)
	if bucket.Acquire() {
		fmt.Println("Acquired")
	} else {
		fmt.Println("Rejected")
	}
}

func Example_fixedWindow() {
	bucket := NewSyncFixedWindow(100, time.Second)
	if bucket.Acquire() {
		fmt.Println("Acquired")
	} else {
		fmt.Println("Rejected")
	}
}