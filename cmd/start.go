package cmd

import (
	"fmt"
	"strings"
	"time"

	"moehl.dev/tt"

	"github.com/AlecAivazis/survey/v2"
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
		project, task, tags, timestamp, copyFrom, err := getStartParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("start: %w", err)
		}
		err = runStart(project, task, tags, timestamp, copyFrom)
		if err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().StringP(flagTimestamp, short(flagTimestamp), "", "manually set the start time for a timer")
	startCmd.Flags().String(flagTags, "", "specify tags for this timer")
	startCmd.Flags().IntP(flagCopy, short(flagCopy), 0, "copy values from a specific timer")
	startCmd.Flags().BoolP(flagResume, short(flagResume), false, "copy values from the previous timer")
	startCmd.Flags().BoolP(flagInteractive, short(flagInteractive), false, "collect values from stdin")

	// TODO: --auto-stop (or something like this) to stop the previous timer automatically and start a new one
	//       how does this relate to the copy option?
}

func runStart(project, task string, tags []string, timestamp time.Time, copyFrom int) error {
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
	printTrackingStartedMsg(timer)
	return nil
}

func getStartParameters(cmd *cobra.Command, args []string) (project, task string, tags []string, timestamp time.Time, copy int, err error) {
	flags, err := flags(cmd, flagQuiet, flagTags, flagTimestamp, flagCopy, flagResume, flagInteractive)
	if err != nil {
		return
	}
	if flags[flagInteractive].(bool) && flags[flagQuiet].(bool) {
		err = fmt.Errorf("interactive and quiet cannot be set together")
		return
	}
	if flags[flagResume].(bool) && flags[flagCopy].(int) <= 0 {
		flags[flagCopy] = 1
	}
	if len(args) > 0 {
		project = args[0]
	}
	if len(args) > 1 {
		task = args[1]
	}
	if flags[flagInteractive].(bool) || (project == "" && flags[flagCopy] == 0) {
		project, task, timestamp, tags, err = getStartParametersInteractive()
		return
	}
	return project, task, flags[flagTags].([]string), flags[flagTimestamp].(time.Time), flags[flagCopy].(int), nil
}

func in(a string, b []string) bool {
	for _, c := range b {
		if a == c {
			return true
		}
	}
	return false
}

func getStartParametersInteractive() (project, task string, timestamp time.Time, tags []string, err error) {
	order := tt.OrderBy{Field: tt.FieldStart, Order: tt.OrderDsc}
	var timers tt.Timers
	err = tt.GetDB().GetTimers(tt.EmptyFilter, order, &timers)
	if err != nil {
		return
	}
	var allProjects []string
	allTasks := make(map[string][]string)
	for _, t := range timers {
		if !in(t.Project, allProjects) {
			allProjects = append(allProjects, t.Project)
		}
		if t.Task != "" && !in(t.Task, allTasks[t.Project]) {
			allTasks[t.Project] = append(allTasks[t.Project], t.Task)
		}
	}

	timestampDefaultStr := ""
	timestampDefault := time.Now().Round(tt.GetConfig().GetRoundStartTime())
	if len(timers) > 0 && timers[0].Stop != nil {
		timestampDefault = *timers[0].Stop
	}

	if datesEqual(time.Now(), timestampDefault) {
		// today, so we can leave out the date
		timestampDefaultStr = fmt.Sprintf("%02d:%02d", timestampDefault.Hour(), timestampDefault.Minute())
	} else {
		timestampDefaultStr = fmt.Sprintf("%04d-%02d-%02d %02d:%02d", timestampDefault.Year(), timestampDefault.Month(), timestampDefault.Day(), timestampDefault.Hour(), timestampDefault.Minute())
	}

	answers := new(struct {
		Project   string
		Task      string
		Timestamp string
		Tags      string
	})

	var suggestedTags []string

	qs := []*survey.Question{
		{
			Name: "timestamp",
			Prompt: &survey.Input{
				Message: "Enter a start timestamp",
				Default: timestampDefaultStr,
			},
		},
		{
			Name: "project",
			Prompt: &survey.Input{
				Message: "Enter a project",
				Default: "",
				Suggest: func(toComplete string) (suggestions []string) {
					for _, project := range allProjects {
						if strings.HasPrefix(strings.ToLower(project), strings.ToLower(toComplete)) {
							suggestions = append(suggestions, project)
						}
					}
					return
				},
			},
			Validate: func(ans interface{}) error {
				if project, ok := ans.(string); !ok {
					return fmt.Errorf("%w: project must be of type string", tt.ErrInvalidParameter)
				} else if project == "" {
					return fmt.Errorf("%w: project cannot be empty", tt.ErrInvalidParameter)
				}
				return nil
			},
		},
		{
			Name: "task",
			Prompt: &survey.Input{
				Message: "Enter a task (optional)",
				Default: "",
				Suggest: func(toComplete string) (suggestions []string) {
					for _, task := range allTasks[answers.Project] {
						if strings.HasPrefix(strings.ToLower(task), strings.ToLower(toComplete)) {
							suggestions = append(suggestions, task)
						}
					}
					return
				},
			},
		},
		{
			Name: "tags",
			Prompt: &survey.Input{
				Message: "Enter tags (optional)",
				Default: "",
				Suggest: func(toComplete string) []string {
					if suggestedTags != nil {
						return suggestedTags
					}
					suggestedTags = make([]string, 0)
					f := tt.NewFilter([]string{answers.Project}, []string{answers.Task}, nil, time.Time{}, time.Time{})
					var timers tt.Timers
					err = tt.GetDB().GetTimers(f, tt.OrderBy{}, &timers)
					if err != nil {
						panic(err.Error())
					}
					for _, t := range timers {
						tags := strings.Join(t.Tags, ",")
						skip := false
						for _, existingTags := range suggestedTags {
							if tags == existingTags {
								skip = true
							}
						}
						if skip {
							continue
						}
						suggestedTags = append(suggestedTags, tags)
					}
					return suggestedTags
				},
			},
		},
	}
	err = survey.Ask(qs, answers)
	if err != nil {
		err = fmt.Errorf("interactive input: %w", err)
		return
	}
	timestamp, err = tt.ParseTime(answers.Timestamp)
	if err != nil {
		return
	}
	tags = strings.Split(answers.Tags, ",")
	return answers.Project, answers.Task, timestamp, tags, nil
}

func printTrackingStartedMsg(t tt.Timer) {
	fmt.Printf("[%02d:%02d] Tracking started!\n", t.Start.Hour(), t.Start.Minute())
	fmt.Printf("  project: %s\n", t.Project)
	if t.Task != "" {
		fmt.Printf("  task   : %s\n", t.Task)
	}
	if len(t.Tags) > 0 {
		fmt.Printf("  tags   : %s\n", strings.Join(t.Tags, ","))
	}
}
