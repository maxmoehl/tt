package cmd

import (
	"fmt"

	"github.com/maxmoehl/tt/storage"
	"github.com/maxmoehl/tt/utils"

	"github.com/spf13/cobra"
)

var breakCmd = &cobra.Command{
	Use:   "break",
	Short: "Start or end a break",
	Long: `This command allows you to take a break from your current activity without
having to manually start and stop the timer. The first time you call break
a break will be started and a short prompt will be printed. The second
time you call it the break will be stopped. Stopping the timer will also
stop any ongoing breaks.`,
	Run: breakR,
}

func init() {
	rootCmd.AddCommand(breakCmd)
}

func breakR(cmd *cobra.Command, args []string) {
	// The weird R is needed because break is a keyword
	silent := getSilent(cmd)
	if len(args) != 0 && !silent {
		utils.PrintWarning(utils.WarningNoArgumentsAccepted)
	}
	openBreak, err := storage.ToggleBreak()
	if err != nil {
		utils.PrintError(err, silent)
	}
	if !silent {
		if openBreak {
			fmt.Println("Break started! Enjoy some time off.")
		} else {
			fmt.Println("Break ended! Time to get things done.")
		}
	}
}
