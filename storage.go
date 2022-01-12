package tt

import (
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
	// UpdateTimer only allows the stop time to be set if it's null,
	// any other changes will be discarded.
	UpdateTimer(timer Timer) Error
}

var storagesMap = map[string]NewStorage{
	"file":   NewFile,
	"sqlite": NewSQLite,
}

var s Storage

func GetStorage() Storage {
	return s
}

func InitStorage() (err Error) {
	c := GetConfig()
	creator, ok := storagesMap[c.StorageType]
	if !ok {
		return NewErrorf("unknown storage type: %s", c.StorageType)
	}
	s, err = creator(c)
	return err
}
