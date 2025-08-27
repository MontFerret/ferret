package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateHour(t *testing.T) {
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
			Name:     "When 7th hour",
			Expected: runtime.NewInt(7),
			Args:     []runtime.Value{mustDefaultLayoutDt("1999-02-07T07:04:05Z")},
		},
		&testCase{
			Name:     "When 15th day",
			Expected: runtime.NewInt(15),
			Args:     []runtime.Value{mustDefaultLayoutDt("1629-02-28T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateHour)
	}
}
