# api-ratelimiter

Why does a rate limiter exist?
Every API has a resource it protects — CPU, memory, database connections, money (third-party API costs). Without a rate limiter, a single misbehaving client can exhaust that resource for everyone else.


 The problem has 5 different shapes
Each algorithm answers that question differently, based on what "too many requests" means to you.
Let me show you the intuition for each one interactively:


- Fixed Window
- Sliding log
- Sliding counter
- Token bucket
- Leaky bucket


```go

// RateLimiter is the common contract all algorithms implement.
type RateLimiter interface {
    Allow() bool
}

// Middleware wires any RateLimiter into an HTTP handler.
func NewMiddleware(limiter RateLimiter, next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !limiter.Allow() {
            http.Error(w, "429 Too Many Requests", http.StatusTooManyRequests)
            return
        }
        next.ServeHTTP(w, r)
    })
}
```


## When to pick which
Scenario                                    Best choice

Simple APIs, low traffic                    ->  Fixed window
Strict per-user fairness                    ->  Sliding window 
logHigh-traffic, memory-sensitive           ->  Sliding window
counterAllow short bursts (e.g. APIs)       ->  Token bucket
Smooth traffic shaping (e.g. queues)        ->  Leaky bucket