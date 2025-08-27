package testing_test

import (
	"context"
	t "testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestDateTime(t *t.T) {
	DateTime := base.NewPositiveAssertion(testing.DateTime)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := DateTime(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not datetime", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := DateTime(context.Background(), runtime.NewString("2023-01-01"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] '2023-01-01' to be datetime")
			})
		})

		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := DateTime(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '1' to be datetime")
			})
		})

		Convey("When arg is array", func() {
			Convey("It should return an error", func() {
				_, err := DateTime(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [array] '[]' to be datetime")
			})
		})
	})

	Convey("When arg is datetime", t, func() {
		Convey("When arg is current datetime", func() {
			Convey("It should not return an error", func() {
				now := time.Now()
				_, err := DateTime(context.Background(), runtime.NewDateTime(now))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is zero datetime", func() {
			Convey("It should not return an error", func() {
				_, err := DateTime(context.Background(), runtime.ZeroDateTime)

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotDateTime(t *t.T) {
	NotDateTime := base.NewNegativeAssertion(testing.DateTime)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotDateTime(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is datetime", t, func() {
		Convey("When arg is current datetime", func() {
			Convey("It should return an error", func() {
				now := time.Now()
				dt := runtime.NewDateTime(now)
				_, err := NotDateTime(context.Background(), dt)

				So(err, ShouldBeError)
				So(err.Error(), ShouldContainSubstring, base.ErrAssertion.Error()+": expected [date_time] '")
				So(err.Error(), ShouldContainSubstring, "' not to be datetime")
			})
		})

		Convey("When arg is zero datetime", func() {
			Convey("It should return an error", func() {
				_, err := NotDateTime(context.Background(), runtime.ZeroDateTime)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [date_time] '0001-01-01 00:00:00 +0000 UTC' not to be datetime")
			})
		})
	})

	Convey("When arg is not datetime", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotDateTime(context.Background(), runtime.NewString("2023-01-01"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotDateTime(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is array", func() {
			Convey("It should not return an error", func() {
				_, err := NotDateTime(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeNil)
			})
		})
	})
}