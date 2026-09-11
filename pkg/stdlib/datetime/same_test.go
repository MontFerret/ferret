package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/datetime"
)

func TestSamePrecision(t *testing.T) {
	base := mustDefaultLayoutDt("2024-06-12T10:20:30.123456789Z")
	for _, tc := range []struct {
		name  string
		right string
		unit  string
		want  bool
	}{
		{"year ignores month", "2024-12-31T23:59:59Z", "year", true},
		{"month ignores day", "2024-06-30T23:59:59Z", "month", true},
		{"month requires year", "2023-06-12T10:20:30.123456789Z", "month", false},
		{"day ignores time", "2024-06-12T23:59:59Z", "day", true},
		{"day requires month", "2024-07-12T10:20:30.123456789Z", "day", false},
		{"hour ignores minute", "2024-06-12T10:59:59Z", "hour", true},
		{"hour requires day", "2024-06-13T10:20:30.123456789Z", "hour", false},
		{"minute ignores second", "2024-06-12T10:20:59Z", "minute", true},
		{"minute requires hour", "2024-06-12T11:20:30.123456789Z", "minute", false},
		{"second ignores fraction", "2024-06-12T10:20:30.999Z", "second", true},
		{"second requires minute", "2024-06-12T10:21:30.123456789Z", "second", false},
		{"millisecond truncates", "2024-06-12T10:20:30.123999999Z", "millisecond", true},
		{"millisecond boundary", "2024-06-12T10:20:30.124Z", "millisecond", false},
		{"millisecond requires second", "2024-06-12T10:20:31.123456789Z", "millisecond", false},
		{"offset ignored", "2024-06-12T10:20:30.123456789+05:30", "millisecond", true},
		{"equal instant different calendar", "2024-06-13T00:20:30.123456789+14:00", "day", false},
		{"same ISO week", "2024-06-16T23:59:59Z", "week", true},
		{"different ISO week", "2024-06-17T10:20:30Z", "week", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := datetime.Same(t.Context(), base, mustDefaultLayoutDt(tc.right), runtime.String(tc.unit))
			if err != nil || got != runtime.Boolean(tc.want) {
				t.Fatalf("Same = %v, %v; want %v", got, err, tc.want)
			}
		})
	}

	for _, tc := range []struct {
		left, right string
		want        bool
	}{
		{"2020-12-31T12:00:00Z", "2021-01-01T12:00:00Z", true},
		{"2021-01-03T12:00:00Z", "2021-01-04T12:00:00Z", false},
		{"2023-01-02T12:00:00Z", "2024-01-01T12:00:00Z", false},
	} {
		got, err := datetime.Same(t.Context(), mustDefaultLayoutDt(tc.left), mustDefaultLayoutDt(tc.right), runtime.String("week"))
		if err != nil || got != runtime.Boolean(tc.want) {
			t.Fatalf("ISO weeks %s / %s = %v, %v", tc.left, tc.right, got, err)
		}
	}
}
