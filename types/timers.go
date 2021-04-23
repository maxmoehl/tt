/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package types

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Timer is the central type that stores a timer and all its relevant
// values
type Timer struct {
	Uuid    uuid.UUID `json:"uuid"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Project string    `json:"project"`
	Task    string    `json:"task,omitempty"`
	Tags    []string  `json:"tags,omitempty"`
}

// Duration returns the duration that the timer has been running,
// excluding any breaks.
func (t Timer) Duration() time.Duration {
	return t.End.Sub(t.Start)
}

// Running indicates whether or not the timer is still running.
func (t Timer) Running() bool {
	return t.End.IsZero()
}

// IsZero checks if the timer has been properly initialized.
func (t Timer) IsZero() bool {
	return t.Start.IsZero() && t.End.IsZero()
}

// Timers stores a list of timers to attach functions to it.
type Timers []Timer

// Running checks if any running timers exist and returns the index if one
// is found.
func (timers Timers) Running() (bool, int) {
	for i, ws := range timers {
		if ws.Running() {
			return true, i
		}
	}
	return false, -1
}

// Last returns the last timer in the list ordered by start time. Running
// timers can be excluded by passing false in.
func (timers Timers) Last(running bool) (t Timer) {
	for _, timer := range timers {
		if !running && timer.End.IsZero() {
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

// Filter applies the given filter to every element and returns those that
// match the filter.
func (timers Timers) Filter(f Filter) (filtered Timers) {
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

func (timers Timers) CSV() (string, error) {
	b := strings.Builder{}
	w := csv.NewWriter(&b)
	err := w.Write([]string{"uuid", "start", "end", "project", "task", "tags"})
	if err != nil {
		return "", err
	}
	for _, t := range timers {
		err = w.Write([]string{t.Uuid.String(), t.Start.String(), t.End.String(), t.Project, t.Task, strings.Join(t.Tags, ", ")})
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

func (timers Timers) SQL() (string, error) {
	b := strings.Builder{}
	b.WriteString("CREATE TABLE IF NOT EXISTS timers (uuid TEXT PRIMARY KEY, start INTEGER NOT NULL, stop INTEGER, project TEXT NOT NULL, task TEXT, tags TEXT);\n")
	for _, t := range timers {
		var stop interface{}
		var task, tags string
		if !t.End.IsZero() {
			stop = t.End.Unix()
		} else {
			stop = "NULL"
		}
		if t.Task != "" {
			task = "'" + t.Task + "'"
		} else {
			task = "NULL"
		}
		if len(t.Tags) > 0 {
			tags = "'" + strings.Join(t.Tags, ",") + "'"
		} else {
			tags = "NULL"
		}
		b.WriteString(fmt.Sprintf("INSERT INTO timers (uuid, start, stop, project, task, tags) VALUES ('%v', %v, %v, '%v', %v, %v);\n",
			t.Uuid.String(), t.Start.Unix(), stop, t.Project, task, tags))
	}
	return b.String(), nil
}
