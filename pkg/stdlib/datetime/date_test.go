package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDate(t *testing.T) {
	now := time.Now()

	tcs := []*testCase{
		{
			Name:     "When more than 2 arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString(time.Now().Format(time.RFC3339)),
				runtime.NewString(time.RFC3339),
				runtime.NewString("foo"),
			},
			ShouldErr: true,
		},
		{
			Name:      "When 0 arguments",
			Expected:  runtime.None,
			Args:      []runtime.Value{},
			ShouldErr: true,
		},
		{
			Name:      "When first argument isn't string",
			Expected:  runtime.None,
			Args:      []runtime.Value{runtime.NewInt(0)},
			ShouldErr: true,
		},
		{
			Name:     "When incorrect timeStrings",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("bla-bla"),
			},
			ShouldErr: true,
		},
		{
			Name:     "When string is in default format",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
			Args: []runtime.Value{
				runtime.NewString("1999-02-07T15:04:05Z"),
			},
		},
		{
			Name:     "When second argument isn't string",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewString("1999-02-07T15:04:05Z"),
				runtime.NewInt(1),
			},
			ShouldErr: true,
		},
		{
			Name:     "When string is in custom format",
			Expected: mustLayoutDt(time.RFC822, now.Format(time.RFC822)),
			Args: []runtime.Value{
				runtime.NewString(now.Format(time.RFC822)),
				runtime.NewString(time.RFC822),
			},
			ShouldErr: false,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Date)
	}
}
