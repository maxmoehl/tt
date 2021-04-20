package types

import (
	"time"

	"github.com/google/uuid"
)

type Break struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

func (b Break) Duration() time.Duration {
	return b.End.Sub(b.Start)
}

func (b Break) Open() bool {
	return b.End.IsZero()
}

type Breaks []Break

func (breaks Breaks) Duration() (d time.Duration) {
	for _, b := range breaks {
		if b.Open() {
			continue
		}
		d += b.Duration()
	}
	return
}

func (breaks Breaks) Open() (bool, int) {
	for i, b := range breaks {
		if b.Open() {
			return true, i
		}
	}
	return false, -1
}

type Timer struct {
	Uuid    uuid.UUID `json:"uuid"`
	Start   time.Time `json:"start"`
	End     time.Time `json:"end"`
	Project string    `json:"project"`
	Task    string    `json:"task,omitempty"`
	Breaks  Breaks    `json:"breaks,omitempty"`
}

func (t Timer) Duration() time.Duration {
	return t.End.Sub(t.Start) - t.Breaks.Duration()
}

func (t Timer) Running() bool {
	return t.End.IsZero()
}

func (t Timer) IsZero() bool {
	return t.Start.IsZero() && t.End.IsZero()
}

type Timers []Timer

func (timers Timers) Running() (bool, int) {
	for i, ws := range timers {
		if ws.Running() {
			return true, i
		}
	}
	return false, -1
}

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
