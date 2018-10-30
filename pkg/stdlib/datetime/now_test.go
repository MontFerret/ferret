package datetime_test

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Name      string
	Expected  core.Value
	TimeArg   time.Time
	Args      []core.Value
	ShouldErr bool
}

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
		tc.Do(t)
	}
}

func (tc *testCase) Do(t *testing.T) {
	Convey(tc.Name, t, func() {
		var expected core.Value

		expected = values.NewDateTime(tc.TimeArg)

		dt, err := datetime.Now(context.Background(), tc.Args...)

		if tc.ShouldErr {
			So(err, ShouldBeError)
			expected = values.None
		} else {
			So(err, ShouldBeNil)
		}

		So(dt, ShouldEqual, expected)
	})
}
