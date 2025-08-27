package main

import (
	"fmt"
	"os"

	"yokiyoki/pkg/formatter"
	"yokiyoki/pkg/interactive"
	"yokiyoki/pkg/models"
	"yokiyoki/pkg/services"

	"github.com/spf13/cobra"
)

var (
	days           int
	startDate      string
	endDate        string
	byUser         bool
	format         string
	sortBy         string
	normalizeUsers bool
	detailedStats  bool
)

var rootCmd = &cobra.Command{
	Use:   "yokiyoki [repositories...]",
	Short: "GitHub metrics collector",
	Long: `A CLI tool to collect and analyze GitHub metrics from repositories.

Metrics include:
- Commit frequency and line changes
- Pull request creation and merge statistics  
- Issue creation and resolution statistics
- User-specific breakdowns

Examples:
  yokiyoki                                    # Interactive mode
  yokiyoki owner/repo1 owner/repo2            # From arguments  
  yokiyoki --days 7 --by-user owner/repo      # Last 7 days, by user
  yokiyoki --start 2024-01-01 --end 2024-01-31 owner/repo  # Date range
  yokiyoki --normalize-users --by-user owner/repo  # Merge similar usernames
  yokiyoki --format csv owner/repo            # CSV output
  yokiyoki --sort-by user,repository owner/repo  # Sort by user then repository
  yokiyoki --detailed-stats owner/repo        # Enable detailed line stats (slower)`,
	Run: runCollect,
}

func main() {
	rootCmd.Flags().IntVarP(&days, "days", "d", 30, "Number of days to analyze (default 30)")
	rootCmd.Flags().StringVar(&startDate, "start", "", "Start date (YYYY-MM-DD format, e.g., 2024-01-01)")
	rootCmd.Flags().StringVar(&endDate, "end", "", "End date (YYYY-MM-DD format, e.g., 2024-01-31)")
	rootCmd.Flags().BoolVarP(&byUser, "by-user", "u", false, "Break down metrics by user")
	rootCmd.Flags().StringVarP(&format, "format", "f", "markdown", "Output format: markdown or csv")
	rootCmd.Flags().StringVarP(&sortBy, "sort-by", "s", "repository", "Sort order: repository, repository,user, user,repository")
	rootCmd.Flags().BoolVarP(&normalizeUsers, "normalize-users", "n", false, "Normalize usernames by removing spaces (merge 'kotaoue' and 'kota oue')")
	rootCmd.Flags().BoolVar(&detailedStats, "detailed-stats", false, "Enable detailed line change statistics (requires individual API calls per commit - slower)")

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func runCollect(cmd *cobra.Command, args []string) {
	fmt.Println("GitHub Metrics Collector")
	fmt.Println("========================")

	repos := collectRepositories(cmd, args)
	if len(repos) == 0 {
		fmt.Println("No repositories selected. Exiting.")
		return
	}

	collectMissingOptions(cmd)

	period := createPeriod()
	allMetrics := processRepositories(repos, period)
	outputResults(allMetrics, period)
}

func collectRepositories(cmd *cobra.Command, args []string) []models.Repository {
	metricsInput := interactive.NewMetrics()

	if len(args) > 0 {
		repos, err := metricsInput.GetRepositoriesFromArgs(args)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return nil
		}
		return repos
	}

	return metricsInput.GetRepositories()
}

func collectMissingOptions(cmd *cobra.Command) {
	metricsInput := interactive.NewMetrics()

	if !cmd.Flags().Changed("detailed-stats") {
		detailedStats = metricsInput.GetDetailedStats()
	}

	if !cmd.Flags().Changed("days") && !cmd.Flags().Changed("start") && !cmd.Flags().Changed("end") {
		var err error
		days, startDate, endDate, err = metricsInput.GetPeriod()
		if err != nil {
			fmt.Printf("Error getting period: %v\n", err)
			os.Exit(1)
		}
	}

	if !cmd.Flags().Changed("by-user") {
		byUser = metricsInput.GetByUser()
	}

	if !cmd.Flags().Changed("format") {
		format = metricsInput.GetFormat()
	}

	if !cmd.Flags().Changed("sort-by") {
		sortBy = metricsInput.GetSortBy()
	}

	if !cmd.Flags().Changed("normalize-users") {
		normalizeUsers = metricsInput.GetNormalizeUsers()
	}
}

func createPeriod() *services.Chronometer {
	var opt services.ChronometerOption

	if startDate != "" && endDate != "" {
		opt = services.ChronometerOption{StartDate: &startDate, EndDate: &endDate}
	} else {
		opt = services.ChronometerOption{Days: &days}
	}

	chronometer, err := services.NewChronometer(opt)
	if err != nil {
		fmt.Printf("Error creating chronometer: %v\n", err)
		os.Exit(1)
	}
	return chronometer
}

func processRepositories(repos []models.Repository, period *services.Chronometer) []models.Metrics {
	var allMetrics []models.Metrics

	fmt.Println()
	for _, repo := range repos {
		fmt.Printf("Processing repository: %s/%s\n", repo.Owner, repo.Name)
		options := services.MetricsOptions{
			Period:         period,
			ByUser:         byUser,
			NormalizeUsers: normalizeUsers,
			DetailedStats:  detailedStats,
			SortBy:         sortBy,
		}
		metrics := services.Execute(repo, options)
		allMetrics = append(allMetrics, metrics...)
	}

	return allMetrics
}

func outputResults(allMetrics []models.Metrics, period *services.Chronometer) {
	fmt.Println("Report")
	fmt.Printf("Analyzing data from %s to %s (%d days)\n\n",
		period.StartTime().Format("2006-01-02"),
		period.EndTime().Format("2006-01-02"),
		days)

	if format == "csv" {
		csv := formatter.NewMetricsCsv(allMetrics)
		csv.Output(byUser, detailedStats)
	} else {
		table := formatter.NewMetricsTable(allMetrics)
		table.Output(byUser, detailedStats)
	}
}
