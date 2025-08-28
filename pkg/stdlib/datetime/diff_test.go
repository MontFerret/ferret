package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

var (
	isFloat        = runtime.NewBoolean(true)
	beginningEpoch = runtime.NewDateTime(time.Time{})
)

func TestDiff(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:      "when less then 3 arguments",
			Expected:  runtime.NewInt(1),
			Args:      []runtime.Value{beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:      "when more then 4 arguments",
			Expected:  runtime.NewInt(1),
			Args:      []runtime.Value{beginningEpoch, beginningEpoch, beginningEpoch, beginningEpoch, beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:      "when wrong type argument",
			Expected:  runtime.NewInt(1),
			Args:      []runtime.Value{beginningEpoch, beginningEpoch, beginningEpoch},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when the difference is 1 year and 1 month (int)",
			Expected: runtime.NewInt(1),
			Args: []runtime.Value{
				beginningEpoch,
				runtime.NewDateTime(
					beginningEpoch.AddDate(1, 1, 0),
				),
				runtime.NewString("y"),
			},
		},
		&testCase{
			Name:     "when the difference is 1 year and 1 month (float)",
			Expected: runtime.NewFloat(1.084931506849315),
			Args: []runtime.Value{
				beginningEpoch,
				runtime.NewDateTime(
					beginningEpoch.AddDate(1, 1, 0),
				),
				runtime.NewString("year"),
				isFloat,
			},
		},
		&testCase{
			Name:     "when date1 after date2 (int)",
			Expected: runtime.NewInt(2),
			Args: []runtime.Value{
				beginningEpoch,
				runtime.NewDateTime(
					beginningEpoch.Add(-time.Hour * 48),
				),
				runtime.NewString("d"),
			},
		},
		&testCase{
			Name:     "when date1 after date2 (float)",
			Expected: runtime.NewFloat(2),
			Args: []runtime.Value{
				beginningEpoch,
				runtime.NewDateTime(
					beginningEpoch.Add(-time.Hour * 48),
				),
				runtime.NewString("d"),
				isFloat,
			},
		},
		&testCase{
			Name:     "when dates are equal (int)",
			Expected: runtime.NewInt(0),
			Args: []runtime.Value{
				beginningEpoch,
				beginningEpoch,
				runtime.NewString("i"),
			},
		},
		&testCase{
			Name:     "when dates are equal (float)",
			Expected: runtime.NewFloat(0),
			Args: []runtime.Value{
				beginningEpoch,
				beginningEpoch,
				runtime.NewString("y"),
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
				Expected: runtime.NewInt(1),
				Args: []runtime.Value{
					beginningEpoch,
					runtime.NewDateTime(
						beginningEpoch.AddDate(dates[0], dates[1], dates[2]),
					),
					runtime.NewString(unit),
				},
			},
			&testCase{
				Name:     "When difference is 1 " + unit + " (float)",
				Expected: runtime.NewFloat(1),
				Args: []runtime.Value{
					beginningEpoch,
					runtime.NewDateTime(
						beginningEpoch.AddDate(dates[0], dates[1], dates[2]),
					),
					runtime.NewString(unit),
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
				Expected: runtime.NewInt(1),
				Args: []runtime.Value{
					beginningEpoch,
					runtime.NewDateTime(
						beginningEpoch.Add(durn),
					),
					runtime.NewString(unit),
				},
			},
			&testCase{
				Name:     "When difference is 1 " + unit + " (int)",
				Expected: runtime.NewFloat(1),
				Args: []runtime.Value{
					beginningEpoch,
					runtime.NewDateTime(
						beginningEpoch.Add(durn),
					),
					runtime.NewString(unit),
					isFloat,
				},
			},
		)
	}

	// Additional test cases to improve coverage
	additionalTests := []*testCase{
		&testCase{
			Name:     "when dates are not equal with invalid unit",
			Expected: runtime.None,
			Args: []runtime.Value{
				beginningEpoch,
				runtime.NewDateTime(beginningEpoch.Add(time.Hour)),
				runtime.NewString("invalid_unit"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when first date is after second date",
			Expected: runtime.NewInt(1),
			Args: []runtime.Value{
				runtime.NewDateTime(beginningEpoch.Add(time.Hour)),
				beginningEpoch,
				runtime.NewString("hour"),
			},
		},
	}

	tcs = append(tcs, additionalTests...)

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDiff)
	}
}
