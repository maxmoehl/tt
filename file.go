package tt

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type file struct {
	config      Config
	timers      Timers
	storagePath string
}

func (f *file) GetTimer(uuid uuid.UUID) (Timer, Error) {
	for _, t := range f.timers {
		if t.Uuid == uuid {
			return t, nil
		}
	}
	return Timer{}, ErrNotFound
}

func (f *file) GetLastTimer(running bool) (Timer, Error) {
	t := f.timers.Last(running)
	if t.IsZero() {
		return Timer{}, ErrNotFound
	}
	return t, nil
}

func (f *file) GetTimers(filter Filter) (Timers, Error) {
	timers := filter.Timers(f.timers)
	return timers, nil
}

func (f *file) StoreTimer(newTimer Timer) Error {
	if newTimer.IsZero() {
		return ErrInvalidData.WithCause(NewErrorf("timer is zero"))
	}
	exists := false
	for _, t := range f.timers {
		if t.Uuid == newTimer.Uuid {
			exists = true
			break
		}
	}
	if exists {
		return ErrInvalidData.WithCause(NewErrorf("timer with uuid %s already exists", newTimer.Uuid.String()))
	}
	f.timers = append(f.timers, newTimer)
	return f.write()
}

func (f *file) UpdateTimer(updatedTimer Timer) Error {
	updated := false
	for i, t := range f.timers {
		if t.Uuid == updatedTimer.Uuid {
			f.timers[i].Stop = updatedTimer.Stop
			updated = true
			break
		}
	}
	if !updated {
		return ErrNotFound
	}
	return f.write()
}

func (f *file) write() Error {
	fileWriter, err := os.Create(f.storagePath)
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	err = json.NewEncoder(fileWriter).Encode(f.timers)
	if err != nil {
		return ErrInternalError.WithCause(err)
	}
	return nil
}

// NewFile initializes and returns a new storage interface that can be used
// to access data.
func NewFile(c Config) (Storage, Error) {
	f := file{
		config:      c,
		storagePath: filepath.Join(c.HomeDir(), "storage.json"),
	}
	fileReader, err := os.Open(f.storagePath)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		return &f, nil
	} else if err != nil {
		return nil, ErrInternalError.WithCause(err)
	}
	err = json.NewDecoder(fileReader).Decode(&f.timers)
	if err != nil {
		return nil, ErrInternalError.WithCause(err)
	}
	return &f, nil
}
