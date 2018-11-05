package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateMillisecond(t *testing.T) {
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
			Name:     "When 129 millisecond",
			Expected: values.NewInt(129),
			Args: []core.Value{
				values.NewDateTime(time.Date(2018, 11, 4, 17, 3, 12, 129410001, time.Local)),
			},
		},
		&testCase{
			Name:     "When 0 milliseconds [0]",
			Expected: values.NewInt(0),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-28T15:59:59Z")},
		},
		// any nanosec < 1000000 equal to 0 milliseconds
		&testCase{
			Name:     "When 0 milliseconds [1]",
			Expected: values.NewInt(0),
			Args: []core.Value{
				values.NewDateTime(time.Date(0, 0, 0, 0, 0, 0, 999999, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateMillisecond)
	}
}
