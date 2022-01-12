package time_clock

import (
	"fmt"
	"testing"
	"time"
)

const statisticsConfigTemplate = `
precision: %s
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

func ExampleStatisticPrintNoProjects() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
	if err != nil {
		panic(err.Error())
	}
	Statistic{
		Worked:     3 * time.Hour,
		Planned:    8 * time.Hour,
		Difference: 5 * time.Hour,
		Percentage: float64(3) / 8,
	}.Print("")
	// Output:
	// worked    : 3h
	// planned   : 8h
	// difference: 5h
	// percentage: 37.50%
}

func ExampleStatisticPrintWithProjects() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
	if err != nil {
		panic(err.Error())
	}
	Statistic{
		Worked:     3 * time.Hour,
		Planned:    8 * time.Hour,
		Difference: 5 * time.Hour,
		Percentage: float64(3) / 8,
		ByProjects: []Project{
			{
				Name:   "testA",
				Worked: time.Hour,
			},
			{
				Name:   "testB",
				Worked: 2 * time.Hour,
			},
		},
	}.Print("")
	// Output:
	// worked    : 3h
	// planned   : 8h
	// difference: 5h
	// percentage: 37.50%
	// by projects:
	//   testA: 1h
	//   testB: 2h
}

func ExampleProjectPrintNoTasks() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "m"))
	if err != nil {
		panic(err.Error())
	}
	d, _ := time.ParseDuration("1h40m")
	Project{
		Name:    "test",
		Worked:  d,
		ByTasks: nil,
	}.Print("")
	// Output:
	// test: 1h40m
}

func ExampleProjectPrintWithTasks() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
	if err != nil {
		panic(err.Error())
	}
	d, _ := time.ParseDuration("1h40m")
	Project{
		Name:   "test",
		Worked: d,
		ByTasks: []Task{
			{
				Name:   "testA",
				Worked: d,
			},
			{
				Name:   "testB",
				Worked: d + time.Hour,
			},
		},
	}.Print("")
	// Output:
	//   test: 1h
	//   by tasks:
	//     testA: 1h
	//     testB: 2h
}

func ExampleTaskPrintSecond() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "s"))
	if err != nil {
		panic(err.Error())
	}
	Task{
		Name:   "test",
		Worked: time.Hour,
	}.Print("")

	// Output:
	// test: 1h0m0s
}

func ExampleTaskPrintMinute() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "m"))
	if err != nil {
		panic(err.Error())
	}
	Task{
		Name:   "test",
		Worked: time.Hour,
	}.Print("")

	// Output:
	// test: 1h0m
}

func ExampleTaskPrintHour() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
	if err != nil {
		panic(err.Error())
	}
	Task{
		Name:   "test",
		Worked: time.Hour,
	}.Print("")

	// Output:
	// test: 1h
}

func ExampleTaskPrintNoTask() {
	err := SetConfig(fmt.Sprintf(statisticsConfigTemplate, "s"))
	if err != nil {
		panic(err.Error())
	}
	Task{
		Name:   "",
		Worked: time.Hour,
	}.Print("")

	// Output:
	// without task: 1h0m0s
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
