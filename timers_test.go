package tt

import (
	"testing"
	"time"
)

func TestTimerDuration(t *testing.T) {
	start := time.Now()
	stop := start.Add(time.Hour)
	timer := Timer{Start: start, Stop: &stop}
	if timer.Duration() != time.Hour {
		t.Fatal("expected duration to be one hour")
	}
}

func TestTimerRunning(t *testing.T) {
	timer := Timer{Start: time.Now()}
	if !timer.Running() {
		t.Fatal("expected timer to be running")
	}
}

func TestTimerNotRunning(t *testing.T) {
	start := time.Now()
	stop := start.Add(time.Minute)
	timer := Timer{Start: start, Stop: &stop}
	if timer.Running() {
		t.Fatal("expected timer to be stopped")
	}
}

func TestTimersNotRunning(t *testing.T) {
	stop := time.Now()
	start := stop.Add(-time.Hour)
	timers := Timers{
		{
			Start: start,
			Stop:  &stop,
		},
	}
	if timers.Running() != -1 {
		t.Fatal("did not expect running timer")
	}
}

func TestTimersRunning(t *testing.T) {
	stop := time.Now().Add(-30 * time.Minute)
	timers := Timers{
		{
			Start: time.Now().Add(-time.Hour),
			Stop:  &stop,
		},
		{
			Start: time.Now(),
		},
	}
	if timers.Running() == -1 {
		t.Fatal("expected running timer")
	}
}

func TestTimersFirst(t *testing.T) {
	stop := time.Now().Add(-30 * time.Minute)
	timer := Timers{
		{
			Start:   time.Now().Add(-time.Hour),
			Stop:    &stop,
			Project: "a",
		},
		{
			Start:   time.Now(),
			Project: "b",
		},
	}.First()
	if timer.Project != "a" {
		t.Fatal("expected to get timer with project a")
	}
}
