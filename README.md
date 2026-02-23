# todome 🗡️

> タスクにトドメを刺せ。

[todo.txt](https://github.com/todotxt/todo.txt) 形式に準拠したターミナル向けタスク管理 CLI ツール。

## インストール

```bash
go install github.com/hengin-eer/todome@latest
```

または、リポジトリをクローンしてビルド:

```bash
git clone https://github.com/hengin-eer/todome.git
cd todome
go build -o todome .
```

## 使い方

### タスクを追加する

```bash
todome add 企画書を書く +仕事 @PC
# 🗡️ タスク #1 を追加した: 企画書を書く +仕事 @PC
```

### タスク一覧を表示する

```bash
todome list          # 未完了タスクのみ
todome list --all    # 完了タスクも含めて表示
todome ls            # list のエイリアス
```

```
 1. 2026-02-23 企画書を書く +仕事 @PC
 2. 2026-02-23 牛乳を買う +買い物
 3. 2026-02-23 レポートを提出する

3 件のタスク
```

### タスクにトドメを刺す（完了）

```bash
todome done 2
# 🗡️ タスク #2 にトドメを刺した！「牛乳を買う +買い物」
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
(A) 2026-02-23 企画書を書く +仕事 @PC
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

## 設定

### ファイルパス

デフォルトではカレントディレクトリの `todo.txt` を使用します。変更するには:

```bash
# 環境変数で指定
export TODOME_FILE=~/todo.txt

# または --file フラグで指定
todome --file ~/todo.txt list
```

## プロジェクト構造

```
todome/
├── main.go                  # エントリーポイント
├── cmd/
│   ├── root.go              # ルートコマンド定義・Store取得
│   ├── add.go               # add サブコマンド
│   ├── list.go              # list サブコマンド
│   ├── done.go              # done サブコマンド
│   └── delete.go            # delete サブコマンド
├── internal/
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
| データ形式 | todo.txt（プレーンテキスト） |
| テスト | 標準 `testing` パッケージ |

## 開発

```bash
# テスト実行
go test ./...

# ビルド
go build -o todome .
```

## ライセンス

MIT
