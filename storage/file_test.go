package storage

import (
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"
)

var testFile *file

func TestMain(m *testing.M) {
	dir := setup()
	defer teardown(dir)

	err := reloadTestFile()
	if err != nil {
		panic(err.Error())
	}

	exitCode := m.Run()
	if exitCode != 0 {
		// os.Exit does not run deferred functions, therefore we run it manually
		teardown(dir)
		os.Exit(exitCode)
	}
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
	id := uuid.Must(uuid.NewRandom())
	timerNew := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		End:     time.Now().Add(time.Hour),
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

func TestFileGetRunningTimer(t *testing.T) {
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
		End:     time.Now().Add(-30 * time.Minute),
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

func TestFileGetLastTimer(t *testing.T) {
	id := uuid.Must(uuid.NewRandom())
	timerNewRunning := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		Project: "test",
	}
	timerNewNotRunning := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		End:     time.Now().Add(-30 * time.Minute),
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
	timers := types.Timers{
		types.Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   time.Now().Add(-time.Hour),
			End:     time.Now().Add(-30 * time.Minute),
			Project: "test",
		},
		types.Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   time.Now().Add(-20 * time.Minute),
			End:     time.Now(),
			Project: "test",
		},
	}
	testFile.timers = timers
	timersStore, err := s.GetTimers(nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if !reflect.DeepEqual(timers, timersStore) {
		t.Fatal("expected both slices to be identical")
	}
}

func TestFileRunningTimerExists(t *testing.T) {
	// place two timers, a running timer and a completed one
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		End:     time.Now().Add(-30 * time.Minute),
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
	id := uuid.Must(uuid.NewRandom())
	timer := types.Timer{
		Uuid:    id,
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}
	testFile.timers = types.Timers{timer}
	timer.End = time.Now()
	err := s.UpdateTimer(timer)
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
	// reset timer storage
	testFile.timers = types.Timers{}
	id := uuid.Must(uuid.NewRandom())
	timerNew := types.Timer{
		Uuid:    id,
		Start:   time.Now(),
		End:     time.Time{},
		Project: "test",
		Task:    "test",
		Tags:    []string{"test"},
	}
	err := s.StoreTimer(timerNew)
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

func TestFileWritesUpdate(t *testing.T) {
	// this test ensures that the update that is passed to the file gets
	// written to disk and can be read again.
	timer := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}
	testFile.timers = types.Timers{timer}
	timer.End = time.Now()
	err := s.UpdateTimer(timer)
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
	testFile.timers = nil
	timer := types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		End:     time.Now(),
		Project: "test",
		Task:    "test",
		Tags:    []string{"a", "b"},
	}
	err := s.StoreTimer(timer)
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
	if !t1.Start.Equal(t2.Start) {
		return fmt.Errorf("start time is not equal")
	}
	if !t1.End.Equal(t2.End) {
		return fmt.Errorf("end time is not equal")
	}
	return nil
}
