package interactive

import (
	"fmt"
	"strconv"
	"strings"

	"yokiyoki/pkg/locale"
	"yokiyoki/pkg/models"
	"yokiyoki/pkg/services"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Metrics handles interactive user input prompting for metrics configuration
type Metrics struct {
	prompt    *services.Prompter
	localizer *i18n.Localizer
}

// NewMetrics creates a new Metrics instance with the given language ("en", "ja", etc.)
func NewMetrics(lang string) *Metrics {
	return &Metrics{
		prompt:    services.NewPrompter(),
		localizer: locale.NewLocalizer(lang),
	}
}

// t localizes a message by its ID. Falls back to English if the current locale
// does not have the message, and to the raw message ID as a last resort.
func (m *Metrics) t(id string) string {
	msg, err := m.localizer.Localize(&i18n.LocalizeConfig{MessageID: id})
	if err != nil {
		if enMsg, enErr := locale.NewLocalizer("en").Localize(&i18n.LocalizeConfig{MessageID: id}); enErr == nil {
			return enMsg
		}
		return id
	}
	return msg
}

// tWithData localizes a message by its ID with template data. Falls back to English
// if the current locale does not have the message, and to the raw message ID as a last resort.
func (m *Metrics) tWithData(id string, data map[string]interface{}) string {
	cfg := &i18n.LocalizeConfig{MessageID: id, TemplateData: data}
	msg, err := m.localizer.Localize(cfg)
	if err != nil {
		if enMsg, enErr := locale.NewLocalizer("en").Localize(cfg); enErr == nil {
			return enMsg
		}
		return id
	}
	return msg
}

// GetLanguage prompts the user to select a display language and returns the language tag.
// This is a package-level function because it must be called before a localizer is created.
// The prompt is always shown bilingually since no language preference is known yet.
func GetLanguage(prompt *services.Prompter) string {
	en := locale.NewLocalizer("en")
	config := services.SingleChoiceConfig{
		Messages: []string{
			en.MustLocalize(&i18n.LocalizeConfig{MessageID: "SelectLanguage"}),
			en.MustLocalize(&i18n.LocalizeConfig{MessageID: "LanguageEnglish"}),
			en.MustLocalize(&i18n.LocalizeConfig{MessageID: "LanguageJapanese"}),
			en.MustLocalize(&i18n.LocalizeConfig{MessageID: "ChoiceDefault1"}),
		},
		Options: []services.PromptOption{
			{Key: "1", Label: "english", Value: "en"},
			{Key: "2", Label: "japanese", Value: "ja"},
		},
		DefaultKey: "1",
	}
	return prompt.PromptSingleChoice(config).(string)
}

// GetRepositories interactively collects repository information from user input
func (m *Metrics) GetRepositories() []models.Repository {
	config := services.MultipleInputConfig{
		HeaderMessages: []string{
			m.t("RepoInputHeader"),
			m.t("RepoInputDone"),
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
			return m.tWithData("RepoAdded", map[string]interface{}{
				"Name": fmt.Sprintf("%s/%s", repo.Owner, repo.Name),
			})
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
		Message:      m.t("DaysInput"),
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
			m.t("PeriodHeader"),
			fmt.Sprintf("%s%s", m.t("Last7Days"), chronometer.GetLast7DaysDescription()),
			fmt.Sprintf("%s%s", m.t("Last30Days"), chronometer.GetLast30DaysDescription()),
			fmt.Sprintf("%s%s", m.t("LastMonth"), chronometer.GetLastMonthDescription()),
			fmt.Sprintf("%s%s", m.t("PreviousHalf"), chronometer.GetPreviousHalfDescription()),
			fmt.Sprintf("%s%s", m.t("PreviousYear"), chronometer.GetPreviousYearDescription()),
			fmt.Sprintf("%s%s", m.t("PreviousFiscalYear"), chronometer.GetPreviousFiscalYearDescription()),
			m.t("CustomPeriod"),
			m.t("ChoiceDefault2"),
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
		HeaderMessages: []string{m.t("CustomPeriodHeader")},
		ParseFunc: func(input string) (any, error) {
			return input, nil
		},
		DoneKeyword: "done",
		Formatter: func(result any) string {
			return result.(string)
		},
	}

	fmt.Print(m.t("StartDatePrompt"))
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
			m.t("ByUserPrompt"),
			"1) Yes",
			"2) No",
			m.t("ChoiceDefault2"),
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
			m.t("FormatHeader"),
			"1) Markdown",
			"2) CSV",
			m.t("ChoiceDefault1"),
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
			m.t("SortHeader"),
			fmt.Sprintf("1) %s", m.t("SortByRepository")),
			fmt.Sprintf("2) %s", m.t("SortByRepositoryUser")),
			fmt.Sprintf("3) %s", m.t("SortByUserRepository")),
			m.t("ChoiceDefault1"),
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
			m.t("NormalizeUsersPrompt"),
			"1) Yes",
			"2) No",
			m.t("ChoiceDefault2"),
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
			m.t("DetailedStatsPrompt"),
			"1) Yes",
			"2) No",
			m.t("ChoiceDefault2"),
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
