package tt

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// Storage is the general storage interface that is used to abstract the
// direct data access.
type Storage interface {
	// GetTimer returns the timer specified by the given uuid. If no timer
	// is found, ErrNotFound is returned
	GetTimer(uuid uuid.UUID) (Timer, Error)
	// GetLastTimer returns either the last timer of all timers if running
	// is true or the last non-running timer if running is false. If no
	// timer is found ErrNotFound is returned or wrapped in the
	// returned error.
	GetLastTimer(running bool) (Timer, Error)
	// GetTimers returns all Timers that match the filter. The result can be
	// nil if no timers are found.
	GetTimers(filter Filter) (Timers, Error)
	// StoreTimer writes the given timer to the configured data source.
	StoreTimer(timer Timer) Error
	// UpdateTimer only allows the stop time to be updated
	// any other changes will be discarded.
	UpdateTimer(timer Timer) Error
}

// StartTimer starts a new timer, and validates that the passed in values.
// If an error is returned no changes have been made, except if the error
// is from writing to the file.
func StartTimer(project, task, timestamp string, tags []string) (Timer, Error) {
	lastTimer, err := s.GetLastTimer(true)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Timer{}, err
	}
	if lastTimer.Running() && !errors.Is(err, ErrNotFound) {
		return Timer{}, ErrBadRequest.WithCause(NewError("running timer found, cannot create a new one"))
	}
	start, err := getStartTime(timestamp)
	if err != nil {
		return Timer{}, err
	}
	t := Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   start,
		Project: project,
		Task:    task,
		Tags:    tags,
	}
	return t, s.StoreTimer(t)
}

// StopTimer stops a timer and validates the given timestamp (if any).
// If an error is returned no changes have been made, except if the error
// is from writing to the file.
func StopTimer(timestamp string) (Timer, Error) {
	runningTimer, err := GetRunningTimer()
	if err != nil {
		return Timer{}, err
	}
	var stop time.Time
	var e error
	if timestamp == "" {
		stop = time.Now()
	} else {
		stop, e = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return Timer{}, ErrBadRequest.WithCause(e)
		}
	}
	if stop.Before(runningTimer.Start) {
		return Timer{}, ErrBadRequest.WithCause(NewError("stop time is before start time"))
	}
	runningTimer.Stop = stop
	return runningTimer, s.UpdateTimer(runningTimer)
}

// ResumeTimer takes the last timer from the storage and copies project,
// task and tags and starts a new timer.
func ResumeTimer() (Timer, Error) {
	runningTimer, err := s.GetLastTimer(true)
	if err != nil {
		return Timer{}, err
	}
	if runningTimer.Running() {
		return Timer{}, ErrBadRequest.WithCause(NewError("found running timer, cannot resume"))
	}
	t := Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: runningTimer.Project,
		Task:    runningTimer.Task,
		Tags:    runningTimer.Tags,
	}
	return t, s.StoreTimer(t)
}

// CheckRunningTimers returns the uuids of all timers that are currently
// running.
func CheckRunningTimers() ([]uuid.UUID, Error) {
	timers, err := s.GetTimers(Filter{})
	if err != nil {
		return nil, err
	}
	var uuids []uuid.UUID
	for _, t := range timers {
		if t.Running() {
			uuids = append(uuids, t.Uuid)
		}
	}
	return uuids, nil
}

// GetRunningTimer returns the running timer or an error if there is no
// running timer or any error that occurred during data access.
func GetRunningTimer() (Timer, Error) {
	timer, err := s.GetLastTimer(true)
	if err != nil {
		return Timer{}, err
	}
	if timer.IsZero() || !timer.Running() {
		return Timer{}, ErrNotFound
	}
	return timer, nil
}

// GetTimers returns all timers after applying the filter
func GetTimers(filter Filter) (Timers, Error) {
	return s.GetTimers(filter)
}

func getStartTime(timestamp string) (time.Time, Error) {
	if timestamp == "" {
		return time.Now(), nil
	}
	start, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, ErrBadRequest.WithCause(err)
	}
	lastTimer, err := s.GetLastTimer(false)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return time.Time{}, ErrInternalError.WithCause(err)
	}
	if lastTimer.Stop.After(start) {
		return time.Time{}, ErrBadRequest.WithCause(NewError("invalid start time, collision with existing timer"))
	}
	return start, nil
}
