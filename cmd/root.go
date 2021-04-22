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
	"github.com/maxmoehl/tt/utils"
	"github.com/spf13/cobra"
)

const flagSilent = "silent"

var rootCmd = &cobra.Command{
	Use:   "tt",
	Short: "tt is a cli application that can be used to track time.",
	Long: `tt is a cli application that can be used to track time.

See the help sections of the individual commands for more details on the
functionality.

Note: Timers cannot overlap. If there are overlapping timers
the application might fail and statistics or analytics may be wrong.

Note: This cli is NOT concurrency safe.`,
	Version: "v0.2.2",
}

func init() {
	rootCmd.PersistentFlags().BoolP(flagSilent, string(flagSilent[0]), false, "Suppress all output")
}

// Execute is the main entry point for the cli.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func getSilent(cmd *cobra.Command) bool {
	silent, err := cmd.Flags().GetBool("silent")
	if err != nil {
		utils.PrintError(err, false)
	}
	return silent
}
