package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestInclude(t *t.T) {
	Include := base.NewPositiveAssertion(testing.Include)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Include(context.Background())

			So(err, ShouldBeError)

			_, err = Include(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When value is a string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := Include(context.Background(), runtime.NewString("Foo"), runtime.NewString("Bar"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' to include [string] 'Bar'").Error())
			})
		})

		Convey("When 'FooBar' and 'Bar'", func() {
			Convey("It should not return an error", func() {
				_, err := Include(context.Background(), runtime.NewString("FooBar"), runtime.NewString("Bar"))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When value is an array", t, func() {
		Convey("When [1,2,3] and 4", func() {
			Convey("It should return an error", func() {
				_, err := Include(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(4),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [array] '[1,2,3]' to include [int] '4'").Error())
			})
		})

		Convey("When [1,2,3] and 2", func() {
			Convey("It should not return an error", func() {
				_, err := Include(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(2),
				)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When value is an object", t, func() {
		Convey("When {a:1,b:2,c:3} and 4", func() {
			Convey("It should return an error", func() {
				_, err := Include(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(4),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [object] '{\"a\":1,\"b\":2,\"c\":3}' to include [int] '4'").Error())
			})
		})

		Convey("When {a:1,b:2,c:3} and 2", func() {
			Convey("It should not return an error", func() {
				_, err := Include(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(2),
				)

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotInclude(t *t.T) {
	NotInclude := base.NewNegativeAssertion(testing.Include)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotInclude(context.Background())

			So(err, ShouldBeError)

			_, err = NotInclude(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When value is a string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should not return an error", func() {
				_, err := NotInclude(context.Background(), runtime.NewString("Foo"), runtime.NewString("Bar"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 'FooBar' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := NotInclude(context.Background(), runtime.NewString("FooBar"), runtime.NewString("Bar"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'FooBar' not to include [string] 'Bar'").Error())
			})
		})
	})

	Convey("When value is an array", t, func() {
		Convey("When [1,2,3] and 4", func() {
			Convey("It should not return an error", func() {
				_, err := NotInclude(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(4),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When [1,2,3] and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotInclude(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
					runtime.NewInt(2),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [array] '[1,2,3]' not to include [int] '2'").Error())
			})
		})
	})

	Convey("When value is an object", t, func() {
		Convey("When {a:1,b:2,c:3} and 4", func() {
			Convey("It should not return an error", func() {
				_, err := NotInclude(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(4),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When {a:1,b:2,c:3} and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotInclude(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
					runtime.NewInt(2),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [object] '{\"a\":1,\"b\":2,\"c\":3}' not to include [int] '2'").Error())
			})
		})
	})
}
