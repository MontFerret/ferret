package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateDateDaysInMonth(t *testing.T) {
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
			Name:     "When Feb and a leap year",
			Expected: values.NewInt(29),
			Args: []core.Value{
				values.NewDateTime(time.Date(1972, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When Feb and not a leap year",
			Expected: values.NewInt(28),
			Args: []core.Value{
				values.NewDateTime(time.Date(1999, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When January",
			Expected: values.NewInt(31),
			Args: []core.Value{
				values.NewDateTime(time.Date(1999, time.January, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When November",
			Expected: values.NewInt(30),
			Args: []core.Value{
				values.NewDateTime(time.Date(1999, time.November, 1, 1, 1, 1, 1, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDaysInMonth)
	}
}
