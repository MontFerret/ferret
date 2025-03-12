package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateDateDaysInMonth(t *testing.T) {
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
			Name:     "When Feb and a leap year",
			Expected: core.NewInt(29),
			Args: []core.Value{
				core.NewDateTime(time.Date(1972, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When Feb and not a leap year",
			Expected: core.NewInt(28),
			Args: []core.Value{
				core.NewDateTime(time.Date(1999, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When January",
			Expected: core.NewInt(31),
			Args: []core.Value{
				core.NewDateTime(time.Date(1999, time.January, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When November",
			Expected: core.NewInt(30),
			Args: []core.Value{
				core.NewDateTime(time.Date(1999, time.November, 1, 1, 1, 1, 1, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDaysInMonth)
	}
}
