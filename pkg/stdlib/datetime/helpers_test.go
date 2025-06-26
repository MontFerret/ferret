package datetime_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
)

type testCase struct {
	Name      string
	Expected  core.Value
	Args      []core.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn core.Function) {
	//Convey(tc.Name, t, func() {
	//	expected := tc.Expected
	//
	//	actual, err := fn(context.Background(), tc.Args...)
	//
	//	if tc.ShouldErr {
	//		So(err, ShouldBeError)
	//		expected = core.None
	//	} else {
	//		So(err, ShouldBeNil)
	//	}
	//
	//	So(actual.Type().Equals(expected.Type()), ShouldBeTrue)
	//	So(actual.Compare(expected), ShouldEqual, 0)
	//})
}

func mustDefaultLayoutDt(timeString string) core.DateTime {
	dt, err := defaultLayoutDt(timeString)

	if err != nil {
		panic(err)
	}

	return dt
}

func mustLayoutDt(layout, value string) core.DateTime {
	dt, err := layoutDt(layout, value)

	if err != nil {
		panic(err)
	}

	return dt
}

func defaultLayoutDt(timeString string) (core.DateTime, error) {
	return layoutDt(core.DefaultTimeLayout, timeString)
}

func layoutDt(layout, value string) (core.DateTime, error) {
	t, err := time.Parse(layout, value)

	if err != nil {
		return core.DateTime{}, err
	}

	return core.NewDateTime(t), nil
}
