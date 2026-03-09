package models

import "time"

type Commit struct {
	Repository string    `json:"repository"`
	SHA        string    `json:"sha"`
	Message    string    `json:"message"`
	Author     string    `json:"author"`
	Date       time.Time `json:"date"`
	URL        string    `json:"url"`
	Additions  int       `json:"additions"`
	Deletions  int       `json:"deletions"`
}
