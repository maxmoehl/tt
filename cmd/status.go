package cmd

import (
	"errors"
	"fmt"
	"time"

	"moehl.dev/tt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints a short notice on the current status",
	Long:  `Reports if you are currently working, taking a break or taking some time off.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		short, err := getStatusParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("status: %w", err)
		}
		err = runStatus(short)
		if err != nil {
			return fmt.Errorf("status: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolP(flagShort, string(flagShort[0]), false, "print the status in short format")
}

func runStatus(short bool) error {
	orderBy := tt.OrderBy{
		Field: tt.FieldStart,
		Order: tt.OrderDsc,
	}
	var lastTimer tt.Timer
	err := tt.GetDB().GetTimer(tt.EmptyFilter, orderBy, &lastTimer)
	if err != nil && !errors.Is(err, tt.ErrNotFound) {
		return err
	}
	if errors.Is(err, tt.ErrNotFound) || !lastTimer.Running() {
		if short {
			fmt.Println("not tracking")
		} else {
			fmt.Println("Currently not tracking. Enjoy your free time :)")
		}
	} else {
		timingFor := tt.FormatDuration(time.Now().Sub(lastTimer.Start))
		if lastTimer.Task != "" {
			if short {
				fmt.Printf("%s / %s / %s\n", lastTimer.Project, lastTimer.Task, timingFor)
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
	return nil
}

func getStatusParameters(cmd *cobra.Command, _ []string) (short bool, err error) {
	flags, err := flags(cmd, flagShort)
	if err != nil {
		return
	}
	return flags[flagShort].(bool), nil
}
