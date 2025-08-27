package testing_test

import (
	"context"
	t "testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestEqual(t *t.T) {
	Equal := base.NewPositiveAssertion(testing.Equal)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Equal(context.Background())

			So(err, ShouldBeError)

			_, err = Equal(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), runtime.NewString("Foo"), runtime.NewString("Bar"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' to be equal to [string] 'Bar'").Error())
			})
		})

		Convey("When 'Foo' and 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), runtime.NewString("Foo"), runtime.NewString("Foo"))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), runtime.NewInt(1), runtime.NewInt(2))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [int] '1' to be equal to [int] '2'").Error())
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), runtime.NewInt(1), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are boolean", t, func() {
		Convey("When False and True", func() {
			Convey("It should return an error", func() {
				_, err := Equal(context.Background(), runtime.False, runtime.True)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [boolean] 'false' to be equal to [boolean] 'true'").Error())
			})
		})

		Convey("When False and False", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(context.Background(), runtime.False, runtime.False)

				So(err, ShouldBeNil)
			})
		})
	})

	Convey("When args are arrays", t, func() {
		Convey("When [1] and [1,2]", func() {
			Convey("It should return an error", func() {
				_, err := Equal(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1)),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[1]' to be equal to [list] '[1,2]'").Error())
			})
		})

		Convey("When [1,2] and [1,2]", func() {
			Convey("It should not return an error", func() {
				_, err := Equal(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
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

			_, err = NotEqual(context.Background(), runtime.NewInt(1))

			So(err, ShouldBeError)
		})
	})

	Convey("When args are string", t, func() {
		Convey("When 'Foo' and 'Bar'", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), runtime.NewString("Foo"), runtime.NewString("Bar"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 'Foo' and 'Foo'", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), runtime.NewString("Foo"), runtime.NewString("Foo"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [string] 'Foo' not to be equal to [string] 'Foo'").Error())
			})
		})
	})

	Convey("When args are numbers", t, func() {
		Convey("When 1 and 2", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), runtime.NewInt(1), runtime.NewInt(2))

				So(err, ShouldBeNil)
			})
		})

		Convey("When 1 and 1", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), runtime.NewInt(1), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [int] '1' not to be equal to [int] '1'").Error())
			})
		})
	})

	Convey("When args are boolean", t, func() {
		Convey("When False and True", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(context.Background(), runtime.False, runtime.True)

				So(err, ShouldBeNil)
			})
		})

		Convey("When False and False", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(context.Background(), runtime.False, runtime.False)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [boolean] 'false' not to be equal to [boolean] 'false'").Error())
			})
		})
	})

	Convey("When args are arrays", t, func() {
		Convey("When [1] and [1,2]", func() {
			Convey("It should return an error", func() {
				_, err := NotEqual(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1)),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				)

				So(err, ShouldBeNil)
			})
		})

		Convey("When [1,2] and [1,2]", func() {
			Convey("It should not return an error", func() {
				_, err := NotEqual(
					context.Background(),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
					runtime.NewArrayWith(runtime.NewInt(1), runtime.NewInt(2)),
				)

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, runtime.Error(base.ErrAssertion, "expected [list] '[1,2]' not to be equal to [list] '[1,2]'").Error())
			})
		})
	})
}
