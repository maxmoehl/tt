package file

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt/types"
)

type file struct {
	timers types.Timers
}

func (f *file) GetTimer(uuid uuid.UUID) (types.Timer, error) {
	for _, t := range f.timers {
		if t.Uuid == uuid {
			return t, nil
		}
	}
	return types.Timer{}, fmt.Errorf("no timer found for uuid %s", uuid.String())
}

func (f *file) GetRunningTimer() (types.Timer, error) {
	for _, t := range f.timers {
		if t.End.IsZero() {
			return t, nil
		}
	}
	return types.Timer{}, fmt.Errorf("no running timer found")
}

func (f *file) GetTimers(filter types.Filter) (types.Timers, error) {
	return f.timers.Filter(filter), nil
}

func (f *file) RunningTimerExists() (bool, error) {
	for _, t := range f.timers {
		if t.Running() {
			return true, nil
		}
	}
	return false, nil
}

func (f *file) StoreTimer(newTimer types.Timer) error {
	exists := false
	for _, t := range f.timers {
		if t.Uuid == newTimer.Uuid {
			exists = true
			break
		}
	}
	if exists {
		return fmt.Errorf("timer with uuid %s already exists", newTimer.Uuid.String())
	}
	f.timers = append(f.timers, newTimer)
	return f.write()
}

func (f *file) UpdateTimer(updatedTimer types.Timer) error {
	for i, t := range f.timers {
		if t.Uuid == updatedTimer.Uuid {
			f.timers[i].Project = updatedTimer.Project
			f.timers[i].Task = updatedTimer.Task
			f.timers[i].Start = updatedTimer.Start
			f.timers[i].End = updatedTimer.End
			f.timers[i].Breaks = updatedTimer.Breaks
			break
		}
	}
	return f.write()
}

func (f *file) write() error {
	fileWriter, err := os.Create(getStorageFile())
	if err != nil {
		return err
	}
	return json.NewEncoder(fileWriter).Encode(f.timers)
}
