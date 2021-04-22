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

	"github.com/maxmoehl/tt/config"
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

// Print prints the Statistic struct to the console indenting everything
// by the given indent and lower levels by multiples of the given indent.
func (s Statistic) Print(indent string) {
	precision := config.Get().GetPrecision()
	f := utils.FormatDuration
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
	fmt.Printf("%s  %s: %s\n", indent, p.Name, utils.FormatDuration(p.Worked, precision))
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
	fmt.Printf("%s    %s: %s\n", indent, name, utils.FormatDuration(t.Worked, precision))
}
