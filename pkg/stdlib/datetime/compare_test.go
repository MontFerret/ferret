package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateCompare(t *testing.T) {
	expectedTrue := core.NewBoolean(true)
	expectedFalse := core.NewBoolean(false)

	tcs := []*testCase{
		&testCase{
			Name:      "When less than 3 arguments",
			Expected:  core.None,
			Args:      []core.Value{core.NewInt(0), core.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When more than 4 arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewInt(0), core.NewInt(0), core.NewInt(0),
				core.NewInt(0), core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewCurrentDateTime(),
				core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when wrong type of optional argument",
			Expected: core.None,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewCurrentDateTime(),
				core.NewString("year"),
				core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when start unit less that end unit",
			Expected: core.None,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewCurrentDateTime(),
				core.NewString("day"),
				core.NewString("year"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "when years are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewCurrentDateTime(),
				core.NewString("year"),
			},
		},
		&testCase{
			Name:     "when years are not equal",
			Expected: expectedFalse,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-07"),
				core.NewString("year"),
				core.NewString("year"),
			},
		},
		&testCase{
			Name:     "when months are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				mustLayoutDt("2006-01-02", "2000-02-09"),
				core.NewString("year"),
				core.NewString("days"),
			},
		},
		&testCase{
			Name:     "when days are equal",
			Expected: expectedTrue,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewCurrentDateTime(),
				core.NewString("days"),
				core.NewString("days"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateCompare)
	}
}
