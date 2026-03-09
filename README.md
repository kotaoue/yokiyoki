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
