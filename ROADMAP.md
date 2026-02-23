# todome ロードマップ

## ゴール

- 個人利用ツールとして完成度を上げる
- 複数端末でファイル同期（Dropbox/Syncthing 等）して使えるようにする
- 友人に紹介・配布できる程度の手軽さを持たせる

---

## Phase 0: プロジェクト基盤整備

> Git管理・開発環境を整える。コードを書く前の土台。

- [ ] `.gitignore` 作成（バイナリ、todo.txt、エディタファイル除外）
- [ ] `git init` & 初回コミット
- [ ] `LICENSE` ファイル追加（MIT）
- [ ] エディタ設定（`.editorconfig`）

## Phase 1: 日常利用の品質向上

> 毎日使うツールとしてのストレスを減らす。

- [ ] `~/.config/todome/config.toml` 設定ファイル対応
  - デフォルトの `todo.txt` パス指定
  - `done.txt` パス指定
  - 言語設定（日本語/英語メッセージ切り替え、将来用）
- [ ] `todome archive` コマンド — 完了タスクを `done.txt` に移動
- [ ] `todome pri <番号> <A-Z>` — 優先度の設定・変更
- [ ] `todome edit <番号>` — タスク内容のインライン編集（`$EDITOR` 起動 or 引数置換）

## Phase 2: 検索・フィルタリング

> タスクが増えても目的のタスクにすぐたどり着ける。

- [ ] `todome list +project` — プロジェクトでフィルタ
- [ ] `todome list @context` — コンテキストでフィルタ
- [ ] `todome list -s priority` — 優先度順ソート
- [ ] `todome list --done` / `--undone` — 状態フィルタ（`--all` の細分化）

## Phase 3: 複数端末対応

> Dropbox/Syncthing でファイル同期した環境で安全に使える。

- [ ] 設定ファイルで `todo.txt` パスを `~/Dropbox/todo/todo.txt` 等に指定可能にする（Phase 1 で対応済み）
- [ ] ファイルロック機構 — 同時書き込み防止（`flock` ベース）
- [ ] `todome backup` コマンド — タイムスタンプ付きバックアップ作成

## Phase 4: UX 改善

> 使い心地を磨く。

- [ ] シェル補完スクリプト生成（Cobra 組み込み: `todome completion bash/zsh/fish`）
- [ ] カラーテーマ設定（設定ファイルで色変更可能に）
- [ ] `todome stats` — 統計表示（完了数、今週の進捗など）
- [ ] `todome undo` — 直前の操作を取り消し

## Phase 5: 配布

> 友人に「`brew install` で入るよ」と言えるようにする。

- [ ] GitHub Actions で CI（テスト・ビルド自動実行）
- [ ] GoReleaser でクロスコンパイル & GitHub Releases に自動公開
- [ ] Homebrew tap リポジトリ作成（`homebrew-tap`）
- [ ] README にインストール手順・バッジ追加

### セキュリティ上の注意（初心者向け）

GoReleaser + GitHub Releases による配布は安全な方法です:
- ソースコードは GitHub 上で公開されているので中身が見える
- バイナリは GitHub Actions 上でビルドされるので改ざんリスクが低い
- Homebrew tap は自分のリポジトリから配布するだけで、公式 Homebrew に登録する必要はない
- **やってはいけないこと**: バイナリに API キーや個人情報をハードコードしない

---

## 実装しない（スコープ外）

- クラウド同期（自前サーバー/API）— ファイル同期で十分
- GUI / TUI（ncurses 的な画面）— ターミナル CLI に徹する
- マルチユーザー機能 — 個人利用ツール
