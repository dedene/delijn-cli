package api

import (
	"testing"
	"time"
)

func TestCircuitBreakerInitial(t *testing.T) {
	cb := NewCircuitBreaker(5, 30*time.Second)

	if cb.IsOpen() {
		t.Error("new circuit breaker should be closed")
	}

	if cb.Failures() != 0 {
		t.Errorf("new circuit breaker should have 0 failures, got %d", cb.Failures())
	}
}

func TestCircuitBreakerOpensAfterMaxFailures(t *testing.T) {
	cb := NewCircuitBreaker(3, 30*time.Second)

	// Record 2 failures - should still be closed
	cb.RecordFailure()
	cb.RecordFailure()

	if cb.IsOpen() {
		t.Error("circuit breaker should still be closed after 2 failures")
	}

	// Third failure should open it
	cb.RecordFailure()

	if !cb.IsOpen() {
		t.Error("circuit breaker should be open after 3 failures")
	}
}

func TestCircuitBreakerSuccessResets(t *testing.T) {
	cb := NewCircuitBreaker(3, 30*time.Second)

	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordSuccess()

	if cb.Failures() != 0 {
		t.Errorf("failures should be 0 after success, got %d", cb.Failures())
	}

	if cb.IsOpen() {
		t.Error("circuit breaker should be closed after success")
	}
}

func TestCircuitBreakerCooldown(t *testing.T) {
	cb := NewCircuitBreaker(2, 50*time.Millisecond)

	cb.RecordFailure()
	cb.RecordFailure()

	if !cb.IsOpen() {
		t.Error("circuit breaker should be open")
	}

	// Wait for cooldown
	time.Sleep(60 * time.Millisecond)

	if cb.IsOpen() {
		t.Error("circuit breaker should be closed after cooldown")
	}

	if cb.Failures() != 0 {
		t.Errorf("failures should be reset after cooldown, got %d", cb.Failures())
	}
}
