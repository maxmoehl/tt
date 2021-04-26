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

package storage

import (
	"errors"
	"fmt"
	"time"

	"github.com/maxmoehl/tt/types"

	"github.com/google/uuid"
)

// StartTimer starts a new timer, and validates that the passed in values.
// If an error is returned no changes have been made, except if the error
// is from writing to the file.
func StartTimer(project, task, timestamp string, tags []string) (types.Timer, error) {
	lastTimer, err := s.GetLastTimer(true)
	if err != nil && !errors.Is(err, types.ErrNotFound) {
		return types.Timer{}, err
	}
	if lastTimer.Running() && !errors.Is(err, types.ErrNotFound) {
		return types.Timer{}, fmt.Errorf("running timer found, cannot create a new one")
	}
	start, err := getStartTime(timestamp)
	if err != nil {
		return types.Timer{}, err
	}
	t := types.Timer{
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
func StopTimer(timestamp string) (types.Timer, error) {
	runningTimer, err := GetRunningTimer()
	if err != nil {
		return types.Timer{}, err
	}
	var stop time.Time
	if timestamp == "" {
		stop = time.Now()
	} else {
		stop, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return types.Timer{}, err
		}
	}
	if stop.Before(runningTimer.Start) {
		return types.Timer{}, fmt.Errorf("stop time is before start time")
	}
	runningTimer.Stop = stop
	return runningTimer, s.UpdateTimer(runningTimer)
}

func ResumeTimer() (types.Timer, error) {
	runningTimer, err := s.GetLastTimer(true)
	if err != nil {
		return types.Timer{}, err
	}
	if runningTimer.Running() {
		return types.Timer{}, fmt.Errorf("found running timer, cannot resume")
	}
	t := types.Timer{
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
func CheckRunningTimers() ([]uuid.UUID, error) {
	timers, err := s.GetTimers(types.Filter{})
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
func GetRunningTimer() (types.Timer, error) {
	timer, err := s.GetLastTimer(true)
	if err != nil {
		return types.Timer{}, err
	}
	if timer.IsZero() || !timer.Running() {
		return types.Timer{}, fmt.Errorf("no running timer found")
	}
	return timer, nil
}

// GetTimers returns all timers after applying the filter
func GetTimers(filter types.Filter) (types.Timers, error) {
	return s.GetTimers(filter)
}

func getStartTime(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Now(), nil
	}
	start, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, err
	}
	lastTimer, err := s.GetLastTimer(false)
	if err != nil && !errors.Is(err, types.ErrNotFound) {
		return time.Time{}, err
	}
	if lastTimer.Stop.After(start) {
		return time.Time{}, fmt.Errorf("invalid start time, collision with existing timer")
	}
	return start, nil
}
