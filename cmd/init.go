package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hengin-eer/todome/internal/config"
	"github.com/spf13/cobra"
)

const defaultConfigContent = `# todome 設定ファイル
# 詳細: https://github.com/hengin-eer/todome

# データディレクトリ（todo.txt, done.txt の保存先）
# デフォルト: ~/.local/share/todome/
# Dropbox/Syncthing で同期する場合はここを変更
# data_dir = "~/Dropbox/todome"

# 個別にファイルパスを指定する場合（data_dir より優先）
# todo_file = "~/Dropbox/todo/todo.txt"
# done_file = "~/Dropbox/todo/done.txt"

# 言語設定（将来用）
# lang = "ja"
`

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "設定ファイルを初期化する",
	Long:  "~/.config/todome/config.toml を作成する。",
	RunE: func(cmd *cobra.Command, args []string) error {
		path := config.FilePath()
		dir := filepath.Dir(path)

		// Check if file already exists
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("設定ファイルが既に存在する: %s\n上書きするか？ [y/N]: ", path)
			reader := bufio.NewReader(os.Stdin)
			answer, _ := reader.ReadString('\n')
			answer = strings.TrimSpace(strings.ToLower(answer))
			if answer != "y" && answer != "yes" {
				fmt.Println("初期化を中止した")
				return nil
			}
		}

		if err := os.MkdirAll(dir, 0o755); err != nil {
			return fmt.Errorf("ディレクトリ作成エラー: %w", err)
		}

		if err := os.WriteFile(path, []byte(defaultConfigContent), 0o644); err != nil {
			return fmt.Errorf("設定ファイル書き込みエラー: %w", err)
		}

		fmt.Printf("🗡️ 設定ファイルを作成した: %s\n", path)
		fmt.Println("必要に応じて編集しろ！")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
