# Phase 4: UX 改善 — 機能アイデア一覧

## 概要

todome の Phase 4「UX 改善」として、既存の ROADMAP にある機能案に加え、CLI の利便性を大幅に向上させるアイデアを幅広くリストアップする。

---

## 既存 ROADMAP の Phase 4 項目（6件）

### 1. `todome edit`（引数なし）で全タスク一括編集
- **利便性**: todo.txt をまとめて整理したいとき、わざわざファイルパスを探さなくていい
- **使い方**: `todome edit` → `$EDITOR` で todo.txt 全体が開き、保存すると反映
- **備考**: 内部的に FileStore.Path をそのまま EDITOR に渡す設計が必要

### 2. 自然言語での期日指定
- **利便性**: `due:2026-03-01` と毎回打つのは面倒。`due:today` `due:tomorrow` `due:fri` で直感的に指定
- **使い方**: `todome add レポート提出 due:today` / `todome add 会議準備 due:+3d` / `todome add 提出 due:fri`
- **対応キーワード案**: `today`, `tomorrow`, `+Nd`（N日後）, `mon`〜`sun`（次の曜日）, `eow`（週末）, `eom`（月末）

### 3. `todome stats` — 統計表示
- **利便性**: 達成感の可視化、習慣化の動機付け。タスク管理のメタ情報を把握
- **使い方**: `todome stats` → 完了数・未完了数・今週の完了数・平均完了時間・プロジェクト別内訳などを表示
- **出力例**:
  ```
  🗡️ 統計
  総タスク: 45 件（完了: 32 / 未完了: 13）
  今週の完了: 8 件
  期限切れ: 2 件
  プロジェクト別: +仕事(5) +個人(3) +買い物(5)
  ```

### 4. `todome undo` — 直前の操作を取り消し
- **利便性**: 誤って done/delete した場合にすぐ戻せる。安全ネット
- **使い方**: `todome undo` → 直前の操作を取り消し
- **実装方針**: 操作前に todo.txt のスナップショットを `~/.local/share/todome/.undo` に保存

### 5. シェル補完スクリプト生成
- **利便性**: Tab 補完でコマンド・フラグを素早く入力。Cobra の組み込み機能で低コスト
- **使い方**: `todome completion bash > ~/.bash_completion.d/todome` / `todome completion zsh > ~/.zfunc/_todome`
- **備考**: Cobra が `completion` サブコマンドを自動生成する機能あり

### 6. カラーテーマ設定
- **利便性**: ターミナルの背景色に合わせて見やすい配色に変更可能
- **使い方**: `config.toml` に `[theme]` セクションを追加:
  ```toml
  [theme]
  priority_a = "red"
  priority_b = "yellow"
  done = "gray"
  overdue = "red_bg"
  ```

---

## 追加アイデア: 入力・操作の効率化（10件）

### 7. バッチ操作（複数タスク一括処理）
- **利便性**: `done 1` `done 2` `done 3` と3回打つ代わりに1回で完了。delete/pri も同様
- **使い方**: `todome done 1 3 5` / `todome delete 2 4 -f` / `todome pri 1 3 5 A`

### 8. コマンドエイリアス短縮形
- **利便性**: 頻繁に使うコマンドを最短で打てる。`a` で add、`l` で list、`d` で done
- **使い方**: `todome a 牛乳買う` / `todome l +仕事` / `todome d 3`
- **備考**: Cobra の Aliases 機能で既存の `ls`/`rm` と同様に追加可能

### 9. `todome add` のインタラクティブモード
- **利便性**: 引数なしで `add` すると対話的にプロジェクト・コンテキスト・期限を聞いてくれる
- **使い方**: `todome add` → プロンプトでタスク名・プロジェクト・期限を順に入力
- **備考**: 初心者にやさしい。既存の引数指定と併存

### 10. `todome add` のテンプレート/プリセット
- **利便性**: よく使うタスクパターンを登録して呼び出せる
- **使い方**: config に `[templates]` 定義 → `todome add --template weekly-review` で展開
  ```toml
  [templates]
  weekly-review = "(B) 週次レビュー +仕事 @PC"
  shopping = "+買い物 @外出"
  ```

### 11. `todome add` でクリップボードから追加
- **利便性**: Web で見つけた情報をコピー → すぐタスク化
- **使い方**: `todome add --clipboard` / `todome add -c`
- **備考**: `xclip`/`pbpaste`/`wl-paste` を検出して利用

