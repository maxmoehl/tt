package types

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt/test"
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
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
				Name: "testA",
				Worked: time.Hour,
			},
			{
				Name: "testB",
				Worked: 2*time.Hour,
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "m"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "s"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "m"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "h"))
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
	err := test.SetConfig(fmt.Sprintf(statisticsConfigTemplate, "s"))
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
