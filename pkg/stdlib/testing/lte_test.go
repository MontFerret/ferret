package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestLte(t *t.T) {
	Lte := base.NewPositiveAssertion(testing.Lte)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Lte(context.Background())

			So(err, ShouldBeError)

			_, err = Lte(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 2 and 1", func() {
			Convey("It should return an error", func() {
				_, err := Lte(context.Background(), values.NewInt(2), values.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [int] '2' to be less than or equal to [int] '1'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := Lte(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 1 and 2", func() {
			Convey("It should not return an error", func() {
				_, err := Lte(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotLte(t *t.T) {
	NotLte := base.NewNegativeAssertion(testing.Lte)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotLte(context.Background())

			So(err, ShouldBeError)

			_, err = NotLte(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotLte(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [int] '1' not to be less than or equal to [int] '2'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should return an error", func() {
				_, err := NotLte(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [int] '1' not to be less than or equal to [int] '1'").Error())
			})
		})

		Convey("When 2 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLte(context.Background(), values.NewInt(2), values.NewInt(1))

				So(err, ShouldBeNil)
			})
		})
	})
}