### 12. 出力フォーマット切り替え（JSON / CSV / TSV）
- **利便性**: 他ツール・スクリプトとの連携。`jq` でフィルタ、スプレッドシートへ貼り付け
- **使い方**: `todome list --format json` / `todome list --format csv` / `todome list -o json`
- **出力例**: `[{"id":1,"text":"企画書を書く","priority":"A","projects":["+仕事"],...}]`

### 13. `--quiet` / `--verbose` モード
- **利便性**: スクリプト内で使うとき余計な出力を抑える / デバッグ時に詳細表示
- **使い方**: `todome add タスク -q` → ID のみ出力（`1`）/ `todome list -v` → 作成日時も表示

### 14. 連続追加モード
- **利便性**: 買い物リストなど、複数タスクを一気に追加したいとき便利
- **使い方**: `todome add --bulk` → 1行1タスクで入力、空行で終了。または `todome add -` でパイプ入力対応
  ```bash
  echo -e "牛乳\n卵\nパン" | todome add --bulk +買い物
  ```

### 15. タスク番号の代わりにキーワード指定
- **利便性**: 番号を覚えなくても文字列の一部でタスクを特定できる
- **使い方**: `todome done "牛乳"` → 「牛乳を買う」を完了（曖昧一致は確認プロンプト）

### 16. `todome move` / `todome swap` — タスク順序の入れ替え
- **利便性**: 表示順を手動で制御。重要タスクを上に持ってくる
- **使い方**: `todome move 5 1` → タスク#5を#1の位置に移動 / `todome swap 1 3`

---

## 追加アイデア: 表示・ナビゲーション（10件）

### 17. `todome today` — 今日のタスク表示
- **利便性**: 毎朝これだけ打てば今日やることがわかる
- **使い方**: `todome today` → 期限が今日以前の未完了タスク + 優先度 A のタスクを表示
- **備考**: `list --overdue` との違いは「今日期限」も含む点

### 18. `todome week` — 今週のタスク表示
- **利便性**: 週間計画の確認に使える
- **使い方**: `todome week` → 今週期限のタスクを曜日ごとにグループ化して表示

### 19. `todome next` — 次にやるべきタスク1件表示
- **利便性**: 迷わず着手できる。優先度→期限→作成日の順で最重要タスクを1件表示
- **使い方**: `todome next` → `🗡️ 次のターゲット: (A) 企画書を書く +仕事 @PC`

### 20. `todome projects` / `todome contexts` — タグ一覧表示
- **利便性**: どんなプロジェクト・コンテキストがあるか把握。タグのタイポ発見にも
- **使い方**: `todome projects` → `+仕事 (5件) / +個人 (3件) / +買い物 (2件)`

### 21. グループ表示（プロジェクト別・コンテキスト別）
- **利便性**: 場所やプロジェクトごとにタスクを俯瞰できる
- **使い方**: `todome list --group project` / `todome list -g context`
  ```
  +仕事:
    1. (A) 企画書を書く @PC
    3. レポート提出
  +買い物:
    2. 牛乳を買う
  ```

### 22. タスク件数のカウント表示
- **利便性**: フィルタ結果の件数だけ知りたいとき。スクリプト連携にも便利
- **使い方**: `todome count` / `todome count +仕事` / `todome count --overdue`
- **備考**: `todome list +仕事 | wc -l` の代替

### 23. カレンダービュー（月間/週間）
- **利便性**: 期限の分布を視覚的に把握。タスクの集中日がわかる
- **使い方**: `todome calendar` / `todome cal`
  ```
  2026年3月
  月  火  水  木  金  土  日
   1   2   3   4   5   6   7
  [2]      [1]          [3]
  ```

### 24. `todome summary` — 1日の振り返りサマリー
- **利便性**: 日報作成の材料に。今日完了したタスクと残タスクを表示
- **使い方**: `todome summary` → 今日完了したタスク一覧 + 残タスク数

### 25. ページネーション / `--limit` フラグ
- **利便性**: タスクが100件超えたとき、画面に収まる量だけ表示
- **使い方**: `todome list --limit 10` / `todome list -l 10 --offset 20`

### 26. `todome show <番号>` — タスク詳細表示
- **利便性**: 1件のタスクの全メタデータ（作成日時、期限、プロジェクト、ノート等）を見やすく表示
- **使い方**: `todome show 3` →
  ```
  タスク #3
  テキスト: レポートを提出する
  優先度:   (B)
  作成日時: 2026-02-23 14:30
  期限:     2026-02-25
  プロジェクト: +仕事
  コンテキスト: @PC
  状態:     未完了
  ```

---

## 追加アイデア: ワークフロー・生産性向上（10件）

