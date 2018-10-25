package datetime_test

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"
)

func TestDateYear(t *testing.T) {
	tcs := []*testCase{
		&testCase{
			Name:     "When more than 1 arguments",
			Expected: values.None,
			Args: []core.Value{
				values.NewString("string"),
				values.NewInt(0),
			},
			ShouldErr: true,
		},
		&testCase{
			Name:      "When 0 arguments",
			Expected:  values.None,
			Args:      []core.Value{},
			ShouldErr: true,
		},
		&testCase{
<<<<<<< HEAD
			Name:      "When argument isn't DateTime",
			Expected:  values.None,
			Args:      []core.Value{values.NewInt(0)},
			ShouldErr: true,
		},
		&testCase{
=======
>>>>>>> db2ab4b... added DATE_YEAR function
			Name:     "When 1999th year",
			Expected: values.NewInt(1999),
			Args:     []core.Value{mustDefaultLayoutDt("1999-02-07T15:04:05Z")},
		},
		&testCase{
			Name:     "When 1629th year",
			Expected: values.NewInt(1629),
			Args:     []core.Value{mustDefaultLayoutDt("1629-02-08T15:04:05Z")},
		},
	}

	for _, tc := range tcs {
		tc.Do(t, datetime.DateYear)
	}
}
