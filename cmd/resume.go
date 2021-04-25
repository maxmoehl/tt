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
	"strings"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/utils"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the last timer.",
	Long: `Resume the last timer.

If no timer is found an error is returned.`,
	Run: func(cmd *cobra.Command, args []string) {
		resume(getResumeParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(resumeCmd)
}

func resume(silent bool) {
	t, err := storage.ResumeTimer()
	if err != nil {
		utils.PrintError(err, silent)
	}
	if !silent {
		fmt.Println("Work tracking started!")
		fmt.Printf("  project: %s\n", t.Project)
		if t.Task != "" {
			fmt.Printf("  task   : %s\n", t.Task)
		}
		if len(t.Tags) > 0 {
			fmt.Printf("  tags   : %s\n", strings.Join(t.Tags, ", "))
		}
	}
}

func getResumeParameters(cmd *cobra.Command, args []string) (silent bool) {
	silent = getSilent(cmd)
	if len(args) != 0 && !silent {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	return
}
