package main

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the last timer.",
	Long: `Resume the last timer.

If no timer is found an error is returned.`,
	Run: func(cmd *cobra.Command, args []string) {
		runResume(getResumeParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func runResume(quiet bool) {
	lastTimer, err := tt.GetStorage().GetLastTimer(true)
	if err != nil && !errors.Is(err, tt.ErrNotFound) {
		tt.PrintError(err, quiet)
	}
	if errors.Is(err, tt.ErrNotFound) {
		tt.PrintError(tt.NewError("no timer found to resume"), quiet)
	}
	newTimer := tt.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: lastTimer.Project,
		Task:    lastTimer.Task,
		Tags:    lastTimer.Tags,
	}
	err = tt.GetStorage().StoreTimer(newTimer)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if !quiet {
		printTrackingStartedMsg(newTimer)
	}
}

func getResumeParameters(cmd *cobra.Command, args []string) (quiet bool) {
	quiet = getQuiet(cmd)
	if len(args) != 0 && !quiet {
		tt.PrintWarning(tt.WarningNoArgumentsAccepted)
	}
	return
}
