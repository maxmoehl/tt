package cmd

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the last timer",
	Long: `Resume the last timer.

If no timer is found an error is returned.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, timestamp, err := getResumeParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("resume: %w", err)
		}
		err = runResume(quiet, timestamp)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "manually set the start time for a timer")
}

func runResume(quiet bool, timestamp time.Time) error {
	timer, err := tt.Resume(timestamp)
	if err != nil {
		return err
	}
	if !quiet {
		printTrackingStartedMsg(timer)
	}
	return nil
}

func getResumeParameters(cmd *cobra.Command, _ []string) (quiet bool, timestamp time.Time, err error) {
	flags, err := flags(cmd, flagQuiet, flagTimestamp)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), flags[flagTimestamp].(time.Time), nil
}
