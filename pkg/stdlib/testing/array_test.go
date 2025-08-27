package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestArray(t *t.T) {
	Array := base.NewPositiveAssertion(testing.Array)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Array(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not array", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := Array(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'test' to be array")
			})
		})

		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := Array(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '1' to be array")
			})
		})

		Convey("When arg is object", func() {
			Convey("It should return an error", func() {
				_, err := Array(context.Background(), runtime.NewObject())

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [map] '{}' to be array")
			})
		})
	})

	Convey("When arg is array", t, func() {
		Convey("When arg is empty array", func() {
			Convey("It should not return an error", func() {
				_, err := Array(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is non-empty array", func() {
			Convey("It should not return an error", func() {
				_, err := Array(context.Background(), runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotArray(t *t.T) {
	NotArray := base.NewNegativeAssertion(testing.Array)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotArray(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is array", t, func() {
		Convey("When arg is empty array", func() {
			Convey("It should return an error", func() {
				_, err := NotArray(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [list] '[]' not to be array")
			})
		})

		Convey("When arg is non-empty array", func() {
			Convey("It should return an error", func() {
				_, err := NotArray(context.Background(), runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [list] '[1,2]' not to be array")
			})
		})
	})

	Convey("When arg is not array", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotArray(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotArray(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is object", func() {
			Convey("It should not return an error", func() {
				_, err := NotArray(context.Background(), runtime.NewObject())

				So(err, ShouldBeNil)
			})
		})
	})
}