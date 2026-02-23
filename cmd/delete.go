package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var forceDelete bool

var deleteCmd = &cobra.Command{
	Use:     "delete <番号>",
	Short:   "タスクを削除する",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		num, err := strconv.Atoi(args[0])
		if err != nil || num < 1 {
			return fmt.Errorf("正しいタスク番号を指定しろ: %s", args[0])
		}

		s := getStore()
		tasks, err := s.Load()
		if err != nil {
			return fmt.Errorf("読み込みエラー: %w", err)
		}

		if num > len(tasks) {
			return fmt.Errorf("タスク #%d は存在しない（全%d件）", num, len(tasks))
		}

		idx := num - 1
		taskText := tasks[idx].Text

		if !forceDelete {
			fmt.Printf("タスク #%d「%s」を削除するか？ [y/N]: ", num, taskText)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("削除を中止した")
				return nil
			}
		}

		tasks = append(tasks[:idx], tasks[idx+1:]...)
		if err := s.Save(tasks); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}

		fmt.Printf("🗑️ タスク #%d「%s」を削除した\n", num, taskText)
		return nil
	},
}

func init() {
	deleteCmd.Flags().BoolVarP(&forceDelete, "force", "f", false, "確認なしで削除する")
	rootCmd.AddCommand(deleteCmd)
}
