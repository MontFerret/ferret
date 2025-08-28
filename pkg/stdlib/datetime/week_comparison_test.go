package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestWeekComparisonBugFix(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When days 7 and 8 are in different weeks - week unit only",
			Expected: runtime.NewBoolean(false),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02T15:04:05Z", "1999-02-07T15:04:05Z"), // Day 7, Week 5
				mustLayoutDt("2006-01-02T15:04:05Z", "2000-03-08T16:05:06Z"), // Day 8, Week 10
				runtime.NewString("week"),
				runtime.NewString("week"),
			},
		},
		&testCase{
			Name:     "When days are in the same ISO week - week unit only",
			Expected: runtime.NewBoolean(true),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02T15:04:05Z", "2023-01-02T15:04:05Z"), // Monday of week 1
				mustLayoutDt("2006-01-02T15:04:05Z", "2023-01-08T16:05:06Z"), // Sunday of week 1
				runtime.NewString("week"),
				runtime.NewString("week"),
			},
		},
		&testCase{
			Name:     "When days are in different ISO weeks but same month - week unit only",
			Expected: runtime.NewBoolean(false),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02T15:04:05Z", "2023-01-01T15:04:05Z"), // Week 52 of previous year
				mustLayoutDt("2006-01-02T15:04:05Z", "2023-01-09T16:05:06Z"), // Week 2 of current year
				runtime.NewString("week"),
				runtime.NewString("week"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateCompare)
	}
}