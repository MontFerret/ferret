package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

func TestFloat(t *t.T) {
	Float := base.NewPositiveAssertion(testing.Float)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Float(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not float", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := Float(context.Background(), runtime.NewString("1.5"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected String '1.5' to be Float")
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should return an error", func() {
				_, err := Float(context.Background(), runtime.True)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Boolean 'true' to be Float")
			})
		})

		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := Float(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Int '1' to be Float")
			})
		})
	})

	Convey("When arg is float", t, func() {
		Convey("When arg is zero", func() {
			Convey("It should not return an error", func() {
				_, err := Float(context.Background(), runtime.NewFloat(0.0))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is positive", func() {
			Convey("It should not return an error", func() {
				_, err := Float(context.Background(), runtime.NewFloat(3.14))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is negative", func() {
			Convey("It should not return an error", func() {
				_, err := Float(context.Background(), runtime.NewFloat(-2.5))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotFloat(t *t.T) {
	NotFloat := base.NewNegativeAssertion(testing.Float)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotFloat(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is float", t, func() {
		Convey("When arg is zero", func() {
			Convey("It should return an error", func() {
				_, err := NotFloat(context.Background(), runtime.NewFloat(0.0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Float '0' not to be Float")
			})
		})

		Convey("When arg is positive", func() {
			Convey("It should return an error", func() {
				_, err := NotFloat(context.Background(), runtime.NewFloat(3.14))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Float '3.14' not to be Float")
			})
		})
	})

	Convey("When arg is not float", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotFloat(context.Background(), runtime.NewString("1.5"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should not return an error", func() {
				_, err := NotFloat(context.Background(), runtime.True)

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotFloat(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})
	})
}
