package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/services"
)

// Metrics handles interactive user input prompting for metrics configuration
type Metrics struct {
	prompt *services.Prompter
}

// NewMetrics creates a new Metrics instance
func NewMetrics() *Metrics {
	return &Metrics{
		prompt: services.NewPrompter(),
	}
}

// GetRepositories interactively collects repository information from user input
func (m *Metrics) GetRepositories() []models.Repository {
	config := services.MultipleInputConfig{
		HeaderMessages: []string{
			"\nリポジトリを入力してください (形式: owner/repo-name)",
			"終了する場合は 'done' と入力:",
		},
		ParseFunc: func(input string) (any, error) {
			repo, err := m.parseRepository(input)
			if err != nil {
				return nil, err
			}
			return repo, nil
		},
		DoneKeyword: "done",
		Formatter: func(result any) string {
			repo := result.(models.Repository)
			return fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
		},
	}

	results := m.prompt.PromptMultipleInput(config)

	repos := make([]models.Repository, len(results))
	for i, result := range results {
		repos[i] = result.(models.Repository)
	}

	return repos
}

// GetRepositoriesFromArgs parses repository strings from command line arguments
func (m *Metrics) GetRepositoriesFromArgs(repos []string) ([]models.Repository, error) {
	var result []models.Repository

	for _, repoStr := range repos {
		repo, err := m.parseRepository(repoStr)
		if err != nil {
			return nil, fmt.Errorf("invalid repository format '%s': %v", repoStr, err)
		}
		result = append(result, repo)
	}

	return result, nil
}

// GetDays prompts user for the number of days to analyze
func (m *Metrics) GetDays() int {
	config := services.SingleInputConfig{
		Message:      "\n分析する日数を入力してください (default 30): ",
		DefaultValue: 30,
		Validator: func(input string) (int, error) {
			days, err := strconv.Atoi(input)
			if err != nil || days <= 0 {
				return 0, fmt.Errorf("invalid number")
			}
			return days, nil
		},
	}

	return m.prompt.PromptSingleInput(config)
}

// GetPeriod prompts user to select a time period for analysis
func (m *Metrics) GetPeriod() (int, string, string, error) {
	chronometer, err := services.NewChronometer(services.ChronometerOption{})
	if err != nil {
		return 0, "", "", err
	}

	config := services.SingleChoiceConfig{
		Messages: []string{
			"\n期間:",
			fmt.Sprintf("1) 過去7日間      %s", chronometer.GetLast7DaysDescription()),
			fmt.Sprintf("2) 過去30日間     %s", chronometer.GetLast30DaysDescription()),
			fmt.Sprintf("3) 先月          %s", chronometer.GetLastMonthDescription()),
			fmt.Sprintf("4) 前半期        %s", chronometer.GetPreviousHalfDescription()),
			fmt.Sprintf("5) 前年(1-12月)  %s", chronometer.GetPreviousYearDescription()),
			fmt.Sprintf("6) 前年度(4-3月) %s", chronometer.GetPreviousFiscalYearDescription()),
			"7) カスタム期間",
			"Choice (default 2): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "7days", Value: chronometer.GetLast7DaysResult},
			{Key: "2", Label: "30days", Value: chronometer.GetLast30DaysResult},
			{Key: "3", Label: "lastmonth", Value: chronometer.GetLastMonthResult},
			{Key: "4", Label: "previoushalf", Value: chronometer.GetPreviousHalfResult},
			{Key: "5", Label: "previousyear", Value: chronometer.GetPreviousYearResult},
			{Key: "6", Label: "previousfiscalyear", Value: chronometer.GetPreviousFiscalYearResult},
			{Key: "7", Label: "custom", Value: func() (int, string, string) {
				return m.getCustomDateRange()
			}},
		},
		DefaultKey: "2",
	}

	result := m.prompt.PromptSingleChoice(config)
	fn := result.(func() (int, string, string))
	days, start, end := fn()
	return days, start, end, nil
}

