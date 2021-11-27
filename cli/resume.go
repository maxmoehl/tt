package main

import (
	"fmt"
	"strings"

	"github.com/maxmoehl/tt"

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
	t, err := tt.ResumeTimer()
	if err != nil {
		PrintError(err, silent)
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
		PrintWarning(WarningNoArgumentsAccepted)
	}
	return
}
