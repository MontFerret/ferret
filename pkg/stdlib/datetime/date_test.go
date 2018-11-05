package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDate(t *testing.T) {
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
			Name:      "When argument isn't DateTime",
			Expected:  values.None,
			Args:      []core.Value{values.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When incorrect timeStrings",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("bla-bla"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When correct timeString in RFC3339 format",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
			Args: []core.Value{
				values.NewString("1999-02-07T15:04:05Z"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Date)
	}
}
