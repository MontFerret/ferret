package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateCompare(t *testing.T) {
	expectedTrue := values.NewBoolean(true)
	expectedFalse := values.NewBoolean(false)

	tcs := []*testCase{
		&testCase{
			Name:      "When less than 3 arguments",
			Expected:  values.None,
			Args:      []core.Value{values.NewInt(0), values.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When more than 4 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0), values.NewInt(0), values.NewInt(0),
				values.NewInt(0), values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewCurrentDateTime(),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of optional argument",
			Expected: values.None,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewCurrentDateTime(),
				values.NewString("year"),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when start unit less that end unit",
			Expected: values.None,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewCurrentDateTime(),
				values.NewString("day"),
				values.NewString("year"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when years are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewCurrentDateTime(),
				values.NewString("year"),
			},
		},
		&testCase{
			Name:     "when years are not equal",
			Expected: expectedFalse,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-07"),
				values.NewString("year"),
				values.NewString("year"),
			},
		},
		&testCase{
			Name:     "when months are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-09"),
				values.NewString("year"),
				values.NewString("days"),
			},
		},
		&testCase{
			Name:     "when days are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewCurrentDateTime(),
				values.NewString("days"),
				values.NewString("days"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateCompare)
	}
}
