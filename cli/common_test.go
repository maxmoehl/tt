package main

import (
	"os"
	"testing"
	"time"

	"github.com/maxmoehl/tt"
)

func ExamplePrintWarning() {
	// redirect stderr to stdout
	stderr := os.Stderr
	os.Stderr = os.Stdout
	defer func() { os.Stderr = stderr }()

	PrintWarning("test")

	// Output:
	// Warning: test
}

func TestFormatDuration(t *testing.T) {
	type args struct {
		d         time.Duration
		precision time.Duration
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"precision hour",
			args{
				time.Hour + 40*time.Minute,
				time.Hour,
			},
			"1h",
		},
		{
			"precision minute",
			args{
				time.Hour + 20*time.Minute,
				time.Minute,
			},
			"1h20m",
		},
		{
			"precision second",
			args{
				time.Hour + 20*time.Minute,
				time.Second,
			},
			"1h20m0s",
		},
		{
			"precision second, negative",
			args{
				-(time.Hour + 20*time.Minute),
				time.Second,
			},
			"-1h20m0s",
		},
		{
			"unknown precision",
			args{
				time.Hour + 20*time.Minute,
				time.Microsecond,
			},
			"unknown precision",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := tt.FormatDuration(test.args.d, test.args.precision); got != test.want {
				t.Errorf("FormatDuration() = %v, want %v", got, test.want)
			}
		})
	}
}
