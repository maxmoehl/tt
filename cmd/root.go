package cmd

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "tt is a cli application that can be used to track time.",
	Long: `tt is a cli application that can be used to track time.

See the help sections of the individual commands for more details on the
functionality.

Note: Timers cannot overlap. If there are overlapping timers the application
might fail and statistics or analytics may be wrong.`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		flags, err := flags(cmd, flagQuiet, flagNoColor)
		if err != nil {
			return err
		}
		if flags[flagQuiet].(bool) {
			os.Stdout, err = os.Open(os.DevNull)
			if err != nil {
				return fmt.Errorf("unable to open /dev/null: %w", err)
			}
		}
		if flags[flagNoColor].(bool) {
			color.NoColor = true
		}
		return nil
	},
	SilenceUsage: true, // do not print usage on error
}

func init() {
	rootCmd.Version = version()
	rootCmd.PersistentFlags().BoolP(flagQuiet, short(flagQuiet), false, "suppress all output to stdout")
	rootCmd.PersistentFlags().BoolP(flagNoColor, short(flagNoColor), false, "disable colored output")
}

func RootCmd() *cobra.Command {
	return rootCmd
}

func version() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		panic("tt must be built with go module support")
	}

	var commit, dirty string
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			commit = setting.Value
		case "vcs.modified":
			dirty = setting.Value
		}
	}

	if dirty == "true" {
		dirty = "dirty"
	} else {
		dirty = ""
	}

	return fmt.Sprintf("%s+%s.%s built using %s", info.Main.Version, commit[:7], dirty, info.GoVersion)
}
