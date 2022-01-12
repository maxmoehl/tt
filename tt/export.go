package main

import (
	"encoding/json"
	"fmt"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

const (
	flagFilter = "filter"

	exportFormatCSV  = "csv"
	exportFormatJSON = "json"
	exportFormatSQL  = "sql"
)

var exportCmd = &cobra.Command{
	Use:   "export <format>",
	Short: "Export data to a given format",
	Run: func(cmd *cobra.Command, args []string) {
		runExport(getExportParameters(cmd, args))
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "Set a filter to apply before exporting")
}

func getExportParameters(cmd *cobra.Command, args []string) (quiet bool, exportFormat string, filter tt.Filter) {
	var err error
	quiet = getQuiet(cmd)
	if len(args) != 1 {
		tt.PrintError(fmt.Errorf("expected one argument"), quiet)
	}
	exportFormat = args[0]
	if exportFormat != exportFormatCSV && exportFormat != exportFormatJSON && exportFormat != exportFormatSQL {
		tt.PrintError(fmt.Errorf("unknown export format %s", exportFormat), quiet)
	}
	rawFilter, err := cmd.LocalFlags().GetString(flagFilter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	filter, err = tt.ParseFilterString(rawFilter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	return
}

func runExport(quiet bool, exportFormat string, filter tt.Filter) {
	if quiet {
		return
	}
	timers, err := tt.GetStorage().GetTimers(filter)
	if err != nil {
		tt.PrintError(err, quiet)
	}
	var out string
	var e error
	switch exportFormat {
	case exportFormatJSON:
		var b []byte
		b, e = json.Marshal(timers)
		out = string(b)
	case exportFormatCSV:
		out, e = timers.CSV()
	case exportFormatSQL:
		out, e = timers.SQL()
	}
	if e != nil {
		tt.PrintError(e, quiet)
	}
	fmt.Println(out)
}
