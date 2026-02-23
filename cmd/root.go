package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hengin-eer/todome/internal/store"
	"github.com/spf13/cobra"
)

var todoFile string

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
	rootCmd.PersistentFlags().StringVar(&todoFile, "file", "", "todo.txtファイルのパス (デフォルト: ./todo.txt)")
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
	dir, err := os.Getwd()
	if err != nil {
		return "todo.txt"
	}
	return filepath.Join(dir, "todo.txt")
}

func exitWithError(msg string) {
	fmt.Fprintln(os.Stderr, "エラー:", msg)
	os.Exit(1)
}
