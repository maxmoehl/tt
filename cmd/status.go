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

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Prints a short notice on the current status",
	Long: `Reports if you are currently working, taking a break or taking some
time off.`,
	Run: status,
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func status(cmd *cobra.Command, args []string) {
	silent := getSilent(cmd)
	if len(args) != 0 && !silent {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	if silent {
		return
	}
	found, timer, err := storage.GetRunningTimer()
	if err != nil {
		utils.PrintError(err, silent)
	}
	if !found {
		fmt.Println("Currently not working. Enjoy your free time :)")
	} else {
		if timer.Task != "" {
			fmt.Printf("Currently timing project %s with task %s, your doing good!\n", timer.Project, timer.Task)
		} else {
			fmt.Printf("Currently timing project %s, your doing good!\n", timer.Project)
		}
	}
}
