package tt

import (
	"encoding/csv"
	"strings"
	"time"
)

// Timers stores a list of timers to attach functions to it.
type Timers []Timer

// Running checks if any running timers exist and returns the index if one
// is found. If the index is -1, no running timer exists.
func (timers Timers) Running() int {
	for i, ws := range timers {
		if ws.Running() {
			return i
		}
	}
	return -1
}

func (timers Timers) Duration() (d time.Duration) {
	for _, t := range timers {
		d += t.Duration()
	}
	return
}

// Last returns the last timer in the list ordered by start time. Running
// timers can be excluded by passing false in.
func (timers Timers) Last(running bool) (t Timer) {
	for _, timer := range timers {
		if !running && timer.Stop.IsZero() {
			continue
		}
		if timer.Start.After(t.Start) {
			t = timer
		}
	}
	return
}

// First returns the first timer from the list, ordered by start time.
func (timers Timers) First() (t Timer) {
	for _, timer := range timers {
		if timer.Start.Before(t.Start) || t.Start.IsZero() {
			t = timer
		}
	}
	return
}

// CSV exports all timers as a csv string
func (timers Timers) CSV() (string, error) {
	b := strings.Builder{}
	w := csv.NewWriter(&b)
	err := w.Write([]string{"uuid", "start", "end", "project", "task", "tags"})
	if err != nil {
		return "", err
	}
	for _, t := range timers {
		err = w.Write([]string{t.ID, t.Start.String(), t.Stop.String(), t.Project, t.Task, strings.Join(t.Tags, ",")})
		if err != nil {
			return "", err
		}
	}
	w.Flush()
	err = w.Error()
	if err != nil {
		return "", err
	}
	return b.String(), nil
}

func (timers Timers) GroupBy(field GroupByOption) map[string]Timers {
	grouped := make(map[string]Timers)
	for _, t := range timers {
		key := t.groupByKey(field)
		if _, ok := grouped[key]; ok {
			grouped[key] = append(grouped[key], t)
		} else {
			grouped[key] = Timers{t}
		}
	}
	return grouped
}
