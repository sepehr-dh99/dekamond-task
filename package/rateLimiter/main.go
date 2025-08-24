package ratelimiter

import (
	"errors"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time // key -> timestamps
}

// NewRateLimiter creates a RateLimiter.
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{requests: make(map[string][]time.Time)}
}

// Allow checks and records a new request; returns error if limit exceeded.
func (rl *RateLimiter) Allow(key string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	now := time.Now()
	windowStart := now.Add(-10 * time.Minute)
	times := rl.requests[key]
	// Filter out old requests
	var recent []time.Time
	for _, t := range times {
		if t.After(windowStart) {
			recent = append(recent, t)
		}
	}
	if len(recent) >= 3 {
		return errors.New("rate limit exceeded")
	}
	// Record this request
	recent = append(recent, now)
	rl.requests[key] = recent
	return nil
}
