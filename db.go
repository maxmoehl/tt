package tt

import (
	"fmt"
)

const (
	FieldStart   Field = "start"
	FieldProject Field = "project"
	FieldTask    Field = "task"
	FieldDay     Field = "day"

	OrderAsc Order = "ASC"
	OrderDsc Order = "DESC"
)

type DB interface {
	// SaveTimer will write a single timer to the database.
	SaveTimer(Timer) error
	// GetTimer reads a single timer from the database by appending `LIMIT 1` to
	// the query. Will return ErrNotFound if no timer matching the filter exists.
	GetTimer(Filter, OrderBy, *Timer) error
	GetTimerById(string, *Timer) error
	// GetTimers returns multiple timers from the database that match the filter.
	// Will never return ErrNotFound but rather just an empty list.
	GetTimers(Filter, OrderBy, *Timers) error
	// UpdateTimer updates existing timers in the db based on the id. If the timer
	// does not exist an ErrNotFound is returned.
	UpdateTimer(Timer) error
	// RemoveTimer removes the timer from the database.
	RemoveTimer(string) error

	SaveVacationDay(VacationDay) error
	GetVacationDay(VacationFilter, *VacationDay) error
	GetVacationDays(OrderBy, *[]VacationDay) error
	RemoveVacationDay(string) error
}

type Field string
type Order string

type OrderBy struct {
	Field Field
	Order Order
}

func (o OrderBy) SQL() string {
	if o.Field == "" {
		return ""
	}
	return fmt.Sprintf("ORDER BY json_extract(`json`, '$.%s') %s", o.Field, o.Order)
}

var db DB

func GetDB() DB {
	if db == nil {
		var err error
		db, err = NewSQLite(GetConfig().DBFile())
		if err != nil {
			panic(err.Error())
		}
	}
	return db
}
