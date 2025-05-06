package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/internal"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateFormat(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 2 arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewString("string"),
				core.NewInt(0),
				internal.NewArray(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 2 arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When first argument is wrong",
			Expected: core.None,
			Args: []core.Value{
				core.NewInt(0),
				core.NewString(time.RFC822),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When second argument is wrong",
			Expected: core.None,
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When DefaultTimeLayout",
			Expected: core.NewString("1999-02-07T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				core.NewString(core.DefaultTimeLayout),
			},
		},
		&testCase{
			Name: "When RFC3339Nano",
			Expected: core.NewString(
				time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local).
					Format(time.RFC3339Nano),
			),
			Args: []core.Value{
				core.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				core.NewString(time.RFC3339Nano),
			},
		},
		&testCase{
			Name: "When custom format",
			Expected: core.NewString(
				time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local).
					Format("2006-01-02"),
			),
			Args: []core.Value{
				core.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				core.NewString("2006-01-02"),
			},
		},
		&testCase{
			Name:     "When empty string",
			Expected: core.NewString(""),
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewString(""),
			},
		},
		&testCase{
			Name:     "When random string without numbers",
			Expected: core.NewString("qwerty"),
			Args: []core.Value{
				core.NewCurrentDateTime(),
				core.NewString("qwerty"),
			},
		},
		&testCase{
			Name:     "When random string with numbers",
			Expected: core.NewString("qwerty2018uio"),
			Args: []core.Value{
				core.NewDateTime(
					time.Date(2018, time.November, 5, 0, 54, 15, 5125, time.Local),
				),
				core.NewString("qwerty2006uio"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateFormat)
	}
}
