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

package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

const (
	flagTask      = "task"
	flagTimestamp = "timestamp"
	flagTags      = "tags"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start project",
	Short: "Starts tracking time.",
	Long: `Starts tracking time.

With this command you can start time tracking and tag it with a project
name and an optional specific task. The project name can be any
alphanumerical identifier, including dashes and underscores. A project
name is required, specifying a task is optional. Tags are also optional
and can be submitted as a comma separated list of strings.

If you want to manually set a start time it has to be in the following
format:
  2020-04-19T08:00:00+02:00
Otherwise an appropriate error will be printed. The cli will check if the
given start time is valid, e.g. if the last timer that ended, ended before
the given start.`,
	Example: "tt start programming -t tt --tags private",
	Run:     start,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP(flagTask, string(flagTask[0]), "", "Set the task your are currently working on.")
	startCmd.Flags().String(flagTimestamp, "", "Manually set the start time for a timer. Format must be RFC3339")
	startCmd.Flags().String(flagTags, "", "Specify tags for this timer")
}

func start(cmd *cobra.Command, args []string) {
	silent := getSilent(cmd)
	if len(args) != 1 {
		utils.PrintError(fmt.Errorf("this command needs exactly one argument"), silent)
	}
	project := args[0]
	silent, task, timestamp, tagsString := getStartFlags(cmd)
	var tags []string
	if tagsString != "" {
		tags = strings.Split(tagsString, ",")
	}
	err := storage.StartTimer(project, task, timestamp, tags)
	if err != nil {
		utils.PrintError(err, silent)
	}
	if !silent {
		fmt.Printf("[%d:%d] Work tracking started!\n", time.Now().Hour(), time.Now().Minute())
		fmt.Printf("  project: %s\n", project)
		if task != "" {
			fmt.Printf("  task   : %s\n", task)
		}
		if len(tags) > 0 {
			fmt.Printf("  tags   : %s\n", strings.Join(tags, ", "))
		}
	}
}

func getStartFlags(cmd *cobra.Command) (silent bool, task, timestamp, tags string) {
	var err error
	silent = getSilent(cmd)
	task, err = cmd.LocalFlags().GetString(flagTask)
	if err != nil {
		utils.PrintError(err, silent)
	}
	timestamp, err = cmd.LocalFlags().GetString(flagTimestamp)
	if err != nil {
		utils.PrintError(err, silent)
	}
	tags, err = cmd.LocalFlags().GetString(flagTags)
	if err != nil {
		utils.PrintError(err, silent)
	}
	return
}
