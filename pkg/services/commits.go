package services

import (
	"fmt"
	"sort"

	"yokiyoki/pkg/models"
	"yokiyoki/pkg/repository"
)

// CommitsOptions represents configuration for commit list collection
type CommitsOptions struct {
	Period        *Chronometer
	DetailedStats bool
}

// ExecuteCommits fetches commits for the given repository, filters them to the
// configured period, tags each commit with its repository name, and returns the list.
func ExecuteCommits(repo models.Repository, opts CommitsOptions) []models.Commit {
	repoFullName := fmt.Sprintf("%s/%s", repo.Owner, repo.Name)
	commits := repository.GetCommits(repo, opts.Period.StartTime(), opts.DetailedStats)

	filtered := filterCommitsInPeriod(commits, opts.Period)
	for i := range filtered {
		filtered[i].Repository = repoFullName
	}

	return filtered
}

// SortCommitsByDate sorts commits by date in descending order (newest first).
func SortCommitsByDate(commits []models.Commit) {
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date.After(commits[j].Date)
	})
}
