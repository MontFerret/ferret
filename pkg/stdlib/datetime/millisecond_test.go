package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateMillisecond(t *testing.T) {
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
			Name:     "When 129 millisecond",
			Expected: runtime.NewInt(129),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(2018, 11, 4, 17, 3, 12, 129410001, time.Local)),
			},
		},
		&testCase{
			Name:     "When 0 milliseconds [0]",
			Expected: runtime.NewInt(0),
			Args:     []runtime.Value{mustDefaultLayoutDt("1629-02-28T15:59:59Z")},
		},
		// any nanosec < 1000000 equal to 0 milliseconds
		&testCase{
			Name:     "When 0 milliseconds [1]",
			Expected: runtime.NewInt(0),
			Args: []runtime.Value{
				runtime.NewDateTime(time.Date(0, 0, 0, 0, 0, 0, 999999, time.Local)),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateMillisecond)
	}
}
