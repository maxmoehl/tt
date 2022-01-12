package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List all existing timers.",
	Run: func(cmd *cobra.Command, args []string) {
		runList(getListParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "Filter results before printing")
}

func runList(quiet bool, filter tt.Filter) {
	if quiet {
		return
	}
	timers, err := tt.GetStorage().GetTimers(filter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	var totalDuration time.Duration = 0
	for _, timer := range timers {
		fmt.Printf("Start   : %s\n", timer.Start.String())
		if !timer.Stop.IsZero() {
			fmt.Printf("Stop    : %s\n", timer.Stop.String())
			fmt.Printf("Duration: %s\n", tt.FormatDuration(timer.Duration(), tt.GetConfig().GetPrecision()))
			totalDuration += timer.Duration()
		}
		fmt.Printf("Project : %s\n", timer.Project)
		if timer.Task != "" {
			fmt.Printf("Task    : %s\n", timer.Task)
		}
		if len(timer.Tags) > 0 {
			fmt.Printf("Tags    : %s\n", strings.Join(timer.Tags, ", "))
		}
		fmt.Printf("--------\n")
	}
	fmt.Printf("Total duration tracked: %s\n", tt.FormatDuration(totalDuration, tt.GetConfig().GetPrecision()))
}

func getListParameters(cmd *cobra.Command, args []string) (quiet bool, filter tt.Filter) {
	quiet = getQuiet(cmd)
	if len(args) != 0 && !quiet {
		tt.PrintWarning(tt.WarningNoArgumentsAccepted)
	}
	rawFilter, err := cmd.Flags().GetString(flagFilter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	filter, err = tt.ParseFilterString(rawFilter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	return
}
