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
	"encoding/json"
	"os"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

const (
	flagJSON      = "json"
	flagByProject = "by-project"
	flagByTask    = "by-task"
	flagFilter    = "filter"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Displays various statistics",
	Long: `This command displays various statistics. The following statistics are
currently available:
  worked    : total time worked (excluding brakes)
  planned   : planned work time as specified by the config
  breaks    : total time that has been spent in breaks
  percentage: percentage of planned time fulfilled
  difference: absolute difference between planned and worked



If no duration is given all days starting from the oldest timer stored
will be considered.`,
	Run: stats,
}

func init() {
	rootCmd.AddCommand(statsCmd)
	statsCmd.Flags().Bool(flagJSON, false, "Enables printing in the json format, any arguments are ignored.")
	statsCmd.Flags().Bool(flagByProject, false, "Print times grouped by project")
	statsCmd.Flags().Bool(flagByTask, false, "Print times grouped by task")
	statsCmd.Flags().String(flagFilter, "", "Filter the data before generating statistics")
}

func stats(cmd *cobra.Command, args []string) {
	silent, j, byProject, byTask, filterString := getStatsFlags(cmd)
	if len(args) != 0 {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	// the only thing we do is provide output, so there is no point in doing anything if
	// no output should be given
	if silent {
		return
	}
	filter, err := types.GetFilter(filterString)
	s, err := storage.GetTimeStatistics(byProject, byTask, filter)
	if err != nil {
		utils.PrintError(err, silent)
	}
	if j {
		printStatsJson(s)
	} else {
		printStatsText(s)
	}
}

func getStatsFlags(cmd *cobra.Command) (silent, jsonFlag, byProject, byTask bool, filter string) {
	var err error
	silent = getSilent(cmd)
	jsonFlag, err = cmd.LocalFlags().GetBool(flagJSON)
	if err != nil {
		utils.PrintError(err, silent)
	}
	byProject, err = cmd.LocalFlags().GetBool(flagByProject)
	if err != nil {
		utils.PrintError(err, silent)
	}
	byTask, err = cmd.LocalFlags().GetBool(flagByTask)
	if err != nil {
		utils.PrintError(err, silent)
	}
	filter, err = cmd.LocalFlags().GetString(flagFilter)
	if err != nil {
		utils.PrintError(err, silent)
	}
	return
}

func printStatsText(s types.Statistic) {
	s.Print()
}

func printStatsJson(s types.Statistic) {
	err := json.NewEncoder(os.Stdout).Encode(s)
	if err != nil {
		utils.PrintError(err, false)
	}
}
