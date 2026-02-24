package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use:   "done <番号> [メッセージ]",
	Short: "タスクにトドメを刺す（完了にする）",
	Args:  cobra.MinimumNArgs(1),
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
		if tasks[idx].Done {
			fmt.Printf("タスク #%d は既にトドメを刺してある\n", num)
			return nil
		}

		tasks[idx].Done = true
		tasks[idx].CompletedAt = time.Now()
		tasks[idx].CompletedHasTime = true

		if len(args) > 1 {
			tasks[idx].Note = strings.Join(args[1:], " ")
		}

		if err := s.Save(tasks); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}

		msg := fmt.Sprintf("🗡️ タスク #%d にトドメを刺した！「%s」", num, tasks[idx].Text)
		if tasks[idx].Note != "" {
			msg += fmt.Sprintf("\n   📝 %s", tasks[idx].Note)
		}
		fmt.Println(msg)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(doneCmd)
}
