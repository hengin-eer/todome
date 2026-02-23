package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/hengin-eer/todome/internal/todo"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <タスク内容>",
	Short: "新しいタスクを追加する",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := getStore()
		tasks, err := s.Load()
		if err != nil {
			return fmt.Errorf("読み込みエラー: %w", err)
		}

		text := strings.Join(args, " ")
		task := todo.Parse(text)
		if task.CreatedAt.IsZero() {
			task.CreatedAt = time.Now()
		}

		tasks = append(tasks, task)
		if err := s.Save(tasks); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}

		fmt.Printf("🗡️ タスク #%d を追加した: %s\n", len(tasks), text)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
