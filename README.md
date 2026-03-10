# yokiyoki

A CLI tool to collect and analyze GitHub metrics for fun.

English | [Japanese](README_JP.md)

## Usage

### Interactive Mode

```bash
$ go run .
GitHub Metrics Collector
========================

Select language / 表示言語を選択してください:
1) English
2) 日本語 (Japanese)
Choice (default 1): 

Mode:
1) Metrics
2) Commit list
Choice (default 1): 

Enter repository (format: owner/repo-name)
Type 'done' to finish:
> kotaoue/chiken
Added: kotaoue/chiken
> kotaoue/gamemo
Added: kotaoue/gamemo
> kotaoue/kota.oue.me
Added: kotaoue/kota.oue.me
> done

Fetch metrics by checking individual PRs? (slower)
1) Yes
2) No
Choice (default 2): 1


Period:
1) Last 7 days    (2025-08-21 to 2025-08-28 JST)
2) Last 30 days   (2025-07-29 to 2025-08-28 JST)
3) Last month     (2025-07-01 to 2025-07-31 JST)
4) First half     (2024-10-01 to 2025-03-31 JST)
5) Last year      (2024-01-01 to 2024-12-31 JST)
6) Last fiscal yr (2024-04-01 to 2025-03-31 JST)
7) Custom range
Choice (default 2): 1

Break down metrics by user?
1) Yes
2) No
Choice (default 2): 1

Output format:
1) Markdown
2) CSV
3) JSON
Choice (default 1):

Sort order:
1) repository
2) repository,user
3) user,repository
Choice (default 1):

Normalize usernames (merge 'kotaoue' and 'kota oue')?
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

### Command-line Mode

```bash
# Same settings as the interactive example above
go run . --days 7 --by-user --normalize-users --detailed-stats --format markdown --sort-by repository kotaoue/chiken kotaoue/gamemo kotaoue/kota.oue.me

# JSON output
go run . --days 7 --by-user --format json kotaoue/chiken
```

### JSON output example

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

## Metrics

| Column               | Description                                                    |
|----------------------|----------------------------------------------------------------|
| Repository           | Repository name                                                |
| User                 | Username (shown when using `--by-user`)                        |
| Commits              | Number of commits                                              |
| PR Merge Rate        | Pull request merge rate (merged / created)                     |
| PR Merge Time        | Average time to merge a pull request (format: 0d 02h 30m)     |
| Issue Resolve Rate   | Issue resolution rate (closed / created)                       |
| Issue Resolve Time   | Average time to resolve an issue (format: 0d 05h 12m)          |
| Active Issues        | Number of currently open issues                                |
| Lines +/-            | Lines added / deleted (shown when using `--detailed-stats`)    |

## Commit List Mode

Select **2) Commit list** at the mode prompt to retrieve commits sorted by date (newest first).

### Interactive example

```bash
$ go run .
GitHub Metrics Collector
========================

Select language / 表示言語を選択してください:
1) English
2) 日本語 (Japanese)
Choice (default 1): 1

Mode:
1) Metrics
2) Commit list
Choice (default 1): 2

Enter repository (format: owner/repo-name)
Type 'done' to finish:
> kotaoue/chiken
Added: kotaoue/chiken
> done

Fetch metrics by checking individual PRs? (slower)
1) Yes
2) No
Choice (default 2): 2


Period:
1) Last 7 days    (2025-08-21 to 2025-08-28 JST)
2) Last 30 days   (2025-07-29 to 2025-08-28 JST)
...
Choice (default 2): 1

Output format:
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

### Commit list columns

| Column     | Description                                                 |
|------------|-------------------------------------------------------------|
| Repository | Repository name                                             |
| SHA        | Short commit hash (7 characters)                            |
| Author     | Commit author                                               |
| Date       | Commit date (JST, format: YYYY-MM-DD HH:mm)                 |
| Message    | First line of the commit message (truncated at 72 chars)    |
| Lines +/-  | Lines added / deleted (shown when using `--detailed-stats`) |
