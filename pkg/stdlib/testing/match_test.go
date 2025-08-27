package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestMatch(t *t.T) {
	Match := base.NewPositiveAssertion(testing.Match)

	Convey("When args are not passed", t, func() {
		Convey("When no args", func() {
			Convey("It should return an error", func() {
				_, err := Match(context.Background())

				So(err, ShouldBeError)
			})
		})

		Convey("When one arg", func() {
			Convey("It should return an error", func() {
				_, err := Match(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeError)
			})
		})
	})

	Convey("When args are valid", t, func() {
		Convey("When string matches pattern", func() {
			Convey("It should not return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewString("hello world"), 
					runtime.NewString("^hello"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When string matches complex pattern", func() {
			Convey("It should not return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewString("abc123def"), 
					runtime.NewString("\\d+"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When string does not match pattern", func() {
			Convey("It should return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewString("hello world"), 
					runtime.NewString("^goodbye"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'hello world' to match regular expression")
			})
		})

		Convey("When value is not string but matches", func() {
			Convey("It should not return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewInt(123), 
					runtime.NewString("\\d+"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When value is not string and does not match", func() {
			Convey("It should return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewInt(123), 
					runtime.NewString("^abc"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '123' to match regular expression")
			})
		})

		Convey("When pattern is invalid", func() {
			Convey("It should return an error", func() {
				_, err := Match(context.Background(), 
					runtime.NewString("test"), 
					runtime.NewString("["))

				So(err, ShouldBeError)
			})
		})
	})
}

func TestNotMatch(t *t.T) {
	NotMatch := base.NewNegativeAssertion(testing.Match)

	Convey("When args are not passed", t, func() {
		Convey("When no args", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background())

				So(err, ShouldBeError)
			})
		})

		Convey("When one arg", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeError)
			})
		})
	})

	Convey("When args are valid", t, func() {
		Convey("When string does not match pattern", func() {
			Convey("It should not return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewString("hello world"), 
					runtime.NewString("^goodbye"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When string matches pattern", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewString("hello world"), 
					runtime.NewString("^hello"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'hello world' not to match regular expression")
			})
		})

		Convey("When string matches complex pattern", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewString("abc123def"), 
					runtime.NewString("\\d+"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'abc123def' not to match regular expression")
			})
		})

		Convey("When value is not string but matches", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewInt(123), 
					runtime.NewString("\\d+"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '123' not to match regular expression")
			})
		})

		Convey("When value is not string and pattern does not match", func() {
			Convey("It should not return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewInt(123), 
					runtime.NewString("^abc"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When pattern is invalid", func() {
			Convey("It should return an error", func() {
				_, err := NotMatch(context.Background(), 
					runtime.NewString("test"), 
					runtime.NewString("["))

				So(err, ShouldBeError)
			})
		})
	})
}