package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/testing/base"
)

func TestObject(t *t.T) {
	Object := base.NewPositiveAssertion(testing.Object)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Object(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not object", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := Object(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected String 'test' to be Object")
			})
		})

		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := Object(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Int '1' to be Object")
			})
		})

		Convey("When arg is array", func() {
			Convey("It should return an error", func() {
				_, err := Object(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Array '[]' to be Object")
			})
		})
	})

	Convey("When arg is object", t, func() {
		Convey("When arg is empty object", func() {
			Convey("It should not return an error", func() {
				_, err := Object(context.Background(), runtime.NewObject())

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is non-empty object", func() {
			Convey("It should not return an error", func() {
				_, err := Object(context.Background(), runtime.NewObjectWith(
					map[string]runtime.Value{
						"key": runtime.NewString("value"),
					},
				))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotObject(t *t.T) {
	NotObject := base.NewNegativeAssertion(testing.Object)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotObject(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is object", t, func() {
		Convey("When arg is empty object", func() {
			Convey("It should return an error", func() {
				_, err := NotObject(context.Background(), runtime.NewObject())

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Object '{}' not to be Object")
			})
		})

		Convey("When arg is non-empty object", func() {
			Convey("It should return an error", func() {
				_, err := NotObject(context.Background(), runtime.NewObjectWith(
					map[string]runtime.Value{
						"key": runtime.NewString("value"),
					},
				))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected Object '{\"key\":\"value\"}' not to be Object")
			})
		})
	})

	Convey("When arg is not object", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotObject(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotObject(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is array", func() {
			Convey("It should not return an error", func() {
				_, err := NotObject(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeNil)
			})
		})
	})
}
