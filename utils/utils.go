/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
)

const (
	// WarningNoArgumentsAccepted is a common warning for commands.
	WarningNoArgumentsAccepted = "this command does not accept any arguments"
	// DateFormat contains the format in which dates are printed.
	DateFormat = "2006-01-02"
)

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

// PrintWarning prints msg in yellow to stderr.
func PrintWarning(msg string) {
	yellow(os.Stderr, "Warning:", msg)
}

// PrintError takes an error and prints the value of error.Error() in red to
// stderr, and exits with os.Exit(1)
func PrintError(err error, silent bool) {
	if !silent {
		red(os.Stderr, "Error:", err.Error())
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
func FormatDuration(d time.Duration, precision time.Duration) string {
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
	switch precision {
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
