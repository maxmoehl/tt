package cmd

import (
	"fmt"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var vacationListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all vacation days",
	Long:    `List all vacation days.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, err := getVacationListParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("vacation list: %w", err)
		}
		err = runVacationList(quiet)
		if err != nil {
			return fmt.Errorf("vacation list: %w", err)
		}
		return nil
	},
}

func init() {
	vacationCmd.AddCommand(vacationListCmd)
}

func runVacationList(quiet bool) error {
	if quiet {
		return nil
	}
	var vacationDays []tt.VacationDay
	order := tt.OrderBy{
		Field: tt.FieldDay,
		Order: tt.OrderAsc,
	}
	err := tt.GetDB().GetVacationDays(order, &vacationDays)
	if err != nil {
		return err
	}
	vacationCount := 0 // one day = 2, half a day = 1
	for _, day := range vacationDays {
		fmt.Println(day.String())
		fmt.Println("----------")
		if !day.Half {
			vacationCount += 2
		} else {
			vacationCount += 1
		}
	}
	fmt.Printf("Total Vacation Days: %.1f\n", float64(vacationCount)/2)
	return nil
}

func getVacationListParameters(cmd *cobra.Command, _ []string) (quiet bool, err error) {
	flags, err := flags(cmd, flagQuiet)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), nil
}
