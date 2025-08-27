package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateCompare(t *testing.T) {
	expectedTrue := runtime.NewBoolean(true)
	expectedFalse := runtime.NewBoolean(false)

	tcs := []*testCase{
		&testCase{
			Name:      "When less than 3 arguments",
			Expected:  runtime.None,
			Args:      []runtime.Value{runtime.NewInt(0), runtime.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When more than 4 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0), runtime.NewInt(0), runtime.NewInt(0),
				runtime.NewInt(0), runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewCurrentDateTime(),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of optional argument",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewCurrentDateTime(),
				runtime.NewString("year"),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when start unit less that end unit",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewCurrentDateTime(),
				runtime.NewString("day"),
				runtime.NewString("year"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when years are equal",
			Expected: expectedTrue,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewCurrentDateTime(),
				runtime.NewString("year"),
			},
		},
		&testCase{
			Name:     "when years are not equal",
			Expected: expectedFalse,
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-07"),
				runtime.NewString("year"),
				runtime.NewString("year"),
			},
		},
		&testCase{
			Name:     "when months are equal",
			Expected: expectedTrue,
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-09"),
				runtime.NewString("year"),
				runtime.NewString("days"),
			},
		},
		&testCase{
			Name:     "when days are equal",
			Expected: expectedTrue,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewCurrentDateTime(),
				runtime.NewString("days"),
				runtime.NewString("days"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateCompare)
	}
}
