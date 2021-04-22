package storage

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
	"github.com/maxmoehl/tt/types"
	"github.com/maxmoehl/tt/utils"
)

const configFile = `precision: m
workHours: 8
workDays:
  monday: true
  tuesday: true
  wednesday: true
  thursday: true
  friday: true
  saturday: false
  sunday: false`

var testFile *file

func setup() string {
	dir, err := setupTempDir()
	if err != nil {
		panic(err.Error())
	}
	// in case the setup fails we teardown right after the setup
	defer func() {
		if err != nil {
			teardown(dir)
		}
	}()
	err = os.Setenv(config.HomeDirEnv, dir)
	if err != nil {
		panic(err.Error())
	}
	err = config.LoadConfig()
	if err != nil {
		panic(err.Error())
	}
	return dir
}

func teardown(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err.Error())
	}
}

func setupTempDir() (string, error) {
	dir, err := os.MkdirTemp("", "tt-testing-*")
	if err != nil {
		return "", err
	}
	f := filepath.Join(dir, "config.yaml")
	err = os.WriteFile(f, []byte(configFile), 0666)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func reloadTestFile() error {
	var err error
	s, err = NewFile()
	if err != nil {
		return err
	}
	var ok bool
	testFile, ok = s.(*file)
	if !ok {
		return fmt.Errorf("expected storage to be of type *file")
	}
	return nil
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
	if !t1.End.Equal(t2.End) {
		return fmt.Errorf("end time is not equal")
	}
	return nil
}
