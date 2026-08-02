package types

import (
	"errors"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToDuration(t *testing.T) {
	t.Parallel()

	tests := []struct {
		input    runtime.Value
		name     string
		expected runtime.Duration
	}{
		{
			name:     "duration",
			input:    runtime.NewDuration(250 * time.Millisecond),
			expected: runtime.NewDuration(250 * time.Millisecond),
		},
		{
			name:     "literal string",
			input:    runtime.NewString("1.5s"),
			expected: runtime.NewDuration(1500 * time.Millisecond),
		},
		{
			name:     "normalized string",
			input:    runtime.NewString("1h30m"),
			expected: runtime.NewDuration(90 * time.Minute),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			actual, err := ToDuration(t.Context(), tt.input)
			if err != nil {
				t.Fatalf("ToDuration() error = %v", err)
			}

			duration, err := runtime.CastDuration(actual)
			if err != nil {
				t.Fatalf("ToDuration() result type = %s, want Duration", runtime.TypeOf(actual))
			}

			if duration != tt.expected {
				t.Fatalf("ToDuration() = %s, want %s", duration, tt.expected)
			}
		})
	}
}

func TestToDurationRejectsUnsupportedType(t *testing.T) {
	t.Parallel()

	_, err := ToDuration(t.Context(), runtime.NewInt(1))
	if !errors.Is(err, runtime.ErrInvalidType) {
		t.Fatalf("ToDuration() error = %v, want invalid type", err)
	}
}

func TestToDurationRejectsInvalidString(t *testing.T) {
	t.Parallel()

	if _, err := ToDuration(t.Context(), runtime.NewString("1fortnight")); err == nil {
		t.Fatal("ToDuration() accepted an invalid duration string")
	}
}
