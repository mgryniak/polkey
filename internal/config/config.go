package config

import (
	"encoding/json"
	"errors"
	"os"
)

// DefaultPath is the default configuration file name in the working directory.
const DefaultPath = "config.json"

// Config stores the user's progression and lesson passing thresholds.
type Config struct {
	LastLevel int `json:"last_level"`
	MaxErrors int `json:"max_errors"`
	MinWPM    int `json:"min_wpm"`
}

// Default returns a Config with sensible beginner defaults.
func Default() Config {
	return Config{
		LastLevel: 1,
		MaxErrors: 5,
		MinWPM:    20,
	}
}

// Normalize ensures all fields stay within valid bounds.
func (c *Config) Normalize() {
	if c.LastLevel < 1 {
		c.LastLevel = 1
	}
	if c.MaxErrors < 0 {
		c.MaxErrors = 0
	}
	if c.MinWPM < 1 {
		c.MinWPM = 1
	}
}

// Load reads configuration from path. If the file does not exist, it creates
// one with default values and returns it.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			cfg := Default()
			if saveErr := Save(path, cfg); saveErr != nil {
				return cfg, saveErr
			}
			return cfg, nil
		}
		return Default(), err
	}

	cfg := Default()
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Default(), err
	}
	cfg.Normalize()
	return cfg, nil
}

// Save writes cfg as formatted JSON to path.
func Save(path string, cfg Config) error {
	cfg.Normalize()
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}
