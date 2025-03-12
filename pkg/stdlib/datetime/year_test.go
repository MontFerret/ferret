package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateYear(t *testing.T) {
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
			Name:     "When 1999th year",
			Expected: core.NewInt(1999),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 1629th year",
			Expected: core.NewInt(1629),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-08T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateYear)
	}
}
