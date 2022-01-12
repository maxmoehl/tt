package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

const (
	flagShort = "short"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints a short notice on the current status",
	Long: `Reports if you are currently working, taking a break or taking some
time off.`,
	Run: func(cmd *cobra.Command, args []string) {
		runStatus(getStatusParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolP(flagShort, string(flagShort[0]), false, "Print the status in short format")
}

func runStatus(quiet, short bool) {
	if quiet {
		return
	}
	lastTimer, err := tt.GetStorage().GetLastTimer(true)
	if err != nil && !errors.Is(err, tt.ErrNotFound) {
		tt.PrintError(err, quiet)
	}
	if errors.Is(err, tt.ErrNotFound) || !lastTimer.Running() {
		if short {
			fmt.Println("not tracking")
		} else {
			fmt.Println("Currently not tracking. Enjoy your free time :)")
		}
	} else {
		timingFor := tt.FormatDuration(time.Now().Sub(lastTimer.Start), tt.GetConfig().GetPrecision())
		if lastTimer.Task != "" {
			if short {
				fmt.Printf("tracking %s, %s for %s\n", lastTimer.Project, lastTimer.Task, timingFor)
			} else {
				fmt.Printf("Currently timing project %s with task %s for %s, you're doing good!\n", lastTimer.Project, lastTimer.Task, timingFor)
			}
		} else {
			if short {
				fmt.Printf("tracking %s for %s\n", lastTimer.Project, timingFor)
			} else {
				fmt.Printf("Currently timing project %s for %s, you're doing good!\n", lastTimer.Project, timingFor)
			}
		}
	}
}

func getStatusParameters(cmd *cobra.Command, args []string) (quiet, short bool) {
	var err error
	quiet = getQuiet(cmd)
	short, err = cmd.LocalFlags().GetBool(flagShort)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if len(args) != 0 && !quiet {
		tt.PrintWarning(tt.WarningNoArgumentsAccepted)
	}
	return
}
