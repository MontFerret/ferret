package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestFalse(t *t.T) {
	False := base.NewPositiveAssertion(testing.False)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := False(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should return an error", func() {
			_, err := False(context.Background(), values.NewString("false"))

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'false' to be [boolean] 'false'")
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should return an error", func() {
			_, err := False(context.Background(), values.True)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'true' to be [boolean] 'false'")
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should not return an error", func() {
			_, err := False(context.Background(), values.False)

			So(err, ShouldBeNil)
		})
	})
}

func TestNotFalse(t *t.T) {
	NotFalse := base.NewNegativeAssertion(testing.False)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotFalse(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should not return an error", func() {
			_, err := NotFalse(context.Background(), values.NewString("false"))

			So(err, ShouldBeNil)
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should return an error", func() {
			_, err := NotFalse(context.Background(), values.False)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'false' not to be [boolean] 'false'")
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should return an error", func() {
			_, err := NotFalse(context.Background(), values.True)

			So(err, ShouldBeNil)
		})
	})
}
