package storage

import (
	"fmt"
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
	_, ok = s.(*file)
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
	_, ok = s.(*file)
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

func TestFileWritesUpdate(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
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
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
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
