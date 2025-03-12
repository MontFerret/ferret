package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateMillisecond(t *testing.T) {
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
			Name:     "When 129 millisecond",
			Expected: core.NewInt(129),
			Args: []core.Value{
				core.NewDateTime(time.Date(2018, 11, 4, 17, 3, 12, 129410001, time.Local)),
			},
		},
		&testCase{
			Name:     "When 0 milliseconds [0]",
			Expected: core.NewInt(0),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-28T15:59:59Z")},
		},
		// any nanosec < 1000000 equal to 0 milliseconds
		&testCase{
			Name:     "When 0 milliseconds [1]",
			Expected: core.NewInt(0),
			Args: []core.Value{
				core.NewDateTime(time.Date(0, 0, 0, 0, 0, 0, 999999, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateMillisecond)
	}
}
