package tt

import (
	"os"

	"github.com/fatih/color"
)

const (
	WarningNoArgumentsAccepted    = "this command does not accept any arguments"
	ErrorNArgumentsAcceptedFormat = "this command accepts %s arguments, but got %d"
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
func PrintError(err error, quiet bool) {
	if !quiet {
		printError(err, "")
	}
	os.Exit(1)
}

func printError(e error, indent string) {
	err, ok := e.(Error)
	if ok {
		red(os.Stderr, indent, "Error:", err.Msg())
		inner := err.Unwrap()
		if inner != nil {
			red(os.Stderr, indent, "Caused by:")
			printError(inner, indent+"  ")
		}
	} else {
		red(os.Stderr, indent, "Error:", e.Error())
	}
}
