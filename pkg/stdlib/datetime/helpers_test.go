package datetime_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
)

type testCase struct {
	Name      string
	Expected  runtime.Value
	Args      []runtime.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn runtime.Function) {
	//Convey(tc.Name, t, func() {
	//	expected := tc.Expected
	//
	//	actual, err := fn(context.Background(), tc.Args...)
	//
	//	if tc.ShouldErr {
	//		So(err, ShouldBeError)
	//		expected = runtime.None
	//	} else {
	//		So(err, ShouldBeNil)
	//	}
	//
	//	So(actual.Type().Equals(expected.Type()), ShouldBeTrue)
	//	So(actual.Compare(expected), ShouldEqual, 0)
	//})
}

func mustDefaultLayoutDt(timeString string) runtime.DateTime {
	dt, err := defaultLayoutDt(timeString)

	if err != nil {
		panic(err)
	}

	return dt
}

func mustLayoutDt(layout, value string) runtime.DateTime {
	dt, err := layoutDt(layout, value)

	if err != nil {
		panic(err)
	}

	return dt
}

func defaultLayoutDt(timeString string) (runtime.DateTime, error) {
	return layoutDt(runtime.DefaultTimeLayout, timeString)
}

func layoutDt(layout, value string) (runtime.DateTime, error) {
	t, err := time.Parse(layout, value)

	if err != nil {
		return runtime.DateTime{}, err
	}

	return runtime.NewDateTime(t), nil
}
