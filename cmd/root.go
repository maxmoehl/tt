package cmd

import (
	"github.com/spf13/cobra"
)

// version to be set during build process:
//   go install \
//     -ldflags "-X main.version=vMAJOR.MINOR.PATCH-label" \
//     github.com/maxmoehl/tt/tt
var version = "dev"

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "tt is a cli application that can be used to track time.",
	Long: `tt is a cli application that can be used to track time.

See the help sections of the individual commands for more details on the
functionality.

Note: Timers cannot overlap. If there are overlapping timers
the application might fail and statistics or analytics may be wrong.`,
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolP(flagQuiet, string(flagQuiet[0]), false, "suppress all output")
}

func GetRootCmd() *cobra.Command {
	return rootCmd
}
