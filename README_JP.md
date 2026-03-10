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

モード:
1) メトリクス取得
2) コミット一覧取得
Choice (default 1): 

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
3) JSON
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

# JSON出力
go run . --days 7 --by-user --format json kotaoue/chiken
```

### JSON出力例

```json
[
  {
    "repository": "kotaoue/chiken",
    "user": "kotaoue",
    "commits": 9,
    "lines_added": 602,
    "lines_deleted": 574,
    "prs_created": 4,
    "prs_merged": 4,
    "pr_merge_rate": "4/4 (100%)",
    "avg_pr_merge_time": "0d 00h 21m",
    "issues_created": 1,
    "issues_closed": 0,
    "issue_resolve_rate": "0/1 (0%)",
    "avg_issue_close_time": "-",
    "open_issues": 1
  }
]
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

## コミット一覧モード

モード選択で **2) コミット一覧取得** を選ぶと、コミット日時の降順 (新しい順) でコミット一覧を取得・表示します。

### インタラクティブ例

```bash
$ go run .
GitHub Metrics Collector
========================

Select language / 表示言語を選択してください:
1) English
2) 日本語 (Japanese)
Choice (default 1): 2

モード:
1) メトリクス取得
2) コミット一覧取得
Choice (default 1): 2

リポジトリを入力してください (形式: owner/repo-name)
終了する場合は 'done' と入力:
> kotaoue/chiken
追加: kotaoue/chiken
> done

個別のPRを確認してメトリクスを取得しますか? (処理が遅くなります)
1) Yes
2) No
Choice (default 2): 2


期間:
1) 過去7日間      (2025-08-21 to 2025-08-28 JST)
2) 過去30日間     (2025-07-29 to 2025-08-28 JST)
...
Choice (default 2): 1

出力フォーマット:
1) Markdown
2) CSV
3) JSON
Choice (default 1): 1


Processing repository: kotaoue/chiken
Found 9 commits for kotaoue/chiken
Report
Analyzing data from 2025-08-21 to 2025-08-28 (7 days)

| Repository     | SHA     | Author  | Date             | Message                         |
|----------------|---------|---------|------------------|---------------------------------|
| kotaoue/chiken | abc1234 | kotaoue | 2025-08-28 10:00 | Fix typo in README              |
| kotaoue/chiken | def5678 | kotaoue | 2025-08-27 15:30 | Add new feature                 |
```

### コミット一覧の列

| Column     | Description                                                 |
|------------|-------------------------------------------------------------|
| Repository | リポジトリ名                                                |
| SHA        | コミットハッシュ (先頭7文字)                                |
| Author     | コミット作者                                                |
| Date       | コミット日時 (形式: YYYY-MM-DD HH:mm)                       |
| Message    | コミットメッセージの1行目 (72文字で切り捨て)                |
| Lines +/-  | 追加・削除行数 (--detailed-stats 使用時)                    |
