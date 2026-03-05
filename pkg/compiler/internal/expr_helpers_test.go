package internal

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
)

func TestParseQueryModifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected queryModifier
	}{
		{
			name:     "exists",
			input:    "exists",
			expected: queryModifierExists,
		},
		{
			name:     "uppercase count",
			input:    "COUNT",
			expected: queryModifierCount,
		},
		{
			name:     "mixed any",
			input:    "AnY",
			expected: queryModifierAny,
		},
		{
			name:     "value",
			input:    "value",
			expected: queryModifierValue,
		},
		{
			name:     "one",
			input:    "one",
			expected: queryModifierOne,
		},
		{
			name:     "unknown",
			input:    "other",
			expected: queryModifierUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := parseQueryModifier(tt.input)
			if actual != tt.expected {
				t.Fatalf("unexpected modifier: got %q want %q", actual, tt.expected)
			}
		})
	}
}

func TestQueryResultTypeForModifier(t *testing.T) {
	tests := []struct {
		name     string
		modifier queryModifier
		expected core.ValueType
	}{
		{
			name:     "exists",
			modifier: queryModifierExists,
			expected: core.TypeBool,
		},
		{
			name:     "count",
			modifier: queryModifierCount,
			expected: core.TypeInt,
		},
		{
			name:     "any",
			modifier: queryModifierAny,
			expected: core.TypeAny,
		},
		{
			name:     "value",
			modifier: queryModifierValue,
			expected: core.TypeAny,
		},
		{
			name:     "one",
			modifier: queryModifierOne,
			expected: core.TypeAny,
		},
		{
			name:     "unknown",
			modifier: queryModifierUnknown,
			expected: core.TypeList,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := queryResultTypeForModifier(tt.modifier)
			if actual != tt.expected {
				t.Fatalf("unexpected value type: got %d want %d", actual, tt.expected)
			}
		})
	}
}
