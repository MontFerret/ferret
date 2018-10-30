package datetime_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/datetime"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Name      string
	Expected  values.DateTime
	Args      []core.Value
	ShouldErr bool
}

func TestNow(t *testing.T) {

}

func (tc *testCase) Do(t *testing.T) {
	Convey(tc.Name, t, func() {
		var expected core.Value

		expected = values.NewDateTime(tc.Expected.Time)

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
