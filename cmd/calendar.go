package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

/*
What should the output look like?

2022 January
Mon       Tue       Wed       Thu       Fri       Sat       Sun
01 08:00  02 08:00  03 Vac    04 08:00  05 08:00  06 00:00  07 00:00
08 08:00  09 08:00  10 08:00  11 08:00  12 08:00  13 00:00  14 00:00
15 08:00  16 08:00  17 08:00  18 08:00  19 08:00  20 00:00  21 00:00
22 08:00  23 08:00  24 08:00  25 08:00  26 08:00  27 00:00  28 00:00
29 08:00  30 08:00  31 08:00

2022 February
*/

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
