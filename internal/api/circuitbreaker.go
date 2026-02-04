package api

import (
	"sync"
	"time"
)

// CircuitBreaker implements a simple circuit breaker pattern.
type CircuitBreaker struct {
	mu               sync.Mutex
	failures         int
	maxFailures      int
	lastFailure      time.Time
	cooldownDuration time.Duration
	open             bool
}

// NewCircuitBreaker creates a new circuit breaker.
// maxFailures is the number of consecutive failures before opening.
// cooldownDuration is how long to wait before attempting again.
func NewCircuitBreaker(maxFailures int, cooldownDuration time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		maxFailures:      maxFailures,
		cooldownDuration: cooldownDuration,
	}
}

// IsOpen returns whether the circuit is open (failing).
func (cb *CircuitBreaker) IsOpen() bool {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if !cb.open {
		return false
	}

	// Check if cooldown has passed
	if time.Since(cb.lastFailure) > cb.cooldownDuration {
		cb.open = false
		cb.failures = 0

		return false
	}

	return true
}

// RecordSuccess records a successful request.
func (cb *CircuitBreaker) RecordSuccess() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures = 0
	cb.open = false
}

// RecordFailure records a failed request.
func (cb *CircuitBreaker) RecordFailure() {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	cb.failures++
	cb.lastFailure = time.Now()

	if cb.failures >= cb.maxFailures {
		cb.open = true
	}
}

// Failures returns the current failure count.
func (cb *CircuitBreaker) Failures() int {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	return cb.failures
}
