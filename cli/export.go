package main

import (
	"encoding/json"
	"fmt"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

const (
	exportFormatCSV  = "csv"
	exportFormatJSON = "json"
	exportFormatSQL  = "sql"
)

var exportCmd = &cobra.Command{
	Use:   "export format",
	Short: "Export data to a given format",
	Run: func(cmd *cobra.Command, args []string) {
		export(getExportParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "Set a filter to apply before exporting")
}

func getExportParameters(cmd *cobra.Command, args []string) (silent bool, filter, exportFormat string) {
	var err error
	silent = getSilent(cmd)
	if len(args) != 1 {
		PrintError(fmt.Errorf("expected one argument"), silent)
	}
	exportFormat = args[0]
	if exportFormat != exportFormatCSV && exportFormat != exportFormatJSON && exportFormat != exportFormatSQL {
		PrintError(fmt.Errorf("unknown export format %s", exportFormat), silent)
	}
	filter, err = cmd.LocalFlags().GetString(flagFilter)
	if err != nil {
		PrintError(err, silent)
	}
	return
}

func export(silent bool, filterString, exportFormat string) {
	if silent {
		return
	}
	filter, err := tt.ParseFilterString(filterString)
	if err != nil {
		PrintError(err, silent)
	}
	timers, err := tt.GetTimers(filter)
	if err != nil {
		PrintError(err, silent)
	}
	var out string
	switch exportFormat {
	case exportFormatJSON:
		var b []byte
		b, err = json.Marshal(timers)
		out = string(b)
	case exportFormatCSV:
		out, err = timers.CSV()
	case exportFormatSQL:
		out, err = timers.SQL()
	}
	if err != nil {
		PrintError(err, silent)
	}
	fmt.Println(out)
}
