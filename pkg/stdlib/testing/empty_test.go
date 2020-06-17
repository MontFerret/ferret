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

func TestEmpty(t *t.T) {
	Empty := base.NewPositiveAssertion(testing.Empty)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background())

			So(err, ShouldBeError)

			_, err = Empty(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should return an error", func() {
				_, err := Empty(context.Background(), values.NewString("Foo"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [string] 'Foo' to be empty").Error())
			})
		})

		Convey("When ''", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(context.Background(), values.NewString(""))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3]", func() {
			Convey("It should return an error", func() {
				_, err := Empty(
					context.Background(),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2), values.NewInt(3)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [array] '[1,2,3]' to be empty").Error())
			})
		})

		Convey("When []", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(
					context.Background(),
					values.NewArray(0),
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
					values.NewObjectWith(
						values.NewObjectProperty("a", values.NewInt(1)),
						values.NewObjectProperty("b", values.NewInt(2)),
						values.NewObjectProperty("c", values.NewInt(3)),
					),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [object] '{\"a\":1,\"b\":2,\"c\":3}' to be empty").Error())
			})
		})

		Convey("When {}", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(
					context.Background(),
					values.NewObject(),
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

			_, err = NotEmpty(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := NotEmpty(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(context.Background(), values.NewString("Foo"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When ''", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(context.Background(), values.NewString(""))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [string] '' not to be empty").Error())
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3]", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2), values.NewInt(3)),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When []", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					values.NewArray(0),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [array] '[]' not to be empty").Error())
			})
		})
	})

	Convey("When arg is an object", t, func() {
		Convey("When { a: 1, b: 2, c: 3 }", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					values.NewObjectWith(
						values.NewObjectProperty("a", values.NewInt(1)),
						values.NewObjectProperty("b", values.NewInt(2)),
						values.NewObjectProperty("c", values.NewInt(3)),
					),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When {}", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					values.NewObject(),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [object] '{}' not to be empty").Error())
			})
		})
	})
}
