package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateQuarter(t *testing.T) {
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
	}

	for month := time.January; month <= time.December; month++ {
		tcs = append(tcs, &testCase{
			Name:     "When " + month.String(),
			Expected: core.NewInt(((int(month) - 1) / 3) + 1),
			Args: []core.Value{
				core.NewDateTime(time.Date(1999, month, 1, 1, 1, 1, 1, time.Local)),
			},
		})
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateQuarter)
	}
}
