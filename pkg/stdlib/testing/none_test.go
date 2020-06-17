package testing_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func TestNone(t *t.T) {
	None := base.NewPositiveAssertion(testing.None)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := None(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not none", t, func() {
		Convey("It should return an error", func() {
			_, err := None(context.Background(), values.NewString("true"))

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'true' to be [none] 'none'")
		})
	})

	Convey("When arg is none", t, func() {
		Convey("It should not return an error", func() {
			_, err := None(context.Background(), values.None)

			So(err, ShouldBeNil)
		})
	})
}

func TestNotNone(t *t.T) {
	NotNone := base.NewNegativeAssertion(testing.None)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotNone(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is none", t, func() {
		Convey("It should return an error", func() {
			_, err := NotNone(context.Background(), values.None)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [none] 'none' not to be [none] 'none'")
		})
	})

	Convey("When arg is not none", t, func() {
		Convey("It should return an error", func() {
			_, err := NotNone(context.Background(), values.False)

			So(err, ShouldBeNil)
		})
	})
}
