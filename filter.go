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

// Filter contains all available filters. If a value is empty (i.e. "" or nil)
// it is assumed that the Filter is not set and all values are included.
type Filter struct {
	// project contains all project names that should be included. Accepts
	// multiple values. Project Filter can be set with the keyword 'project'.
	project []string
	// task contains all task names that should be included. Accepts
	// multiple values. Task Filter can be set with the keyword 'task'.
	task []string
	// since stores the date from which on the data should be included. Since
	// is inclusive and only accepts a single value in the following form:
	//   yyyy-MM-dd
	// Since filter can be set with the keyword 'since'.
	since time.Time
	// until stores the last date that should be included. Until is inclusive
	// and only accepts a single value in the following form:
	//	 yyyy-MM-dd
	// Until Filter can be set with the keyword 'until'.
	until time.Time
	// tags contains all tags that should be included. They are parsed as a
	// comma separated list and can be set with the keyword 'tags'.
	tags []string
}

// Match checks if a given Timer matches this Filter. An empty filter matches
// everything
func (f Filter) Match(t Timer) bool {
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

func (f Filter) Timers(timers Timers) (filtered Timers) {
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
func (f Filter) SQL() string {
	var filters []string
	projects := convertFilterToSql("project", f.project, sqlOperatorEquals)
	if projects != "" {
		filters = append(filters, projects)
	}
	tasks := convertFilterToSql("task", f.task, sqlOperatorEquals)
	if tasks != "" {
		filters = append(filters, tasks)
	}
	if !f.since.IsZero() {
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.start') >= '%s'", f.since.Format(DateFormat)))
	}
	if !f.until.IsZero() {
		filters = append(filters, fmt.Sprintf("json_extract(`json`, '$.start') < '%s'", f.until.Format(DateFormat)))
	}
	tags := convertFilterToSql("tags", f.tags, sqlOperatorLike)
	if tags != "" {
		filters = append(filters, tags)
	}
	// if there are no filters return TRUE to match all values
	if len(filters) == 0 {
		return ""
	}
	return "WHERE " + strings.Join(filters, " AND ")
}

// ParseFilterString takes a string and creates a Filter from it. The Filter
// string has to be in the following format:
//
//   filterName=values;filterName=values;...
//
// each filterName consists of a string, values contains the Filter value. Some
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
// TODO: increase usability by allowing more relaxed values (e.g. 'today', 'yesterday' for since/until)
func ParseFilterString(filterString string) (Filter, error) {
	var f Filter
	if len(filterString) == 0 {
		return f, nil
	}
	for _, fSlice := range strings.Split(filterString, filtersSeparator) {
		key, values, err := parseFilter(fSlice)
		if err != nil {
			return f, err
		}
		err = parseValuesInto(key, values, &f)
		if err != nil {
			return f, err
		}
	}
	// To make until inclusive we have to add 24h to it since Filter.Match
	// checks if the end of a timer is before until
	// TODO: this might not be necessary anymore
	if !f.until.IsZero() {
		f.until = f.until.Add(time.Hour * 24)
	}
	return f, nil
}

// NewFilter allows to create a Filter that is not parsed from a Filter string.
func NewFilter(projects, tasks, tags []string, since, until time.Time) Filter {
	return Filter{
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

func parseValuesInto(key, values string, f *Filter) (err error) {
	switch key {
	case filterProject:
		if f.project != nil {
			err = fmt.Errorf("%w: redeclared filter project", ErrInvalidData)
			return
		}
		f.project = strings.Split(values, valuesSeparator)
	case filterTask:
		if f.task != nil {
			err = fmt.Errorf("%w: redeclared filter task", ErrInvalidData)
			return
		}
		f.task = strings.Split(values, valuesSeparator)
	case filterSince:
		if !f.since.IsZero() {
			err = fmt.Errorf("%w: redeclared filter since", ErrInvalidData)
			return
		}
		f.since, err = time.Parse(DateFormat, values)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
			return
		}
	case filterUntil:
		if !f.until.IsZero() {
			err = fmt.Errorf("%w: redeclared filter until", ErrInvalidData)
			return
		}
		f.until, err = time.Parse(DateFormat, values)
		if err != nil {
			err = fmt.Errorf("%w: %s", ErrInvalidData, err.Error())
			return
		}
	case filterTags:
		if f.tags != nil {
			err = fmt.Errorf("%w: redeclared filter tags", ErrInvalidData)
			return
		}
		f.tags = strings.Split(values, valuesSeparator)
	}
	return
}

// stringContainsAny checks if any of the valid values is inside the searchable
// list.
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

func convertFilterToSql(key string, values []string, operator string) string {
	if len(values) == 0 {
		return ""
	}
	b := strings.Builder{}
	for i, v := range values {
		if i == 0 {
			b.WriteString("(")
		} else {
			b.WriteString(" OR ")
		}
		switch operator {
		case sqlOperatorEquals:
			b.WriteString(fmt.Sprintf("json_extract(`json`,'$.%s')='%s'", key, v))
		case sqlOperatorLike:
			b.WriteString(fmt.Sprintf("json_extract(`json`,'$.%s') LIKE '%%%s%%'", key, v))
		}
	}
	b.WriteString(")")
	return b.String()
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
