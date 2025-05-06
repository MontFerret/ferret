package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestTrue(t *t.T) {
	True := base.NewPositiveAssertion(testing.True)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background(), core.NewString("true"))

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'true' to be [boolean] 'true'")
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should return an error", func() {
			_, err := True(context.Background(), core.False)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'false' to be [boolean] 'true'")
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should not return an error", func() {
			_, err := True(context.Background(), core.True)

			So(err, ShouldBeNil)
		})
	})
}

func TestNotTrue(t *t.T) {
	NotTrue := base.NewNegativeAssertion(testing.True)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not boolean", t, func() {
		Convey("It should not return an error", func() {
			_, err := NotTrue(context.Background(), core.NewString("true"))

			So(err, ShouldBeNil)
		})
	})

	Convey("When arg is true", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background(), core.True)

			So(err, ShouldBeError)
			So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'true' not to be [boolean] 'true'")
		})
	})

	Convey("When arg is false", t, func() {
		Convey("It should return an error", func() {
			_, err := NotTrue(context.Background(), core.False)

			So(err, ShouldBeNil)
		})
	})
}
