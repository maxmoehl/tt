package tt

import (
	"testing"
	"time"
)

func TestTimerDuration(t *testing.T) {
	now := time.Now()
	timer := Timer{Start: now, Stop: now.Add(time.Hour)}
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
	timer := Timer{Start: time.Now(), Stop: time.Now().Add(time.Minute)}
	if timer.Running() {
		t.Fatal("expected timer to be stopped")
	}
}

func TestTimerIsZero(t *testing.T) {
	if !(Timer{}).IsZero() {
		t.Fatal("expected initial value of timer to be zero")
	}
}

func TestTimerNotZero(t *testing.T) {
	timer := Timer{Start: time.Now()}
	if timer.IsZero() {
		t.Fatal("expected timer to not be zero")
	}
}

func TestTimersNotRunning(t *testing.T) {
	timers := Timers{
		{
			Start: time.Now().Add(-time.Hour),
			Stop:  time.Now(),
		},
	}
	if running, _ := timers.Running(); running {
		t.Fatal("did not expect running timer")
	}
}

func TestTimersRunning(t *testing.T) {
	timers := Timers{
		{
			Start: time.Now().Add(-time.Hour),
			Stop:  time.Now().Add(-30 * time.Minute),
		},
		{
			Start: time.Now(),
		},
	}
	if running, _ := timers.Running(); !running {
		t.Fatal("expected running timer")
	}
}

func TestTimersFirst(t *testing.T) {
	timer := Timers{
		{
			Start:   time.Now().Add(-time.Hour),
			Stop:    time.Now().Add(-30 * time.Minute),
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
