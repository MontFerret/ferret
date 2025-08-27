package datetime_test

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

var (
	utcLoc, _ = time.LoadLocation("UTC")
)

func TestDateAdd(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 3 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
				runtime.NewInt(0),
				runtime.NewInt(0),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 3 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When incorrect arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("bla-bla"),
				runtime.NewInt(0),
				runtime.NewString("be-be"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When wrong unit given",
			Expected: runtime.None,
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(5),
				runtime.NewString("not_exist"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name: "When argument have correct types",
			Expected: func() runtime.Value {
				expected, _ := datetime.DateAdd(
					context.Background(),
					mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
					runtime.NewInt(1),
					runtime.NewString("day"),
				)
				return expected
			}(),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(1),
				runtime.NewString("day"),
			},
		},
		&testCase{
			Name:     "-1 day",
			Expected: mustDefaultLayoutDt("1999-02-06T15:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(-1),
				runtime.NewString("day"),
			},
		},
		&testCase{
			Name:     "+3 months",
			Expected: mustDefaultLayoutDt("1999-05-07T15:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(3),
				runtime.NewString("months"),
			},
		},
		&testCase{
			Name:     "+5 years",
			Expected: mustLayoutDt("2006-01-02", "2004-02-07"),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(5),
				runtime.NewString("y"),
			},
		},
		&testCase{
			Name: "1999 minus 2000 years",
			Expected: runtime.NewDateTime(
				time.Date(-1, 2, 7, 0, 0, 0, 0, utcLoc),
			),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(-2000),
				runtime.NewString("year"),
			},
		},
		&testCase{
			Name:     "+2 hours",
			Expected: mustDefaultLayoutDt("1999-02-07T17:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(2),
				runtime.NewString("h"),
			},
		},
		&testCase{
			Name:     "+20 minutes",
			Expected: mustDefaultLayoutDt("1999-02-07T15:24:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(20),
				runtime.NewString("i"),
			},
		},
		&testCase{
			Name:     "+30 seconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:35Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(30),
				runtime.NewString("s"),
			},
		},
		&testCase{
			Name:     "+1000 milliseconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:06Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(1000),
				runtime.NewString("f"),
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
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
				runtime.NewInt(0),
				runtime.NewInt(0),
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When less than 3 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When incorrect arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("bla-bla"),
				runtime.NewInt(0),
				runtime.NewString("be-be"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:     "When wrong unit given",
			Expected: runtime.None,
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(5),
				runtime.NewString("not_exist"),
			},
			ShouldErr: true,
		},
		&testCase{
			Name: "When argument have correct types",
			Expected: func() runtime.Value {
				expected, _ := datetime.DateSubtract(
					context.Background(),
					mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
					runtime.NewInt(1),
					runtime.NewString("day"),
				)
				return expected
			}(),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(1),
				runtime.NewString("day"),
			},
		},
		&testCase{
			Name:     "-1 day",
			Expected: mustDefaultLayoutDt("1999-02-08T15:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(-1),
				runtime.NewString("day"),
			},
		},
		&testCase{
			Name:     "+3 months",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-05-07T15:04:05Z"),
				runtime.NewInt(3),
				runtime.NewString("months"),
			},
		},
		&testCase{
			Name:     "+5 years",
			Expected: mustLayoutDt("2006-01-02", "1994-02-07"),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(5),
				runtime.NewString("y"),
			},
		},
		&testCase{
			Name: "1999 minus 2000 years",
			Expected: runtime.NewDateTime(
				time.Date(-1, 2, 7, 0, 0, 0, 0, utcLoc),
			),
			Args: []runtime.Value{
				mustLayoutDt("2006-01-02", "1999-02-07"),
				runtime.NewInt(2000),
				runtime.NewString("year"),
			},
		},
		&testCase{
			Name:     "-2 hours",
			Expected: mustDefaultLayoutDt("1999-02-07T13:04:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(2),
				runtime.NewString("h"),
			},
		},
		&testCase{
			Name:     "-20 minutes",
			Expected: mustDefaultLayoutDt("1999-02-07T14:44:05Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(20),
				runtime.NewString("i"),
			},
		},
		&testCase{
			Name:     "-30 seconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:03:35Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(30),
				runtime.NewString("s"),
			},
		},
		&testCase{
			Name:     "-1000 milliseconds",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:04Z"),
			Args: []runtime.Value{
				mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
				runtime.NewInt(1000),
				runtime.NewString("f"),
			},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateSubtract)
	}
}
