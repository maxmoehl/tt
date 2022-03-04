package tt

import (
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
	if d.Vacation == nil {
		return fmt.Sprintf("%02d %s", d.Time.Day(), FormatDuration(d.Timers.Duration(), time.Minute))
	} else if d.Vacation.Half {
		// TODO: account for half a day vacation
		return fmt.Sprintf("%02d %s", d.Time.Day(), FormatDuration(d.Timers.Duration(), time.Minute))
	} else {
		return "vacation "
	}
}

// correctedWeekday adjusts from Sunday as first day of the week
// to monday as first day of the week.
func correctedWeekday(weekday time.Weekday) int {
	if weekday == time.Sunday {
		return int(time.Saturday)
	} else {
		return int(weekday) - 1
	}
}
