package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/maxmoehl/tt"

	"github.com/spf13/cobra"
)

const (
	exportFormatCSV  = "csv"
	exportFormatJSON = "json"
)

var exportCmd = &cobra.Command{
	Use:       "export <format>",
	Short:     "Export data to a given format",
	Example:   "tt export json -f since=today",
	ValidArgs: []string{exportFormatCSV, exportFormatJSON},
	RunE: func(cmd *cobra.Command, args []string) error {
		exportFormat, filter, err := getExportParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}
		err = runExport(exportFormat, filter)
		if err != nil {
			return fmt.Errorf("export: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(exportCmd)
	exportCmd.Flags().StringP(flagFilter, string(flagFilter[0]), "", "set a filter to apply before exporting")
}

func runExport(exportFormat string, filter tt.Filter) error {
	var timers tt.Timers
	err := tt.GetDB().GetTimers(filter, tt.OrderBy{}, &timers)
	if err != nil {
		return err
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
	default:
		err = fmt.Errorf("unknown format: %s", exportFormat)
	}
	if e != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

func getExportParameters(cmd *cobra.Command, args []string) (exportFormat string, filter tt.Filter, err error) {
	flags, err := flags(cmd, flagFilter)
	if err != nil {
		return
	}
	if len(args) != 1 {
		err = fmt.Errorf("expected one argument")
		return
	}
	return args[0], flags[flagFilter].(tt.Filter), nil
}