func (m *Metrics) getCustomDateRange() (int, string, string) {
	config := services.MultipleInputConfig{
		HeaderMessages: []string{"カスタム期間を入力してください:"},
		ParseFunc: func(input string) (any, error) {
			return input, nil
		},
		DoneKeyword: "done",
		Formatter: func(result any) string {
			return result.(string)
		},
	}

	fmt.Print("開始日 (YYYY-MM-DD JST, 例: 2024-01-01): ")
	results := m.prompt.PromptMultipleInput(config)

	if len(results) >= 2 {
		startDate := results[0].(string)
		endDate := results[1].(string)
		return 0, startDate, endDate
	}

	return 30, "", ""
}

// GetByUser prompts user whether to break down metrics by user
func (m *Metrics) GetByUser() bool {
	config := services.SingleChoiceConfig{
		Messages: []string{
			"ユーザー別にメトリクスを表示しますか?",
			"1) Yes",
			"2) No",
			"Choice (default 2): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "yes", Value: true},
			{Key: "y", Label: "yes", Value: true},
			{Key: "2", Label: "no", Value: false},
		},
		DefaultKey: "2",
	}
	return m.prompt.PromptSingleChoice(config).(bool)
}

// GetFormat prompts user to select output format (markdown or csv)
func (m *Metrics) GetFormat() string {
	config := services.SingleChoiceConfig{
		Messages: []string{
			"出力フォーマット:",
			"1) Markdown",
			"2) CSV",
			"Choice (default 1): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "markdown", Value: "markdown"},
			{Key: "2", Label: "csv", Value: "csv"},
		},
		DefaultKey: "1",
	}
	return m.prompt.PromptSingleChoice(config).(string)
}

// GetSortBy prompts user to select sort order for output
func (m *Metrics) GetSortBy() string {
	config := services.SingleChoiceConfig{
		Messages: []string{
			"ソート順:",
			"1) リポジトリ",
			"2) リポジトリ,ユーザー",
			"3) ユーザー,リポジトリ",
			"Choice (default 1): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "repository", Value: "repository"},
			{Key: "2", Label: "repository,user", Value: "repository,user"},
			{Key: "3", Label: "user,repository", Value: "user,repository"},
		},
		DefaultKey: "1",
	}
	return m.prompt.PromptSingleChoice(config).(string)
}

// GetNormalizeUsers prompts user whether to normalize usernames
func (m *Metrics) GetNormalizeUsers() bool {
	config := services.SingleChoiceConfig{
		Messages: []string{
			"ユーザー名を正規化しますか ('kotaoue' と 'kota oue' をマージ)?",
			"1) Yes",
			"2) No",
			"Choice (default 2): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "yes", Value: true},
			{Key: "y", Label: "yes", Value: true},
			{Key: "2", Label: "no", Value: false},
		},
		DefaultKey: "2",
	}
	return m.prompt.PromptSingleChoice(config).(bool)
}

// GetDetailedStats prompts user whether to enable detailed statistics collection
func (m *Metrics) GetDetailedStats() bool {
	config := services.SingleChoiceConfig{
		Messages: []string{
			"個別のPRを確認してメトリクスを取得しますか? (処理が遅くなります)",
			"1) Yes",
			"2) No",
			"Choice (default 2): ",
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "yes", Value: true},
			{Key: "y", Label: "yes", Value: true},
			{Key: "2", Label: "no", Value: false},
		},
		DefaultKey: "2",
	}
	return m.prompt.PromptSingleChoice(config).(bool)
}

func (m *Metrics) parseRepository(input string) (models.Repository, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return models.Repository{}, fmt.Errorf("invalid format. Please use: owner/repo-name")
	}

	owner := strings.TrimSpace(parts[0])
	name := strings.TrimSpace(parts[1])

	if owner == "" || name == "" {
		return models.Repository{}, fmt.Errorf("owner and repository name cannot be empty")
	}

	return models.Repository{
		Owner: owner,
		Name:  name,
	}, nil
}
