package testing_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestEmpty(t *t.T) {
	Empty := base.NewPositiveAssertion(testing.Empty)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background())

			So(err, ShouldBeError)

			_, err = Empty(context.Background(), core.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := Empty(context.Background(), core.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should return an error", func() {
				_, err := Empty(context.Background(), core.NewString("Foo"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [string] 'Foo' to be empty").Error())
			})
		})

		Convey("When ''", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(context.Background(), core.NewString(""))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When arg is an array", t, func() {
		Convey("When [1,2,3]", func() {
			Convey("It should return an error", func() {
				_, err := Empty(
					context.Background(),
					internal.NewArrayWith(core.NewInt(1), core.NewInt(2), core.NewInt(3)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [array] '[1,2,3]' to be empty").Error())
			})
		})

		Convey("When []", func() {
			Convey("It should not return an error", func() {
				_, err := Empty(
					context.Background(),
					internal.NewArray(0),
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
					internal.NewObjectWith(
						internal.NewObjectProperty("a", core.NewInt(1)),
						internal.NewObjectProperty("b", core.NewInt(2)),
						internal.NewObjectProperty("c", core.NewInt(3)),
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
					internal.NewObject(),
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

			_, err = NotEmpty(context.Background(), core.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg are not measurable", t, func() {
		Convey("It should return an error", func() {
			_, err := NotEmpty(context.Background(), core.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is a string", t, func() {
		Convey("When 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(context.Background(), core.NewString("Foo"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When ''", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(context.Background(), core.NewString(""))

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
					internal.NewArrayWith(core.NewInt(1), core.NewInt(2), core.NewInt(3)),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When []", func() {
			Convey("It should return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					internal.NewArray(0),
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
					internal.NewObjectWith(
						internal.NewObjectProperty("a", core.NewInt(1)),
						internal.NewObjectProperty("b", core.NewInt(2)),
						internal.NewObjectProperty("c", core.NewInt(3)),
					),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When {}", func() {
			Convey("It should not return an error", func() {
				_, err := NotEmpty(
					context.Background(),
					internal.NewObject(),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [object] '{}' not to be empty").Error())
			})
		})
	})
}
