package ratelimiter

import (
	"context"
	"sync"
	"time"
)

// is my leaky bucket implementation correct?
// Answer
type LeakyBucket struct {
	limitermu sync.Mutex
	count     int

	limit    int
	leakTime time.Duration
	leakSize int
}

func NewLeakyBucket(ctx context.Context, limit int, leakTime time.Duration, leakSize int) RateLimiter {
	if leakTime <= 0 {
		// use default
	}
	if leakSize <= 0 {
		// use default
	}
	if limit <= 0 {
		// use default
	}
	lb := &LeakyBucket{
		limit:    limit,
		leakTime: leakTime,
		leakSize: leakSize,
		count:    0,
	}
	ticker := time.NewTicker(leakTime)
	go func(ctx context.Context, ticker *time.Ticker) {
		defer ticker.Stop()
		lb.leak(ctx, ticker)
	}(ctx, ticker)
	return lb
}

func (l *LeakyBucket) leak(ctx context.Context, ticker *time.Ticker) {
	for {
		select {
		case <-ticker.C:
			l.limitermu.Lock()
			if l.count > 0 {
				l.count -= l.leakSize
				if l.count < 0 {
					l.count = 0
				}
			}
			l.limitermu.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

func (l *LeakyBucket) Allow() bool {
	l.limitermu.Lock()
	defer l.limitermu.Unlock()

	if l.count < l.limit {
		l.count++
		return true
	}
	return false
}
