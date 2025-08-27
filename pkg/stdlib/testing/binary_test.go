package testing_test

import (
	"context"
	t "testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/testing"
	"github.com/MontFerret/ferret/pkg/stdlib/testing/base"
)

func TestBinary(t *t.T) {
	Binary := base.NewPositiveAssertion(testing.Binary)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := Binary(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is not binary", t, func() {
		Convey("When arg is string", func() {
			Convey("It should return an error", func() {
				_, err := Binary(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [string] 'test' to be binary")
			})
		})

		Convey("When arg is int", func() {
			Convey("It should return an error", func() {
				_, err := Binary(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [int] '1' to be binary")
			})
		})

		Convey("When arg is array", func() {
			Convey("It should return an error", func() {
				_, err := Binary(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [list] '[]' to be binary")
			})
		})
	})

	Convey("When arg is binary", t, func() {
		Convey("When arg is empty binary", func() {
			Convey("It should not return an error", func() {
				_, err := Binary(context.Background(), runtime.NewBinary([]byte{}))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is non-empty binary", func() {
			Convey("It should not return an error", func() {
				_, err := Binary(context.Background(), runtime.NewBinary([]byte{1, 2, 3}))

				So(err, ShouldBeNil)
			})
		})
	})
}

func TestNotBinary(t *t.T) {
	NotBinary := base.NewNegativeAssertion(testing.Binary)

	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := NotBinary(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When arg is binary", t, func() {
		Convey("When arg is empty binary", func() {
			Convey("It should return an error", func() {
				_, err := NotBinary(context.Background(), runtime.NewBinary([]byte{}))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [binary] '' not to be binary")
			})
		})

		Convey("When arg is non-empty binary", func() {
			Convey("It should return an error", func() {
				_, err := NotBinary(context.Background(), runtime.NewBinary([]byte{1, 2, 3}))

				So(err, ShouldBeError)
				So(err.Error(), ShouldEqual, base.ErrAssertion.Error()+": expected [binary] '\x01\x02\x03' not to be binary")
			})
		})
	})

	Convey("When arg is not binary", t, func() {
		Convey("When arg is string", func() {
			Convey("It should not return an error", func() {
				_, err := NotBinary(context.Background(), runtime.NewString("test"))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is int", func() {
			Convey("It should not return an error", func() {
				_, err := NotBinary(context.Background(), runtime.NewInt(1))

				So(err, ShouldBeNil)
			})
		})

		Convey("When arg is array", func() {
			Convey("It should not return an error", func() {
				_, err := NotBinary(context.Background(), runtime.NewArray(0))

				So(err, ShouldBeNil)
			})
		})
	})
}