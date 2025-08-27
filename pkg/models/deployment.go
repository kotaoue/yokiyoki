package models

import "time"

type Deployment struct {
	ID          int       `json:"id"`
	Environment string    `json:"environment"`
	State       string    `json:"state"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	SHA         string    `json:"sha"`
}
