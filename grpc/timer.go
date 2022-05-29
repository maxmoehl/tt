package grpc

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt"
)

func fromTt(timer tt.Timer) *Timer {
	stop := ""
	if !timer.Running() {
		stop = timer.Stop.Format(time.RFC3339)
	}
	return &Timer{
		Id:      timer.ID,
		Project: timer.Project,
		Task:    timer.Task,
		Tags:    timer.Tags,
		Start:   timer.Start.Format(time.RFC3339),
		Stop:    stop,
	}
}

func (x *Timer) StartTime() (time.Time, error) {
	if x == nil {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, x.Start)
}

func (x *Timer) StopTime() (time.Time, error) {
	if x == nil || x.Stop == "" {
		return time.Time{}, nil
	}
	return time.Parse(time.RFC3339, x.Stop)
}

func (x *Timer) Tt() (tt.Timer, error) {
	if x == nil {
		return tt.Timer{}, fmt.Errorf("%w: nil", tt.ErrInvalidTimer)
	}
	if _, err := uuid.Parse(x.Id); err != nil {
		return tt.Timer{}, fmt.Errorf("%w: id is not a valid uuid", tt.ErrInvalidTimer)
	}
	// can never be zero
	start, err := x.StartTime()
	if err != nil {
		return tt.Timer{}, fmt.Errorf("%w: %s", tt.ErrInvalidTimer, err.Error())
	} else if start.IsZero() {
		return tt.Timer{}, fmt.Errorf("%w: start time is zero", tt.ErrInvalidTimer)
	}
	// can be zero, but needs to be valid if it's present
	var stop *time.Time
	*stop, err = x.StopTime()
	if err != nil {
		return tt.Timer{}, fmt.Errorf("%w: %s", tt.ErrInvalidTimer, err.Error())
	} else if !stop.IsZero() && stop.Unix() <= start.Unix() {
		return tt.Timer{}, fmt.Errorf("%w: stop is before or equal to start", tt.ErrInvalidTimer)
	} else if stop.IsZero() {
		stop = nil
	}
	if x.Project == "" {
		return tt.Timer{}, fmt.Errorf("%w: project is an empty string", tt.ErrInvalidTimer)
	}
	return tt.Timer{
		ID:      x.Id,
		Start:   start,
		Stop:    nil,
		Project: x.Project,
		Task:    x.Task,
		Tags:    x.Tags,
	}, nil
}
