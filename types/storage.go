package types

import (
	"github.com/google/uuid"
)

// Storage is the general storage interface that is used to abstract the
// direct data access.
type Storage interface {
	GetTimer(uuid uuid.UUID) (Timer, error)
	GetRunningTimer() (Timer, error)
	GetTimers(filter Filter) (Timers, error)
	RunningTimerExists() (bool, error)
	StoreTimer(timer Timer) error
	UpdateTimer(timer Timer) error
}
