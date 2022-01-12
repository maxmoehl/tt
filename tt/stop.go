package main

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"end"},
	Short:   "Stops a timer.",
	Long: `This command stops the current timer. If the timer is currently in a break
the break is also ended without further notice.

If you want to manually set a stop time it should look something like
this:
  2020-04-19 08:00
you can also omit the date, the current date will be used:
  08:00
or add seconds if that's your thing:
  09:32:42
You can also supply a full RFC3339 date-time string.

Otherwise an appropriate error will be printed. The cli will check if the
given stop time is valid, e.g. if the last timer and break that were
started, started before the given stop.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStop(getStopParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "Manually set the stop time for a timer")
}

func runStop(quiet bool, timestamp time.Time) {
	timer, err := tt.GetStorage().GetLastTimer(true)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	timer.Stop = timestamp
	err = tt.GetStorage().UpdateTimer(timer)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if !quiet {
		fmt.Printf("You worked for %s! Good job.\n", timer.Duration().Round(time.Second).String())
	}
}

func getStopParameters(cmd *cobra.Command, args []string) (quiet bool, timestamp time.Time) {
	var err error
	quiet = getQuiet(cmd)
	rawTimestamp, err := cmd.Flags().GetString(flagTimestamp)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if rawTimestamp == "" {
		timestamp = time.Now()
	} else {
		timestamp, err = tt.ParseDate(rawTimestamp)
		if err != nil {
			tt.PrintError(err, quiet)
		}
	}
	if len(args) != 0 && !quiet {
		tt.PrintWarning(tt.WarningNoArgumentsAccepted)
	}
	return
}
