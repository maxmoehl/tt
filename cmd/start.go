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

With this command you can start time tracking and tag it with a project
name and an optional specific task. The project name can be any
alphanumerical identifier, including dashes and underscores. A project
name is required, specifying a task is optional. Tags are also optional
and can be submitted as a comma separated list of strings.

If you want to manually set a start time it should look something like
this:
  2020-04-19 08:00
you can also omit the date, the current date will be used:
  08:00
or add seconds if that's your thing:
  09:32:42
You can also supply a full RFC3339 date-time string.

The cli will check if the given start time is valid, e.g. if the last timer
that ended, ended before the given start.`,
	Example: "tt start programming tt --tags private",
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, project, task, tags, timestamp, err := getStartParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
		err = runStart(quiet, project, task, tags, timestamp)
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
}

func runStart(quiet bool, project, task string, tags []string, timestamp time.Time) error {
	timer, err := tt.Start(project, task, tags, timestamp)
	if err != nil {
		return err
	}
	if !quiet {
		printTrackingStartedMsg(timer)
	}
	return nil
}

func getStartParameters(cmd *cobra.Command, args []string) (quiet bool, project, task string, tags []string, timestamp time.Time, err error) {
	flags, err := flags(cmd, flagQuiet, flagTags, flagTimestamp)
	if err != nil {
		return
	}
	if len(args) < 1 || len(args) > 2 {
		err = fmt.Errorf("expected 1-2 arguments")
		return
	}
	project = args[0]
	if len(args) > 1 {
		task = args[1]
	}
	return flags[flagQuiet].(bool), project, task, flags[flagTags].([]string), flags[flagTimestamp].(time.Time), nil
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
