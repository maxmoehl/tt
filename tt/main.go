package main

import (
	"os"

	"moehl.dev/tt/cmd"
)

func main() {
	err := cmd.RootCmd().Execute()
	if err != nil {
		// the error will be printed by cobra
		os.Exit(1)
	}
}
