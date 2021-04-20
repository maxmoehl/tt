package types

import (
	"fmt"
	"time"

	"github.com/maxmoehl/tt/utils"
)

type Statistic struct {
	Worked     time.Duration `json:"worked"`
	Planned    time.Duration `json:"planned"`
	Breaks     time.Duration `json:"breaks"`
	Difference time.Duration `json:"difference"`
	Percentage float64       `json:"percentage"`
	ByProjects []Project     `json:"by_projects,omitempty"`
}

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

type Project struct {
	Name    string        `json:"name"`
	Worked  time.Duration `json:"worked"`
	Breaks  time.Duration `json:"breaks"`
	ByTasks []Task        `json:"by_tasks,omitempty"`
}

type Task struct {
	Name   string        `json:"name"`
	Worked time.Duration `json:"worked"`
	Breaks time.Duration `json:"breaks"`
}
