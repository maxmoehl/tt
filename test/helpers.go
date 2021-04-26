package test

import (
	"os"
	"path/filepath"

	"github.com/maxmoehl/tt/config"
)

var testDir string

func Main(run func() int) {
	dir := setup()
	defer teardown(dir)

	exitCode := run()
	if exitCode != 0 {
		// os.Exit does not run deferred functions, therefore we run it manually
		teardown(dir)
		os.Exit(exitCode)
	}
}

func setup() string {
	var err error
	testDir, err = os.MkdirTemp("", "tt-testing-*")
	if err != nil {
		panic(err.Error())
	}
	err = os.Setenv(config.HomeDirEnv, testDir)
	if err != nil {
		panic(err.Error())
	}
	return testDir
}

func teardown(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		panic(err.Error())
	}
}

// SetConfig sets the given string as configfile and reloads the config.
// This ensures every test has the correct config when being executed.
func SetConfig(configFile string) error {
	err := os.WriteFile(filepath.Join(testDir, "config.yaml"), []byte(configFile), 0666)
	if err != nil {
		return err
	}
	err = config.LoadConfig()
	if err != nil {
		return err
	}
	return nil
}