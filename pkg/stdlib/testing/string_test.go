package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestString(t *t.T) {
	String := base.NewPositiveAssertion(testing.String)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := String(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not string", t, func() {
		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := String(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '1' to be string")
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should return an error", func() {
				_, err := String(context.Background(), runtime.True)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [boolean] 'true' to be string")
			})
		})

		Convey("When arg is array", func() {
			Convey("It should return an error", func() {
				_, err := String(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [array] '[]' to be string")
			})
		})
	})

	Convey("When arg is string", t, func() {
		Convey("When arg is empty string", func() {
			Convey("It should not return an error", func() {
				_, err := String(context.Background(), runtime.NewString(""))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is non-empty string", func() {
			Convey("It should not return an error", func() {
				_, err := String(context.Background(), runtime.NewString("hello"))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotString(t *t.T) {
	NotString := base.NewNegativeAssertion(testing.String)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotString(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is string", t, func() {
		Convey("When arg is empty string", func() {
			Convey("It should return an error", func() {
				_, err := NotString(context.Background(), runtime.NewString(""))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] '' not to be string")
			})
		})

		Convey("When arg is non-empty string", func() {
			Convey("It should return an error", func() {
				_, err := NotString(context.Background(), runtime.NewString("hello"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'hello' not to be string")
			})
		})
	})

	Convey("When arg is not string", t, func() {
		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotString(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is boolean", func() {
			Convey("It should not return an error", func() {
				_, err := NotString(context.Background(), runtime.True)

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is array", func() {
			Convey("It should not return an error", func() {
				_, err := NotString(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeNil)
			})
		})
	})
}