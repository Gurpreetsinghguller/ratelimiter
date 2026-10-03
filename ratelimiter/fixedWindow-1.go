package ratelimiter

import (
	"sync"
	"time"
)

type FixedWindowRateLimiter struct {
	mu                  sync.Mutex
	currentRequestCount int
	totalRequestCount   int
	windowStart         time.Time
	windowDuration      time.Duration
}

func NewFixedWindowRateLimiter() RateLimiter {
	return &FixedWindowRateLimiter{
		currentRequestCount: 0,
		totalRequestCount:   10,
		windowStart:         time.Now(),
		windowDuration:      1 * time.Second,
	}
}

func (fx *FixedWindowRateLimiter) Allow() bool {
	fx.mu.Lock()
	defer fx.mu.Unlock()

	now := time.Now()

	// check if this request is in the same window as previous requests
	if now.Sub(fx.windowStart) >= fx.windowDuration { // current window passed
		// reset the window and request count
		fx.windowStart = now
		fx.currentRequestCount = 0
	}

	fx.currentRequestCount++
	return fx.currentRequestCount <= fx.totalRequestCount
}
