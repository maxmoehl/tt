package main

import (
	"os"

	"github.com/maxmoehl/tt/cmd"
)

func main() {
	err := cmd.GetRootCmd().Execute()
	if err != nil {
		// the error will be printed by cobra
		os.Exit(1)
	}
}
