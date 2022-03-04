package cmd

import (
	"errors"
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
			return fmt.Errorf("resume: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
	resumeCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "manually set the start time for a timer")
}

func runResume(quiet bool, timestamp time.Time) error {
	var db = tt.GetDB()
	orderBy := tt.OrderBy{
		Field: tt.FieldStart,
		Order: tt.OrderDsc,
	}
	var lastTimer tt.Timer
	err := db.GetTimer(tt.Filter{}, orderBy, &lastTimer)
	if err != nil && !errors.Is(err, tt.ErrNotFound) {
		return err
	}
	if errors.Is(err, tt.ErrNotFound) {
		return fmt.Errorf("no timer found")
	}
	newTimer := tt.Timer{
		Start:   timestamp,
		Project: lastTimer.Project,
		Task:    lastTimer.Task,
		Tags:    lastTimer.Tags,
	}
	err = db.SaveTimer(newTimer)
	if err != nil {
		return err
	}
	if !quiet {
		printTrackingStartedMsg(newTimer)
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
