package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateFormat(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 2 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("string"),
				runtime.NewInt(0),
				runtime.NewArray(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 2 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When first argument is wrong",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
				runtime.NewString(time.RFC822),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When second argument is wrong",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When DefaultTimeLayout",
			Expected: runtime.NewString("1999-02-07T15:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewString(runtime.DefaultTimeLayout),
			},
		},
		&testCase{
			Name: "When RFC3339Nano",
			Expected: runtime.NewString(
				time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local).
					Format(time.RFC3339Nano),
			),
			Args: []runtime.Value{
				runtime.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				runtime.NewString(time.RFC3339Nano),
			},
		},
		&testCase{
			Name: "When custom format",
			Expected: runtime.NewString(
				time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local).
					Format("2006-01-02"),
			),
			Args: []runtime.Value{
				runtime.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				runtime.NewString("2006-01-02"),
			},
		},
		&testCase{
			Name:     "When empty string",
			Expected: runtime.NewString(""),
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewString(""),
			},
		},
		&testCase{
			Name:     "When random string without numbers",
			Expected: runtime.NewString("qwerty"),
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
				runtime.NewString("qwerty"),
			},
		},
		&testCase{
			Name:     "When random string with numbers",
			Expected: runtime.NewString("qwerty2018uio"),
			Args: []runtime.Value{
				runtime.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				runtime.NewString("qwerty2006uio"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateFormat)
	}
}
