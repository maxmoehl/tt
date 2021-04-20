package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt/utils"
)

const (
	filterProject = "project"
	filterTask    = "task"
	filterSince   = "since"
	filterUntil   = "until"

	filtersSeparator = ";"
	valuesSeparator  = ","

	dateFormat = "2006-01-02"
)

type Filter interface {
	Match(Timer) bool
	Filter(Timers) Timers
}

// filter contains all available filters. If a value is empty (i.e. ""
// or nil) it is assumed that the filter is not set and all values are
// included.
type filter struct {
	// Project contains all project names that should be included. Accepts
	// multiple values. Project filter can be set with the keyword 'project'.
	Project []string
	// Task contains all task names that should be included. Accepts
	// multiple values. Task filter can be set with the keyword 'task'.
	Task []string
	// Since stores the date from which on the data should be included. Since
	// is inclusive and only accepts a single value in the following form:
	//   yyyy-MM-dd
	// Since filter can be set with the keyword 'since'.
	Since time.Time
	// Until stores the last date that should be included. Until is inclusive
	// and only accepts a single value in the following form:
	//	 yyyy-MM-dd
	// Until filter can be set with the keyword 'until'.
	Until time.Time
}

func (f filter) Match(t Timer) bool {
	if f.Project != nil && !utils.StringSliceContains(f.Project, t.Project) {
		return false
	}
	if f.Task != nil && !utils.StringSliceContains(f.Task, t.Task) {
		return false
	}
	if !f.Since.IsZero() && t.Start.Before(f.Since) {
		return false
	}
	if !f.Until.IsZero() && (t.Start.After(f.Until) || t.End.After(f.Until)) {
		return false
	}
	return true
}

func (f filter) Filter(timers Timers) (filtered Timers) {
	for _, t := range timers {
		if f.Match(t) {
			filtered = append(filtered, t)
		}
	}
	return
}

// GetFilter takes a string and creates a filter struct from it. The
// filter string has to be in the following format:
//
//   filterName=values;filterName=values;...
//
// each filterName consists only of strings. values contains the filter
// value. Some filters only accept a single value, others accept multiple
// values separated by commas.
//
// Example:
//   projectName=work,school;since=2020-01-01;until=2020-02-01
//
// Available filters are:
//
//   project: accepts multiple string values
//   task   : accepts multiple string values
//   since  : accepts a single string int the format of yyyy-MM-dd
//   until  : accepts a single string int the format of yyyy-MM-dd
//
// since and until are inclusive, both dates will be included in filtered
// data.
func GetFilter(filterString string) (Filter, error) {
	var F filter
	for _, f := range strings.Split(filterString, filtersSeparator) {
		key, values, err := parseFilter(f)
		if err != nil {
			return nil, err
		}
		err = parseValuesInto(key, values, &F)
		if err != nil {
			return nil, err
		}
	}
	// To make until inclusive we have to add 24h to it since filter.Match
	// checks if the end of a timer is before until
	if !F.Until.IsZero() {
		F.Until = F.Until.Add(time.Hour * 24)
	}
	return F, nil
}

func parseFilter(in string) (key, values string, err error) {
	filterSplit := strings.Split(in, "=")
	if len(filterSplit) != 2 {
		return "", "", fmt.Errorf("expected one '=' per filter but got %d: [%s]", len(filterSplit), in)
	}
	return filterSplit[0], filterSplit[1], nil
}

func parseValuesInto(key, values string, f *filter) (err error) {
	switch key {
	case filterProject:
		if f.Project != nil {
			err = fmt.Errorf("redeclared filter project")
		}
		f.Project = strings.Split(values, valuesSeparator)
	case filterTask:
		if f.Task != nil {
			err = fmt.Errorf("redeclared filter task")
		}
		f.Task = strings.Split(values, valuesSeparator)
	case filterSince:
		if !f.Since.IsZero() {
			err = fmt.Errorf("redeclared filter since")
		}
		f.Since, err = time.Parse(dateFormat, values)
	case filterUntil:
		if !f.Until.IsZero() {
			err = fmt.Errorf("redeclared filter until")
		}
		f.Until, err = time.Parse(dateFormat, values)
	}
	return
}
