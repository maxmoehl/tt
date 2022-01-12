package tt

import (
	"fmt"
	"net/url"
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

// Filter contains all available filters. If a value is empty (i.e. ""
// or nil) it is assumed that the Filter is not set and all values are
// included.
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
	// Since Filter can be set with the keyword 'since'.
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
	if !f.since.IsZero() && t.Start.Before(f.since) {
		return false
	}
	if !f.until.IsZero() && (t.Start.After(f.until) || t.Stop.After(f.until)) {
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
	tags := convertFilterToSql("tags", f.tags, sqlOperatorLike)
	if tags != "" {
		filters = append(filters, tags)
	}
	if !f.since.IsZero() {
		filters = append(filters, fmt.Sprintf("start > %d", f.since.Unix()))
	}
	if !f.until.IsZero() {
		filters = append(filters, fmt.Sprintf("stop < %d", f.until.Unix()))
	}
	// if there are no filters return TRUE to match all values
	if len(filters) == 0 {
		return "TRUE"
	}
	return strings.Join(filters, " AND ")
}

// ParseFilterString takes a string and creates a Filter from it. The
// Filter string has to be in the following format:
//
//   filterName=values;filterName=values;...
//
// each filterName consists of a string, values contains the Filter
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
//   tags   : accepts multiple string values
//
// since and until are inclusive, both dates will be included in filtered
// data.
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
	if !f.until.IsZero() {
		f.until = f.until.Add(time.Hour * 24)
	}
	return f, nil
}

func ParseFilterQuery(q url.Values) (Filter, error) {
	var f Filter
	for qKey, qValues := range q {
		key, err := url.QueryUnescape(qKey)
		if err != nil {
			return f, ErrInvalidData.WithCause(err)
		}
		if len(qValues) != 1 {
			return f, ErrInvalidData.WithCause(NewError("got an empty query parameter or duplicate keys"))
		}
		values, err := url.QueryUnescape(qValues[0])
		if err != nil {
			return f, ErrInvalidData.WithCause(err)
		}
		err = parseValuesInto(key, values, &f)
		if err != nil {
			return f, err
		}
	}
	return f, nil
}

// NewFilter allows to create a Filter that is not parsed from a Filter
// string.
func NewFilter(projects, tasks, tags []string, since, until time.Time) Filter {
	return Filter{
		project: projects,
		task:    tasks,
		since:   since,
		until:   until,
		tags:    tags,
	}
}

func parseFilter(in string) (key, values string, err error) {
	filterSplit := strings.Split(in, "=")
	if len(filterSplit) != 2 {
		return "", "", ErrInvalidData.WithCause(NewErrorf("expected one '=' per Filter but got %d: [%s]", len(filterSplit), in))
	}
	return filterSplit[0], filterSplit[1], nil
}

func parseValuesInto(key, values string, f *Filter) (err Error) {
	var e error
	switch key {
	case filterProject:
		if f.project != nil {
			err = ErrInvalidData.WithCause(NewError("redeclared filter project"))
			return
		}
		f.project = strings.Split(values, valuesSeparator)
	case filterTask:
		if f.task != nil {
			err = ErrInvalidData.WithCause(NewError("redeclared filter task"))
			return
		}
		f.task = strings.Split(values, valuesSeparator)
	case filterSince:
		if !f.since.IsZero() {
			err = ErrInvalidData.WithCause(NewError("redeclared filter since"))
			return
		}
		f.since, e = time.Parse(DateFormat, values)
		if err != nil {
			err = ErrInvalidData.WithCause(e)
			return
		}
		f.since = time.Date(f.since.Year(), f.since.Month(), f.since.Day(), 0, 0, 0, 0, time.Local)
	case filterUntil:
		if !f.until.IsZero() {
			err = ErrInvalidData.WithCause(NewError("redeclared filter until"))
			return
		}
		f.until, e = time.Parse(DateFormat, values)
		if err != nil {
			err = ErrInvalidData.WithCause(e)
			return
		}
		f.until = time.Date(f.until.Year(), f.until.Month(), f.until.Day(), 0, 0, 0, 0, time.Local)
	case filterTags:
		if f.tags != nil {
			err = ErrInvalidData.WithCause(NewError("redeclared filter tags"))
			return
		}
		f.tags = strings.Split(values, valuesSeparator)
	}
	return
}

// stringSliceContainsAny checks if any of the valid values is inside the
// searchable list.
func stringSliceContainsAny(validValues, searchable []string) bool {
	for _, v := range validValues {
		if stringSliceContains(searchable, v) {
			return true
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
		if i > 0 {
			b.WriteString(" OR ")
		} else {
			b.WriteString("(")
		}
		switch operator {
		case sqlOperatorEquals:
			b.WriteString(key)
			b.WriteString("='")
			b.WriteString(v)
			b.WriteString("'")
		case sqlOperatorLike:
			b.WriteString(key)
			b.WriteString(" LIKE '%")
			b.WriteString(v)
			b.WriteString("%'")
		}
	}
	b.WriteString(")")
	return b.String()
}

// stringSliceContains checks if the given string is contained in the
// given string slice.
func stringSliceContains(strings []string, s string) bool {
	for _, t := range strings {
		if t == s {
			return true
		}
	}
	return false
}
