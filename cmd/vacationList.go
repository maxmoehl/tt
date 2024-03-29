package cmd

import (
	"fmt"

	"moehl.dev/tt"

	"github.com/spf13/cobra"
)

var vacationListCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all vacation days",
	Long:    `List all vacation days.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		err := runVacationList()
		if err != nil {
			return fmt.Errorf("vacation list: %w", err)
		}
		return nil
	},
}

func init() {
	vacationCmd.AddCommand(vacationListCmd)
}

func runVacationList() error {
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
