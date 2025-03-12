package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

func TestNow(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When too many arguments",
			Expected: core.None,
			Args: []core.Value{
				core.NewCurrentDateTime(),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Now)
	}
}
