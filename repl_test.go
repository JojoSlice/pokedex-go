package main

import (
	"slices"
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "basic split",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "extra whitespace",
			input:    "  hello   world  ",
			expected: []string{"hello", "world"},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := cleanInput(c.input)
			if !slices.Equal(actual, c.expected) {
				t.Errorf("expected: %v, got: %v", c.expected, actual)
			}
		})
	}
}
