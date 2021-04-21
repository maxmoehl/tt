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

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

var inspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspects the data and looks for inconsistencies to report.",
	Long: `Inspects the data and looks for inconsistencies to report.

If any timers are printed and the command reports an error, this means
that there is more than one running timer. Please try to remove the
timer form the data source or finish it by setting a valid end time.`,
	Run: inspect,
}

func init() {
	rootCmd.AddCommand(inspectCmd)
}

func inspect(cmd *cobra.Command, args []string) {
	silent := getSilent(cmd)
	if len(args) != 0 && !silent {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	if silent {
		return
	}
	// TODO: add check if there are running timers that are older than the most recent timer
	fmt.Println("checking for running timers...")
	timerUuids, err := storage.CheckRunningTimers()
	if err != nil {
		utils.PrintError(err, silent)
	}
	if len(timerUuids) == 0 {
		fmt.Println("found no running timers")
		fmt.Println("ok")
	} else if len(timerUuids) == 1 {
		fmt.Println("found one running timer:")
		fmt.Printf("\t%s\n", timerUuids[0].String())
		fmt.Println("ok")
	} else {
		fmt.Println("found more then one running timer:")
		for _, u := range timerUuids {
			fmt.Printf("\t%s\n", u.String())
		}
		fmt.Println("ERROR")
	}
}
