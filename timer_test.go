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
