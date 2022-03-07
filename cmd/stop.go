package cmd

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"end"},
	Short:   "Stops a timer",
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
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, timestamp, err := getStopParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("stop: %w", err)
		}
		err = runStop(quiet, timestamp)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "manually set the stop time for a timer")
}

func runStop(quiet bool, timestamp time.Time) error {
	timer, err := tt.Stop(timestamp)
	if err != nil {
		return err
	}
	if !quiet {
		fmt.Printf("You worked for %s. Good job!\n", timer.Duration().Round(time.Second).String())
	}
	return nil
}

func getStopParameters(cmd *cobra.Command, _ []string) (quiet bool, timestamp time.Time, err error) {
	flags, err := flags(cmd, flagQuiet, flagTimestamp)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), flags[flagTimestamp].(time.Time), nil
}
