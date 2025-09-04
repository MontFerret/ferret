package parser

import (
	"testing"
)

func TestPreprocessStepSyntax(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple increment",
			input:    "FOR i = 0 WHILE i < 10 STEP i++ RETURN i",
			expected: "FOR i = 0 WHILE i < 10 STEP i = i + 1 RETURN i",
		},
		{
			name:     "Simple decrement",
			input:    "FOR i = 10 WHILE i > 0 STEP i-- RETURN i",
			expected: "FOR i = 10 WHILE i > 0 STEP i = i - 1 RETURN i",
		},
		{
			name:     "Different variable name increment",
			input:    "FOR counter = 0 WHILE counter < 5 STEP counter++ RETURN counter",
			expected: "FOR counter = 0 WHILE counter < 5 STEP counter = counter + 1 RETURN counter",
		},
		{
			name:     "Different variable name decrement",
			input:    "FOR index = 5 WHILE index > 0 STEP index-- RETURN index",
			expected: "FOR index = 5 WHILE index > 0 STEP index = index - 1 RETURN index",
		},
		{
			name:     "Case insensitive STEP",
			input:    "FOR i = 0 WHILE i < 5 step i++ RETURN i",
			expected: "FOR i = 0 WHILE i < 5 STEP i = i + 1 RETURN i",
		},
		{
			name:     "Extra spaces around ++",
			input:    "FOR i = 0 WHILE i < 5 STEP i ++ RETURN i",
			expected: "FOR i = 0 WHILE i < 5 STEP i = i + 1 RETURN i",
		},
		{
			name:     "Extra spaces around --",
			input:    "FOR i = 5 WHILE i > 0 STEP i -- RETURN i",
			expected: "FOR i = 5 WHILE i > 0 STEP i = i - 1 RETURN i",
		},
		{
			name:     "Underscore in variable name",
			input:    "FOR var_name = 0 WHILE var_name < 3 STEP var_name++ RETURN var_name",
			expected: "FOR var_name = 0 WHILE var_name < 3 STEP var_name = var_name + 1 RETURN var_name",
		},
		{
			name:     "No transformation needed",
			input:    "FOR i = 0 WHILE i < 5 STEP i = i + 2 RETURN i",
			expected: "FOR i = 0 WHILE i < 5 STEP i = i + 2 RETURN i",
		},
		{
			name:     "Multiple loops in same query",
			input:    "LET a = (FOR i = 0 WHILE i < 2 STEP i++ RETURN i) LET b = (FOR j = 2 WHILE j > 0 STEP j-- RETURN j) RETURN [a, b]",
			expected: "LET a = (FOR i = 0 WHILE i < 2 STEP i = i + 1 RETURN i) LET b = (FOR j = 2 WHILE j > 0 STEP j = j - 1 RETURN j) RETURN [a, b]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := preprocessStepSyntax(tt.input)
			if result != tt.expected {
				t.Errorf("preprocessStepSyntax() = %q, want %q", result, tt.expected)
			}
		})
	}
}