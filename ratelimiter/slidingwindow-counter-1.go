package ratelimiter

import (
	"sync"
	"time"
)

type SlidingWindowCounter1 struct {
	mu             sync.Mutex
	limit          int
	currentCount   int
	prevCount      int
	windowStart    time.Time
	windowDuration time.Duration
}

func (sl *SlidingWindowCounter1) Allow() bool {
	sl.mu.Lock()
	defer sl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(sl.windowStart)
	// check if the request is in current window
	if elapsed >= sl.windowDuration { // not in current window
		// check how many windows have passed since the last request
		windowsPassed := int(elapsed / sl.windowDuration)
		if windowsPassed == 1 {
			sl.prevCount = sl.currentCount
		} else {
			sl.prevCount = 0
		}
		sl.currentCount = 0
		sl.windowStart = sl.windowStart.Add(
			time.Duration(windowsPassed) * sl.windowDuration,
		)

		elapsed = now.Sub(sl.windowStart)
	}
	// How much of the previous window still overlaps
	// the last `windowDuration` seconds.
	overlap := sl.windowDuration - elapsed
	if overlap < 0 {
		overlap = 0
	}

	prevWeight := float64(overlap) / float64(sl.windowDuration)
	estimate := float64(sl.prevCount)*prevWeight + float64(sl.currentCount)
	if estimate+1 >= float64(sl.limit) {
		return false
	}
	sl.currentCount++
	return true
}
