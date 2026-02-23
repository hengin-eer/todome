package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds application settings loaded from config.toml.
type Config struct {
	TodoFile string `toml:"todo_file"`
	DoneFile string `toml:"done_file"`
	Lang     string `toml:"lang"`
}

// DefaultConfig returns a Config with default values.
func DefaultConfig() Config {
	return Config{
		TodoFile: "",
		DoneFile: "",
		Lang:     "ja",
	}
}

// Load reads the config file from the standard path.
// Returns default config if the file does not exist.
func Load() (Config, error) {
	cfg := DefaultConfig()
	path := FilePath()

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return cfg, nil
	}

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return cfg, err
	}
	return cfg, nil
}

// FilePath returns the config file path, respecting XDG_CONFIG_HOME.
func FilePath() string {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "todome", "config.toml")
}
