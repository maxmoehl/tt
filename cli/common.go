package main

import (
	"os"

	"github.com/maxmoehl/tt/config"

	"github.com/fatih/color"
)

const (
	// WarningNoArgumentsAccepted is a common warning for commands.
	WarningNoArgumentsAccepted = "this command does not accept any arguments"
)

var (
	red    = color.New(color.FgRed).FprintlnFunc()
	yellow = color.New(color.FgYellow).FprintlnFunc()
)

func init() {
	err := config.Load()
	if err != nil {
		PrintError(err, false)
	}
}

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
