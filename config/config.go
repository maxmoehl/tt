package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"
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

func init() {
	ttHomeDir := HomeDir()
	file, err := os.Open(filepath.Join(ttHomeDir, "config.yaml"))
	if err != nil {
		// If the occurred error is related to the file not existing we set the defaults
		// by calling validate and return.
		if errors.Is(err, os.ErrNotExist) {
			_ = validate()
			return
		}
		panic(err.Error())
	}
	err = yaml.NewDecoder(file).Decode(&c)
	if err != nil {
		panic(err.Error())
	}
	err = validate()
	if err != nil {
		panic(err.Error())
	}
}

// Get returns the current Config
func Get() Config {
	return c
}

// HomeDir returns the path to the directory that contains for storage
// files and configuration files.
func HomeDir() string {
	ttHomeDir := os.Getenv("TT_HOME_DIR")
	if ttHomeDir == "" {
		ttHomeDir = filepath.Join(os.Getenv("HOME"), ".tt")
	}
	return ttHomeDir
}

// validate checks if the values are inside valid ranges and applies defaults if
// values have not been set.
func validate() error {
	if c.WorkHours == 0 {
		c.WorkHours = 8
	} else if c.WorkHours < 0 || c.WorkHours > 24 {
		return fmt.Errorf("workHours has be bigger than 0 and smaller than 25 but is %d", c.WorkHours)
	}
	return nil
}
