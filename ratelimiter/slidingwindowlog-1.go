package ratelimiter

import "time"

type SlidingWindowLogAgain struct {
	limit          int
	windowDuration time.Duration
	timelog        []time.Time
}

func (s *SlidingWindowLogAgain) Allow() bool {
	now := time.Now()
	cutOff := now.Add(-s.windowDuration)
	i := 0
	for i < len(s.timelog) && !s.timelog[i].After(cutOff) {
		i++
	}
	s.timelog = s.timelog[i:]
	if len(s.timelog) >= s.limit {
		return false
	}
	s.timelog = append(s.timelog, now)
	return true
}
