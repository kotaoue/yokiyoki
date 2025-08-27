package models

import (
	"strings"
)

// UserName represents a normalized username
type UserName struct {
	original   string
	normalized string
}

// NewUserName creates a new UserName with normalization
func NewUserName(username string) UserName {
	return UserName{
		original:   username,
		normalized: normalizeUserName(username),
	}
}

// Original returns the original username
func (u UserName) Original() string {
	return u.original
}

// String returns the normalized username for easy usage
func (u UserName) String() string {
	return u.normalized
}

// Name returns the username based on normalize flag
func (u UserName) Name(normalize bool) string {
	if normalize {
		return u.normalized
	}
	return u.original
}

func normalizeUserName(username string) string {
	normalized := strings.ToLower(strings.ReplaceAll(username, " ", ""))
	return strings.TrimSpace(normalized)
}