### 27. `todome focus <番号>` — フォーカスモード
- **利便性**: 1つのタスクに集中するための表示。他タスクを一時非表示に
- **使い方**: `todome focus 3` → タスク #3 のみ表示 + 経過時間カウント

### 28. 繰り返しタスク（`rec:daily` / `rec:weekly`）
- **利便性**: 毎日・毎週のルーティンを自動再生成
- **使い方**: `todome add 筋トレ +習慣 rec:daily` → done 後に翌日の期限で新タスク自動追加
- **キーワード**: `rec:daily`, `rec:weekly`, `rec:monthly`, `rec:weekdays`

### 29. `todome review` — GTD スタイルレビュー
- **利便性**: 全タスクを1件ずつ確認して整理する対話的フロー
- **使い方**: `todome review` → タスクを1件ずつ表示し、`[d]one / [s]kip / [e]dit / [p]ri / [D]elete` を選択

### 30. `todome waiting` — 待機中タスク表示
- **利便性**: 他人の返事待ちなど、自分では進められないタスクを分離管理
- **使い方**: `todome add 上司の承認待ち @waiting` → `todome waiting` で `@waiting` タスク一覧
- **備考**: 既存のコンテキストフィルタ `@waiting` の省略コマンド

### 31. タスクへのメモ追加（完了時以外も）
- **利便性**: 進捗メモ・参考 URL など、タスクに補足情報を付けたい
- **使い方**: `todome note 3 "参考: https://example.com"` → `memo:` タグで保存
- **備考**: 現在 `note:` は完了時のみ。拡張して任意タイミングで付与可能に

### 32. タスクの開始マーク（`todome start`）
- **利便性**: 「作業中」状態を表現。GTD の "Next Action" 的な使い方
- **使い方**: `todome start 3` → テキストに `started:2026-03-01` を付与 / `todome list --started`

### 33. 延期コマンド（`todome defer` / `todome snooze`）
- **利便性**: 期限を手軽にずらせる。edit で日付を書き換えるより直感的
- **使い方**: `todome defer 3 +3d` → 期限を3日後に変更 / `todome defer 3 tomorrow`

### 34. タスクの複製（`todome copy` / `todome dup`）
- **利便性**: 似たタスクを少し変えて追加したいとき便利
- **使い方**: `todome copy 3` → タスク #3 のコピーを末尾に追加（新しい作成日時で）

### 35. `todome inbox` — 素早いキャプチャ
- **利便性**: プロジェクトやコンテキストを考えず、とりあえずメモ。後で整理
- **使い方**: `todome inbox ふと思いついたこと` → プロジェクトなし・優先度なしで即追加
- **備考**: `todome add` との違いは「整理前」を明示する点。`+inbox` タグ自動付与

### 36. `todome delegate <番号> <名前>` — 委任マーク
- **利便性**: 誰に頼んだかを記録。`@delegated` + `to:名前` で追跡
- **使い方**: `todome delegate 3 田中` → `@delegated to:田中` を付与

---

## 追加アイデア: システム連携・拡張性（8件）

### 37. Git-style エイリアス設定
- **利便性**: 自分だけのショートカットを config に定義できる
- **使い方**: config.toml に定義:
  ```toml
  [aliases]
  t = "today"
  w = "week"
  s = "stats"
  lp = "list -s priority"
  ```

### 38. フック機能（pre/post コマンド実行）
- **利便性**: タスク完了時にスクリプト実行（通知送信、ログ記録など）
- **使い方**: config.toml に定義:
  ```toml
  [hooks]
  post_done = "notify-send 'タスク完了！'"
  post_add = "echo 'Added' >> ~/task-log.txt"
  ```

### 39. Markdown / HTML エクスポート
- **利便性**: 日報・週報の素材として、タスク一覧をフォーマット出力
- **使い方**: `todome export --format md` / `todome export --format html --done --since 2026-02-24`

### 40. 他フォーマットからのインポート
- **利便性**: 既存のタスクリストを todome に移行
- **使い方**: `todome import --from markdown tasks.md` / `todome import --from csv tasks.csv`

### 41. `todome version` コマンド
- **利便性**: バグ報告時にバージョン確認。ビルド情報の表示
- **使い方**: `todome version` → `todome v0.4.0 (go1.22, linux/amd64, commit: abc1234)`
- **備考**: `go build -ldflags` でビルド時に埋め込み

