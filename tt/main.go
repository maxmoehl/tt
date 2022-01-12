package main

import (
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

const flagQuiet = "quiet"

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
the application might fail and statistics or analytics may be wrong.

Note: This cli is NOT concurrency safe.`,
	Version: version,
}

func init() {
	rootCmd.PersistentFlags().BoolP(flagQuiet, string(flagQuiet[0]), false, "Suppress all output")
}

func main() {
	// to make sure everything works as desired we have to obey a certain order:
	//   1. load the configuration (will be done implicitly on the first call to tt.GetConfig())
	//   2. load all storages from plugins
	//   3. initialize the storage
	//   4. load all commands and pass them the storage instance
	err := tt.LoadPlugins()
	if err != nil {
		tt.PrintError(err, false)
	}
	tt.LoadPluginStorages()
	err = tt.InitStorage()
	if err != nil {
		tt.PrintError(err, false)
	}
	tt.LoadPluginCmds(rootCmd)
	cobra.CheckErr(rootCmd.Execute())
}

func getQuiet(cmd *cobra.Command) bool {
	quiet, err := cmd.Flags().GetBool(flagQuiet)
	if err != nil {
		tt.PrintError(err, false)
	}
	return quiet
}
