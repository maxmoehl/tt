package cmd

import (
	"fmt"
	"sort"
	"time"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

var timeclockCmd = &cobra.Command{
	Use:     "timeclock",
	Aliases: []string{"tc"},
	Short:   "Track your time and compare it to planned time",
	Long: `Track your time and compare it to planned time.

This command prints planned vs. worked time.

See subcommands for more details.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, day, filter, err := getTimeclockParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("timeclock: %w", err)
		}
		err = runTimeclock(quiet, day, filter)
		if err != nil {
			return fmt.Errorf("timeclock: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeclockCmd)
	timeclockCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "filter timers before showing statistics")
	timeclockCmd.Flags().BoolP(flagDay, string(flagDay[0]), false, "show time per day")
}

func runTimeclock(quiet, day bool, filter tt.Filter) error {
	if quiet {
		return nil
	}
	orderBy := tt.OrderBy{
		Field: tt.FieldStart,
		Order: tt.OrderAsc,
	}
	var timers tt.Timers
	err := tt.GetDB().GetTimers(filter, orderBy, &timers)
	if err != nil {
		return err
	}
	if day {
		e := statsByDay(timers)
		if err != nil {
			return e
		}
		fmt.Println("\nOverall statistics:")
	}
	return overallStats(timers)
}

func getTimeclockParameters(cmd *cobra.Command, _ []string) (quiet, day bool, filter tt.Filter, err error) {
	flags, err := flags(cmd, flagQuiet, flagDay, flagFilter)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), flags[flagDay].(bool), flags[flagFilter].(tt.Filter), nil
}

func firstAndLast(timers tt.Timers) (first, last time.Time, err error) {
	groupedTimers := timers.GroupBy(tt.GroupByDay)
	var days []string
	for day := range groupedTimers {
		days = append(days, day)
	}
	sort.Strings(days)
	first, err = tt.ParseDayString(days[0])
	if err != nil {
		return
	}
	last, err = tt.ParseDayString(days[len(days)-1])
	if err != nil {
		return
	}
	return
}

func statsByDay(timers tt.Timers) error {
	from, to, err := firstAndLast(timers)
	if err != nil {
		return err
	}
	to = to.AddDate(0, 0, 1)
	precision := tt.GetConfig().GetPrecision()
	for ; !datesEqual(from, to); from = from.AddDate(0, 0, 1) {
		dayTimers := tt.NewFilter(nil, nil, nil, from, from).Timers(timers)
		worked := dayTimers.Duration()
		planned, err := plannedTime(from, from)
		if err != nil {
			return err
		}
		if worked == 0 && planned == 0 {
			continue
		}
		fmt.Printf("%s: %s / %s\n", dayString(from), tt.FormatDuration(worked, precision), tt.FormatDuration(planned, precision))
	}
	return nil
}

func overallStats(timers tt.Timers) error {
	worked := timers.Duration()
	from, to, err := firstAndLast(timers)
	if err != nil {
		return err
	}
	planned, err := plannedTime(from, to)
	if err != nil {
		return err
	}
	precision := tt.GetConfig().GetPrecision()
	f := tt.FormatDuration
	fmt.Printf("worked    : %s\n", f(worked, precision))
	fmt.Printf("planned   : %s\n", f(planned, precision))
	fmt.Printf("difference: %s\n", f(worked-planned, precision))
	fmt.Printf("percentage: %.2f%%\n", float64(worked)/float64(planned)*100)
	return nil
}

func datesEqual(one time.Time, two time.Time) bool {
	return one.Year() == two.Year() && one.Month() == two.Month() && one.Day() == two.Day()
}

func plannedTime(from time.Time, to time.Time) (time.Duration, error) {
	if from.IsZero() || to.IsZero() {
		return 0, fmt.Errorf("from and to must be non-zero times")
	}
	to = to.AddDate(0, 0, 1)
	var d time.Duration
	for ; !datesEqual(from, to); from = from.AddDate(0, 0, 1) {
		t, err := tt.PlannedTime(from)
		if err != nil {
			return 0, err
		}
		d += t
	}
	return d, nil
}
