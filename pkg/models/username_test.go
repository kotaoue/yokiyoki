package models_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"yokiyoki/pkg/models"
)

func TestUserName_Name(t *testing.T) {
	tests := []struct {
		name      string
		username  string
		normalize bool
		want      string
	}{
		{
			name:      "normalize true returns normalized",
			username:  "Kota Oue",
			normalize: true,
			want:      "kotaoue",
		},
		{
			name:      "normalize false returns original",
			username:  "Kota Oue",
			normalize: false,
			want:      "Kota Oue",
		},
		{
			name:      "normalize true with spaces",
			username:  "  KOTA  OUE  ",
			normalize: true,
			want:      "kotaoue",
		},
		{
			name:      "normalize false with spaces",
			username:  "  KOTA  OUE  ",
			normalize: false,
			want:      "  KOTA  OUE  ",
		},
		{
			name:      "empty string normalize true",
			username:  "",
			normalize: true,
			want:      "",
		},
		{
			name:      "empty string normalize false",
			username:  "",
			normalize: false,
			want:      "",
		},
		{
			name:      "whitespace only normalize true",
			username:  "   ",
			normalize: true,
			want:      "",
		},
		{
			name:      "whitespace only normalize false",
			username:  "   ",
			normalize: false,
			want:      "   ",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			username := models.NewUserName(tt.username)
			got := username.Name(tt.normalize)
			assert.Equal(t, tt.want, got)
		})
	}
}
