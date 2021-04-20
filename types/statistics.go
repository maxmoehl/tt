package types

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt/utils"
)

// Statistic holds all values that are generated as part of the stats
// command.
type Statistic struct {
	Worked     time.Duration `json:"worked"`
	Planned    time.Duration `json:"planned"`
	Breaks     time.Duration `json:"breaks"`
	Difference time.Duration `json:"difference"`
	Percentage float64       `json:"percentage"`
	ByProjects []Project     `json:"by_projects,omitempty"`
}

// Print prints the Statistic struct to the console indenting lower
// levels with two spaceds.
func (s Statistic) Print() {
	f := utils.FormatDuration
	fmt.Printf(
		"worked    : %s\n"+
			"planned   : %s\n"+
			"breaks    : %s\n"+
			"difference: %s\n"+
			"percentage: %.2f%%\n", f(s.Worked), f(s.Planned), f(s.Breaks), f(s.Difference), s.Percentage*100)
	if s.ByProjects != nil {
		fmt.Println("by projects:")
		for _, p := range s.ByProjects {
			fmt.Printf(
				"  %s:\n"+
					"    worked: %s\n"+
					"    breaks: %s\n", p.Name, f(p.Worked), f(p.Breaks))
			if p.ByTasks != nil {
				fmt.Println("    by tasks:")
				for _, t := range p.ByTasks {
					fmt.Printf(
						"      %s:\n"+
							"        worked: %s\n"+
							"        breaks: %s\n", t.Name, f(t.Worked), f(t.Breaks))
				}
			}
		}
	}
}

// Project is a part of Statistic and contains data for the stat command
type Project struct {
	Name    string        `json:"name"`
	Worked  time.Duration `json:"worked"`
	Breaks  time.Duration `json:"breaks"`
	ByTasks []Task        `json:"by_tasks,omitempty"`
}

// Task is a part of Statistic and contains data for the stat command
type Task struct {
	Name   string        `json:"name"`
	Worked time.Duration `json:"worked"`
	Breaks time.Duration `json:"breaks"`
}
