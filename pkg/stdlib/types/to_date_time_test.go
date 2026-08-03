package types

import (
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestToDateTime(t *testing.T) {
	t.Parallel()

	expected := runtime.NewDateTime(time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC))

	for _, input := range []runtime.Value{expected, runtime.NewString("2026-08-01T12:00:00Z")} {
		actual, err := ToDateTime(t.Context(), input)
		if err != nil {
			t.Fatalf("ToDateTime(%s): %v", input, err)
		}

		dateTime, err := runtime.CastDateTime(actual)
		if err != nil || !dateTime.Equal(expected.Time) {
			t.Fatalf("ToDateTime(%s) = %v, %v", input, actual, err)
		}
	}

	for _, input := range []runtime.Value{runtime.NewInt(0), runtime.NewFloat(0)} {
		if _, err := ToDateTime(t.Context(), input); !errors.Is(err, runtime.ErrInvalidArgument) ||
			!strings.Contains(err.Error(), "requires an explicit epoch unit") {
			t.Fatalf("ToDateTime(%s) error = %v", runtime.TypeOf(input), err)
		}
	}
}

func TestToDateTimeEpoch(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected time.Time
		value    runtime.Value
		unit     runtime.Value
		name     string
	}{
		{
			name:     "integer epoch",
			value:    runtime.NewInt(1),
			unit:     runtime.NewString("s"),
			expected: time.Unix(1, 0).UTC(),
		},
		{
			name:     "float epoch",
			value:    runtime.NewFloat(1.5),
			unit:     runtime.NewString("s"),
			expected: time.Unix(1, 500_000_000).UTC(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := toDateTime2(t.Context(), test.value, test.unit)
			if err != nil {
				t.Fatal(err)
			}

			actual, ok := value.(runtime.DateTime)
			if !ok || !actual.Equal(test.expected) {
				t.Fatalf("toDateTime2(%v, %v) = %v, want %v", test.value, test.unit, value, test.expected)
			}
		})
	}
}

func TestToDateTimeRegistration(t *testing.T) {
	t.Parallel()

	library := runtime.NewLibrary()
	RegisterLib(library)

	functions, err := library.Build()
	if err != nil {
		t.Fatal(err)
	}
	if functions.Var().Has("TO_DATETIME") {
		t.Fatal("TO_DATETIME unexpectedly remains registered as variadic")
	}
	if !functions.A1().Has("TO_DATETIME") || !functions.A2().Has("TO_DATETIME") {
		t.Fatal("TO_DATETIME is not registered at arities 1 and 2")
	}
}
