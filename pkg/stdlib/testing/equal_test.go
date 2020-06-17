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

func TestEqual(t *t.T) {
	Equal := base.NewPositiveAssertion(testing.Equal)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Equal(context.Background())

			So(err, ShouldBeError)

			_, err = Equal(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), values.NewString("Foo"), values.NewString("Bar"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [string] 'Foo' to be equal to [string] 'Bar'").Error())
			})
		})

		Convey("When 'Foo' and 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), values.NewString("Foo"), values.NewString("Foo"))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [int] '1' to be equal to [int] '2'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are boolean", t, func() {
		Convey("When False and True", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), values.False, values.True)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [boolean] 'false' to be equal to [boolean] 'true'").Error())
			})
		})

		Convey("When False and False", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), values.False, values.False)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are arrays", t, func() {
		Convey("When [1] and [1,2]", func() {
			Convey("It should return an error", func() {
				_, err := Equal(
					context.Background(),
					values.NewArrayWith(values.NewInt(1)),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [array] '[1]' to be equal to [array] '[1,2]'").Error())
			})
		})

		Convey("When [1,2] and [1,2]", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(
					context.Background(),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
				)

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotEqual(t *t.T) {
	NotEqual := base.NewNegativeAssertion(testing.Equal)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotEqual(context.Background())

			So(err, ShouldBeError)

			_, err = NotEqual(context.Background(), values.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), values.NewString("Foo"), values.NewString("Bar"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 'Foo' and 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), values.NewString("Foo"), values.NewString("Foo"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [string] 'Foo' not to be equal to [string] 'Foo'").Error())
			})
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), values.NewInt(1), values.NewInt(2))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), values.NewInt(1), values.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [int] '1' not to be equal to [int] '1'").Error())
			})
		})
	})

	Convey("When args are boolean", t, func() {
		Convey("When False and True", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), values.False, values.True)

				So(err, ShouldBeNil)
			})
		})

		Convey("When False and False", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), values.False, values.False)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [boolean] 'false' not to be equal to [boolean] 'false'").Error())
			})
		})
	})

	Convey("When args are arrays", t, func() {
		Convey("When [1] and [1,2]", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(
					context.Background(),
					values.NewArrayWith(values.NewInt(1)),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When [1,2] and [1,2]", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(
					context.Background(),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
					values.NewArrayWith(values.NewInt(1), values.NewInt(2)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, core.Error(base.ErrAssertion, "expected [array] '[1,2]' not to be equal to [array] '[1,2]'").Error())
			})
		})
	})
}
