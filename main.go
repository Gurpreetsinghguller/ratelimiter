package main

import rm "github.com/Gurpreetsinghguller/ratelimiter/ratelimiter"

type RatelimiterConfig struct {
	Type           string
	Limit          int
	WindowDuration int
}
type Config struct {
	RateLimiter RatelimiterConfig
}

func main() {
	// we can read config from a yml file
	// for now we will hardcode the config
	config := Config{
		RateLimiter: RatelimiterConfig{
			Type:           "fixedWindow",
			Limit:          5,
			WindowDuration: 10,
		},
	}
	rm.RateLimiterFactory(config.RateLimiter.Type)

}
