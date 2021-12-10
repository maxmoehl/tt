package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints a short notice on the current status",
	Long: `Reports if you are currently working, taking a break or taking some
time off.`,
	Run: func(cmd *cobra.Command, args []string) {
		status(getStatusParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func status(silent bool) {
	if silent {
		return
	}
	runningTimers, err := tt.CheckRunningTimers()
	if err != nil && !errors.Is(err, tt.ErrNotFound) {
		PrintError(err, silent)
	}
	if len(runningTimers) == 0 {
		fmt.Println("Currently not working. Enjoy your free time :)")
		return
	}
	timer, err := tt.GetRunningTimer()
	if err != nil {
		PrintError(err, silent)
	}
	timingFor := time.Now().Sub(timer.Start).Round(time.Second).String()
	if timer.Task != "" {
		fmt.Printf("Currently timing project %s with task %s for %s, your doing good!\n", timer.Project, timer.Task, timingFor)
	} else {
		fmt.Printf("Currently timing project %s for %s, your doing good!\n", timer.Project, timingFor)
	}
}

func getStatusParameters(cmd *cobra.Command, args []string) (silent bool) {
	silent = getSilent(cmd)
	if len(args) != 0 && !silent {
		PrintWarning(WarningNoArgumentsAccepted)
	}
	return
}
