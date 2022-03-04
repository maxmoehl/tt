package cmd

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var vacationCmd = &cobra.Command{
	Use:     "vacation",
	Aliases: []string{"vac"},
	Short:   "Modify vacation days for timeclock",
	Long: `Modify vacation days for timeclock.

Allows you to specify days which will be considered vacation and therefore
not be included in the statistics.`,
}

func init() {
	rootCmd.AddCommand(vacationCmd)
}

func dayString(day time.Time) string {
	return fmt.Sprintf("%04d-%02d-%02d", day.Year(), day.Month(), day.Day())
}
