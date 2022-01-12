package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

const (
	flagTimestamp = "timestamp"
	flagTags      = "tags"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:     "start <project> [<task>]",
	Aliases: []string{"begin"},
	Short:   "Starts tracking time.",
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

The cli will check if the
given start time is valid, e.g. if the last timer that ended, ended before
the given start.`,
	Example: "tt start programming tt --tags private",
	Run: func(cmd *cobra.Command, args []string) {
		runStart(getStartParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP(flagTimestamp, string(flagTimestamp[0]), "", "Manually set the start time for a timer")
	startCmd.Flags().String(flagTags, "", "Specify tags for this timer")
}

func runStart(quiet bool, project, task string, tags []string, start time.Time) {
	newTimer := tt.Timer{
		Uuid:    uuid.Must(uuid.NewRandom()),
		Start:   start,
		Project: project,
		Task:    task,
		Tags:    tags,
	}
	err := tt.GetStorage().StoreTimer(newTimer)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if !quiet {
		printTrackingStartedMsg(newTimer)
	}
}

func getStartParameters(cmd *cobra.Command, args []string) (quiet bool, project, task string, tags []string, timestamp time.Time) {
	var err error
	quiet = getQuiet(cmd)
	rawTimestamp, err := cmd.LocalFlags().GetString(flagTimestamp)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if rawTimestamp != "" {
		timestamp, err = tt.ParseDate(rawTimestamp)
		if err != nil {
			tt.PrintError(err, quiet)
		}
	} else {
		timestamp = time.Now()
	}
	rawTags, err := cmd.LocalFlags().GetString(flagTags)
	if rawTags != "" {
		tags = strings.Split(rawTags, ",")
	}
	if err != nil {
		tt.PrintError(err, quiet)
	}
	if len(args) < 1 || len(args) > 2 {
		tt.PrintError(tt.NewErrorf(tt.ErrorNArgumentsAcceptedFormat, "1-2", len(args)), quiet)
	}
	project = args[0]
	if len(args) > 1 {
		task = args[1]
	}
	return
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
