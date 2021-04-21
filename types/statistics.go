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
	Difference time.Duration `json:"difference"`
	Percentage float64       `json:"percentage"`
	ByProjects []Project     `json:"by_projects,omitempty"`
}

// Print prints the Statistic struct to the console indenting lower
// levels with two spaces.
func (s Statistic) Print() {
	f := utils.FormatDuration
	fmt.Printf(
		"worked    : %s\n"+
			"planned   : %s\n"+
			"difference: %s\n"+
			"percentage: %.2f%%\n", f(s.Worked), f(s.Planned), f(s.Difference), s.Percentage*100)
	if s.ByProjects != nil {
		fmt.Println("by projects:")
		for _, p := range s.ByProjects {
			fmt.Printf(
				"  %s: %s\n", p.Name, f(p.Worked))
			if p.ByTasks != nil {
				fmt.Println("    by tasks:")
				for _, t := range p.ByTasks {
					fmt.Printf(
						"      %s: %s\n", t.Name, f(t.Worked))
				}
			}
		}
	}
}

// Project is a part of Statistic and contains data for the stat command
type Project struct {
	Name    string        `json:"name"`
	Worked  time.Duration `json:"worked"`
	ByTasks []Task        `json:"by_tasks,omitempty"`
}

// Task is a part of Statistic and contains data for the stat command
type Task struct {
	Name   string        `json:"name"`
	Worked time.Duration `json:"worked"`
}
