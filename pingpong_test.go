package models

import "testing"

func TestPing(t *testing.T) {
	result := ping()
	expected := "pong"

	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func ping() string {
	return "pong"
}
