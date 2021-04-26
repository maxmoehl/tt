package storage

import (
	"fmt"
	"testing"

	"github.com/maxmoehl/tt/test"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"
)

func TestMain(m *testing.M) {
	test.Main(m.Run)
}

func timersEqual(t1, t2 types.Timer) error {
	if t1.Uuid != t2.Uuid {
		return fmt.Errorf("uuids are not equal")
	}
	if t1.Project != t2.Project {
		return fmt.Errorf("projects are not equal")
	}
	if t1.Task != t2.Task {
		return fmt.Errorf("tasks are not equal")
	}
	if len(t1.Tags) != len(t2.Tags) {
		return fmt.Errorf("tags are not equal")
	}
	for _, tag := range t1.Tags {
		if !utils.StringSliceContains(t2.Tags, tag) {
			return fmt.Errorf("tags are not equal")
		}
	}
	if !t1.Start.Equal(t2.Start) {
		return fmt.Errorf("start time is not equal")
	}
	if !t1.Stop.Equal(t2.Stop) {
		return fmt.Errorf("end time is not equal")
	}
	return nil
}
