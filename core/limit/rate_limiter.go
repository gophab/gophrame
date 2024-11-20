package limit

import (
	"sync"
	"time"
)

type RateLimiter struct {
	limit     int
	interval  time.Duration
	tokens    chan struct{}
	lastReset time.Time
	mux       sync.Mutex
}

func NewRateLimiter(limit int, interval time.Duration) *RateLimiter {
	rl := &RateLimiter{
		limit:    limit,
		interval: interval,
		tokens:   make(chan struct{}, limit),
	}
	rl.lastReset = time.Now()
	for i := 0; i < limit; i++ {
		rl.tokens <- struct{}{}
	}
	go rl.refillTokens()
	return rl
}

func (rl *RateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.interval)
	for range ticker.C {
		rl.mux.Lock()
		for i := 0; i < rl.limit; i++ {
			rl.tokens <- struct{}{}
		}
		rl.lastReset = time.Now()
		rl.mux.Unlock()
	}
}

func (rl *RateLimiter) Wait() {
	select {
	case <-rl.tokens:
		// Got a token, continue
	case <-time.After(time.Until(rl.lastReset.Add(rl.interval))):
		// Waited for the next interval
	}
}
