package ratelimiter

import (
	"sync"
	"time"
)

// The hidden bug: a user can fire 5 requests at 4.99s and 5 more at 5.01s — 10 requests in 20ms, all within "limit". This is the edge burst problem.
type FixedWindow struct {
	limitMu        sync.Mutex
	limit          int
	count          int
	windowStart    time.Time
	windowDuration time.Duration
}

func NewFixedWindow(limit int, windowDuration time.Duration) RateLimiter {
	return &FixedWindow{
		limitMu:        sync.Mutex{},
		limit:          limit,
		count:          0,
		windowStart:    time.Now(),
		windowDuration: windowDuration,
	}
}

func (fw *FixedWindow) Allow() bool {
	fw.limitMu.Lock()
	defer fw.limitMu.Unlock()
	now := time.Now()
	if now.Sub(fw.windowStart) >= fw.windowDuration {
		fw.windowStart = now
		fw.count = 0
	}
	if fw.count >= fw.limit {
		return false
	}
	fw.count++
	return true
}
