package datetime_test

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

var (
	utcLoc, _ = time.LoadLocation("UTC")
)

func TestDateAdd(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 3 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
				values.NewInt(0),
				values.NewInt(0),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 3 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When incorrect arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("bla-bla"),
				values.NewInt(0),
				values.NewString("be-be"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When wrong unit given",
			Expected: values.None,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(5),
				values.NewString("not_exist"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name: "When argument have correct types",
			Expected: func() core.Value {
				expected, _ := datetime.DateAdd(
					context.Background(),
					mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
					values.NewInt(1),
					values.NewString("day"),
				)
				return expected
			}(),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(1),
				values.NewString("day"),
			},
		},
		&testCase{
			Name:     "-1 day",
			Expected: mustDefaultLayoutDt("1999-02-06T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(-1),
				values.NewString("day"),
			},
		},
		&testCase{
			Name:     "+3 months",
			Expected: mustDefaultLayoutDt("1999-05-07T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(3),
				values.NewString("months"),
			},
		},
		&testCase{
			Name:     "+5 years",
			Expected: mustLayoutDt("2006-01-02", "2004-02-07"),
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(5),
				values.NewString("y"),
			},
		},
		&testCase{
			Name: "1999 minus 2000 years",
			Expected: values.NewDateTime(
				time.Date(-1, 2, 7, 0, 0, 0, 0, utcLoc),
			),
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(-2000),
				values.NewString("year"),
			},
		},
		&testCase{
			Name:     "+2 hours",
			Expected: mustDefaultLayoutDt("1999-02-07T17:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(2),
				values.NewString("h"),
			},
		},
		&testCase{
			Name:     "+20 minutes",
			Expected: mustDefaultLayoutDt("1999-02-07T15:24:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(20),
				values.NewString("i"),
			},
		},
		&testCase{
			Name:     "+30 seconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:35Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(30),
				values.NewString("s"),
			},
		},
		&testCase{
			Name:     "+1000 milliseconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:06Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(1000),
				values.NewString("f"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateAdd)
	}
}

func TestDateSubtract(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 3 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
				values.NewInt(0),
				values.NewInt(0),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 3 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When incorrect arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("bla-bla"),
				values.NewInt(0),
				values.NewString("be-be"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When wrong unit given",
			Expected: values.None,
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(5),
				values.NewString("not_exist"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name: "When argument have correct types",
			Expected: func() core.Value {
				expected, _ := datetime.DateSubtract(
					context.Background(),
					mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
					values.NewInt(1),
					values.NewString("day"),
				)
				return expected
			}(),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(1),
				values.NewString("day"),
			},
		},
		&testCase{
			Name:     "-1 day",
			Expected: mustDefaultLayoutDt("1999-02-08T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(-1),
				values.NewString("day"),
			},
		},
		&testCase{
			Name:     "+3 months",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-05-07T15:04:05Z"),
				values.NewInt(3),
				values.NewString("months"),
			},
		},
		&testCase{
			Name:     "+5 years",
			Expected: mustLayoutDt("2006-01-02", "1994-02-07"),
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(5),
				values.NewString("y"),
			},
		},
		&testCase{
			Name: "1999 minus 2000 years",
			Expected: values.NewDateTime(
				time.Date(-1, 2, 7, 0, 0, 0, 0, utcLoc),
			),
			Args: []core.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				values.NewInt(2000),
				values.NewString("year"),
			},
		},
		&testCase{
			Name:     "-2 hours",
			Expected: mustDefaultLayoutDt("1999-02-07T13:04:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(2),
				values.NewString("h"),
			},
		},
		&testCase{
			Name:     "-20 minutes",
			Expected: mustDefaultLayoutDt("1999-02-07T14:44:05Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(20),
				values.NewString("i"),
			},
		},
		&testCase{
			Name:     "-30 seconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:03:35Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(30),
				values.NewString("s"),
			},
		},
		&testCase{
			Name:     "-1000 milliseconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:04Z"),
			Args: []core.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				values.NewInt(1000),
				values.NewString("f"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateSubtract)
	}
}