### 42. `todome doctor` — 環境チェック
- **利便性**: 設定ファイルの問題や todo.txt の整合性を診断
- **使い方**: `todome doctor` → 設定パス確認・ファイル存在チェック・パースエラー検出
  ```
  ✓ config: ~/.config/todome/config.toml
  ✓ todo.txt: ~/.local/share/todome/todo.txt (13 tasks)
  ✓ done.txt: ~/.local/share/todome/done.txt (32 tasks)
  ⚠ 行 7: パース警告 — 不明なタグ "prio:A"
  ```

### 43. マンページ生成
- **利便性**: `man todome` で本格的なドキュメント参照。Cobra の doc 生成機能で自動化
- **使い方**: `todome man` でマンページ生成 / Makefile に `make man` ターゲット追加

### 44. パイプ入出力の完全対応
- **利便性**: Unix 哲学に沿った他コマンドとの連携
- **使い方**:
  ```bash
  todome list --format plain | grep "仕事" | wc -l
  cat tasks.txt | todome add --bulk
  todome list --format json | jq '.[] | select(.priority == "A")'
  ```

---

## 追加アイデア: データ安全性・信頼性（4件）

### 45. アトミック書き込み（一時ファイル + rename）
- **利便性**: 書き込み中のクラッシュでデータ消失を防ぐ
- **使い方**: 内部実装の改善。ユーザーには透過的
- **備考**: 現在の FileStore.Save() は直接上書き。temp → rename パターンに変更

### 46. `todome validate` — データ整合性チェック
- **利便性**: todo.txt のフォーマットエラーを検出・修正提案
- **使い方**: `todome validate` → パースエラー行を報告 / `todome validate --fix` で自動修正

### 47. 自動バックアップ（操作前スナップショット）
- **利便性**: undo の基盤にもなる。破壊的操作前に自動保存
- **使い方**: 内部実装。`~/.local/share/todome/.backup/` に世代管理で保存
- **備考**: `todome undo` と連携

### 48. 操作ログ（audit trail）
- **利便性**: いつ何をしたか振り返れる。stats の精度も向上
- **使い方**: `todome log` → 操作履歴を表示
  ```
  2026-03-01 09:00 add "企画書を書く +仕事"
  2026-03-01 10:30 done 1
  2026-03-01 10:31 undo
  ```

---

## 追加アイデア: モチベーション・ゲーミフィケーション（5件）

### 49. 完了時のランダム褒め言葉
- **利便性**: 完了時に毎回同じメッセージだと飽きる。バリエーションで達成感アップ
- **使い方**: `todome done 3` → ランダムで表示:
  - `🗡️ 見事にトドメを刺した！`
  - `🗡️ 一刀両断！お見事！`
  - `🗡️ ターゲット撃破！残り 5 件`

### 50. ストリーク表示（連続達成日数）
- **利便性**: 毎日タスクを完了する習慣づけ
- **使い方**: `todome stats` 内に `🔥 連続 7 日タスク完了中！` と表示

### 51. 完了数に応じた称号・ランク
- **利便性**: 長期的なモチベーション維持
- **使い方**: `todome stats` に表示: `称号: 見習い剣士（50件完了）→ 次: 熟練剣士（100件）`
- **備考**: todome の「トドメを刺す」テーマに合わせた称号体系

### 52. 今日の完了数プログレスバー
- **利便性**: 今日の目標タスク数に対する進捗を可視化
- **使い方**: `todome today` の出力に追加: `進捗: ████████░░ 8/10 (80%)`
- **備考**: 目標数は config で設定 `daily_goal = 10`

### 53. `todome celebrate` — 全タスク完了時の演出
- **利便性**: 全タスク完了のご褒美。小さな楽しみ
- **使い方**: 自動検出 or `todome celebrate` で ASCII アート表示:
  ```
  🗡️✨ 全てのタスクにトドメを刺した！ ✨🗡️
       今日の戦い、見事な勝利だ！
  ```

---

## 実装優先度の提案

| 優先度 | アイデア | 理由 |
|:---:|---|---|
| ★★★ | #2 自然言語期日, #5 シェル補完, #7 バッチ操作 | 日常の入力効率が劇的に向上 |
| ★★★ | #41 version, #45 アトミック書き込み | 基本品質・安全性 |
| ★★☆ | #1 全タスク編集, #4 undo, #17 today | 頻出ユースケースのショートカット |
| ★★☆ | #3 stats, #19 next, #20 projects/contexts | タスク俯瞰・意思決定支援 |
| ★★☆ | #12 JSON出力, #49 ランダム褒め言葉 | 連携性・使っていて楽しい |
| ★☆☆ | #28 繰り返し, #29 review, #37 エイリアス | パワーユーザー向け |
| ★☆☆ | 残りすべて | 必要に応じて段階的に |
