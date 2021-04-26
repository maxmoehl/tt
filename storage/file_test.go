package storage

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/maxmoehl/tt/test"
	"github.com/maxmoehl/tt/types"

	"github.com/google/uuid"
)

const fileConfig = `precision: m
workHours: 8
storageType: file
workDays:
  monday: true
  tuesday: true
  wednesday: true
  thursday: true
  friday: true
  saturday: false
  sunday: false`

var testFile *file

func setupFileTest() error {
	err := test.SetConfig(fileConfig)
	if err != nil {
		return err
	}
	err = initStorage()
	if err != nil {
		return err
	}
	var ok bool
	testFile, ok = s.(*file)
	if !ok {
		return fmt.Errorf("expected storage to be of type *file")
	}
	return nil
}

func reloadTestFile() error {
	var err error
	s, err = NewFile()
	if err != nil {
		return err
	}
	var ok bool
	testFile, ok = s.(*file)
	if !ok {
		return fmt.Errorf("expected storage to be of type *file")
	}
	return nil
}

func TestNewFile(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	s, err := NewFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	_, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
}

func TestFileGetTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	id := uuid.Must(uuid.NewRandom())
	timerNew := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Stop:    time.Now().Add(time.Hour),
		Project: "test",
	}
	testFile.timers = types.Timers{timerNew}
	timerStore, err := s.GetTimer(id)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timerNew, timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileGetTimerNoExist(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Stop:    time.Now().Add(time.Hour),
		Project: "test",
	}}
	_, err = s.GetTimer(uuid.Must(uuid.Parse("00000000-0000-0000-0000-000000000000")))
	if err == nil {
		t.Fatal("expected error because timer does not exist")
	}
}

func TestFileGetRunningTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	id := uuid.Must(uuid.NewRandom())
	timerNew := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Project: "test",
	}
	// place two timers, the running timer and a completed one
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now().Add(-30 * time.Minute),
		Project: "test",
	}, timerNew}
	timerStore, err := s.GetRunningTimer()
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timerNew, timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileGetRunningTimerNoExist(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	// place two timers, the running timer and a completed one
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now().Add(-30 * time.Minute),
		Project: "test",
	}}
	_, err = s.GetRunningTimer()
	if err == nil {
		t.Fatal("expected error because of no running timer")
	}
}

func TestFileGetLastTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	id := uuid.Must(uuid.NewRandom())
	timerNewRunning := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Project: "test",
	}
	timerNewNotRunning := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now().Add(-30 * time.Minute),
		Project: "test",
	}
	// place two timers, the running timer and a completed one
	testFile.timers = types.Timers{timerNewNotRunning, timerNewRunning}
	timerStore, err := s.GetLastTimer(false)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timerNewNotRunning, timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
	timerStore, err = s.GetLastTimer(true)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = timersEqual(timerNewRunning, timerStore)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileGetTimers(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	timers := types.Timers{
		types.Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   time.Now().Add(-time.Hour),
			Stop:    time.Now().Add(-30 * time.Minute),
			Project: "test",
		},
		types.Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   time.Now().Add(-20 * time.Minute),
			Stop:    time.Now(),
			Project: "test",
		},
	}
	testFile.timers = timers
	timersStore, err := s.GetTimers(types.Filter{})
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(timers, timersStore) {
		t.Fatal("expected both slices to be identical")
	}
}

func TestFileRunningTimerExists(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	// place two timers, a running timer and a completed one
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now().Add(-30 * time.Minute),
		Project: "test",
	}}
	runningTimerExists, err := s.RunningTimerExists()
	if err != nil {
		t.Fatal(err.Error())
	}
	if runningTimerExists {
		t.Fatal("expected to find no running timer")
	}
	testFile.timers = append(testFile.timers, types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	})
	runningTimerExists, err = s.RunningTimerExists()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !runningTimerExists {
		t.Fatal("expected to find running timer")
	}
}

func TestFileUpdateTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	id := uuid.Must(uuid.NewRandom())
	timer := types.Timer{
		Uuid:    id,
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}
	testFile.timers = types.Timers{timer}
	timer.Stop = time.Now()
	err = s.UpdateTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected testFile to contain one timer but contains %d", len(testFile.timers))
	}
	err = timersEqual(timer, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileStoreTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	// reset timer storage
	testFile.timers = types.Timers{}
	id := uuid.Must(uuid.NewRandom())
	timerNew := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Stop:    time.Time{},
		Project: "test",
		Task:    "test",
		Tags:    []string{"test"},
	}
	err = s.StoreTimer(timerNew)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected exactly one timer but got %d", len(testFile.timers))
	}
	err = timersEqual(timerNew, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileStoreTimerNoDuplicate(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	id := uuid.Must(uuid.NewRandom())
	// reset timer storage
	testFile.timers = types.Timers{types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Project: "test",
		Task:    "test",
		Tags:    []string{"test"},
	}}
	err = s.StoreTimer(types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Project: "test",
	})
	if err == nil {
		t.Fatal("expected error because of duplicate uuids")
	}
}

func TestFileWritesUpdate(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	// this test ensures that the update that is passed to the file gets
	// written to disk and can be read again.
	timer := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}
	testFile.timers = types.Timers{timer}
	timer.Stop = time.Now()
	err = s.UpdateTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = reloadTestFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected exactly one timer but got %d", len(testFile.timers))
	}
	err = timersEqual(timer, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestFileWritesStore(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile.timers = nil
	timer := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "test",
		Task:    "test",
		Tags:    []string{"a", "b"},
	}
	err = s.StoreTimer(timer)
	if err != nil {
		t.Fatal(err.Error())
	}
	err = reloadTestFile()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected to find one timer but got %d", len(testFile.timers))
	}
	err = timersEqual(timer, testFile.timers[0])
	if err != nil {
		t.Fatal(err.Error())
	}
}
