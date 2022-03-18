package tt

import (
	"encoding/csv"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	groupByProject GroupByOption = "project"
	groupByTask    GroupByOption = "task"
	groupByDay     GroupByOption = "day"
)

type GroupByOption string

// Timer is the central type that stores a timer and all its relevant values.
type Timer struct {
	ID      string     `json:"id"`
	Start   time.Time  `json:"start"`
	Stop    *time.Time `json:"stop,omitempty"`
	Project string     `json:"project"`
	Task    string     `json:"task,omitempty"`
	Tags    []string   `json:"tags,omitempty"`
}

func (t Timer) Validate() error {
	if _, err := uuid.Parse(t.ID); err != nil {
		return fmt.Errorf("%w: id is not a valid uuid", ErrInvalidTimer)
	}
	if t.Start.IsZero() {
		return fmt.Errorf("%w: start time is zero", ErrInvalidTimer)
	}
	if t.Stop != nil && t.Stop.IsZero() {
		return fmt.Errorf("%w: stop is non-nil but zero", ErrInvalidTimer)
	}
	if t.Project == "" {
		return fmt.Errorf("%w: project is an empty string", ErrInvalidTimer)
	}
	return nil
}

// Duration returns the duration that the timer has been running. If the timer
// is still running it will return the time it has run until now.
func (t Timer) Duration() time.Duration {
	if t.Stop == nil {
		return time.Now().Sub(t.Start)
	}
	return t.Stop.Sub(t.Start)
}

// Running indicates whether the timer is still running.
func (t Timer) Running() bool {
	return t.Stop == nil
}

func (t Timer) String() string {
	b := strings.Builder{}

	b.WriteString("ID      : ")
	b.WriteString(t.ID)
	b.WriteRune('\n')

	b.WriteString("Start   : ")
	b.WriteString(t.Start.String())
	b.WriteRune('\n')

	if !t.Running() {
		b.WriteString("Stop    : ")
		b.WriteString(t.Stop.String())
		b.WriteRune('\n')
	}

	b.WriteString("Duration: ")
	b.WriteString(FormatDuration(t.Duration()))
	b.WriteRune('\n')

	b.WriteString("Project : ")
	b.WriteString(t.Project)
	b.WriteRune('\n')

	if t.Task != "" {
		b.WriteString("Task    : ")
		b.WriteString(t.Task)
		b.WriteRune('\n')
	}
	if len(t.Tags) > 0 {
		b.WriteString("Tags    : ")
		b.WriteString(strings.Join(t.Tags, ", "))
	}

	return b.String()
}

func (t Timer) groupByKey(f GroupByOption) string {
	switch f {
	case groupByProject:
		return t.Project
	case groupByTask:
		if t.Task == "" {
			return "no-task"
		}
		return t.Task
	case groupByDay:
		return fmt.Sprintf("%04d-%02d-%02d", t.Start.Year(), t.Start.Month(), t.Start.Day())
	default:
		panic(fmt.Sprintf("%s is not a group by field", f))
	}
}

// Timers stores a list of timers to attach functions to it.
type Timers []Timer

func (timers Timers) Duration() (d time.Duration) {
	for _, t := range timers {
		d += t.Duration()
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
		stop := ""
		if t.Stop != nil {
			stop = t.Stop.String()
		}
		err = w.Write([]string{t.ID, t.Start.String(), stop, t.Project, t.Task, strings.Join(t.Tags, ",")})
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

// GroupByTask groups all timers by project and task.
func (timers Timers) GroupByTask() map[string]map[string]Timers {
	grouped := make(map[string]map[string]Timers)
	for k, v := range timers.GroupByProject() {
		grouped[k] = v.groupBy(groupByTask)
	}
	return grouped
}

func (timers Timers) GroupByProject() map[string]Timers {
	return timers.groupBy(groupByProject)
}

func (timers Timers) GroupByDay() map[string]Timers {
	return timers.groupBy(groupByDay)
}

func (timers Timers) groupBy(field GroupByOption) map[string]Timers {
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
