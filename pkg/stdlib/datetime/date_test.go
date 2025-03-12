package datetime_test

import (
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDate(t *testing.T) {
	now := time.Now()

	tcs := []*testCase{
		{
			Name:     "When more than 2 arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewString(time.Now().Format(time.RFC3339)),
				core.NewString(time.RFC3339),
				core.NewString("foo"),
			},
			ShouldErr: true,
		},
		{
			Name:      "When 0 arguments",
			Expected:  core.None,
			Args:      []core.Value{},
			ShouldErr: true,
		},
		{
			Name:      "When first argument isn't string",
			Expected:  core.None,
			Args:      []core.Value{core.NewInt(0)},
			ShouldErr: true,
		},
		{
			Name:     "When incorrect timeStrings",
			Expected: core.None,
			Args: []core.Value{
				core.NewString("bla-bla"),
			},
			ShouldErr: true,
		},
		{
			Name:     "When string is in default format",
			Expected: mustDefaultLayoutDt("1999-02-07T15:04:05Z"),
			Args: []core.Value{
				core.NewString("1999-02-07T15:04:05Z"),
			},
		},
		{
			Name:     "When second argument isn't string",
			Expected: core.None,
			Args: []core.Value{
				core.NewString("1999-02-07T15:04:05Z"),
				core.NewInt(1),
			},
			ShouldErr: true,
		},
		{
			Name:     "When string is in custom format",
			Expected: mustLayoutDt(time.RFC822, now.Format(time.RFC822)),
			Args: []core.Value{
				core.NewString(now.Format(time.RFC822)),
				core.NewString(time.RFC822),
			},
			ShouldErr: false,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Date)
	}
}
