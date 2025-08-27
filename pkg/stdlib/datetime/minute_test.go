package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateMinute(t *testing.T) {
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
			Name:     "When 4th minute",
			Expected: runtime.NewInt(4),
			Args:     []runtime.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 59th minute",
			Expected: runtime.NewInt(59),
			Args:     []runtime.Value{mustDefaultLayoutDt("1629-02-28T15:59:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateMinute)
	}
}
