# todome 🗡️

> タスクにトドメを刺せ。

[todo.txt](https://github.com/todotxt/todo.txt) 形式に準拠したターミナル向けタスク管理 CLI ツール。

## インストール

```bash
git clone https://github.com/hengin-eer/todome.git
cd todome
make install    # ビルドして ~/.bin/ にインストール
```

または:

```bash
go install github.com/hengin-eer/todome@latest
```

## 使い方

### タスクを追加する

```bash
todome add 企画書を書く +仕事 @PC
# 🗡️ タスク #1 を追加した: 企画書を書く +仕事 @PC

# 期限付きタスク
todome add 請求書処理 +仕事 due:2026-03-01
```

### タスク一覧を表示する

```bash
todome list          # 未完了タスクのみ
todome list --all    # 完了タスクも含めて表示
todome list --done   # 完了タスクのみ
todome list --overdue # 期限切れタスクのみ
todome ls            # list のエイリアス
```

```
 1. (A) 2026-02-23 企画書を書く +仕事 @PC due:2026-02-20  [期限切れ!] 
 2. 2026-02-23 牛乳を買う +買い物
 3. (B) 2026-02-23 レポート due:2026-02-25  [あと2日] 

3 件のタスク
```

### フィルタリング

```bash
todome list +仕事              # +仕事 プロジェクトのタスク
todome list @PC                # @PC コンテキストのタスク
todome list +仕事 @PC          # AND: +仕事 かつ @PC
todome list --or +仕事 +個人   # OR: +仕事 または +個人
todome list -n +仕事           # NOT: +仕事 を除外
```

### ソート

```bash
todome list -s priority        # 優先度順（A→B→...→なし）
todome list -s created         # 作成日順（新しい順）
todome list -s due             # 期限順（近い順、期限なしは末尾）
todome list -s priority -r     # 逆順
```

### タスクにトドメを刺す（完了）

```bash
todome done 2
# 🗡️ タスク #2 にトドメを刺した！「牛乳を買う +買い物」
```

### 優先度を設定する

```bash
todome pri 1 A       # 優先度 (A) を設定
todome pri 1 none    # 優先度をクリア
```

### タスクを編集する

```bash
todome edit 1 "新しいタスク内容"    # テキストを直接指定
todome edit 1                       # $EDITOR で編集
```

### 完了タスクをアーカイブする

```bash
todome archive
# 🗡️ 2 件の完了タスクをアーカイブした → ~/.local/share/todome/done.txt
```

### タスクを削除する

```bash
todome delete 1        # 確認プロンプトあり
todome delete 1 -f     # 確認なしで削除
todome rm 1            # delete のエイリアス
```

## todo.txt フォーマット

データは [todo.txt 形式](https://github.com/todotxt/todo.txt) のプレーンテキストで保存されます。

```
(A) 2026-02-23 企画書を書く +仕事 @PC due:2026-03-01
x 2026-02-23 2026-02-20 牛乳を買う +買い物
2026-02-23 レポートを提出する
```

| 要素 | 例 | 説明 |
|------|-----|------|
| 完了マーク | `x ` | 行頭に付くと完了タスク |
| 優先度 | `(A)` | A〜Z で優先度を表す |
| 日付 | `2026-02-23` | 作成日（完了時は完了日 + 作成日） |
| プロジェクト | `+仕事` | `+` に続くタグ |
| コンテキスト | `@PC` | `@` に続くタグ |
| 期限 | `due:2026-03-01` | key:value 形式の期限日 |

## 設定

### 初期化

```bash
todome init
# 🗡️ 設定ファイルを作成した: ~/.config/todome/config.toml
```

### 設定ファイル (`~/.config/todome/config.toml`)

```toml
# データディレクトリ（todo.txt, done.txt の保存先）
# デフォルト: ~/.local/share/todome/
# Dropbox/Syncthing で同期する場合はここを変更
# data_dir = "~/Dropbox/todome"

# 個別にファイルパスを指定する場合（data_dir より優先）
# todo_file = "~/Dropbox/todo/todo.txt"
# done_file = "~/Dropbox/todo/done.txt"

# 言語設定（将来用）
# lang = "ja"
```

### ファイルパスの優先順位

| 優先順位 | ソース | 例 |
|:---:|------|-----|
| 1 | `--file` フラグ | `todome --file /tmp/todo.txt list` |
| 2 | `TODOME_FILE` 環境変数 | `export TODOME_FILE=~/todo.txt` |
| 3 | `config.toml` の `todo_file` | `todo_file = "~/custom/todo.txt"` |
| 4 | `config.toml` の `data_dir` | `data_dir = "~/Dropbox/todome"` |
| 5 | XDG デフォルト | `~/.local/share/todome/todo.txt` |

### ファイル配置

```
~/.config/todome/config.toml     # 設定ファイル（todome init で作成）
~/.local/share/todome/todo.txt   # タスクデータ（デフォルト）
~/.local/share/todome/done.txt   # アーカイブデータ（デフォルト）
```

### 複数端末での同期（Dropbox/Syncthing）

`config.toml` で `data_dir` を同期フォルダに向けるだけ:

```toml
data_dir = "~/Dropbox/todome"
```

## プロジェクト構造

```
todome/
├── main.go                  # エントリーポイント
├── Makefile                 # build / test / install / uninstall
├── cmd/
│   ├── root.go              # ルートコマンド定義・設定読み込み
│   ├── add.go               # add サブコマンド
│   ├── list.go              # list サブコマンド
│   ├── done.go              # done サブコマンド
│   ├── delete.go            # delete サブコマンド
│   ├── archive.go           # archive サブコマンド
│   ├── pri.go               # pri サブコマンド
│   ├── edit.go              # edit サブコマンド
│   └── init.go              # init サブコマンド
├── internal/
│   ├── config/
│   │   ├── config.go        # 設定ファイル読み込み
│   │   └── config_test.go
│   ├── todo/
│   │   ├── task.go          # Task 構造体
│   │   ├── parser.go        # todo.txt パーサー・シリアライザー
│   │   └── parser_test.go
│   ├── store/
│   │   ├── store.go         # Store インターフェース
│   │   ├── file.go          # FileStore 実装
│   │   └── file_test.go
│   └── ui/
│       └── display.go       # 色付き表示フォーマット
└── testdata/
```

## 技術スタック

| 項目 | 選定 |
|------|------|
| 言語 | Go |
| CLI フレームワーク | [spf13/cobra](https://github.com/spf13/cobra) |
| 設定ファイル | [BurntSushi/toml](https://github.com/BurntSushi/toml) |
| データ形式 | todo.txt（プレーンテキスト） |
| テスト | 標準 `testing` パッケージ |

## 開発

```bash
make test       # テスト実行
make build      # ビルド
make install    # ビルド & ~/.bin/ にインストール
make uninstall  # アンインストール
```

## ライセンス

MIT
