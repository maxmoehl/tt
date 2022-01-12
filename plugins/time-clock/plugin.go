package time_clock

import (
	"github.com/maxmoehl/tt"
	"github.com/spf13/cobra"
)

type statsIndex struct {
}

func (i statsIndex) Storages() tt.Storage {

}

func (i statsIndex) Commands() *cobra.Command {
	return nil
}
