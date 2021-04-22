package storage

import (
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
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
