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
	tracked := d.Timers.Duration()
	if d.Vacation == nil && !IsWorkDay(d.Time) && tracked == 0 {
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
	err := GetDB().GetVacationDay(VacationFilter(date), &vac)
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

func BuildCalendar() ([]Year, error) {
	db := GetDB()
	var timers Timers
	err := db.GetTimers(Filter{}, OrderBy{}, &timers)
	if err != nil {
		return nil, err
	}
	groupedTimers := timers.GroupBy(GroupByDay)
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
	years := make(map[int]Year)
	for y := start.Year(); y <= stop.Year(); y++ {
		var months [12]Month
		for m := 0; !(y == stop.Year() && time.Month(m+1) > stop.Month()) && m < 12; m++ {
			if y == start.Year() && time.Month(m+1) < start.Month() {
				continue
			}
			var days []Day
			for d := 0; isValidDate(y, m, d); d++ {
				key := fmt.Sprintf("%04d-%02d-%02d", y, m+1, d+1)
				t := time.Date(y, time.Month(m+1), d+1, 0, 0, 0, 0, time.UTC)
				var vac VacationDay
				var pVac *VacationDay
				err = db.GetVacationDay(VacationFilter(t), &vac)
				if err == nil {
					pVac = &vac
				} else if err != nil && !errors.Is(err, ErrNotFound) {
					return nil, err
				}
				days = append(days, Day{
					Time:     t,
					Timers:   groupedTimers[key],
					Vacation: pVac,
				})
			}
			months[m] = Month{
				Days: days,
			}
		}
		years[y] = Year{
			Year:   y,
			Months: months,
		}
	}
	return nil, nil
}

func isValidDate(year, month, day int) bool {
	d := time.Date(year, time.Month(month+1), day+1, 0, 0, 0, 0, time.UTC)
	return d.Year() == year && month+1 == int(d.Month()) && day+1 == d.Day()
}
