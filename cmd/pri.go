package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var priCmd = &cobra.Command{
	Use:   "pri <番号> <A-Z|none>",
	Short: "タスクの優先度を設定・変更する",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		num, err := strconv.Atoi(args[0])
		if err != nil || num < 1 {
			return fmt.Errorf("正しいタスク番号を指定しろ: %s", args[0])
		}

		pri := strings.ToUpper(args[1])
		if pri != "NONE" {
			if len(pri) != 1 || pri[0] < 'A' || pri[0] > 'Z' {
				return fmt.Errorf("優先度はA-Zまたはnoneで指定しろ: %s", args[1])
			}
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
		if tasks[idx].Done {
			return fmt.Errorf("完了タスクには優先度を設定できない")
		}

		if pri == "NONE" {
			tasks[idx].Priority = ""
			fmt.Printf("🗡️ タスク #%d の優先度をクリアした「%s」\n", num, tasks[idx].Text)
		} else {
			tasks[idx].Priority = pri
			fmt.Printf("🗡️ タスク #%d の優先度を (%s) に設定した「%s」\n", num, pri, tasks[idx].Text)
		}

		if err := s.Save(tasks); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(priCmd)
}
