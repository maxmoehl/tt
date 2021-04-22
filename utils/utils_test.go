package utils

import (
	"os"
	"testing"
	"time"
)

func ExamplePrintWarning() {
	// redirect stderr to stdout
	stderr := os.Stderr
	os.Stderr = os.Stdout
	defer func() {os.Stderr = stderr}()

	PrintWarning("test")

	// Output:
	// Warning: test
}

func TestStringSliceContains(t *testing.T) {
	type args struct {
		strings []string
		s       string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"string is inside slice",
			args{
				[]string{"test", "a", "b", "c"},
				"a",
			},
			true,
		},
		{
			"string is not inside slice",
			args{
				[]string{"test", "a", "b", "c"},
				"d",
			},
			false,
		},
		{
			"nil slice",
			args{
				nil,
				"a",
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringSliceContains(tt.args.strings, tt.args.s); got != tt.want {
				t.Errorf("StringSliceContains() = %v, want %v", got, tt.want)
			}
		})
	}
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
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatDuration(tt.args.d, tt.args.precision); got != tt.want {
				t.Errorf("FormatDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}
