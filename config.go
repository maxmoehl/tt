package tt

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const (
	// HomeDirEnv stores the name of the environment variable that contains the
	// path to the home directory of this cli.
	HomeDirEnv = "TT_HOME_DIR"
)

var (
	c *Config
)

// Config holds all available configuration values.
type Config struct {
	// Precision sets how precise the stats should be evaluated. Available
	// values are: [s second m minute h hour]
	// Default: second
	Precision string `json:"precision"`
	AutoStop  bool   `json:"autoStop"`
	// RoundStartTime will take the start time and round by the factor given.
	// Example:
	//   60s: 13:45:23 -> 13:45:00
	//   60s: 09:01:59 -> 09:02:00
	//   5m : 23:32:29 -> 23:30:00
	// Refer to time.Time.Round on how it works
	RoundStartTime string `json:"roundStartTime"`
	Timeclock      struct {
		HoursPerDay int `json:"hoursPerDay"`
		DaysPerWeek struct {
			Monday    bool `json:"monday"`
			Tuesday   bool `json:"tuesday"`
			Wednesday bool `json:"wednesday"`
			Thursday  bool `json:"thursday"`
			Friday    bool `json:"friday"`
			Saturday  bool `json:"saturday"`
			Sunday    bool `json:"sunday"`
		} `json:"daysPerWeek"`
	} `json:"timeclock"`
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

func (c Config) Validate() error {
	_, err := time.ParseDuration(c.RoundStartTime)
	if err != nil {
		return fmt.Errorf("config: validate: %w", err)
	}
	return nil
}

func (c Config) GetRoundStartTime() time.Duration {
	d, err := time.ParseDuration(c.RoundStartTime)
	if err != nil {
		panic(err.Error())
	}
	return d
}

// GetConfig returns the current Config and lazy loads it if necessary.
func GetConfig() Config {
	if c == nil {
		err := LoadConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	return *c
}

// HomeDir returns the path to the directory that contains storage and
// configuration files.
func (Config) HomeDir() string {
	ttHomeDir := os.Getenv(HomeDirEnv)
	if ttHomeDir == "" {
		ttHomeDir = filepath.Join(os.Getenv("HOME"), ".tt")
	}

	return ttHomeDir
}

func (c Config) DBFile() string {
	return filepath.Join(c.HomeDir(), "storage.db")
}

// LoadConfig allows to manually load the configuration file.
func LoadConfig() error {
	c = &Config{}
	ttHomeDir := c.HomeDir()

	file, err := os.Open(filepath.Join(ttHomeDir, "config.json"))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return err
	}

	err = json.NewDecoder(file).Decode(c)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	if c.RoundStartTime == "" {
		c.RoundStartTime = "0"
	}

	err = c.Validate()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	return nil
}
