package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"moehl.dev/tt"

	"github.com/spf13/cobra"
)

// TODO: we need a few more options here:
//       - limit amount of printed timers
//       - specify order of timers

const (
	groupByProject = "project"
	groupByTask    = "task"
	groupByDay     = "day"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all existing timers",
	Long: `List all existing timers.

The filter (if provided) has to be in the following format:
  filterName=values;filterName=values;...

Each filterName consists of a string, values contains the filter value. Some
filters only accept a single value, others accept multiple values separated by
commas.

Example:
  projectName=work,school;since=2020-01-01;until=2020-02-01

Available filters are:
  project: accepts multiple string values
  task   : accepts multiple string values
  since  : accepts a single string int the format of yyyy-MM-dd
  until  : accepts a single string int the format of yyyy-MM-dd
  tags   : accepts multiple string values

since and until are inclusive, both dates will be included in filtered data.

The group-by string is a comma separated list of groups that should be formed.
Available values are:
  project: Show time by project
  tasks  : Show time by task, automatically sets project
  days   : Show report for each day`,
	RunE: func(cmd *cobra.Command, args []string) error {
		filter, groupBy, short, err := getListParameters(cmd, args)
		if err != nil {
			return err
		}
		err = runList(filter, groupBy, short)
		if err != nil {
			return fmt.Errorf("list: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "filter results before printing")
	listCmd.Flags().StringP(flagGroupBy, string(flagGroupBy[0]), "", "group results before printing")
	listCmd.Flags().BoolP(flagShort, string(flagShort[0]), false, "shorten the output")
}

func runList(filter tt.Filter, groupBy string, short bool) error {
	orderBy := tt.OrderBy{
		Field: tt.FieldStart,
		Order: tt.OrderAsc,
	}

	timers, err := tt.List(filter, orderBy)
	if err != nil {
		return err
	}

	switch groupBy {
	case "":
		printTimers(timers, short)
	case groupByDay:
		printTimersGrouped(timers.GroupByDay(), short)
	case groupByProject:
		printTimersGrouped(timers.GroupByProject(), short)
	case groupByTask:
		printTimersByTask(timers.GroupByTask(), short)
	default:
		return fmt.Errorf("list: unknown group by option: %s", groupBy)
	}

	fmt.Printf("Overall total duration tracked: %s\n", tt.FormatDuration(timers.Duration()))
	return nil
}

func getListParameters(cmd *cobra.Command, _ []string) (filter tt.Filter, groupBy string, short bool, err error) {
	flags, err := flags(cmd, flagFilter, flagGroupBy, flagShort)
	if err != nil {
		return
	}
	return flags[flagFilter].(tt.Filter), flags[flagGroupBy].(string), flags[flagShort].(bool), nil
}

func printTimers(timers tt.Timers, short bool) {
	var totalDuration time.Duration = 0
	for _, t := range timers {
		if short {
			task := "no-task"
			if t.Task != "" {
				task = t.Task
			}
			fmt.Printf("%s (%s) %s / %s\n", t.Start.Format(tt.TimeFormat), tt.FormatDuration(t.Duration()), t.Project, task)
		} else {
			fmt.Println(t.String())
			fmt.Println("--------")
		}
		totalDuration += t.Duration()
	}
	fmt.Printf("Total duration tracked: %s\n\n", tt.FormatDuration(totalDuration))
}

func printTimersGrouped(groupedTimers map[string]tt.Timers, short bool) {
	var keys []string
	for key := range groupedTimers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if short {
			fmt.Printf("%s: %s\n", key, tt.FormatDuration(groupedTimers[key].Duration()))
		} else {
			fmt.Printf("### %s ###\n", key)
			printTimers(groupedTimers[key], false)
			fmt.Println()
		}
	}
}

func printTimersByTask(groupedTimers map[string]map[string]tt.Timers, short bool) {
	var keys []string
	for key := range groupedTimers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		var d time.Duration
		for _, t := range groupedTimers[key] {
			d += t.Duration()
		}
		if !short {
			fmt.Printf("####%s####\n", strings.Repeat("#", len(key)))
		}
		fmt.Printf("### %s (%s) ###\n", key, tt.FormatDuration(d))
		if !short {
			fmt.Printf("####%s####\n\n", strings.Repeat("#", len(key)))
		}
		printTimersGrouped(groupedTimers[key], short)
		fmt.Println()
	}
}
