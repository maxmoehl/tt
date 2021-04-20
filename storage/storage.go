package storage

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt/types"
)

func StartTimer(project, task, timestamp string, tags []string) error {
	if running, err := s.RunningTimerExists(); err != nil {
		return err
	} else if running {
		return fmt.Errorf("running timer found, cannot create a new one")
	}
	var start time.Time
	if timestamp == "" {
		start = time.Now()
	} else {
		var err error
		start, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return err
		}
		validStartTime, err := isValidStartTime(start)
		if err != nil {
			return err
		}
		if !validStartTime {
			return fmt.Errorf("the given start time is not valid, collision with other timer")
		}
	}
	return s.StoreTimer(types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   start,
		End:     time.Time{},
		Project: project,
		Task:    task,
		Tags:    tags,
		Breaks:  nil,
	})
}

func StopTimer(timestamp string) error {
	runningTimer, err := s.GetRunningTimer()
	if err != nil {
		return err
	}
	openBreak, breakIdx := runningTimer.Breaks.Open()
	var stop time.Time
	if timestamp == "" {
		stop = time.Now()
	} else {
		var err error
		stop, err = time.Parse(time.RFC3339, timestamp)
		if err != nil {
			return err
		}
		if !openBreak && len(runningTimer.Breaks) > 0 {
			breakIdx = getMostRecentBreak(runningTimer.Breaks)
		}
		if !isValidStopTime(runningTimer, breakIdx, stop) {
			return fmt.Errorf("the given stop time is not valid for this timer")
		}
	}
	if openBreak {
		runningTimer.Breaks[breakIdx].End = stop
	}
	runningTimer.End = stop
	return s.UpdateTimer(runningTimer)
}

// ToggleBreak starts a new break if no break is open or ends and open
// break. If there is no running timer, an error is returned. The
// returned bool indicates whether a break is running or not after
// the toggle has been executed.
func ToggleBreak() (bool, error) {
	runningTimer, err := s.GetRunningTimer()
	if err != nil {
		return false, err
	}
	openBreak, breakIdx := runningTimer.Breaks.Open()
	var breakOpenAfter bool
	if openBreak {
		runningTimer.Breaks[breakIdx].End = time.Now()
		breakOpenAfter = false
		err = s.UpdateTimer(runningTimer)
	} else {
		runningTimer.Breaks = append(runningTimer.Breaks, types.Break{
			Start: time.Now(),
		})
		breakOpenAfter = true
		err = s.UpdateTimer(runningTimer)
	}
	if err != nil {
		return !breakOpenAfter, err
	}
	return breakOpenAfter, nil
}

func CheckRunningTimers() ([]uuid.UUID, error) {
	timers, err := s.GetTimers(nil)
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

func CheckTimersOpenBreaks() ([]uuid.UUID, error) {
	timers, err := s.GetTimers(nil)
	if err != nil {
		return nil, err
	}
	var uuids []uuid.UUID
	for _, t := range timers {
		if t.Running() {
			openBreaks := 0
			for _, b := range t.Breaks {
				if b.Open() {
					openBreaks++
				}
			}
			if openBreaks > 1 {
				uuids = append(uuids, t.Uuid)
			}
		} else if openBreaks, _ := t.Breaks.Open(); openBreaks {
			uuids = append(uuids, t.Uuid)
		}
	}
	return uuids, nil
}

func GetRunningTimer() (bool, types.Timer, error) {
	exists, err := s.RunningTimerExists()
	if err != nil {
		return false, types.Timer{}, err
	}
	if !exists {
		return false, types.Timer{}, nil
	}
	t, err := s.GetRunningTimer()
	if err != nil {
		return false, types.Timer{}, err
	}
	return exists, t, nil
}

// isValidStartTime checks if there are other timers that are more recent
// than the time passed in.
func isValidStartTime(t time.Time) (bool, error) {
	if t.IsZero() {
		return false, nil
	}
	timer, err := getMostRecentTimer()
	if err != nil {
		return false, err
	}
	if timer.IsZero() {
		return true, nil
	}
	return t.After(timer.End), nil
}

// isValidStopTime checks if the given timer and break started before the
// timestamp. If breakIdx is -1, it is ignored.
func isValidStopTime(timer types.Timer, breakIdx int, time time.Time) bool {
	// if there is some sort of break passed in, check if that break started before
	// the given time
	if breakIdx != -1 {
		b := timer.Breaks[breakIdx]
		return timer.Start.Before(time) && b.Start.Before(time)
	}
	return timer.Start.Before(time)
}

// getMostRecentTimer returns the most recent timer. If no timer is found
// an empty timer is returned which can be identified by IsZero()
func getMostRecentTimer() (types.Timer, error) {
	timers, err := s.GetTimers(nil)
	if err != nil {
		return types.Timer{}, err
	}
	mostRecent := types.Timer{}
	for _, t := range timers {
		if t.Start.After(mostRecent.Start) {
			mostRecent = t
		}
	}
	return mostRecent, nil
}

func getMostRecentBreak(breaks types.Breaks) int {
	mostRecent := 0
	for i, b := range breaks {
		if b.Start.After(breaks[mostRecent].Start) {
			mostRecent = i
		}
	}
	return mostRecent
}
