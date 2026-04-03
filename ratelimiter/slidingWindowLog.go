package ratelimiter

import "time"

//No edge burst problem. But memory grows with traffic — O(requests) per user.
type SlidingWindowLog struct {
	limit      int
	windowSize time.Duration
	log        []time.Time
}

func NewSlidingWindowLog(limit int, windowSize time.Duration) RateLimiter {
	return &SlidingWindowLog{
		limit:      limit,
		windowSize: windowSize,
		log:        []time.Time{},
	}
}

func (sw *SlidingWindowLog) Allow() bool {

	now := time.Now()
	cutOff := now.Add(-sw.windowSize)

	i := 0
	for i < len(sw.log) && sw.log[i].Before(cutOff) {
		i++
	}

	sw.log = sw.log[i:]

	if len(sw.log) > sw.limit {
		return false
	}
	sw.log = append(sw.log, now)
	return true
}
