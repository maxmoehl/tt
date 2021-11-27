package tt

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestTest(t *testing.T) {
	var f Filter
	f.Match(Timer{})
}

func TestStorageInterface(t *testing.T) {
	storageTypes := map[string]func() error{
		"sqlite": setupSqliteTest,
		"file":   setupFileTest,
	}

	testCases := map[string]map[string]func(t *testing.T){
		"GetTimer": {
			"timer not found":          GetTimerNotFound,
			"no error on timer exists": GetTimersExist,
			"check correct timer":      GetTimerCorrectTimer,
		},
		"GetLastTimer": {
			"true: w/ running timer":            GetLastTimerRunningTrueExists,
			"true w/o running timer":            GetLastTimerRunningTrueNotExists,
			"true: no timer found":              GetLastTimerTrueNotFound,
			"false: w/ running timer":           GetLastTimerRunningFalseExists,
			"false: w/o running timer":          GetLastTimerRunningFalseNotExists,
			"false: no timer found":             GetLastTimerFalseNotFound,
			"false: not found w/ running timer": GetLastTimerFalseNotFoundWithRunning,
		},
		"GetTimers": {
			"not found error": GetTimersNotFound,
			"no filter":       GetTimersNoFilter,
		},
		"StoreTimer": {
			"normal timer":         StoreTimer,
			"zero timer":           StoreTimerZeroTimer,
			"timer already exists": StoreTimerDuplicateTimer,
		},
		"UpdateTimer": {
			"timer does not exist": UpdateTimerNotExist,
			"timer gets updated":   UpdateTimer,
		},
	}

	// for each available storage type
	for storageType, setupFunc := range storageTypes {
		t.Run(storageType, func(t *testing.T) {
			err := setupFunc()
			if err != nil {
				t.Fatal(err.Error())
			}
			// for each test group
			for testGroup, testFuncs := range testCases {
				t.Run(testGroup, func(t *testing.T) {
					// execute all test cases
					for testName, testFunc := range testFuncs {
						// and clear the storage each time to make sure all tests
						// have the same starting conditions
						err = clearStorage()
						if err != nil {
							t.Fatal(err.Error())
						}
						t.Run(testName, testFunc)
					}
				})
			}
		})
	}

}

func GetTimerNotFound(t *testing.T) {
	err := writeTimerToStorage(Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	})
	if err != nil {
		t.Fatalf("unable to write test timer: %s", err.Error())
	}
	_, err = s.GetTimer(uuid.Must(uuid.NewRandom()))
	if err == nil {
		t.Fatal("expected error because of missing timer")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected error to wrap ErrNotFound")
	}
}

func GetTimersExist(t *testing.T) {
	timer := Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "test-project",
		Task:    "task",
		Tags:    []string{"a", "b", "c"},
	}
	err := writeTimerToStorage(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStorage, err := s.GetTimer(timer.Uuid)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timer, timerStorage)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetTimerCorrectTimer(t *testing.T) {
	timers := getRandomTimers(4)
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timer := timers[2]
	timerStorage, err := s.GetTimer(timer.Uuid)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timer, timerStorage)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetLastTimerRunningTrueExists(t *testing.T) {
	timers := getRandomTimers(5)
	// Set most recent timer to running
	timers[0].Stop = time.Time{}
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStore, err := s.GetLastTimer(true)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timers[0], timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetLastTimerRunningTrueNotExists(t *testing.T) {
	timers := getRandomTimers(5)
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStore, err := s.GetLastTimer(true)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timers[0], timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetLastTimerRunningFalseExists(t *testing.T) {
	timers := getRandomTimers(5)
	// Set most recent timer to running
	timers[0].Stop = time.Time{}
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStore, err := s.GetLastTimer(false)
	if err != nil {
		t.Fatal(err.Error())
	}
	// we expect the second most recent timer since the last timer is
	// still running
	err = timersEqual(timers[1], timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetLastTimerRunningFalseNotExists(t *testing.T) {
	timers := getRandomTimers(5)
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStore, err := s.GetLastTimer(false)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timers[0], timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func GetLastTimerTrueNotFound(t *testing.T) {
	_, err := s.GetLastTimer(true)
	if err == nil {
		t.Fatal("expected error because of missing timer")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected error to wrap types.ErrNotFound")
	}
}

func GetLastTimerFalseNotFound(t *testing.T) {
	_, err := s.GetLastTimer(false)
	if err == nil {
		t.Fatal("expected error because of missing timer")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected error to wrap types.ErrNotFound")
	}
}

func GetLastTimerFalseNotFoundWithRunning(t *testing.T) {
	timer := getRandomTimers(1)[0]
	timer.Stop = time.Time{}
	err := writeTimerToStorage(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	_, err = s.GetLastTimer(false)
	if err == nil {
		t.Fatal("expected error because of missing timer")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected error to wrap types.ErrNotFound")
	}
}

func GetTimersNotFound(t *testing.T) {
	timers, err := s.GetTimers(Filter{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(timers) != 0 {
		t.Fatalf("expected to get 0 timers but got %d", len(timers))
	}
}

func GetTimersNoFilter(t *testing.T) {
	timers := getRandomTimers(5)
	err := writeTimersToStorage(timers)
	if err != nil {
		t.Fatal(err.Error())
	}
	timersStorage, err := s.GetTimers(Filter{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(timers) != len(timersStorage) {
		t.Fatal("expected same number of timers")
	}
	for i := range timers {
		err = timersEqual(timers[i], timersStorage[i])
		if err != nil {
			t.Fatal(err.Error())
		}
	}
}

func StoreTimer(t *testing.T) {
	timer := getRandomTimers(1)[0]
	err := s.StoreTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	exists, e := storageContainsTimer(timer)
	if e != nil {
		t.Fatal(err.Error())
	}
	if !exists {
		t.Fatal("expected timer to be in storage")
	}
}

func StoreTimerZeroTimer(t *testing.T) {
	err := s.StoreTimer(Timer{})
	if err == nil {
		t.Fatal("expected error because of zero timer")
	}
}

func StoreTimerDuplicateTimer(t *testing.T) {
	timer := getRandomTimers(1)[0]
	err := writeTimerToStorage(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = s.StoreTimer(timer)
	if err == nil {
		t.Fatal("expected error because of duplicate uuid")
	}
}

func UpdateTimerNotExist(t *testing.T) {
	err := s.UpdateTimer(getRandomTimers(1)[0])
	if err == nil {
		t.Fatal("expected error because of non existent timer")
	}
	if !errors.Is(err, ErrNotFound) {
		t.Fatal("expected error to wrap types.ErrNotFound")
	}
}

func UpdateTimer(t *testing.T) {
	timer := getRandomTimers(1)[0]
	timer.Stop = time.Time{}
	err := writeTimerToStorage(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	timer.Stop = time.Now()
	err = s.UpdateTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	exists, err := storageContainsTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !exists {
		t.Fatal("expected updated timer to be in storage")
	}
}
