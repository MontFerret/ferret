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

func TestToDateTimeVariadic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		expected time.Time
		name     string
		args     []runtime.Value
	}{
		{
			name:     "RFC3339",
			args:     []runtime.Value{runtime.NewString("2026-08-02T12:00:00Z")},
			expected: time.Date(2026, time.August, 2, 12, 0, 0, 0, time.UTC),
		},
		{
			name:     "integer epoch",
			args:     []runtime.Value{runtime.NewInt(1), runtime.NewString("s")},
			expected: time.Unix(1, 0).UTC(),
		},
		{
			name:     "float epoch",
			args:     []runtime.Value{runtime.NewFloat(1.5), runtime.NewString("s")},
			expected: time.Unix(1, 500_000_000).UTC(),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			value, err := toDateTime(t.Context(), test.args...)
			if err != nil {
				t.Fatal(err)
			}

			actual, ok := value.(runtime.DateTime)
			if !ok || !actual.Equal(test.expected) {
				t.Fatalf("toDateTime(%v) = %v, want %v", test.args, value, test.expected)
			}
		})
	}

	for name, args := range map[string][]runtime.Value{
		"missing": nil,
		"extra":   {runtime.NewInt(1), runtime.NewString("s"), runtime.True},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := toDateTime(t.Context(), args...); !errors.Is(err, runtime.ErrInvalidArgumentNumber) {
				t.Fatalf("toDateTime(%v) error = %v, want ErrInvalidArgumentNumber", args, err)
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
	if !functions.Var().Has("TO_DATETIME") {
		t.Fatal("TO_DATETIME is not registered as variadic")
	}
	if functions.A1().Has("TO_DATETIME") {
		t.Fatal("TO_DATETIME unexpectedly remains registered as fixed arity")
	}
}
