package cmd

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var vacationAddCmd = &cobra.Command{
	Use:   "add <day>",
	Short: "Add a vacation day",
	Long: `Add a vacation day.

<day> should be in the format of: YYYY-MM-DD.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		half, day, err := getVacationAddParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("vacation add: %w", err)
		}
		err = runVacationAdd(half, day)
		if err != nil {
			return fmt.Errorf("vacation add: %w", err)
		}
		return nil
	},
}

func init() {
	vacationCmd.AddCommand(vacationAddCmd)
	vacationAddCmd.Flags().Bool(flagHalf, false, "only add half a day instead of a full day.")
}

func runVacationAdd(half bool, day time.Time) error {
	v := tt.VacationDay{
		Day:  day,
		Half: half,
	}
	return tt.GetDB().SaveVacationDay(v)
}

func getVacationAddParameters(cmd *cobra.Command, args []string) (half bool, day time.Time, err error) {
	flags, err := flags(cmd, flagHalf)
	if err != nil {
		return
	}
	if len(args) != 1 {
		err = fmt.Errorf("expected one argument")
		return
	}
	day, err = tt.ParseDayString(args[0])
	return flags[flagHalf].(bool), day, nil
}
