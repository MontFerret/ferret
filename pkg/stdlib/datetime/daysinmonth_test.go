package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateDateDaysInMonth(t *testing.T) {
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
		&testCase{
			Name:     "When Feb and a leap year",
			Expected: runtime.NewInt(29),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(1972, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When Feb and not a leap year",
			Expected: runtime.NewInt(28),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(1999, time.February, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When January",
			Expected: runtime.NewInt(31),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(1999, time.January, 1, 1, 1, 1, 1, time.Local)),
			},
		},
		&testCase{
			Name:     "When November",
			Expected: runtime.NewInt(30),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(1999, time.November, 1, 1, 1, 1, 1, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDaysInMonth)
	}
}
