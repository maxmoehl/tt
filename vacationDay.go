package tt

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type VacationDay struct {
	ID   string    `json:"id"`
	Day  time.Time `json:"day"`
	Half bool      `json:"half"`
}

func (v VacationDay) String() string {
	return fmt.Sprintf("ID  : %s\nDay : %s\nHalf: %t", v.ID, v.Day.String(), v.Half)
}

// ParseDayString parses the given day string and expects the format YYYY-MM-DD.
// The returned time is always in timezone UTC to avoid daylight-saving-time
// issues when adding/subtracting days.
func ParseDayString(dayStr string) (time.Time, error) {
	dayParts := strings.Split(dayStr, "-")
	if len(dayParts) != 3 {
		return time.Time{}, ErrInvalidFormat
	}
	year, err := strconv.ParseInt(dayParts[0], 10, 0)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidFormat, err.Error())
	}
	month, err := strconv.ParseInt(dayParts[1], 10, 0)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidFormat, err.Error())
	}
	day, err := strconv.ParseInt(dayParts[2], 10, 0)
	if err != nil {
		return time.Time{}, fmt.Errorf("%w: %s", ErrInvalidFormat, err.Error())
	}
	return time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.UTC), nil
}

type VacationFilter time.Time

func (f VacationFilter) SQL() string {
	t := time.Time(f)
	return fmt.Sprintf("json_extract(`json`, '$.day') LIKE '%04d-%02d-%02d%%'", t.Year(), t.Month(), t.Day())
}
