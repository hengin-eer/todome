package cmd

import (
	"fmt"

	"github.com/hengin-eer/todome/internal/ui"
	"github.com/spf13/cobra"
)

var listAll bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "タスク一覧を表示する",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		s := getStore()
		tasks, err := s.Load()
		if err != nil {
			return fmt.Errorf("読み込みエラー: %w", err)
		}

		if len(tasks) == 0 {
			fmt.Println("タスクがないぞ。todome add で追加しろ！")
			return nil
		}

		count := 0
		for i, t := range tasks {
			if !listAll && t.Done {
				continue
			}
			fmt.Println(ui.FormatTask(i+1, t))
			count++
		}

		if count == 0 {
			fmt.Println("未完了のタスクはない。全部にトドメを刺した！ 🎉")
		} else {
			fmt.Printf("\n%d 件のタスク\n", count)
		}
		return nil
	},
}

func init() {
	listCmd.Flags().BoolVarP(&listAll, "all", "a", false, "完了タスクも表示する")
	rootCmd.AddCommand(listCmd)
}
