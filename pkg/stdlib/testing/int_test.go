package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestInt(t *t.T) {
	Int := base.NewPositiveAssertion(testing.Int)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Int(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not int", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := Int(context.Background(), runtime.NewString("1"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] '1' to be int")
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should return an error", func() {
				_, err := Int(context.Background(), runtime.True)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'true' to be int")
			})
		})

		Convey("When arg is float", func() {
			Convey("It should return an error", func() {
				_, err := Int(context.Background(), runtime.NewFloat(1.5))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [float] '1.5' to be int")
			})
		})
	})

	Convey("When arg is int", t, func() {
		Convey("When arg is zero", func() {
			Convey("It should not return an error", func() {
				_, err := Int(context.Background(), runtime.NewInt(0))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is positive", func() {
			Convey("It should not return an error", func() {
				_, err := Int(context.Background(), runtime.NewInt(42))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is negative", func() {
			Convey("It should not return an error", func() {
				_, err := Int(context.Background(), runtime.NewInt(-10))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotInt(t *t.T) {
	NotInt := base.NewNegativeAssertion(testing.Int)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotInt(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is int", t, func() {
		Convey("When arg is zero", func() {
			Convey("It should return an error", func() {
				_, err := NotInt(context.Background(), runtime.NewInt(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '0' not to be int")
			})
		})

		Convey("When arg is positive", func() {
			Convey("It should return an error", func() {
				_, err := NotInt(context.Background(), runtime.NewInt(42))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '42' not to be int")
			})
		})
	})

	Convey("When arg is not int", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotInt(context.Background(), runtime.NewString("1"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should not return an error", func() {
				_, err := NotInt(context.Background(), runtime.True)

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is float", func() {
			Convey("It should not return an error", func() {
				_, err := NotInt(context.Background(), runtime.NewFloat(1.5))

				So(err, ShouldBeNil)
			})
		})
	})
}
