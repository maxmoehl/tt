package tt

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func Resume(timestamp time.Time) (Timer, error) {
	db := GetDB()
	orderBy := OrderBy{
		Field: FieldStart,
		Order: OrderDsc,
	}
	var timer Timer
	err := db.GetTimer(Filter{}, orderBy, &timer)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return Timer{}, fmt.Errorf("resume: %w", err)
	}
	timer.ID = uuid.Must(uuid.NewRandom()).String()
	timer.Start = timestamp
	timer.Stop = nil
	err = db.SaveTimer(timer)
	if err != nil {
		return Timer{}, fmt.Errorf("resume: %w", err)
	}
	return timer, nil
}

func Start(project, task string, tags []string, timestamp time.Time) (Timer, error) {
	t := Timer{
		ID:      uuid.Must(uuid.NewRandom()).String(),
		Start:   timestamp,
		Project: project,
		Task:    task,
		Tags:    tags,
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
