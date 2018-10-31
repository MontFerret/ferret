package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/datetime"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

func TestNow(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When too many arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewCurrentDateTime(),
			},
			ShouldErr: true,
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.Now)
	}
}
