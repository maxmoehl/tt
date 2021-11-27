/*
Copyright © 2021 Maximilian Moehl contact@moehl.eu

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
)

const (
	// HomeDirEnv stores the name of the environment variable that
	// contains the path to the home directory of this cli.
	HomeDirEnv = "TT_HOME_DIR"

	// StorageTypeFile is the identifier for file storage configuration
	StorageTypeFile   = "file"
	// StorageTypeSQLite is the identifier for SQLite storage configuration
	StorageTypeSQLite = "sqlite"
)

var c Config

// Config holds all available configuration values.
type Config struct {
	// WorkHours specifies the hours you want to work per day.
	// Default: 8
	WorkHours int `yaml:"workHours"`
	// WorkDays specifies on which days of the week you work.
	WorkDays struct {
		Monday    bool `yaml:"monday"`
		Tuesday   bool `yaml:"tuesday"`
		Wednesday bool `yaml:"wednesday"`
		Thursday  bool `yaml:"thursday"`
		Friday    bool `yaml:"friday"`
		Saturday  bool `yaml:"saturday"`
		Sunday    bool `yaml:"sunday"`
	} `yaml:"workDays"`
	// Precision sets how precise the stats should be evaluated. Available values are: [s second m minute h hour]
	// Default: second
	Precision string `yaml:"precision"`
	// StorageType indicates which type of storage tt should use. Available values are: [file sqlite]
	// Default: file
	StorageType string `yaml:"storageType"`
}

// GetPrecision returns the precision as a duration.
func (c Config) GetPrecision() time.Duration {
	switch c.Precision {
	case "h":
		fallthrough
	case "hour":
		return time.Hour
	case "m":
		fallthrough
	case "minute":
		return time.Minute
	case "s":
		fallthrough
	case "second":
		fallthrough
	default:
		return time.Second
	}
}

// Get returns the current Config
func Get() Config {
	return c
}

// HomeDir returns the path to the directory that contains for storage
// files and configuration files.
func HomeDir() string {
	ttHomeDir := os.Getenv(HomeDirEnv)
	if ttHomeDir == "" {
		ttHomeDir = filepath.Join(os.Getenv("HOME"), ".tt")
	}
	return ttHomeDir
}

// Load allows to manually load the configuration file.
func Load() error {
	ttHomeDir := HomeDir()
	file, err := os.Open(filepath.Join(ttHomeDir, "config.yaml"))
	if err != nil {
		// If the occurred error is related to the file not existing we set the defaults
		// by calling validate and return.
		if errors.Is(err, os.ErrNotExist) {
			return validate()
		}
		return err
	}
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		return err
	}
	return validate()
}

// validate checks if the values are inside valid ranges and applies defaults if
// values have not been set.
func validate() error {
	if c.WorkHours == 0 {
		c.WorkHours = 8
	} else if c.WorkHours < 0 || c.WorkHours > 24 {
		return fmt.Errorf("workHours has to be bigger than 0 and smaller than 25 but is %d", c.WorkHours)
	}
	if c.StorageType == "" {
		c.StorageType = StorageTypeFile
	} else if c.StorageType != StorageTypeFile && c.StorageType != StorageTypeSQLite {
		return fmt.Errorf("invalid storage type %s", c.StorageType)
	}
	return nil
}
