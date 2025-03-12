package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateMonth(t *testing.T) {
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
			Name:     "When 2th month",
			Expected: core.NewInt(2),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 12th month",
			Expected: core.NewInt(12),
			Args:     []core.Value{mustDefaultLayoutDt("1999-12-07T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateMonth)
	}
}
