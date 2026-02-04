package output

import (
	"fmt"

	"github.com/muesli/termenv"
)

var (
	profile = termenv.ColorProfile()
	noColor = false
)

// SetNoColor disables color output.
func SetNoColor(disable bool) {
	noColor = disable
	if disable {
		profile = termenv.Ascii
	} else {
		profile = termenv.ColorProfile()
	}
}

// Style creates a styled string.
func Style(s string) termenv.Style {
	if noColor {
		return termenv.String(s)
	}

	return termenv.String(s)
}

// Dim returns dimmed text.
func Dim(s string) string {
	if noColor {
		return s
	}

	return Style(s).Faint().String()
}

// Bold returns bold text.
func Bold(s string) string {
	if noColor {
		return s
	}

	return Style(s).Bold().String()
}

// Green returns green text.
func Green(s string) string {
	if noColor {
		return s
	}

	return Style(s).Foreground(profile.Color("2")).String()
}

// Yellow returns yellow text.
func Yellow(s string) string {
	if noColor {
		return s
	}

	return Style(s).Foreground(profile.Color("3")).String()
}

// Red returns red text.
func Red(s string) string {
	if noColor {
		return s
	}

	return Style(s).Foreground(profile.Color("1")).String()
}

// Cyan returns cyan text.
func Cyan(s string) string {
	if noColor {
		return s
	}

	return Style(s).Foreground(profile.Color("6")).String()
}

// FormatDelay returns a colored delay string.
// Negative = early (green), 0 = on time (dim), positive = late (red/yellow).
func FormatDelay(seconds int) string {
	switch {
	case seconds < 0:
		return Green("-" + formatDuration(-seconds))
	case seconds == 0:
		return Dim("on time")
	case seconds < 120:
		return Yellow("+" + formatDuration(seconds))
	default:
		return Red("+" + formatDuration(seconds))
	}
}

func formatDuration(seconds int) string {
	if seconds < 60 {
		return "<1m"
	}

	m := seconds / 60

	return fmt.Sprintf("%dm", m)
}
