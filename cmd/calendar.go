package cmd

import (
	"fmt"

	"moehl.dev/tt"

	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:     "calendar",
	Aliases: []string{"cal"},
	Short:   "Show all data in a nice calendar format",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runCalendar()
		if err != nil {
			return fmt.Errorf("calendar: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(calendarCmd)
	// TODO: add flags for --abs and --rel that either show absolute values (current implementation)
	//       or the relative percentage indicating the fulfilment and something like `-%` for days
	//       where planned time == 0
}

func runCalendar() error {
	years, err := tt.BuildCalendar()
	if err != nil {
		return err
	}
	for _, year := range years {
		fmt.Printf(year.String())
	}
	return nil
}
