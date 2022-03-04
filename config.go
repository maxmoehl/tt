package tt

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

const (
	// HomeDirEnv stores the name of the environment variable that
	// contains the path to the home directory of this cli.
	HomeDirEnv = "TT_HOME_DIR"
)

var (
	c *Config
)

// Config holds all available configuration values.
type Config struct {
	// Precision sets how precise the stats should be evaluated. Available values are: [s second m minute h hour]
	// Default: second
	Precision string `json:"precision"`
	Timeclock struct {
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

// GetConfig returns the current Config
func GetConfig() Config {
	if c == nil {
		err := LoadConfig()
		if err != nil {
			panic(err.Error())
		}
	}
	return *c
}

// HomeDir returns the path to the directory that contains storage
// and configuration files.
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
	return json.NewDecoder(file).Decode(c)
}
