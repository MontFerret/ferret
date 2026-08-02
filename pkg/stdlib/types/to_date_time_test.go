package types

import (
	"errors"
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
		if _, err := ToDateTime(t.Context(), input); !errors.Is(err, runtime.ErrInvalidType) {
			t.Fatalf("ToDateTime(%s) error = %v", runtime.TypeOf(input), err)
		}
	}
}
