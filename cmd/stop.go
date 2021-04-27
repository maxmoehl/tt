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
	"time"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:     "stop",
	Aliases: []string{"end"},
	Short:   "Stops a timer.",
	Long: `This command stops the current timer. If the timer is currently in a break
the break is also ended without further notice.

If you want to manually set a stop time it has to be in the following
format:

  2020-04-19T08:00:00+02:00

Otherwise an appropriate error will be printed. The cli will check if the
given stop time is valid, e.g. if the last timer and break that were
started, started before the given stop.`,
	Run: func(cmd *cobra.Command, args []string) {
		stop(getStopParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
	stopCmd.Flags().String(flagTimestamp, "", "Manually set the stop time for a timer. Format must be RFC3339")
}

func stop(silent bool, timestamp string) {
	timer, err := storage.StopTimer(timestamp)
	if err != nil {
		utils.PrintError(err, silent)
	}
	if !silent {
		fmt.Printf("You worked for %s! Good job.\n", timer.Duration().Round(time.Second).String())
	}
}

func getStopParameters(cmd *cobra.Command, args []string) (silent bool, timestamp string) {
	var err error
	silent = getSilent(cmd)
	timestamp, err = cmd.Flags().GetString(flagTimestamp)
	if err != nil {
		utils.PrintError(err, silent)
	}
	if len(args) != 0 && !silent {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	return
}
