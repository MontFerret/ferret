package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"

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
			_, err := None(context.Background(), core.NewString("true"))

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'true' to be [none] 'none'")
		})
	})

	Convey("When arg is none", t, func() {
		Convey("It should not return an error", func() {
			_, err := None(context.Background(), core.None)

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
			_, err := NotNone(context.Background(), core.None)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [none] 'none' not to be [none] 'none'")
		})
	})

	Convey("When arg is not none", t, func() {
		Convey("It should return an error", func() {
			_, err := NotNone(context.Background(), core.False)

			So(err, ShouldBeNil)
		})
	})
}
