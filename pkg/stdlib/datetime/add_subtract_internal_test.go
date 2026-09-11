package datetime

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func TestAddUnitInvalidUnit(t *testing.T) {
	for _, tc := range []struct {
		name string
		unit unit
	}{
		{"below minimum", millisecond - 1},
		{"above maximum", year + 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := addUnit(time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC), 1, tc.unit)
			if !got.IsZero() || !errors.Is(err, runtime.ErrUnexpected) {
				t.Fatalf("addUnit(%d) = %v, %v; want zero time and an unexpected error", tc.unit, got, err)
			}

			if want := fmt.Sprintf("unsupported datetime unit %d", tc.unit); !strings.Contains(err.Error(), want) {
				t.Fatalf("addUnit(%d) error = %q; want message containing %q", tc.unit, err, want)
			}
		})
	}
}
