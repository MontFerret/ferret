package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateFormat(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 2 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("string"),
				values.NewInt(0),
				values.NewArray(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 2 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When first argument is wrong",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
				values.NewString(time.RFC822),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When second argument is wrong",
			Expected: values.None,
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When DefaultTimeLayout",
			Expected: values.NewString("1999-02-07T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewString(values.DefaultTimeLayout),
			},
		},
		&testCase{
			Name:     "When RFC3339Nano",
			Expected: values.NewString("2018-11-05T00:54:15.000005125+03:00"),
			Args: []core.Value{
				values.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				values.NewString(time.RFC3339Nano),
			},
		},
		&testCase{
			Name:     "When custom format",
			Expected: values.NewString("2018-11-05"),
			Args: []core.Value{
				values.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				values.NewString("2006-01-02"),
			},
		},
		&testCase{
			Name:     "When empty string",
			Expected: values.NewString(""),
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewString(""),
			},
		},
		&testCase{
			Name:     "When random string without numbers",
			Expected: values.NewString("qwerty"),
			Args: []core.Value{
				values.NewCurrentDateTime(),
				values.NewString("qwerty"),
			},
		},
		&testCase{
			Name:     "When random string with numbers",
			Expected: values.NewString("qwerty2018uio"),
			Args: []core.Value{
				values.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				values.NewString("qwerty2006uio"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateFormat)
	}
}
