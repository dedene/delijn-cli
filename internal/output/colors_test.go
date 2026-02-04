package output

import (
	"testing"
)

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"zero", 0, "<1m"},
		{"30 seconds", 30, "<1m"},
		{"59 seconds", 59, "<1m"},
		{"60 seconds", 60, "1m"},
		{"90 seconds", 90, "1m"},
		{"120 seconds", 120, "2m"},
		{"5 minutes", 300, "5m"},
		{"10 minutes", 600, "10m"},
		{"100 minutes", 6000, "100m"},
		{"999 minutes", 59940, "999m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.seconds)
			if got != tt.expected {
				t.Errorf("formatDuration(%d) = %q, want %q", tt.seconds, got, tt.expected)
			}
		})
	}
}

func TestFormatDelayNoColor(t *testing.T) {
	// Ensure consistent output with noColor=true
	SetNoColor(true)
	defer SetNoColor(false)

	tests := []struct {
		name     string
		seconds  int
		expected string
	}{
		{"early 2 min", -120, "-2m"},
		{"early 30 sec", -30, "-<1m"},
		{"on time", 0, "on time"},
		{"late 1 min", 60, "+1m"},
		{"late 2 min", 120, "+2m"},
		{"late 5 min", 300, "+5m"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatDelay(tt.seconds)
			if got != tt.expected {
				t.Errorf("FormatDelay(%d) = %q, want %q", tt.seconds, got, tt.expected)
			}
		})
	}
}
