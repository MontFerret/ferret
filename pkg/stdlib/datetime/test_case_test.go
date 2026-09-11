package datetime_test

import (
	"context"
	"reflect"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type testCase struct {
	Name      string
	Expected  runtime.Value
	Args      []runtime.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn runtime.Function) {
	Convey(tc.Name, t, func() {
		actual, err := fn(context.Background(), tc.Args...)

		if tc.ShouldErr {
			So(err, ShouldNotBeNil)
			So(actual, ShouldEqual, runtime.None)
		} else {
			So(err, ShouldBeNil)
			So(reflect.TypeOf(actual), ShouldEqual, reflect.TypeOf(tc.Expected))
			equal, equalityErr := runtime.EqualValues(context.Background(), actual, tc.Expected)
			So(equalityErr, ShouldBeNil)
			So(equal, ShouldEqual, runtime.True)
		}
	})
}
