/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tt

import (
	"fmt"
	"math"
	"time"

	"github.com/maxmoehl/tt/config"
)

// DateFormat contains the format in which dates are printed.
var DateFormat = "2006-01-02"

// Statistic holds all values that are generated as part of the stats
// command.
type Statistic struct {
	Worked     time.Duration `json:"worked"`
	Planned    time.Duration `json:"planned"`
	Difference time.Duration `json:"difference"`
	Percentage float64       `json:"percentage"`
	ByProjects []Project     `json:"by_projects,omitempty"`
}

// Print prints the Statistic struct to the console indenting everything
// by the given indent and lower levels by multiples of the given indent.
func (s Statistic) Print(indent string) {
	precision := config.Get().GetPrecision()
	f := FormatDuration
	fmt.Printf("%sworked    : %s\n", indent, f(s.Worked, precision))
	fmt.Printf("%splanned   : %s\n", indent, f(s.Planned, precision))
	fmt.Printf("%sdifference: %s\n", indent, f(s.Difference, precision))
	fmt.Printf("%spercentage: %.2f%%\n", indent, s.Percentage*100)
	if len(s.ByProjects) > 0 {
		fmt.Printf("%sby projects:\n", indent)
		for _, p := range s.ByProjects {
			p.Print(indent)
		}
	}
}

// Project is a part of Statistic and contains data for the stat command
type Project struct {
	Name    string        `json:"name"`
	Worked  time.Duration `json:"worked"`
	ByTasks []Task        `json:"by_tasks,omitempty"`
}

// Print prints an indented representation of Project and Tasks if ByTasks
// contains at least one element.
func (p Project) Print(indent string) {
	precision := config.Get().GetPrecision()
	fmt.Printf("%s  %s: %s\n", indent, p.Name, FormatDuration(p.Worked, precision))
	if len(p.ByTasks) > 0 {
		fmt.Printf("%s  by tasks:\n", indent)
		for _, t := range p.ByTasks {
			t.Print(indent)
		}
	}
}

// Task is a part of Statistic and contains data for the stat command
type Task struct {
	Name   string        `json:"name"`
	Worked time.Duration `json:"worked"`
}

// Print prints an indented representation of Task
func (t Task) Print(indent string) {
	precision := config.Get().GetPrecision()
	name := t.Name
	if name == "" {
		name = "without task"
	}
	fmt.Printf("%s    %s: %s\n", indent, name, FormatDuration(t.Worked, precision))
}

// FormatDuration formats a duration in the precision defined by the
// config.
func FormatDuration(d time.Duration, precision time.Duration) string {
	h := d / time.Hour
	m := (d - (h * time.Hour)) / time.Minute
	s := (d - (h * time.Hour) - (m * time.Minute)) / time.Second
	sign := ""
	if d < 0 {
		sign = "-"
		h *= -1
		m *= -1
		s *= -1
	}
	switch precision {
	case time.Second:
		return fmt.Sprintf("%s%dh%dm%ds", sign, h, m, s)
	case time.Minute:
		return fmt.Sprintf("%s%dh%dm", sign, h, m)
	case time.Hour:
		return fmt.Sprintf("%s%dh", sign, h)
	default:
		return "unknown precision"
	}
}

// GetTimeStatistics generates a types.Statistic struct for all timers
// that match the filter. The data is grouped as defined by byTask and
// byProject.
func GetTimeStatistics(byProject, byTask bool, filter Filter) (statistic Statistic, err error) {
	var timers Timers
	timers, err = s.GetTimers(filter)
	if err != nil {
		return
	}
	return getTimeStatisticsForTimers(timers, byProject, byTask)
}

