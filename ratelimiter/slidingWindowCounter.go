package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter struct {
	mu         sync.Mutex
	limit      int
	windowSize time.Duration

	prevCount   int       // requests in the completed previous window
	currCount   int       // requests so far in the current window
	windowStart time.Time // when the current window began
}

func NewSlidingWindowCounter(limit int, windowSize time.Duration) RateLimiter {
	return &SlidingWindowCounter{
		mu:          sync.Mutex{},
		limit:       limit,
		windowSize:  windowSize,
		currCount:   0,
		prevCount:   0,
		windowStart: time.Now(),
	}
}

func (sc *SlidingWindowCounter) Allow() bool {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	now := time.Now()
	// window start hone k baad kab request aayi hai, uska time nikalna
	elapsed := now.Sub(sc.windowStart) // elapsed — how far into the current window are we right now?

	// Has the current window fully expired?
	if elapsed >= sc.windowSize {
		if elapsed >= 2*sc.windowSize {
			// So much time passed that even prevCount is irrelevant
			sc.prevCount = 0
		} else {
			sc.prevCount = sc.currCount
		}
		sc.currCount = 0
		sc.windowStart = now
		elapsed = 0
	}

	prevWeight := 1.0 - elapsed.Seconds()/sc.windowSize.Seconds()
	estimate := float64(sc.prevCount)*prevWeight + float64(sc.currCount)

	if estimate >= float64(sc.limit) {
		return false
	}

	sc.currCount++
	return true

}
