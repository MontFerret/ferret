package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateLeapYear(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 1 arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewString("string"),
				core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When 0 arguments",
			Expected:  core.None,
			Args:      []core.Value{},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When argument isn't DateTime",
			Expected:  core.None,
			Args:      []core.Value{core.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When not a leap year",
			Expected: core.NewBoolean(false),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When a leap year",
			Expected: core.NewBoolean(true),
			Args:     []core.Value{mustDefaultLayoutDt("1972-12-07T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateLeapYear)
	}
}
