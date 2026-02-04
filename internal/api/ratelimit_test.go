package api

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiterInitial(t *testing.T) {
	rl := NewRateLimiter(10, time.Second)

	available := rl.Available()
	if available != 10 {
		t.Errorf("initial available should be 10, got %d", available)
	}
}

func TestRateLimiterConsumesTokens(t *testing.T) {
	rl := NewRateLimiter(5, time.Second)
	ctx := context.Background()

	// Consume all tokens
	for range 5 {
		if err := rl.Wait(ctx); err != nil {
			t.Fatalf("Wait() error: %v", err)
		}
	}

	available := rl.Available()
	if available != 0 {
		t.Errorf("available should be 0 after consuming all tokens, got %d", available)
	}
}

func TestRateLimiterRefills(t *testing.T) {
	// 100 tokens per second = 1 token per 10ms
	rl := NewRateLimiter(100, time.Second)
	ctx := context.Background()

	// Consume all tokens
	for range 100 {
		_ = rl.Wait(ctx)
	}

	// Wait for refill (should get ~5 tokens in 50ms at 100/sec)
	time.Sleep(50 * time.Millisecond)

	available := rl.Available()
	if available < 3 || available > 7 {
		t.Errorf("expected ~5 tokens after 50ms, got %d", available)
	}
}

func TestRateLimiterContextCancel(t *testing.T) {
	rl := NewRateLimiter(1, time.Second)
	ctx, cancel := context.WithCancel(context.Background())

	// Consume the only token
	_ = rl.Wait(ctx)

	// Cancel context before next wait completes
	cancel()

	err := rl.Wait(ctx)
	if err == nil {
		t.Error("expected error when context is cancelled")
	}
}
