package tt

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var monthNames = []string{
	"January",
	"February",
	"March",
	"April",
	"May",
	"June",
	"July",
	"August",
	"September",
	"October",
	"November",
	"December",
}

type Year struct {
	Year   int
	Months [12]Month
}

func (y Year) String() string {
	b := strings.Builder{}
	for m, month := range y.Months {
		if len(month.Days) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("%04d %s\n", y.Year, monthNames[m]))
		b.WriteString(month.String())
		b.WriteString("\n\n")
	}
	return b.String()
}

type Month struct {
	// Days contains the days of the month but starts at 0
	Days []Day
}

func (m Month) String() string {
	b := strings.Builder{}
	b.WriteString("Mon        Tue        Wed        Thu        Fri        Sat        Sun\n")

	// we need to insert space to account for months that do not start on mondays
	dayOfWeek := correctedWeekday(m.Days[0].Time.Weekday())
	b.WriteString(strings.Repeat("           ", dayOfWeek))
	for i, day := range m.Days {
		b.WriteString(day.String())
		if day.Time.Weekday() == time.Sunday && i < len(m.Days)-1 {
			// after a sunday we need a new line
			b.WriteRune('\n')
		} else {
			b.WriteString("  ")
		}
	}
	return b.String()
}

type Day struct {
	Time     time.Time
	Timers   Timers
	Vacation *VacationDay
}

func (d Day) String() string {
	// TODO: this could probably be simplified:
	//       simply compare planned to tracked time
	//         if planned == 0 leave empty, print vac, or plus time in green
	//         if planned > 0 show fulfilment using color coding (tbd)
	dur := d.Timers.Duration()
	if d.Vacation == nil && !IsWorkDay(d.Time) && dur == 0 {
		return fmt.Sprintf("%02d       ", d.Time.Day())
	} else if d.Vacation == nil && IsWorkDay(d.Time) {
		// TODO: color code depending on fulfillment
		return fmt.Sprintf("%02d %s", d.Time.Day(), FormatDuration(d.Timers.Duration(), time.Minute))
	} else if d.Vacation.Half {
		// TODO: color code depending on fulfilment and account for half a day vacation
		return fmt.Sprintf("%02d %s", d.Time.Day(), FormatDuration(d.Timers.Duration(), time.Minute))
	} else {
		return fmt.Sprintf("%02d vac.  ", d.Time.Day())
	}
}

// correctedWeekday adjusts from Sunday as first day of the week to monday as
// first day of the week.
func correctedWeekday(weekday time.Weekday) int {
	if weekday == time.Sunday {
		return int(time.Saturday)
	} else {
		return int(weekday) - 1
	}
}

// IsWorkDay checks if the day should have been worked on. It does not take into
// account vacation days, but only the configured weekdays.
func IsWorkDay(d time.Time) bool {
	days := GetConfig().Timeclock.DaysPerWeek
	switch d.Weekday() {
	case time.Monday:
		return days.Monday
	case time.Tuesday:
		return days.Tuesday
	case time.Wednesday:
		return days.Wednesday
	case time.Thursday:
		return days.Thursday
	case time.Friday:
		return days.Friday
	case time.Saturday:
		return days.Saturday
	case time.Sunday:
		return days.Sunday
	default:
		panic(fmt.Sprintf("unknown day of week %d", d.Weekday()))
	}
}

// PlannedTime returns the duration that was planned for the given date.
func PlannedTime(date time.Time) (time.Duration, error) {
	if !IsWorkDay(date) {
		return 0, nil
	}
	workTime := time.Duration(GetConfig().Timeclock.HoursPerDay) * time.Hour
	var vac VacationDay
	err := GetDB().GetVacationDay(date, &vac)
	if errors.Is(err, ErrNotFound) {
		return workTime, nil
	} else if err != nil {
		return 0, err
	}
	if vac.Half {
		return workTime / 2, nil
	}
	return 0, nil
}
