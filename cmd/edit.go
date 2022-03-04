package cmd

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit existing timers, even after they are closed",
	RunE: func(cmd *cobra.Command, args []string) error {
		quiet, remove, id, err := getEditParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("edit: %w", err)
		}
		err = runEdit(quiet, remove, id)
		if err != nil {
			return fmt.Errorf("edit: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
	editCmd.Flags().BoolP(flagRemove, short(flagRemove), false, "remove the given timer")
}

func runEdit(quiet, remove bool, id string) error {
	db := tt.GetDB()
	var t tt.Timer
	err := db.GetTimerById(id, &t)
	if err != nil {
		return err
	}
	if remove {
		err := db.RemoveTimer(t.ID)
		if err != nil {
			return err
		}
		if !quiet {
			fmt.Printf("removed timer with id %s\n", t.ID)
		}
		return nil
	} else {
		return fmt.Errorf("no operation provided")
	}
}

func getEditParameters(cmd *cobra.Command, args []string) (quiet, remove bool, id string, err error) {
	flags, err := flags(cmd, flagQuiet, flagRemove)
	if err != nil {
		return
	}
	if len(args) < 1 {
		err = fmt.Errorf("expected one argument")
		return
	}
	id = args[0]
	_, err = uuid.Parse(id)
	if err != nil {
		return
	}
	return flags[flagQuiet].(bool), flags[flagRemove].(bool), id, nil
}
