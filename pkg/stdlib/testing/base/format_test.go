package base

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestFormatValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		value    runtime.Value
		expected string
	}{
		{
			name:     "plain string",
			value:    runtime.NewString("test"),
			expected: "String 'test'",
		},
		{
			name:     "none",
			value:    runtime.None,
			expected: "None 'none'",
		},
		{
			name:     "apostrophe and backslash",
			value:    runtime.NewString(`can't\stop`),
			expected: `String 'can\'t\\stop'`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if actual := FormatValue(tt.value); actual != tt.expected {
				t.Fatalf("FormatValue() = %q, want %q", actual, tt.expected)
			}
		})
	}
}
