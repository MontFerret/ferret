package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateQuarter(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 1 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("string"),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When 0 arguments",
			Expected:  runtime.None,
			Args:      []runtime.Value{},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When argument isn't DateTime",
			Expected:  runtime.None,
			Args:      []runtime.Value{runtime.NewInt(0)},
			ShouldErr: true,
		},
	}

	for month := time.January; month <= time.December; month++ {
		tcs = append(tcs, &testCase{
			Name:     "When " + month.String(),
			Expected: runtime.NewInt(((int(month) - 1) / 3) + 1),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(1999, month, 1, 1, 1, 1, 1, time.Local)),
			},
		})
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateQuarter)
	}
}
