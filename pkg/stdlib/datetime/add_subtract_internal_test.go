package datetime

import (
	"testing"
	"time"
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
			defer func() {
				if got := recover(); got != "unreachable" {
					t.Fatalf("addUnit(%d) panic = %v; want unreachable", tc.unit, got)
				}
			}()

			addUnit(time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC), 1, tc.unit)
		})
	}
}
