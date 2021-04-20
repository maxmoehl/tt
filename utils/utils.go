package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/maxmoehl/tt/config"
)

const (
	WarningNoArgumentsAccepted = "this command does not accept any arguments"

	DateFormat = "2006-01-02"
)

// PrintWarning prints msg in yellow.
func PrintWarning(msg string) {
	fmt.Println("\x1b[33mWarning: " + msg + "\x1b[0m")
}

// PrintError takes an error and prints the value of error.Error() in red to the
// console, and exits with os.Exit(1)
func PrintError(err error, silent bool) {
	if !silent {
		fmt.Println("\x1b[31mError: " + err.Error() + "\x1b[0m")
	}
	os.Exit(1)
}

// StringSliceContains checks if the given string is contained in the
// given string slice.
func StringSliceContains(strings []string, s string) bool {
	for _, t := range strings {
		if t == s {
			return true
		}
	}
	return false
}

// FormatDuration formats a duration in the precision defined by the
// config.
func FormatDuration(d time.Duration) string {
	h := d / time.Hour
	m := (d - (h * time.Hour)) / time.Minute
	s := (d - (h * time.Hour) - (m * time.Minute)) / time.Second
	sign := ""
	if d < 0 {
		sign = "-"
		h *= -1
		m *= -1
		s *= -1
	}
	switch config.Get().GetPrecision() {
	case time.Second:
		return fmt.Sprintf("%s%dh%dm%ds", sign, h, m, s)
	case time.Minute:
		return fmt.Sprintf("%s%dh%dm", sign, h, m)
	case time.Hour:
		return fmt.Sprintf("%s%dh", sign, h)
	default:
		return "unknown precision"
	}
}
