package tt

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var datePattern = regexp.MustCompile(`((\d{4})[-/.](\d{1,2})[-/.](\d{1,2}))?([T ])?(\d{1,2}):(\d{1,2}):?(\d{1,2})?`)

// ParseDate will take in a string that contains a time in some
// format and try to guess the missing parts. Currently, the following
// cases are supported:
// - 15:04
// - 15:04:05
// - 2006/01/02 15:04
// - 2006/01/02 15:04:05
// valid date separators: dash, dot, slash
// valid time separators: colon
// valid separators between date and time: space, upper-case t
//
// more general information will be taken from time.Now() (e.g. day or
// year) and more specific information (e.g. seconds) will be set to zero.
func ParseDate(in string) (time.Time, error) {
	// if we can parse as RFC3339 we just return it
	t, err := time.Parse(time.RFC3339, in)
	if err == nil {
		return t, nil
	}

	matches := datePattern.FindSubmatch([]byte(in))
	if len(matches) == 0 {
		return time.Time{},
			ErrInvalidData.WithCause(NewError("timestamp is not RFC3339 compliant and does not match custom format"))
	}

	now := time.Now()
	var year, month, day, hour, min, sec int

	yearStr := string(matches[2])
	if yearStr != "" {
		year, err = strconv.Atoi(yearStr)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		year = now.Year()
	}

	monthStr := string(matches[3])
	if monthStr != "" {
		month, err = strconv.Atoi(monthStr)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		month = int(now.Month())
	}

	dayStr := string(matches[4])
	if dayStr != "" {
		day, err = strconv.Atoi(dayStr)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		day = now.Day()
	}

	hour, err = strconv.Atoi(string(matches[6]))
	if err != nil {
		return time.Time{}, err
	}

	min, err = strconv.Atoi(string(matches[7]))
	if err != nil {
		return time.Time{}, err
	}

	secStr := string(matches[8])
	if secStr != "" {
		sec, err = strconv.Atoi(secStr)
		if err != nil {
			return time.Time{}, err
		}
	} else {
		sec = 0
	}

	return time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local), nil
}

// FormatDuration formats a duration in the precision defined by the
// Config.
func FormatDuration(d time.Duration, precision time.Duration) string {
	h := d / time.Hour
	m := (d - (h * time.Hour)) / time.Minute
	s := (d - (h * time.Hour) - (m * time.Minute)) / time.Second
	sign := ""
	if d < 0 {
		sign = "-"
		h *= -1
		m *= -1
		s *= -1
	}
	switch precision {
	case time.Second:
		return fmt.Sprintf("%s%dh%dm%ds", sign, h, m, s)
	case time.Minute:
		return fmt.Sprintf("%s%dh%dm", sign, h, m)
	case time.Hour:
		return fmt.Sprintf("%s%dh", sign, h)
	default:
		return "unknown precision"
	}
}