// GetTimeStatisticsByDay generates a similar report to GetTimeStatistics
// but does the analysis on a daily basis.
func GetTimeStatisticsByDay(byProject, byTask bool, filter Filter) (map[string]Statistic, error) {
	timers, err := s.GetTimers(filter)
	if err != nil {
		return nil, err
	}

	b := timers.First().Start
	e := timers.Last(false).Stop

	currentDay := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(e.Year(), e.Month(), e.Day()+1, 0, 0, 0, 0, time.Local)

	statistics := make(map[string]Statistic)

	for currentDay.Before(end) {
		f := NewFilter(nil, nil, nil, currentDay, currentDay.Add(time.Hour*24))
		statistic, err := getTimeStatisticsForTimers(f.Timers(timers), byProject, byTask)
		if err != nil {
			return nil, err
		}
		statistics[currentDay.Format(DateFormat)] = statistic

		currentDay = currentDay.Add(time.Hour * 24)
	}
	return statistics, nil
}

func getTimeStatisticsForTimers(timers Timers, byProject, byTask bool) (statistic Statistic, err error) {
	statistic.Worked = workTime(timers)
	statistic.Planned, err = plannedTime(timers)
	if err != nil {
		return Statistic{}, err
	}
	statistic.Difference = statistic.Worked - statistic.Planned
	statistic.Percentage = float64(statistic.Worked) / float64(statistic.Planned)
	if math.IsNaN(statistic.Percentage) || math.IsInf(statistic.Percentage, 0) {
		statistic.Percentage = 0
	}
	if byProject {
		statistic.ByProjects = getTimeByProjects(timers, byTask)
	}
	return
}

func getTimeByProjects(timers Timers, byTasks bool) (projects []Project) {
	projectsMap := getTimersByProjects(timers)
	for name, timers := range projectsMap {
		p := Project{
			Name:   name,
			Worked: workTime(timers),
		}
		if byTasks {
			p.ByTasks = getTimeByTasks(timers)
		}
		projects = append(projects, p)
	}
	return
}

func getTimeByTasks(timers Timers) (tasks []Task) {
	tasksMap := getTimersByTasks(timers)
	for name, timers := range tasksMap {
		tasks = append(tasks, Task{
			Name:   name,
			Worked: workTime(timers),
		})
	}
	return
}

func getTimersByProjects(timers Timers) map[string]Timers {
	res := make(map[string]Timers)
	for _, t := range timers {
		if _, ok := res[t.Project]; ok {
			res[t.Project] = append(res[t.Project], t)
		} else {
			res[t.Project] = Timers{t}
		}
	}
	return res
}

func getTimersByTasks(timers Timers) map[string]Timers {
	res := make(map[string]Timers)
	for _, t := range timers {
		if _, ok := res[t.Task]; ok {
			res[t.Task] = append(res[t.Task], t)
		} else {
			res[t.Task] = Timers{t}
		}
	}
	return res
}

func workTime(timers Timers) (d time.Duration) {
	for _, t := range timers {
		if t.Running() {
			continue
		}
		d += t.Duration()
	}
	return
}

func plannedTime(timers Timers) (time.Duration, error) {
	if len(timers) == 0 {
		return 0, nil
	}
	// get oldest timer
	start := timers[0].Start
	for _, t := range timers {
		if t.Start.Unix() < start.Unix() {
			start = t.Start
		}
	}

	// get most recent timer
	end := timers[0].Stop
	for _, t := range timers {
		if t.Stop.Unix() > end.Unix() {
			end = t.Stop
		}
	}

	// we want to get all dates before tomorrow
	var d time.Duration
	end = time.Date(end.Year(), end.Month(), end.Day()+1, 0, 0, 0, 0, time.Local)
	c := config.Get()

	// once for every day add the hours that would have been worked on that day
	for currentTime := start; currentTime.Before(end); currentTime = currentTime.Add(time.Hour * 24) {
		switch currentTime.Weekday() {
		case time.Monday:
			if c.WorkDays.Monday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Tuesday:
			if c.WorkDays.Tuesday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Wednesday:
			if c.WorkDays.Wednesday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Thursday:
			if c.WorkDays.Thursday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Friday:
			if c.WorkDays.Friday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Saturday:
			if c.WorkDays.Saturday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		case time.Sunday:
			if c.WorkDays.Sunday {
				d += time.Duration(c.WorkHours) * time.Hour
			}
		}
	}
	return d, nil
}
