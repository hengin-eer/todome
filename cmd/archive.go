package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/hengin-eer/todome/internal/todo"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive",
	Short: "完了タスクをdone.txtに移動する",
	RunE: func(cmd *cobra.Command, args []string) error {
		s := getStore()
		tasks, err := s.Load()
		if err != nil {
			return fmt.Errorf("読み込みエラー: %w", err)
		}

		var done []todo.Task
		var remaining []todo.Task
		for _, t := range tasks {
			if t.Done {
				done = append(done, t)
			} else {
				remaining = append(remaining, t)
			}
		}

		if len(done) == 0 {
			fmt.Println("アーカイブするタスクがない")
			return nil
		}

		// Append done tasks to done.txt
		donePath := getDoneFile()
		f, err := os.OpenFile(donePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
		if err != nil {
			return fmt.Errorf("done.txtを開けない: %w", err)
		}
		defer f.Close()

		w := bufio.NewWriter(f)
		for _, t := range done {
			w.WriteString(todo.Serialize(t))
			w.WriteString("\n")
		}
		if err := w.Flush(); err != nil {
			return fmt.Errorf("done.txt書き込みエラー: %w", err)
		}

		// Save remaining tasks
		if err := s.Save(remaining); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}

		fmt.Printf("🗡️ %d 件の完了タスクをアーカイブした → %s\n", len(done), donePath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(archiveCmd)
}
