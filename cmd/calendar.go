package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var calendarCmd = &cobra.Command{
	Use:     "calendar",
	Aliases: []string{"cal"},
	Short:   "Show all data in a nice calendar format",
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, err := getCalendarParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("calendar: %w", err)
		}
		err = runCalendar(quiet)
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

func runCalendar(quiet bool) error {
	if quiet {
		return nil
	}
	db := tt.GetDB()
	var timers tt.Timers
	err := db.GetTimers(tt.Filter{}, tt.OrderBy{}, &timers)
	if err != nil {
		return err
	}
	groupedTimers := timers.GroupBy(tt.GroupByDay)
	var start time.Time
	var stop time.Time
	for _, timers := range groupedTimers {
		for _, timer := range timers {
			if start.IsZero() || timer.Start.Unix() < start.Unix() {
				start = timer.Start
			}
			if stop.IsZero() || timer.Start.Unix() > stop.Unix() {
				stop = timer.Start
			}
		}
	}
	years := make(map[int]tt.Year)
	for y := start.Year(); y <= stop.Year(); y++ {
		var months [12]tt.Month
		for m := 0; !(y == stop.Year() && time.Month(m+1) > stop.Month()) && m < 12; m++ {
			if y == start.Year() && time.Month(m+1) < start.Month() {
				continue
			}
			var days []tt.Day
			for d := 0; isValidDate(y, m, d); d++ {
				key := fmt.Sprintf("%04d-%02d-%02d", y, m+1, d+1)
				t := time.Date(y, time.Month(m+1), d+1, 0, 0, 0, 0, time.UTC)
				var vac tt.VacationDay
				var pVac *tt.VacationDay
				err = tt.GetDB().GetVacationDay(t, &vac)
				if err == nil {
					pVac = &vac
				} else if err != nil && !errors.Is(err, tt.ErrNotFound) {
					return err
				}
				days = append(days, tt.Day{
					Time:     t,
					Timers:   groupedTimers[key],
					Vacation: pVac,
				})
			}
			months[m] = tt.Month{
				Days: days,
			}
		}
		years[y] = tt.Year{
			Year:   y,
			Months: months,
		}
	}
	for _, year := range years {
		fmt.Printf(year.String())
	}
	return nil
}

func getCalendarParameters(cmd *cobra.Command, _ []string) (quiet bool, err error) {
	flags, err := flags(cmd, flagQuiet)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), nil
}

func isValidDate(year, month, day int) bool {
	d := time.Date(year, time.Month(month+1), day+1, 0, 0, 0, 0, time.UTC)
	return d.Year() == year && month+1 == int(d.Month()) && day+1 == d.Day()
}
