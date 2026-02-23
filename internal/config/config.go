package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds application settings loaded from config.toml.
type Config struct {
	DataDir  string `toml:"data_dir"`
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

// DataDirPath returns the data directory path.
// Priority: config data_dir > XDG_DATA_HOME/todome/
func (c Config) DataDirPath() string {
	if c.DataDir != "" {
		return ExpandHome(c.DataDir)
	}
	dir := os.Getenv("XDG_DATA_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dir, "todome")
}

// ExpandHome replaces a leading ~/ with the user's home directory.
func ExpandHome(path string) string {
	if len(path) >= 2 && path[:2] == "~/" {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}
