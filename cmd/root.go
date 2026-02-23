package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	return store.NewFileStore(path)
}

func defaultTodoPath() string {
	if env := os.Getenv("TODOME_FILE"); env != "" {
		return env
	}
	if appCfg.TodoFile != "" {
		return expandHome(appCfg.TodoFile)
	}
	dir, err := os.Getwd()
	if err != nil {
		return "todo.txt"
	}
	return filepath.Join(dir, "todo.txt")
}

func getDoneFile() string {
	if appCfg.DoneFile != "" {
		return expandHome(appCfg.DoneFile)
	}
	// done.txt in the same directory as todo.txt
	todoPath := defaultTodoPath()
	return filepath.Join(filepath.Dir(todoPath), "done.txt")
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, "エラー:", msg)
	os.Exit(1)
}
