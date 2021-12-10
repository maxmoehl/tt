package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

const (
	flagJSON    = "json"
	flagGroupBy = "group-by"
	flagFilter  = "filter"

	groupByProject = "project"
	groupByTask    = "task"
	groupByDay     = "day"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays various statistics.",
	Long: `Displays various statistics.

This command displays various statistics. The following statistics are
currently available:

  worked    : total time worked
  planned   : planned work time as specified by the config
  percentage: percentage of planned time fulfilled
  difference: absolute difference between planned and worked

The filter string has to be in the following format:

  filterName=values;filterName=values;...

Each filterName consists of a string, values contains the filter value.
Some filters only accept a single value, others accept multiple values
separated by commas.

Example:

  projectName=work,school;since=2020-01-01;until=2020-02-01

Available filters are:

  project: accepts multiple string values
  task   : accepts multiple string values
  since  : accepts a single string int the format of yyyy-MM-dd
  until  : accepts a single string int the format of yyyy-MM-dd
  tags   : accepts multiple string values

since and until are inclusive, both dates will be included in filtered
data.

The group-by string is a comma separated list of groups that should be
formed. Available values are:

  project: Show time by project
  tasks  : Show time by task, automatically sets project
  days   : Show report for each day

Invalid values are ignored.`,
	Run: func(cmd *cobra.Command, args []string) {
		stats(getStatsParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().BoolP(flagJSON, string(flagJSON[0]), false, "Enables printing in the json format, any arguments are ignored.")
	statsCmd.Flags().StringP(flagGroupBy, string(flagGroupBy[0]), "", "Group output by certain aspects: project task day")
	statsCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "Filter the data before generating statistics")
}

func stats(silent, jsonFlag bool, groupBy, filterString string) {
	byProject, byTask, byDay := getGroupByFields(groupBy)
	// the only thing we do is provide output, so there is no point in doing anything if
	// no output should be given
	if silent {
		return
	}
	filter, err := tt.ParseFilterString(filterString)
	if err != nil {
		PrintError(err, silent)
	}
	if byDay {
		statistics, err := tt.GetTimeStatisticsByDay(byProject, byTask, filter)
		if err != nil {
			PrintError(err, silent)
		}
		printStatsStatistics(statistics, jsonFlag)
		if !jsonFlag {
			statistic, err := tt.GetTimeStatistics(false, false, filter)
			if err != nil {
				PrintError(err, silent)
			}
			fmt.Println("Summary:")
			statistic.Print("  ")
		}
	} else {
		statistic, err := tt.GetTimeStatistics(byProject, byTask, filter)
		if err != nil {
			PrintError(err, silent)
		}
		if jsonFlag {
			err = json.NewEncoder(os.Stdout).Encode(statistic)
			if err != nil {
				PrintError(err, false)
			}
		} else {
			statistic.Print("")
		}
	}
}

func printStatsStatistics(statistics map[string]tt.Statistic, j bool) {
	var dates []string

	if j {
		err := json.NewEncoder(os.Stdout).Encode(statistics)
		if err != nil {
			PrintError(err, false)
		}
		return
	}

	for d := range statistics {
		dates = append(dates, d)
	}

	sort.Strings(dates)

	for _, date := range dates {
		fmt.Println(date)
		statistics[date].Print("  ")
		fmt.Println("----------")
	}
}

func getStatsParameters(cmd *cobra.Command, args []string) (silent, jsonFlag bool, groupBy, filter string) {
	var err error
	silent = getSilent(cmd)
	jsonFlag, err = cmd.LocalFlags().GetBool(flagJSON)
	if err != nil {
		PrintError(err, silent)
	}
	groupBy, err = cmd.LocalFlags().GetString(flagGroupBy)
	if err != nil {
		PrintError(err, silent)
	}
	filter, err = cmd.LocalFlags().GetString(flagFilter)
	if err != nil {
		PrintError(err, silent)
	}
	if len(args) != 0 {
		PrintWarning(WarningNoArgumentsAccepted)
	}
	return
}

func getGroupByFields(s string) (byProject, byTask, byDay bool) {
	bys := strings.Split(s, ",")
	for _, b := range bys {
		switch b {
		case groupByProject:
			byProject = true
		case groupByTask:
			byProject = true
			byTask = true
		case groupByDay:
			byDay = true
		}
	}
	return
}
