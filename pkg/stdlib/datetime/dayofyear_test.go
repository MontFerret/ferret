package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateDayOfYear(t *testing.T) {
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
			Name:     "When 38th day of the year",
			Expected: core.NewInt(38),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 59th day of the year",
			Expected: core.NewInt(59),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-28T15:59:05Z")},
		},
		&testCase{
			Name:     "When 366th day of the year",
			Expected: core.NewInt(366),
			Args: []core.Value{
				core.NewDateTime(time.Date(1972, time.December, 31, 0, 0, 0, 0, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateDayOfYear)
	}
}
