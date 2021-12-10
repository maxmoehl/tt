package main

import (
	"fmt"
	"strings"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

const (
	flagTask      = "task"
	flagTimestamp = "timestamp"
	flagTags      = "tags"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start <project>",
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
	Run: func(cmd *cobra.Command, args []string) {
		start(getStartParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP(flagTask, string(flagTask[0]), "", "Set the task your are currently working on.")
	startCmd.Flags().String(flagTimestamp, "", "Manually set the start time for a timer. Format must be RFC3339")
	startCmd.Flags().String(flagTags, "", "Specify tags for this timer")
}

func start(silent bool, project, task, timestamp, tagsString string) {
	var tags []string
	if tagsString != "" {
		tags = strings.Split(tagsString, ",")
	}
	timer, err := tt.StartTimer(project, task, timestamp, tags)
	if err != nil {
		PrintError(err, silent)
	}
	if !silent {
		fmt.Printf("[%02d:%02d] Work tracking started!\n", timer.Start.Hour(), timer.Start.Minute())
		fmt.Printf("  project: %s\n", timer.Project)
		if task != "" {
			fmt.Printf("  task   : %s\n", timer.Task)
		}
		if len(tags) > 0 {
			fmt.Printf("  tags   : %s\n", strings.Join(timer.Tags, ", "))
		}
	}
}

func getStartParameters(cmd *cobra.Command, args []string) (silent bool, project, task, timestamp, tags string) {
	var err error
	silent = getSilent(cmd)
	task, err = cmd.LocalFlags().GetString(flagTask)
	if err != nil {
		PrintError(err, silent)
	}
	timestamp, err = cmd.LocalFlags().GetString(flagTimestamp)
	if err != nil {
		PrintError(err, silent)
	}
	tags, err = cmd.LocalFlags().GetString(flagTags)
	if err != nil {
		PrintError(err, silent)
	}
	if len(args) != 1 {
		PrintError(fmt.Errorf(WarningNArgumentsAcceptedFormat, len(args), 1), silent)
	}
	project = args[0]
	return
}
