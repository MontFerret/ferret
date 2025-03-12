package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

var (
	isFloat        = core.NewBoolean(true)
	beginningEpoch = core.NewDateTime(time.Time{})
)

func TestDiff(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:      "when less then 3 arguments",
			Expected:  core.NewInt(1),
			Args:      []core.Value{beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:      "when more then 4 arguments",
			Expected:  core.NewInt(1),
			Args:      []core.Value{beginningEpoch, beginningEpoch, beginningEpoch, beginningEpoch, beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:      "when wrong type argument",
			Expected:  core.NewInt(1),
			Args:      []core.Value{beginningEpoch, beginningEpoch, beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when the difference is 1 year and 1 month (int)",
			Expected: core.NewInt(1),
			Args: []core.Value{
				beginningEpoch,
				core.NewDateTime(
					beginningEpoch.AddDate(1, 1, 0),
				),
				core.NewString("y"),
			},
		},
		&testCase{
			Name:     "when the difference is 1 year and 1 month (float)",
			Expected: core.NewFloat(1.084931506849315),
			Args: []core.Value{
				beginningEpoch,
				core.NewDateTime(
					beginningEpoch.AddDate(1, 1, 0),
				),
				core.NewString("year"),
				isFloat,
			},
		},
		&testCase{
			Name:     "when date1 after date2 (int)",
			Expected: core.NewInt(2),
			Args: []core.Value{
				beginningEpoch,
				core.NewDateTime(
					beginningEpoch.Add(-time.Hour * 48),
				),
				core.NewString("d"),
			},
		},
		&testCase{
			Name:     "when date1 after date2 (float)",
			Expected: core.NewFloat(2),
			Args: []core.Value{
				beginningEpoch,
				core.NewDateTime(
					beginningEpoch.Add(-time.Hour * 48),
				),
				core.NewString("d"),
				isFloat,
			},
		},
		&testCase{
			Name:     "when dates are equal (int)",
			Expected: core.NewInt(0),
			Args: []core.Value{
				beginningEpoch,
				beginningEpoch,
				core.NewString("i"),
			},
		},
		&testCase{
			Name:     "when dates are equal (float)",
			Expected: core.NewFloat(0),
			Args: []core.Value{
				beginningEpoch,
				beginningEpoch,
				core.NewString("y"),
				isFloat,
			},
		},
	}

	bigUnits := map[string][3]int{
		"y": [3]int{1, 0, 0}, "year": [3]int{1, 0, 0}, "years": [3]int{1, 0, 0},
		"m": [3]int{0, 1, 0}, "month": [3]int{0, 1, 0}, "months": [3]int{0, 1, 0},
		"w": [3]int{0, 0, 7}, "week": [3]int{0, 0, 7}, "weeks": [3]int{0, 0, 7},
		"d": [3]int{0, 0, 1}, "day": [3]int{0, 0, 1}, "days": [3]int{0, 0, 1},
	}

	for unit, dates := range bigUnits {
		tcs = append(tcs,
			&testCase{
				Name:     "When difference is 1 " + unit + " (int)",
				Expected: core.NewInt(1),
				Args: []core.Value{
					beginningEpoch,
					core.NewDateTime(
						beginningEpoch.AddDate(dates[0], dates[1], dates[2]),
					),
					core.NewString(unit),
				},
			},
			&testCase{
				Name:     "When difference is 1 " + unit + " (float)",
				Expected: core.NewFloat(1),
				Args: []core.Value{
					beginningEpoch,
					core.NewDateTime(
						beginningEpoch.AddDate(dates[0], dates[1], dates[2]),
					),
					core.NewString(unit),
					isFloat,
				},
			},
		)
	}

	units := map[string]time.Duration{
		"h": time.Hour, "hour": time.Hour, "hours": time.Hour,
		"i": time.Minute, "minute": time.Minute, "minutes": time.Minute,
		"s": time.Second, "second": time.Second, "seconds": time.Second,
		"f": time.Millisecond, "millisecond": time.Millisecond, "milliseconds": time.Millisecond,
	}

	for unit, durn := range units {
		tcs = append(tcs,
			&testCase{
				Name:     "When difference is 1 " + unit + " (int)",
				Expected: core.NewInt(1),
				Args: []core.Value{
					beginningEpoch,
					core.NewDateTime(
						beginningEpoch.Add(durn),
					),
					core.NewString(unit),
				},
			},
			&testCase{
				Name:     "When difference is 1 " + unit + " (int)",
				Expected: core.NewFloat(1),
				Args: []core.Value{
					beginningEpoch,
					core.NewDateTime(
						beginningEpoch.Add(durn),
					),
					core.NewString(unit),
					isFloat,
				},
			},
		)
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDiff)
	}
}
