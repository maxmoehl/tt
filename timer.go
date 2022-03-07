package tt

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Timer is the central type that stores a timer and all its relevant
// values
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

// Duration returns the duration that the timer has been running.
// If the timer is still running it will return the time it has run
// until now.
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
	b.WriteString(FormatDuration(t.Duration(), GetConfig().GetPrecision()))
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
	case GroupByProject:
		return t.Project
	case GroupByTask:
		if t.Task == "" {
			return "no-task"
		}
		return t.Task
	case GroupByDay:
		return fmt.Sprintf("%04d-%02d-%02d", t.Start.Year(), t.Start.Month(), t.Start.Day())
	default:
		panic(fmt.Sprintf("%s is not a group by field", f))
	}
}
