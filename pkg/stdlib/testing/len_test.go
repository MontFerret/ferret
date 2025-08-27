package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
)

func TestLen(t *t.T) {
	Len := base.NewPositiveAssertion(testing.Len)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Len(context.Background())

			So(err, ShouldBeError)

			_, err = Len(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := Len(context.Background(), runtime.NewInt(1), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo' should have length 1", func() {
			Convey("It should return an error", func() {
				_, err := Len(context.Background(), runtime.NewString("Foo"), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' to has size 1").Error())
			})
		})

		Convey("When 'Foo' should have length 3", func() {
			Convey("It should not return an error", func() {
				_, err := Len(context.Background(), runtime.NewString("Foo"), runtime.NewInt(3))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3] should have length 1", func() {
			Convey("It should return an error", func() {
				_, err := Len(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(1),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[1,2,3]' to has size 1").Error())
			})
		})

		Convey("When [1,2,3] should have length 3", func() {
			Convey("It should not return an error", func() {
				_, err := Len(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(3),
				)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an object", t, func() {
		Convey("When { a: 1, b: 2, c: 3 } should have length 1", func() {
			Convey("It should return an error", func() {
				_, err := Len(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(1),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [map] '{\"a\":1,\"b\":2,\"c\":3}' to has size 1").Error())
			})
		})

		Convey("When [1,2,3] should have length 3", func() {
			Convey("It should not return an error", func() {
				_, err := Len(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(3),
				)

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotLen(t *t.T) {
	NotLen := base.NewNegativeAssertion(testing.Len)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotLen(context.Background())

			So(err, ShouldBeError)

			_, err = NotLen(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := NotLen(context.Background(), runtime.NewInt(1), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo' should not have length 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLen(context.Background(), runtime.NewString("Foo"), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 'Foo' should not have length 3", func() {
			Convey("It should return an error", func() {
				_, err := NotLen(context.Background(), runtime.NewString("Foo"), runtime.NewInt(3))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' not to has size 3").Error())
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3] should not have length 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLen(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(1),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When [1,2,3] should have length 3", func() {
			Convey("It should return an error", func() {
				_, err := NotLen(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(3),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[1,2,3]' not to has size 3").Error())
			})
		})
	})

	Convey("When arg is an object", t, func() {
		Convey("When { a: 1, b: 2, c: 3 } should have length 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotLen(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(1),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When [1,2,3] should have length 3", func() {
			Convey("It should not return an error", func() {
				_, err := NotLen(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(3),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [map] '{\"a\":1,\"b\":2,\"c\":3}' not to has size 3").Error())
			})
		})
	})
}
