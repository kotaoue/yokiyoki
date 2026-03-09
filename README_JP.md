# yokiyoki
GitHubの情報を取得してあれこれ分析して楽しむためのコマンド

[English](README.md) | Japanese

## Usage

### インタラクティブモード

```bash
$ go run .     
GitHub Metrics Collector
========================

Select language / 表示言語を選択してください:
1) English
2) 日本語 (Japanese)
Choice (default 1): 2

リポジトリを入力してください (形式: owner/repo-name)
終了する場合は 'done' と入力:
> kotaoue/chiken
追加: kotaoue/chiken
> kotaoue/gamemo
追加: kotaoue/gamemo
> kotaoue/kota.oue.me
追加: kotaoue/kota.oue.me
> done

個別のPRを確認してメトリクスを取得しますか? (処理が遅くなります)
1) Yes
2) No
Choice (default 2): 1


期間:
1) 過去7日間      (2025-08-21 to 2025-08-28 JST)
2) 過去30日間     (2025-07-29 to 2025-08-28 JST)
3) 先月          (2025-07-01 to 2025-07-31 JST)
4) 前半期        (2024-10-01 to 2025-03-31 JST)
5) 前年(1-12月)  (2024-01-01 to 2024-12-31 JST)
6) 前年度(4-3月) (2024-04-01 to 2025-03-31 JST)
7) カスタム期間
Choice (default 2): 1

ユーザー別にメトリクスを表示しますか?
1) Yes
2) No
Choice (default 2): 1

出力フォーマット:
1) Markdown
2) CSV
Choice (default 1):  

ソート順:
1) リポジトリ
2) リポジトリ,ユーザー
3) ユーザー,リポジトリ
Choice (default 1): 

ユーザー名を正規化しますか ('kotaoue' と 'kota oue' をマージ)?
1) Yes
2) No
Choice (default 2): 1


Processing repository: kotaoue/chiken
Found 9 commits for kotaoue/chiken
Found 4 pull requests for kotaoue/chiken
Found 5 issues for kotaoue/chiken
Processing repository: kotaoue/gamemo
Found 0 commits for kotaoue/gamemo
Found 0 pull requests for kotaoue/gamemo
Found 0 issues for kotaoue/gamemo
Processing repository: kotaoue/kota.oue.me
Found 2 commits for kotaoue/kota.oue.me
Found 1 pull requests for kotaoue/kota.oue.me
Found 4 issues for kotaoue/kota.oue.me
Report
Analyzing data from 2025-08-21 to 2025-08-28 (7 days)

| Repository          | User    | Commits | PR Merge Rate | PR Merge Time | Issue Resolve Rate | Issue Resolve Time | Active Issues | Lines +/- |
|---------------------|---------|---------|---------------|---------------|--------------------|--------------------|---------------|-----------|
| kotaoue/chiken      | kotaoue |       9 | 4/4 (100%)    | 0d 00h 21m    | 0/1 (0%)           | -                  |             1 | +602/-574 |
| kotaoue/gamemo      | -       |       0 | -/-           | -             | -/-                | -                  |             0 | +0/-0     |
| kotaoue/kota.oue.me | kotaoue |       2 | 0/1 (0%)      | -             | 0/1 (0%)           | -                  |             1 | +6/-2     |
```

### コマンドラインモード

```bash
# インタラクティブモードと同じ設定
go run . --days 7 --by-user --normalize-users --detailed-stats --format markdown --sort-by repository kotaoue/chiken kotaoue/gamemo kotaoue/kota.oue.me
```

## メトリクス

| Column               | Description                                           |
|----------------------|-------------------------------------------------------|
| Repository           | リポジトリ名                                          |
| User                 | ユーザー名 (--by-user 使用時)                        |
| Commits              | コミット数                                            |
| PR Merge Rate        | プルリクエストのマージ率 (マージ数/作成数)            |
| PR Merge Time        | プルリクエストの平均マージ時間 (形式: 0d 02h 30m)    |
| Issue Resolve Rate   | イシューの解決率 (クローズ数/作成数)                  |
| Issue Resolve Time   | イシューの平均解決時間 (形式: 0d 05h 12m)            |
| Active Issues        | 現在のオープンイシュー数                              |
| Lines +/-            | 追加・削除行数 (--detailed-stats 使用時)             |

---

## conversation サブコマンド

指定した期間内のPRおよびIssueのコメントスレッドを取得・表示します。

### 使い方

```bash
go run . conversation owner/repo
go run . conversation --days 7 owner/repo1 owner/repo2
go run . conversation --start 2024-01-01 --end 2024-01-31 owner/repo
```

### オプション

| フラグ        | 説明                                                  |
|---------------|-------------------------------------------------------|
| `-d, --days`  | 分析する日数 (デフォルト 30)                          |
| `--start`     | 開始日 (YYYY-MM-DD 形式、例: 2024-01-01)             |
| `--end`       | 終了日 (YYYY-MM-DD 形式、例: 2024-01-31)             |

### 出力例

```
kotaoue/chiken (3 conversations)
---------------------------------

PR #42 [merged] by alice - Fix login bug (2024-01-10)
  bob (2024-01-11 09:30): LGTM, nice cleanup
  alice (2024-01-11 10:00): Thanks, merging

PR #43 [open] by carol - Add dark mode (2024-01-12)
  (no comments)

Issue #57 [closed] by bob - Button misaligned on mobile (2024-01-08)
  alice (2024-01-09 14:22): Confirmed, will fix in next PR
  bob (2024-01-10 11:05): Fixed in #42, closing
```
