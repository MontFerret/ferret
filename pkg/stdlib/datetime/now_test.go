package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestNow(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When too many arguments",
			Expected: runtime.None,
			Args: []runtime.Value{
				runtime.NewCurrentDateTime(),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Now)
	}
}
