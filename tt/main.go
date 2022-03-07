package main

import (
	"os"

	"github.com/maxmoehl/tt/cmd"
)

var version = "dev"

func main() {
	err := cmd.RootCmd(version).Execute()
	if err != nil {
		// the error will be printed by cobra
		os.Exit(1)
	}
}
