package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hengin-eer/todome/internal/config"
	"github.com/hengin-eer/todome/internal/store"
	"github.com/spf13/cobra"
)

var (
	todoFile string
	appCfg   config.Config
)

var rootCmd = &cobra.Command{
	Use:   "todome",
	Short: "todome — タスクにトドメを刺せ 🗡️",
	Long:  "todo.txt形式のタスク管理CLI。タスクを追加し、片付け、トドメを刺せ。",
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(loadConfig)
	rootCmd.PersistentFlags().StringVar(&todoFile, "file", "", "todo.txtファイルのパス (デフォルト: ./todo.txt)")
}

func loadConfig() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "警告: 設定ファイル読み込みエラー: %v\n", err)
		cfg = config.DefaultConfig()
	}
	appCfg = cfg
}

func getStore() *store.FileStore {
	path := todoFile
	if path == "" {
		path = defaultTodoPath()
	}
	// Ensure directory exists
	if dir := filepath.Dir(path); dir != "." {
		os.MkdirAll(dir, 0o755)
	}
	return store.NewFileStore(path)
}

func defaultTodoPath() string {
	if env := os.Getenv("TODOME_FILE"); env != "" {
		return env
	}
	if appCfg.TodoFile != "" {
		return config.ExpandHome(appCfg.TodoFile)
	}
	return filepath.Join(appCfg.DataDirPath(), "todo.txt")
}

func getDoneFile() string {
	if appCfg.DoneFile != "" {
		return config.ExpandHome(appCfg.DoneFile)
	}
	return filepath.Join(appCfg.DataDirPath(), "done.txt")
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, "エラー:", msg)
	os.Exit(1)
}
