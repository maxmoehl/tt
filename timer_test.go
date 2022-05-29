package tt

import (
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestTimer_Validate(t1 *testing.T) {
	tests := []struct {
		name    string
		timer   Timer
		wantErr bool
	}{
		{
			"simple valid timer",
			Timer{
				ID:      uuid.Must(uuid.NewRandom()).String(),
				Start:   time.Now(),
				Stop:    nil,
				Project: "foo",
				Task:    "",
				Tags:    nil,
			},
			false,
		},
		{
			"zero start time",
			Timer{
				ID:      uuid.Must(uuid.NewRandom()).String(),
				Start:   time.Time{},
				Stop:    nil,
				Project: "foo",
				Task:    "",
				Tags:    nil,
			},
			true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			if err := tt.timer.Validate(); (err != nil) != tt.wantErr {
				t1.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
