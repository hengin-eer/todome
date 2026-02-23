package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.TodoFile != "" {
		t.Errorf("expected empty TodoFile, got %q", cfg.TodoFile)
	}
	if cfg.DoneFile != "" {
		t.Errorf("expected empty DoneFile, got %q", cfg.DoneFile)
	}
	if cfg.Lang != "ja" {
		t.Errorf("expected lang 'ja', got %q", cfg.Lang)
	}
}

func TestLoadMissing(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TodoFile != "" {
		t.Errorf("expected empty TodoFile, got %q", cfg.TodoFile)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", dir)

	configDir := filepath.Join(dir, "todome")
	os.MkdirAll(configDir, 0o755)

	content := `todo_file = "~/Dropbox/todo/todo.txt"
done_file = "~/Dropbox/todo/done.txt"
lang = "en"
`
	os.WriteFile(filepath.Join(configDir, "config.toml"), []byte(content), 0o644)

	cfg, err := Load()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.TodoFile != "~/Dropbox/todo/todo.txt" {
		t.Errorf("expected '~/Dropbox/todo/todo.txt', got %q", cfg.TodoFile)
	}
	if cfg.DoneFile != "~/Dropbox/todo/done.txt" {
		t.Errorf("expected '~/Dropbox/todo/done.txt', got %q", cfg.DoneFile)
	}
	if cfg.Lang != "en" {
		t.Errorf("expected 'en', got %q", cfg.Lang)
	}
}

func TestFilePath(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "/custom/config")
	got := FilePath()
	expected := "/custom/config/todome/config.toml"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestFilePathDefault(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", "")
	got := FilePath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".config", "todome", "config.toml")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDataDirPathDefault(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "")
	cfg := DefaultConfig()
	got := cfg.DataDirPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".local", "share", "todome")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDataDirPathXDG(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/custom/data")
	cfg := DefaultConfig()
	got := cfg.DataDirPath()
	if got != "/custom/data/todome" {
		t.Errorf("expected '/custom/data/todome', got %q", got)
	}
}

func TestDataDirPathFromConfig(t *testing.T) {
	t.Setenv("XDG_DATA_HOME", "/should/not/use")
	cfg := Config{DataDir: "~/Dropbox/todome"}
	got := cfg.DataDirPath()
	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "Dropbox", "todome")
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}
