package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/AlecAivazis/survey/v2"
	"github.com/google/uuid"
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:     "edit <id>",
	Aliases: []string{"e"},
	Short:   "Edit existing timers, even after they are closed",
	RunE: func(cmd *cobra.Command, args []string) error {
		remove, id, err := getEditParameters(cmd, args)
		if err != nil {
			return fmt.Errorf("edit: %w", err)
		}
		err = runEdit(remove, id)
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

func runEdit(remove bool, id string) error {
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
		fmt.Printf("removed timer with id %s\n", t.ID)
		return nil
	} else {
		t, err = openEditor(t)
		if err != nil {
			return err
		}
		return db.UpdateTimer(t)
	}
}

func openEditor(timer tt.Timer) (tt.Timer, error) {
	content, err := json.MarshalIndent(timer, "", "\t")
	if err != nil {
		return tt.Timer{}, fmt.Errorf("editor: %w", err)
	}
	editor := survey.Editor{
		Message:       "Press enter to edit the timer",
		Default:       string(content),
		AppendDefault: true,
	}
	var resp string
	err = survey.AskOne(&editor, &resp)
	if err != nil {
		return tt.Timer{}, fmt.Errorf("editor: %w", err)
	}
	var updatedTimer tt.Timer
	err = json.Unmarshal([]byte(resp), &updatedTimer)
	if err != nil {
		return tt.Timer{}, fmt.Errorf("editor: %w", err)
	}
	if updatedTimer.ID != timer.ID {
		return tt.Timer{}, fmt.Errorf("editor: ID of timer changed")
	}
	return updatedTimer, nil
}

func getEditParameters(cmd *cobra.Command, args []string) (remove bool, id string, err error) {
	flags, err := flags(cmd, flagRemove)
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
	return flags[flagRemove].(bool), id, nil
}
