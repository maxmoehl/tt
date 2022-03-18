package tt

import (
	"fmt"
	"strings"
	"time"
)

const (
	filterProject = "project"
	filterTask    = "task"
	filterSince   = "since"
	filterUntil   = "until"
	filterTags    = "tags"

	filtersSeparator = ";"
	valuesSeparator  = ","

	sqlOperatorLike   = "LIKE"
	sqlOperatorEquals = "="

	DateFormat = "2006-01-02"
)

var EmptyFilter *filter

type Filter interface {
	DatabaseFilter
	Match(Timer) bool
	Timers(Timers) Timers
}

// filter contains all available filters. If a value is empty (i.e. "" or nil)
// it is assumed that the filter is not set and all values are included.
type filter struct {
	// project contains all project names that should be included. Accepts
	// multiple values. Project filter can be set with the keyword 'project'.
	project []string
	// task contains all task names that should be included. Accepts
	// multiple values. Task filter can be set with the keyword 'task'.
	task []string
	// since stores the date from which on the data should be included. Since
	// is inclusive and only accepts a single value in the following form:
	//   yyyy-MM-dd
	// Since filter can be set with the keyword 'since'.
	since time.Time
	// until stores the last date that should be included. Until is inclusive
	// and only accepts a single value in the following form:
	//	 yyyy-MM-dd
	// Until filter can be set with the keyword 'until'.
	until time.Time
	// tags contains all tags that should be included. They are parsed as a
	// comma separated list and can be set with the keyword 'tags'.
	tags []string
}

// Match checks if a given Timer matches this filter. An empty filter matches
// everything
func (f *filter) Match(t Timer) bool {
	if f == nil {
		return true
	}
	if f.project != nil && !stringSliceContains(f.project, t.Project) {
		return false
	}
	if f.task != nil && !stringSliceContains(f.task, t.Task) {
		return false
	}
	if !f.since.IsZero() && beforeDate(t.Start, f.since) {
		return false
	}
	if !f.until.IsZero() && beforeDate(f.until, t.Start) {
		return false
	}
	if f.tags != nil && !stringSliceContainsAny(f.tags, t.Tags) {
		return false
	}
	return true
}

func (f *filter) Timers(timers Timers) (filtered Timers) {
	if f == nil {
		return timers
	}
	for _, t := range timers {
		if f.Match(t) {
			filtered = append(filtered, t)
		}
	}
	return
}

// SQL returns the WHERE clause that can be used to match timers using this
// filter. A known limitation is that filtering for tags is not supported
// because of the way the tags are stored in the database.
func (f *filter) SQL() string {
	if f == nil {
		return ""
	}
	var filters []string
	if len(f.project) > 0 {
		// f.project = ["a", "b", "c"] => "json_extract(`json`, '$.project') IN ('a', 'b', 'c')"
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.project') IN ('%s')", strings.Join(f.project, "', '")))
	}
	if len(f.task) > 0 {
		// f.task = ["a", "b", "c"] => "json_extract(`json`, '$.task') IN ('a', 'b', 'c')"
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.task') IN ('%s')", strings.Join(f.task, "', '")))
	}
	if len(f.tags) > 0 {
		// f.tags = ["a", "b", "c"] => "`uuid` IN (SELECT `uuid` FROM `timers`, json_each(json_extract(timers.json, '$.tags')) WHERE value IN ('a', 'b', 'c'))"
		filters = append(filters, fmt.Sprintf("`uuid` IN (SELECT `uuid` FROM `timers`, json_each(json_extract(timers.json, '$.tags')) WHERE value IN ('%s'))", strings.Join(f.tags, "', '")))
	}
	if !f.since.IsZero() {
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.start') >= '%s'", f.since.Format(DateFormat)))
	}
	if !f.until.IsZero() {
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.start') < '%s'", f.until.AddDate(0, 0, 1).Format(DateFormat)))
	}
	// if there are no filters return TRUE to match all values
	if len(filters) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(filters, " AND ")
}

// ParseFilterString takes a string and creates a filter from it. The filter
// string has to be in the following format:
//
//   filterName=values;filterName=values;...
//
// each filterName consists of a string, values contains the filter value. Some
// filters only accept a single value, others accept multiple values separated
// by commas.
//
// Example:
//   project=work,school;since=2020-01-01;until=2020-02-01
//
// Available filters are:
//
//   project: accepts multiple string values
//   task   : accepts multiple string values
//   since  : accepts a single string int the format of yyyy-MM-dd
//   until  : accepts a single string int the format of yyyy-MM-dd
//   tags   : accepts multiple string values
//
// since and until are inclusive, both dates will be included in filtered data.
func ParseFilterString(filterString string) (Filter, error) {
	var f *filter
	if len(filterString) == 0 {
		return f, nil
	}
	f = new(filter)
	for _, fSlice := range strings.Split(filterString, filtersSeparator) {
		key, values, err := parseFilter(fSlice)
		if err != nil {
			return nil, err
		}
		err = parseValuesInto(key, values, f)
		if err != nil {
			return nil, err
		}
	}
	return f, nil
}

// NewFilter allows to create a filter that is not parsed from a filter string.
func NewFilter(projects, tasks, tags []string, since, until time.Time) Filter {
	return &filter{
		project: projects,
		task:    tasks,
		since:   time.Date(since.Year(), since.Month(), since.Day(), 0, 0, 0, 0, time.UTC),
		until:   time.Date(until.Year(), until.Month(), until.Day(), 0, 0, 0, 0, time.UTC),
		tags:    tags,
	}
}

func parseFilter(in string) (key, values string, err error) {
	filterSplit := strings.Split(in, "=")
	if len(filterSplit) != 2 {
		return "", "", ErrInvalidData
	}
	return filterSplit[0], filterSplit[1], nil
}

func parseValuesInto(key, values string, f *filter) (err error) {
	switch key {
	case filterProject:
		if f.project != nil {
			err = fmt.Errorf("redeclared filter project")
			break
		}
		f.project = strings.Split(values, valuesSeparator)
	case filterTask:
		if f.task != nil {
			err = fmt.Errorf("redeclared filter task")
			break
		}
		f.task = strings.Split(values, valuesSeparator)
	case filterSince:
		if !f.since.IsZero() {
			err = fmt.Errorf("redeclared filter since")
			break
		}
		f.since, err = ParseDate(values)
	case filterUntil:
		if !f.until.IsZero() {
			err = fmt.Errorf("redeclared filter until")
			break
		}
		f.until, err = ParseDate(values)
	case filterTags:
		if f.tags != nil {
			err = fmt.Errorf("redeclared filter tags")
			break
		}
		f.tags = strings.Split(values, valuesSeparator)
	}
	if err != nil {
		return fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
	}
	return
}

// stringSliceContainsAny checks if any of the valid values is inside the
// searchable list.
func stringSliceContainsAny(shouldContain []string, toSearch []string) bool {
	for _, c := range shouldContain {
		for _, t := range toSearch {
			if c == t {
				return true
			}
		}
	}
	return false
}

// stringSliceContains checks if the given string is contained in the given
// string slice.
func stringSliceContains(strings []string, s string) bool {
	for _, t := range strings {
		if t == s {
			return true
		}
	}
	return false
}

// beforeDate returns whether date one is before date two.
func beforeDate(one, two time.Time) bool {
	if one.Year() < two.Year() {
		return true
	} else if one.Year() > two.Year() {
		return false
	}
	if one.Month() < two.Month() {
		return true
	} else if one.Month() > two.Month() {
		return false
	}
	if one.Day() < two.Day() {
		return true
	} else if one.Day() > two.Day() {
		return false
	}
	return false
}
