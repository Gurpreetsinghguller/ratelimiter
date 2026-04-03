package ratelimiter

import (
	"fmt"
	"time"
)

// RateLimiter is the single contract all algorithms implement.
// This is the Strategy pattern from GoF — the algorithm is interchangeable.
type RateLimiter interface {
	Allow() bool
}

// basically we can pass the configs for the rate limiter here and it will return the instance of the rate limiter
func RateLimiterFactory(rltype string) RateLimiter {
	switch rltype {
	case "fixedWindow":
		return NewFixedWindow(5, 10*time.Second)
	default:
		fmt.Println("Invalid Rate Limiter Type")
		return nil
	}
}
