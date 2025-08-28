package models

// Metrics represents GitHub repository metrics
type Metrics struct {
	Repository        string
	User              string // "" for repository-wide metrics
	Commits           int
	LinesAdded        int
	LinesDeleted      int
	PRsCreated        int
	PRsMerged         int
	PRMergeRate       string
	AvgPRMergeTime    string
	IssuesCreated     int
	IssuesClosed      int
	IssueResolveRate  string
	AvgIssueCloseTime string
	OpenIssues        int
}
