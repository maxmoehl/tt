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
	"fmt"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

const (
	exportFormatCSV  = "csv"
	exportFormatJSON = "json"
)

var exportCmd = &cobra.Command{
	Use:   "export format",
	Short: "Export data to a given format",
	Run:   export,
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "Set a filter to apply before exporting")
}

func export(cmd *cobra.Command, args []string) {
	silent, filter, exportFormat := getExportParameters(cmd, args)
	if silent {
		return
	}
	timers, err := storage.GetTimers(filter)
	if err != nil {
		utils.PrintError(err, silent)
	}
	var out string
	switch exportFormat {
	case exportFormatJSON:
		var b []byte
		b, err = json.Marshal(timers)
		out = string(b)
	case exportFormatCSV:
		out, err = timers.CSV()
	}
	if err != nil {
		utils.PrintError(err, silent)
	}
	fmt.Println(out)
}

func getExportParameters(cmd *cobra.Command, args []string) (silent bool, filter types.Filter, exportFormat string) {
	silent = getSilent(cmd)
	if len(args) != 1 {
		utils.PrintError(fmt.Errorf("expected one argument"), silent)
	}
	exportFormat = args[0]
	if exportFormat != exportFormatCSV && exportFormat != exportFormatJSON {
		utils.PrintError(fmt.Errorf("unknown export format %s", exportFormat), silent)
	}
	filterString, err := cmd.LocalFlags().GetString(flagFilter)
	if err != nil {
		utils.PrintError(err, silent)
	}
	filter, err = types.GetFilter(filterString)
	if err != nil {
		utils.PrintError(err, silent)
	}
	return
}
