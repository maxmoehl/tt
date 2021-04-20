package storage

import (
	"math"
	"time"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"
)

// GetTimeStatistics generates a types.Statistic struct for all timers
// that match the filter. The data is grouped as defined by byTask and
// byProject.
func GetTimeStatistics(byProject, byTask bool, filter types.Filter) (statistic types.Statistic, err error) {
	var timers types.Timers
	timers, err = s.GetTimers(filter)
	if err != nil {
		return
	}
	return getTimeStatisticsForTimers(timers, byProject, byTask)
}

// GetTimeStatisticsByDay generates a similar report to GetTimeStatistics
// but does the analysis on a daily basis.
func GetTimeStatisticsByDay(byProject, byTask bool, filter types.Filter) (map[string]types.Statistic, error) {
	timers, err := s.GetTimers(filter)
	if err != nil {
		return nil, err
	}

	b := timers.First().Start
	e := timers.Last(false).End

	currentDay := time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.Local)
	end := time.Date(e.Year(), e.Month(), e.Day()+1, 0, 0, 0, 0, time.Local)

	statistics := make(map[string]types.Statistic)

	for currentDay.Before(end) {
		f := types.NewFilter(nil, nil, nil, currentDay, currentDay.Add(time.Hour*24))
		statistic, err := getTimeStatisticsForTimers(timers.Filter(f), byProject, byTask)
		if err != nil {
			return nil, err
		}
		statistics[currentDay.Format(utils.DateFormat)] = statistic

		currentDay = currentDay.Add(time.Hour * 24)
	}
	return statistics, nil
}

func getTimeStatisticsForTimers(timers types.Timers, byProject, byTask bool) (statistic types.Statistic, err error) {
	statistic.Worked = workTime(timers)
	statistic.Breaks = breakTime(timers)
	statistic.Planned, err = plannedTime(timers)
	if err != nil {
		return types.Statistic{}, err
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

func getTimeByProjects(timers types.Timers, byTasks bool) (projects []types.Project) {
	projectsMap := getTimersByProjects(timers)
	for name, timers := range projectsMap {
		p := types.Project{
			Name:   name,
			Worked: workTime(timers),
			Breaks: breakTime(timers),
		}
		if byTasks {
			p.ByTasks = getTimeByTasks(timers)
		}
		projects = append(projects, p)
	}
	return
}

func getTimeByTasks(timers types.Timers) (tasks []types.Task) {
	tasksMap := getTimersByTasks(timers)
	for name, timers := range tasksMap {
		tasks = append(tasks, types.Task{
			Name:   name,
			Worked: workTime(timers),
			Breaks: breakTime(timers),
		})
	}
	return
}

func getTimersByProjects(timers types.Timers) map[string]types.Timers {
	res := make(map[string]types.Timers)
	for _, t := range timers {
		if _, ok := res[t.Project]; ok {
			res[t.Project] = append(res[t.Project], t)
		} else {
			res[t.Project] = types.Timers{t}
		}
	}
	return res
}

func getTimersByTasks(timers types.Timers) map[string]types.Timers {
	res := make(map[string]types.Timers)
	for _, t := range timers {
		if _, ok := res[t.Task]; ok {
			res[t.Task] = append(res[t.Task], t)
		} else {
			res[t.Task] = types.Timers{t}
		}
	}
	return res
}

func workTime(timers types.Timers) (d time.Duration) {
	for _, t := range timers {
		if t.Running() {
			continue
		}
		d += t.Duration()
	}
	return
}

func breakTime(timers types.Timers) (d time.Duration) {
	for _, t := range timers {
		d += t.Breaks.Duration()
	}
	return
}

func plannedTime(timers types.Timers) (time.Duration, error) {
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
	end := timers[0].End
	for _, t := range timers {
		if t.End.Unix() > end.Unix() {
			end = t.End
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
