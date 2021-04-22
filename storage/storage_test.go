package storage

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt/types"
)

func TestStartTimer(t *testing.T) {
	testFile.timers = nil
	err := StartTimer("test", "", "", nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected file to contain one timer but got %d", len(testFile.timers))
	}
	timer := testFile.timers[0]
	if timer.Project != "test" || timer.Task != "" || len(timer.Tags) != 0 ||
		!timer.Start.Before(time.Now()) || !timer.Start.After(time.Now().Add(-time.Minute)) {
		t.Fatal("timer does not match expectations")
	}
}

func TestStartTimerRunningTimer(t *testing.T) {
	testFile.timers = types.Timers{types.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	err := StartTimer("test", "", "", nil)
	if err == nil {
		t.Fatal("expected error since a timer is already running")
	}
}

func TestStartTimerValidTimestamp(t *testing.T) {
	testFile.timers = nil
	// we have to round since .Format does not add fractions of seconds
	startTime := time.Now().Round(time.Second)
	err := StartTimer("test", "", startTime.Format(time.RFC3339), nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected file to contain one timer but got %d", len(testFile.timers))
	}
	timer := testFile.timers[0]
	if timer.Project != "test" || timer.Task != "" || len(timer.Tags) != 0 ||
		!timer.Start.Equal(startTime) {
		fmt.Printf("%s\n", timer.Start.Format(time.RFC3339))
		fmt.Printf("%s\n", startTime.Format(time.RFC3339))
		t.Fatal("timer does not match expectations")
	}
}

func TestStartTimerInvalidTimestamp(t *testing.T) {
	testFile.timers = nil
	err := StartTimer("test", "", "invalid timestamp", nil)
	if err == nil {
		t.Fatal("expected error because of invalid timestamp")
	}
}

