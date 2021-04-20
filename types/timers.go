package types

import (
	"time"

	"github.com/google/uuid"
)

// Break stores a single break.
type Break struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Duration returns the duration of the break.
func (b Break) Duration() time.Duration {
	return b.End.Sub(b.Start)
}

// Open returns whether or not the break is still open.
func (b Break) Open() bool {
	return b.End.IsZero()
}

// Breaks stores a list of breaks. This type exists to attach functions
// to it.
type Breaks []Break

// Duration returns the accumulated duration of all breaks contained
func (breaks Breaks) Duration() (d time.Duration) {
	for _, b := range breaks {
		if b.Open() {
			continue
		}
		d += b.Duration()
	}
	return
}

// Open indicates whether an open break exists and returns the index if so.
func (breaks Breaks) Open() (bool, int) {
	for i, b := range breaks {
		if b.Open() {
			return true, i
		}
	}
	return false, -1
}

// Timer is the central type that stores a timer and all its relevant
// values
type Timer struct {
	Uuid    uuid.UUID `json:"uuid"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Project string    `json:"project"`
	Task    string    `json:"task,omitempty"`
	Tags    []string  `json:"tags,omitempty"`
	Breaks  Breaks    `json:"breaks,omitempty"`
}

// Duration returns the duration that the timer has been running,
// excluding any breaks.
func (t Timer) Duration() time.Duration {
	return t.End.Sub(t.Start) - t.Breaks.Duration()
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
