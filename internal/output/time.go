package output

import (
	"fmt"
	"time"
)

var brusselsTZ *time.Location

func init() {
	var err error

	brusselsTZ, err = time.LoadLocation("Europe/Brussels")
	if err != nil {
		brusselsTZ = time.Local
	}
}

// BrusselsTimezone returns the Europe/Brussels timezone.
func BrusselsTimezone() *time.Location {
	return brusselsTZ
}

// ParseAPITime parses a time string from the De Lijn API.
// Format: "2006-01-02T15:04:05"
func ParseAPITime(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", s, brusselsTZ)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse time %q: %w", s, err)
	}

	return t, nil
}

// FormatTime formats a time for display (HH:MM).
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return "--:--"
	}

	return t.In(brusselsTZ).Format("15:04")
}

// FormatTimeWithSeconds formats a time for display (HH:MM:SS).
func FormatTimeWithSeconds(t time.Time) string {
	if t.IsZero() {
		return "--:--:--"
	}

	return t.In(brusselsTZ).Format("15:04:05")
}

// FormatRelative formats a time relative to now.
func FormatRelative(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	now := time.Now()
	diff := t.Sub(now)

	if diff < 0 {
		return "now"
	}

	minutes := int(diff.Minutes())
	if minutes == 0 {
		return "now"
	}

	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60

	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, mins)
}
