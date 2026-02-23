package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <番号> [新しいテキスト]",
	Short: "タスク内容を編集する",
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
		oldText := tasks[idx].Text

		var newText string
		if len(args) > 1 {
			// Inline replacement
			newText = strings.Join(args[1:], " ")
		} else {
			// Open $EDITOR
			editor := os.Getenv("EDITOR")
			if editor == "" {
				return fmt.Errorf("$EDITOR が設定されていない。テキストを引数で指定するか $EDITOR を設定しろ")
			}

			tmpFile, err := os.CreateTemp("", "todome-edit-*.txt")
			if err != nil {
				return fmt.Errorf("一時ファイル作成エラー: %w", err)
			}
			tmpPath := tmpFile.Name()
			defer os.Remove(tmpPath)

			tmpFile.WriteString(oldText)
			tmpFile.Close()

			c := exec.Command(editor, tmpPath)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("エディタ実行エラー: %w", err)
			}

			data, err := os.ReadFile(tmpPath)
			if err != nil {
				return fmt.Errorf("一時ファイル読み込みエラー: %w", err)
			}
			newText = strings.TrimSpace(string(data))
		}

		if newText == "" {
			return fmt.Errorf("空のタスクは設定できない")
		}

		if newText == oldText {
			fmt.Println("変更なし")
			return nil
		}

		tasks[idx].Text = newText
		if err := s.Save(tasks); err != nil {
			return fmt.Errorf("保存エラー: %w", err)
		}

		fmt.Printf("🗡️ タスク #%d を編集した\n  前: %s\n  後: %s\n", num, oldText, newText)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(editCmd)
}
