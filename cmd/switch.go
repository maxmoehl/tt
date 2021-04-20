package cmd

import (
	"github.com/spf13/cobra"
)

var switchCmd = &cobra.Command{
	Use: "switch project",
	Short: "Switch allows you to quickly switch between different projects and tasks",
	Long: `This command is a shorthand for:
tt stop
tt start project [-t task]

It stops the current timer (and break of open) and starts a new timer with
the given project and task.

Note: switch does not support manually setting a start time.`,
	Run: switchR,
}

func init() {
	rootCmd.AddCommand(switchCmd)
	switchCmd.Flags().StringP(flagTask, string(flagTask[0]), "", "Set the task your are currently working on.")
}

func switchR(cmd *cobra.Command, args []string) {
	// The weird R is needed because switch is a keyword
	stop(cmd, nil)
	start(cmd, args)
}
