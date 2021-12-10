package main

import (
	"github.com/spf13/cobra"
)

const flagSilent = "silent"

// version to be set during build process:
//   go build \
//     -ldflags "-X main.version=vMAJOR.MINOR.PATCH-label" \
//     github.com/maxmoehl/tt/cli
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "tt is a cli application that can be used to track time.",
	Long: `tt is a cli application that can be used to track time.

See the help sections of the individual commands for more details on the
functionality.

Note: Timers cannot overlap. If there are overlapping timers
the application might fail and statistics or analytics may be wrong.

Note: This cli is NOT concurrency safe.`,
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolP(flagSilent, string(flagSilent[0]), false, "Suppress all output")
}

func main() {
	cobra.CheckErr(rootCmd.Execute())
}

func getSilent(cmd *cobra.Command) bool {
	silent, err := cmd.Flags().GetBool("silent")
	if err != nil {
		PrintError(err, false)
	}
	return silent
}
