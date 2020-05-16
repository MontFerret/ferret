package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func TestLt(t *t.T) {
	Lt := testing.NewPositive(testing.Lt)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Lt(context.Background())

			So(err, ShouldBeError)

			_, err = Lt(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 2 and 1", func() {
			Convey("It should return an error", func() {
				_, err := Lt(context.Background(), values.NewInt(2), values.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(testing.ErrAssertion, "expected [int] '2' to be less than [int] '1'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should return an error", func() {
				_, err := Lt(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(testing.ErrAssertion, "expected [int] '1' to be less than [int] '1'").Error())
			})
		})

		Convey("When 1 and 2", func() {
			Convey("It should not return an error", func() {
				_, err := Lt(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotLt(t *t.T) {
	NotLt := testing.NewNegative(testing.Lt)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotLt(context.Background())

			So(err, ShouldBeError)

			_, err = NotLt(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotLt(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(testing.ErrAssertion, "expected [int] '1' not to be less than [int] '2'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLt(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 2 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLt(context.Background(), values.NewInt(2), values.NewInt(1))

				So(err, ShouldBeNil)
			})
		})
	})
}
