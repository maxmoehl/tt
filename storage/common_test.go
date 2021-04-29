package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt/test"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"
)

func TestMain(m *testing.M) {
	test.Main(m.Run)
}

func clearStorage() error {
	var err error
	if fileStorage, ok := s.(*file); ok {
		fileStorage.timers = nil
	} else if sqliteStorage, ok := s.(*sqlite); ok {
		_, err = sqliteStorage.db.Exec("DELETE FROM timers WHERE TRUE")
	} else {
		err = fmt.Errorf("unknown storage type: %T", s)
	}
	return err
}

func writeTimersToStorage(timers types.Timers) (err error) {
	for _, t := range timers {
		err = writeTimerToStorage(t)
		if err != nil {
			return
		}
	}
	return
}

func writeTimerToStorage(timer types.Timer) error {
	var err error
	if fileStorage, ok := s.(*file); ok {
		fileStorage.timers = append(fileStorage.timers, timer)
	} else if sqliteStorage, ok := s.(*sqlite); ok {
		id := timer.Uuid.String()
		start := timer.Start.Unix()
		var stop, task, tags interface{}
		if !timer.Stop.IsZero() {
			stop = timer.Stop.Unix()
		}
		if timer.Task != "" {
			task = timer.Task
		}
		if len(timer.Tags) > 0 {
			tags = strings.Join(timer.Tags, ",")
		}
		_, err = sqliteStorage.db.Exec(
			"INSERT INTO timers (uuid, start, stop, project, task, tags) VALUES (?, ?, ?, ?, ?, ?);",
			id, start, stop, timer.Project, task, tags)
	} else {
		err = fmt.Errorf("unknown storage type: %T", s)
	}
	return err
}

func storageContainsTimer(timer types.Timer) (bool, error) {
	if fileStorage, ok := s.(*file); ok {
		found := false
		for _, timerStorage := range fileStorage.timers {
			if err := timersEqual(timer, timerStorage); err == nil {
				found = true
			}
		}
		return found, nil
	} else if sqliteStorage, ok := s.(*sqlite); ok {
		row := sqliteStorage.db.QueryRow("SELECT uuid, start, stop, project, task, tags FROM timers WHERE uuid = ?",
			timer.Uuid.String())
		timerStorage, err := sqliteStorage.scanRow(row)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return false, nil
			}
			return false, err
		}
		err = timersEqual(timer, timerStorage)
		return err == nil, nil
	} else {
		return false, fmt.Errorf("unknown storage type: %T", s)
	}
}

// getRandomTimers returns n randomly created timers. The timers are
// guaranteed not to collide with each other and be finished. The first
// timer returned is the most recent.
func getRandomTimers(n int) (timers types.Timers) {
	/*
		i | start time     | stop time
		0 | now - 1h + 10m | now - 0h - 10m
		1 | now - 2h + 10m | now - 1h - 10m
		2 | now - 3h + 10m | now - 2h - 10m
		3 | ...
	*/
	for i := 0; i < n; i++ {
		start := time.Now().Add(-time.Duration(i+1)*time.Hour + 10*time.Minute)
		stop := time.Now().Add(-time.Duration(i)*time.Hour + (-10 * time.Minute))
		timers = append(timers, types.Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   start,
			Stop:    stop,
			Project: fmt.Sprintf("test-%d", i),
			Task:    fmt.Sprintf("test-task-%d", i),
			Tags:    []string{"a", "b"},
		})
	}
	return
}

func timersEqual(t1, t2 types.Timer) error {
	if t1.Uuid != t2.Uuid {
		return fmt.Errorf("uuids are not equal")
	}
	if t1.Project != t2.Project {
		return fmt.Errorf("projects are not equal")
	}
	if t1.Task != t2.Task {
		return fmt.Errorf("tasks are not equal")
	}
	if len(t1.Tags) != len(t2.Tags) {
		return fmt.Errorf("tags are not equal")
	}
	for _, tag := range t1.Tags {
		if !utils.StringSliceContains(t2.Tags, tag) {
			return fmt.Errorf("tags are not equal")
		}
	}
	// cut anything beyond second precision since our storage interface only requires
	// precision up to seconds
	if !time.Unix(t1.Start.Unix(), 0).Equal(time.Unix(t2.Start.Unix(), 0)) {
		return fmt.Errorf("start time is not equal")
	}
	if !time.Unix(t1.Stop.Unix(), 0).Equal(time.Unix(t2.Stop.Unix(), 0)) {
		return fmt.Errorf("end time is not equal")
	}
	return nil
}
