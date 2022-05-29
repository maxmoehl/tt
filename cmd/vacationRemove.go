package cmd

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var vacationRemoveCmd = &cobra.Command{
	Use:     "remove <day>",
	Aliases: []string{"rm"},
	Short:   "Remove a vacation day",
	Long: `Remove a vacation day.

<day> should be in the format of: YYYY-MM-DD.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		day, err := getVacationRemoveParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("vacation remove: %w", err)
		}
		err = runVacationRemove(day)
		if err != nil {
			return fmt.Errorf("vacation remove: %w", err)
		}
		return nil
	},
}

func init() {
	vacationCmd.AddCommand(vacationRemoveCmd)
}

func runVacationRemove(day time.Time) error {
	var vac tt.VacationDay
	err := tt.GetDB().GetVacationDay(tt.VacationFilter(day), &vac)
	if err != nil {
		return err
	}
	return tt.GetDB().RemoveVacationDay(vac.ID)
}

func getVacationRemoveParameters(_ *cobra.Command, args []string) (day time.Time, err error) {
	if len(args) != 1 {
		err = fmt.Errorf("expected one argument")
		return
	}
	day, err = tt.ParseDayString(args[0])
	return day, nil
}
