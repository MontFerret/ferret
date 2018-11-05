package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateQuarter(t *testing.T) {
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
	}

	for month := time.January; month <= time.December; month++ {
		tcs = append(tcs, &testCase{
			Name:     "When " + month.String(),
			Expected: values.NewInt(((int(month) - 1) / 3) + 1),
			Args: []core.Value{
				values.NewDateTime(time.Date(1999, month, 1, 1, 1, 1, 1, time.Local)),
			},
		})
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateQuarter)
	}
}
