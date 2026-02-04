package api

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// RateLimiter implements a token bucket rate limiter.
type RateLimiter struct {
	mu       sync.Mutex
	tokens   float64
	maxRate  float64 // tokens per second
	capacity float64
	lastTime time.Time
}

// NewRateLimiter creates a new rate limiter.
// rate is the number of requests allowed, period is the time window.
func NewRateLimiter(rate int, period time.Duration) *RateLimiter {
	maxRate := float64(rate) / period.Seconds()

	return &RateLimiter{
		tokens:   float64(rate),
		maxRate:  maxRate,
		capacity: float64(rate),
		lastTime: time.Now(),
	}
}

// Wait blocks until a token is available or the context is cancelled.
func (r *RateLimiter) Wait(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refill()

	if r.tokens >= 1 {
		r.tokens--

		return nil
	}

	// Calculate wait time for next token
	waitDuration := time.Duration((1 - r.tokens) / r.maxRate * float64(time.Second))

	select {
	case <-ctx.Done():
		return fmt.Errorf("rate limiter wait: %w", ctx.Err())
	case <-time.After(waitDuration):
		r.refill()
		r.tokens--

		return nil
	}
}

// refill adds tokens based on elapsed time.
func (r *RateLimiter) refill() {
	now := time.Now()
	elapsed := now.Sub(r.lastTime).Seconds()
	r.lastTime = now

	r.tokens += elapsed * r.maxRate
	if r.tokens > r.capacity {
		r.tokens = r.capacity
	}
}

// Available returns the number of tokens currently available.
func (r *RateLimiter) Available() int {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.refill()

	return int(r.tokens)
}
