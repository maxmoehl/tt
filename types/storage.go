package types

import (
	"github.com/google/uuid"
)

type Interface interface {
	GetTimer(uuid uuid.UUID) (Timer, error)
	GetRunningTimer() (Timer, error)
	GetTimers(filter Filter) (Timers, error)
	RunningTimerExists() (bool, error)
	StoreTimer(timer Timer) error
	UpdateTimer(timer Timer) error
}
