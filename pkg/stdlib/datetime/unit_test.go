package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

// Test UnitFromString function which is exported
func TestUnitFromString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		hasError bool
	}{
		{"Valid year", "year", false},
		{"Valid y", "y", false},
		{"Valid years", "years", false},
		{"Valid month", "month", false},
		{"Valid m", "m", false},
		{"Valid months", "months", false},
		{"Valid week", "week", false},
		{"Valid w", "w", false},
		{"Valid weeks", "weeks", false},
		{"Valid day", "day", false},
		{"Valid d", "d", false},
		{"Valid days", "days", false},
		{"Valid hour", "hour", false},
		{"Valid h", "h", false},
		{"Valid hours", "hours", false},
		{"Valid minute", "minute", false},
		{"Valid i", "i", false},
		{"Valid minutes", "minutes", false},
		{"Valid second", "second", false},
		{"Valid s", "s", false},
		{"Valid seconds", "seconds", false},
		{"Valid millisecond", "millisecond", false},
		{"Valid f", "f", false},
		{"Valid milliseconds", "milliseconds", false},
		{"Invalid unit", "invalid_unit", true},
		{"Case insensitive", "YEAR", false},
		{"Case insensitive", "Hour", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := datetime.UnitFromString(tt.input)
			if tt.hasError && err == nil {
				t.Error("expected error but got none")
			}
			if !tt.hasError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}