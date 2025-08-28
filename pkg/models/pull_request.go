package models

import "time"

type PullRequest struct {
	Number    int        `json:"number"`
	Title     string     `json:"title"`
	State     string     `json:"state"`
	Author    string     `json:"author"`
	CreatedAt time.Time  `json:"created_at"`
	MergedAt  *time.Time `json:"merged_at"`
	ClosedAt  *time.Time `json:"closed_at"`
	URL       string     `json:"url"`
	Additions int        `json:"additions"`
	Deletions int        `json:"deletions"`
}
