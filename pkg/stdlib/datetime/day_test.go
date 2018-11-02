package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateDay(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 1 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("string"),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When 0 arguments",
			Expected:  values.None,
			Args:      []core.Value{},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When 7th day",
			Expected: values.NewInt(7),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 28th day",
			Expected: values.NewInt(28),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-28T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDay)
	}
}
