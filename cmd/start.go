package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start <project> [<task>]",
	Aliases: []string{"begin"},
	Short:   "Starts tracking time",
	Long: `Starts tracking time.

With this command you can start time tracking and tag it with a project name
and an optional specific task. The project name can be any alphanumerical
identifier, including dashes and underscores. A project name is required,
specifying a task is optional. Tags are also optional and can be submitted as a
comma separated list of strings.

If you want to manually set a start time it should look something like this:
  2020-04-19 08:00

you can also omit the date, the current date will be used:
  08:00

or add seconds if that's your thing:
  09:32:42

You can also supply a full RFC3339 date-time string.

The two options --resume and --copy <timer> help to reduce typing by copying
values from previous timers, unless provided explicitly.Resume automatically
picks the last timer that was stopped. Copy needs an integer indicating how
many timers it should go back (1 being the same as resume). Copy ignores values
of zero and below. If you copy/resume the syntax of the command changes
slightly to:
  tt start [<task>] [flags]
  tt start [<project>] [<task>] [flags]

This is to enable you to set the task without having to redefine the project
because that is most likely the more frequent use case (compared to keeping the
task and only changing the project).

The cli will check if the given start time is valid, e.g. if the last timer
that ended, ended before the given start.`,
	Example: "tt start programming tt --tags private",
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, project, task, tags, timestamp, copyFrom, resume, err := getStartParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
		err = runStart(quiet, project, task, tags, timestamp, copyFrom, resume)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "manually set the start time for a timer")
	startCmd.Flags().String(flagTags, "", "specify tags for this timer")
	startCmd.Flags().IntP(flagCopy, string(flagCopy[0]), 0, "copy values from a specific timer")
	startCmd.Flags().BoolP(flagResume, string(flagResume[0]), false, "copy values from the previous timer")

	// TODO: --auto-stop (or something like this) to stop the previous timer automatically and start a new one
	//       how does this relate to the copy option?
}

func runStart(quiet bool, project, task string, tags []string, timestamp time.Time, copyFrom int, resume bool) error {
	if resume && copyFrom > 0 {
		return fmt.Errorf("start: %w: cannot have copy and resume", tt.ErrInvalidParameters)
	}
	if resume {
		copyFrom = 1
	}
	// if we are copying, and we only have a project, the order is reversed
	// so the project becomes the task.
	if copyFrom > 0 && task == "" {
		task = project
		project = ""
	}
	timer, err := tt.Start(project, task, tags, timestamp, copyFrom)
	if err != nil {
		return err
	}
	if !quiet {
		printTrackingStartedMsg(timer)
	}
	return nil
}

func getStartParameters(cmd *cobra.Command, args []string) (quiet bool, project, task string, tags []string, timestamp time.Time, copy int, resume bool, err error) {
	// TODO: could we also add a --interactive | -i mode that collects the information from stdin with autocomplete?
	flags, err := flags(cmd, flagQuiet, flagTags, flagTimestamp, flagCopy, flagResume)
	if err != nil {
		return
	}
	if len(args) > 0 {
		project = args[0]
	}
	if len(args) > 1 {
		task = args[1]
	}
	return flags[flagQuiet].(bool), project, task, flags[flagTags].([]string), flags[flagTimestamp].(time.Time), flags[flagCopy].(int), flags[flagResume].(bool), nil
}

func printTrackingStartedMsg(t tt.Timer) {
	fmt.Printf("[%02d:%02d] Tracking started!\n", t.Start.Hour(), t.Start.Minute())
	fmt.Printf("  project: %s\n", t.Project)
	if t.Task != "" {
		fmt.Printf("  task   : %s\n", t.Task)
	}
	if len(t.Tags) > 0 {
		fmt.Printf("  tags   : %s\n", strings.Join(t.Tags, ", "))
	}
}
