package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// version to be set during build process:
//   go install \
//     -ldflags "-X main.version=vMAJOR.MINOR.PATCH-label" \
//     -tags "json1"
//     github.com/maxmoehl/tt/tt

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
	rootCmd.PersistentFlags().BoolP(flagQuiet, short(flagQuiet), false, "suppress all output to stdout")
	rootCmd.PersistentFlags().BoolP(flagNoColor, short(flagNoColor), false, "disable colored output")
}

func RootCmd(version string) *cobra.Command {
	rootCmd.Version = version
	return rootCmd
}
