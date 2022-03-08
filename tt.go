package tt

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

func Start(project, task string, tags []string, timestamp time.Time, copy int) (Timer, error) {
	db := GetDB()
	orderBy := OrderBy{
		Field: FieldStart,
		Order: OrderDsc,
	}
	var baseTimer Timer
	if copy > 0 {
		var timers Timers
		err := db.GetTimers(Filter{}, orderBy, &timers)
		if err != nil {
			return Timer{}, fmt.Errorf("start: %w", err)
		}
		if len(timers) < copy {
			return Timer{}, fmt.Errorf("start: copy from timer: %w", ErrNotFound)
		}
		baseTimer = timers[copy-1]
	}

	t := Timer{
		ID:      uuid.Must(uuid.NewRandom()).String(),
		Start:   timestamp,
		Project: project,
		Task:    task,
		Tags:    tags,
	}

	// copy values if they haven't been provided, and we should copy
	if copy > 0 && t.Project == "" {
		t.Project = baseTimer.Project
	}
	if copy > 0 && t.Task == "" {
		t.Task = baseTimer.Task
	}
	if copy > 0 && len(t.Tags) == 0 {
		t.Tags = baseTimer.Tags
	}

	err := t.Validate()
	if err != nil {
		return Timer{}, fmt.Errorf("start: %w", err)
	}
	err = GetDB().SaveTimer(t)
	if err != nil {
		return Timer{}, fmt.Errorf("start: %w", err)
	}
	return t, nil
}

func Stop(timestamp time.Time) (Timer, error) {
	db := GetDB()
	orderBy := OrderBy{
		Field: FieldStart,
		Order: OrderDsc,
	}
	var timer Timer
	err := db.GetTimer(Filter{}, orderBy, &timer)
	if err != nil {
		return Timer{}, fmt.Errorf("stop: %w", err)
	}
	if !timer.Running() {
		return Timer{}, fmt.Errorf("stop: %w: running timer", ErrNotFound)
	}
	timer.Stop = &timestamp
	err = db.UpdateTimer(timer)
	if err != nil {
		return Timer{}, fmt.Errorf("stop: %w", err)
	}
	return timer, nil
}
