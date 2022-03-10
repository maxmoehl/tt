package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"end"},
	Short:   "Stops the current timer",
	Long: `Stops the current timer.

If you want to manually set a stop time it should look something like this:
  2020-04-19 08:00
  2020-04-19T08:00

you can also omit the date, the current date will be used:
  08:00

or add seconds if that's your thing:
  09:32:42

You can also supply a full RFC3339 date-time string.

Otherwise an appropriate error will be printed.`,
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
		printTrackingStoppedMsg(timer)
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

func printTrackingStoppedMsg(t tt.Timer) {
	fmt.Printf("[%02d:%02d] Tracking stopped!\n", t.Stop.Hour(), t.Stop.Minute())
	fmt.Printf("  start   : %02d:%02d\n", t.Start.Hour(), t.Start.Minute())
	fmt.Printf("  project : %s\n", t.Project)
	if t.Task != "" {
		fmt.Printf("  task    : %s\n", t.Task)
	}
	if len(t.Tags) > 0 {
		fmt.Printf("  tags    : %s\n", strings.Join(t.Tags, ","))
	}
	fmt.Printf("  duration: %s\n", tt.FormatDuration(t.Duration(), tt.GetConfig().GetPrecision()))
}
