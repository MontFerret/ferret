package testing_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func TestTrue(t *t.T) {
	True := testing.NewPositive(testing.True)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background(), values.NewString("true"))

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, testing.ErrAssertion.Error()+": expected [string] true to be [boolean] true")
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background(), values.False)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, testing.ErrAssertion.Error()+": expected [boolean] false to be [boolean] true")
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should not return an error", func() {
			_, err := True(context.Background(), values.True)

			So(err, ShouldBeNil)
		})
	})
}

func TestNotTrue(t *t.T) {
	NotTrue := testing.NewNegative(testing.True)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should not return an error", func() {
			_, err := NotTrue(context.Background(), values.NewString("true"))

			So(err, ShouldBeNil)
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background(), values.True)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, testing.ErrAssertion.Error()+": expected [boolean] true not to be [boolean] true")
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background(), values.False)

			So(err, ShouldBeNil)
		})
	})
}
