package datetime_test

import (
	"context"
	"testing"
	"time"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Name      string
	Expected  core.Value
<<<<<<< HEAD
=======
	TimeArg   time.Time
>>>>>>> 94dc1c6... add helpers functions
	Args      []core.Value
	ShouldErr bool
}

func (tc *testCase) Do(t *testing.T, fn core.Function) {
	Convey(tc.Name, t, func() {
		expected := tc.Expected

		actual, err := fn(context.Background(), tc.Args...)

		if tc.ShouldErr {
			So(err, ShouldBeError)
			expected = values.None
		} else {
			So(err, ShouldBeNil)
		}

<<<<<<< HEAD
		So(actual.Type().Equals(expected.Type()), ShouldBeTrue)
=======
>>>>>>> 94dc1c6... add helpers functions
		So(actual.Compare(expected), ShouldEqual, 0)
	})
}

func mustDefaultLayoutDt(timeString string) values.DateTime {
	dt, err := defaultLayoutDt(timeString)
<<<<<<< HEAD

	if err != nil {
		panic(err)
	}

=======
	if err != nil {
		panic(err)
	}
>>>>>>> 94dc1c6... add helpers functions
	return dt
}

func mustLayoutDt(layout, value string) values.DateTime {
	dt, err := layoutDt(layout, value)
<<<<<<< HEAD

	if err != nil {
		panic(err)
	}

=======
	if err != nil {
		panic(err)
	}
>>>>>>> 94dc1c6... add helpers functions
	return dt
}

func defaultLayoutDt(timeString string) (values.DateTime, error) {
	return layoutDt(values.DefaultTimeLayout, timeString)
}

func layoutDt(layout, value string) (values.DateTime, error) {
	t, err := time.Parse(layout, value)
<<<<<<< HEAD

	if err != nil {
		return values.DateTime{}, err
	}

=======
	if err != nil {
		return values.DateTime{}, err
	}
>>>>>>> 94dc1c6... add helpers functions
	return values.NewDateTime(t), nil
}
