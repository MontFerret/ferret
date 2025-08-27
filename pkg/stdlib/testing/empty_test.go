package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestEmpty(t *t.T) {
	Empty := base.NewPositiveAssertion(testing.Empty)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background())

			So(err, ShouldBeError)

			_, err = Empty(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should return an error", func() {
				_, err := Empty(context.Background(), runtime.NewString("Foo"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' to be empty").Error())
			})
		})

		Convey("When ''", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(context.Background(), runtime.NewString(""))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3]", func() {
			Convey("It should return an error", func() {
				_, err := Empty(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[1,2,3]' to be empty").Error())
			})
		})

		Convey("When []", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(
					context.Background(),
					runtime.NewArray(0),
				)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an object", t, func() {
		Convey("When { a: 1, b: 2, c: 3 }", func() {
			Convey("It should return an error", func() {
				_, err := Empty(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [map] '{\"a\":1,\"b\":2,\"c\":3}' to be empty").Error())
			})
		})

		Convey("When {}", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(
					context.Background(),
					runtime.NewObject(),
				)

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotEmpty(t *t.T) {
	NotEmpty := base.NewNegativeAssertion(testing.Empty)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotEmpty(context.Background())

			So(err, ShouldBeError)

			_, err = NotEmpty(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := NotEmpty(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(context.Background(), runtime.NewString("Foo"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When ''", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(context.Background(), runtime.NewString(""))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] '' not to be empty").Error())
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3]", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2), runtime.NewInt(3)),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When []", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					runtime.NewArray(0),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[]' not to be empty").Error())
			})
		})
	})

	Convey("When arg is an object", t, func() {
		Convey("When { a: 1, b: 2, c: 3 }", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					runtime.NewObjectWith(
						runtime.NewObjectProperty("a", runtime.NewInt(1)),
						runtime.NewObjectProperty("b", runtime.NewInt(2)),
						runtime.NewObjectProperty("c", runtime.NewInt(3)),
					),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When {}", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					runtime.NewObject(),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [map] '{}' not to be empty").Error())
			})
		})
	})
}
