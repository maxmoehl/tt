package tt

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestStartTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	_, err = StartTimer("test", "", "", nil)
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
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	_, err = StartTimer("test", "", "", nil)
	if err == nil {
		t.Fatal("expected error since a timer is already running")
	}
}

func TestStartTimerValidTimestamp(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	// we have to round since .Format does not add fractions of seconds
	startTime := time.Now().Round(time.Second)
	_, err = StartTimer("test", "", startTime.Format(time.RFC3339), nil)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatalf("expected file to contain one timer but got %d", len(testFile.timers))
	}
	timer := testFile.timers[0]
	if timer.Project != "test" || timer.Task != "" || len(timer.Tags) != 0 ||
		!timer.Start.Equal(startTime) {
		t.Fatal("timer does not match expectations")
	}
}

func TestStartTimerInvalidTimestamp(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	_, err = StartTimer("test", "", "invalid timestamp", nil)
	if err == nil {
		t.Fatal("expected error because of invalid timestamp")
	}
}

func TestStartTimerTimestampCollision(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "test",
	}}
	_, err = StartTimer("test", "", time.Now().Add(-30*time.Minute).Format(time.RFC3339), nil)
	if err == nil {
		t.Fatal("expected error because of collision with existing timer")
	}
}

func TestStopTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	_, err = StopTimer("")
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestStopTimerNoRunningTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	_, err = StopTimer("")
	if err == nil {
		t.Fatal("expected error because of no running timer")
	}
}

func TestStopTimerValidTimestamp(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}}
	stopTime := time.Now().Round(time.Second)
	_, err = StopTimer(stopTime.Format(time.RFC3339))
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 1 {
		t.Fatal("expected exactly one timer")
	}
	timer := testFile.timers[0]
	if timer.Project != "test" || timer.Task != "" || len(timer.Tags) != 0 ||
		!timer.Stop.Equal(stopTime) {
		t.Fatal("timer does not match expectations")
	}
}

func TestStopTimerInvalidTimestamp(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Project: "test",
	}}
	_, err = StopTimer("invalid timestamp")
	if err == nil {
		t.Fatal("expected error because of invalid timestamp")
	}
}

func TestStopTimerTimestampBeforeStart(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	_, err = StopTimer(time.Now().Add(-time.Hour).Format(time.RFC3339))
	if err == nil {
		t.Fatal("expected error because of stop time before start time")
	}
}

func TestResumeTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now().Add(-time.Hour),
		Stop:    time.Now(),
		Project: "testA",
		Task:    "testB",
		Tags:    []string{"a", "b"},
	}}
	timerResume, err := ResumeTimer()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(testFile.timers) != 2 {
		t.Fatalf("expected two timers but got %d", len(testFile.timers))
	}
	timerFile := testFile.timers[1]
	if !reflect.DeepEqual(timerResume, timerFile) {
		t.Fatal("expected returned timer and timer in storage to be equal")
	}
	if timerResume.Project != "testA" ||
		timerResume.Task != "testB" ||
		len(timerResume.Tags) != 2 ||
		!stringSliceContains(timerResume.Tags, "a") ||
		!stringSliceContains(timerResume.Tags, "b") {
		t.Fatal("expected properties to be copied but found different values")
	}
}

func TestResumeTimerNoTimerToResume(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	_, err = ResumeTimer()
	if err == nil {
		t.Fatal("expected error because there is no timer to resume")
	}
}

func TestResumeTimerRunningTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	_, err = ResumeTimer()
	if err == nil {
		t.Fatal("expected error because there is a running timer")
	}
}

func TestCheckRunningTimers(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	runningId := uuid.Must(uuid.NewRandom())
	testFile.timers = Timers{
		Timer{
			Uuid:    uuid.Must(uuid.NewRandom()),
			Start:   time.Now().Add(-time.Hour),
			Stop:    time.Now().Add(-45 * time.Minute),
			Project: "test",
		},
		Timer{
			Uuid:    runningId,
			Start:   time.Now().Add(-30 * time.Minute),
			Project: "test",
		},
	}
	ids, err := CheckRunningTimers()
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(ids) != 1 {
		t.Fatalf("expected to find one running timer but got %d", len(ids))
	}
	if ids[0] != runningId {
		t.Fatal("returned id is not id of running timer")
	}
}

func TestGetRunningTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = Timers{Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   time.Now(),
		Project: "test",
	}}
	timer, err := GetRunningTimer()
	if err != nil {
		t.Fatal(err.Error())
	}
	if !timer.Running() {
		t.Fatal("expected to get running timer")
	}
	if timer.Project != "test" ||
		timer.Task != "" ||
		timer.Tags != nil {
		t.Fatal("timer does not match expectations")
	}
}

func TestGetRunningTimerNoTimer(t *testing.T) {
	err := setupFileTest()
	if err != nil {
		t.Fatal(err.Error())
	}
	testFile, ok := s.(*file)
	if !ok {
		t.Fatal("expected storage to be of type *file")
	}
	testFile.timers = nil
	_, err = GetRunningTimer()
	if err == nil {
		t.Fatal("expected error because of no running timer")
	}
}
