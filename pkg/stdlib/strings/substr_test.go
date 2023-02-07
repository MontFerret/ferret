package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestSubstring(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Substring(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Substring(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Substring('foobar', 3) should return 'bar'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			values.NewString("foobar"),
			values.NewInt(3),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, "bar")
	})

	Convey("Substring('foobar', 3, 2) should return 'ba'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			values.NewString("foobar"),
			values.NewInt(3),
			values.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, "ba")
	})

	Convey("Substring('foobar', 3, 5) should return 'bar'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			values.NewString("foobar"),
			values.NewInt(3),
			values.NewInt(5),
		)

		So(err, ShouldBeNil)
		So(out, ShouldEqual, "bar")
	})
}

func TestLeft(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Left(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Left(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Left('foobarfoobar', 3) should return 'foo'", t, func() {
		out, _ := strings.Left(
			context.Background(),
			values.NewString("foobarfoobar"),
			values.NewInt(3),
		)

		So(out, ShouldEqual, "foo")
	})

	Convey("Left('foobar', 10) should return 'foobar'", t, func() {
		out, _ := strings.Left(
			context.Background(),
			values.NewString("foobar"),
			values.NewInt(10),
		)

		So(out, ShouldEqual, "foobar")
	})
}

func TestRight(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Right(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Right(context.Background(), values.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Right('foobarfoobar', 3) should return 'bar'", t, func() {
		out, _ := strings.Right(
			context.Background(),
			values.NewString("foobarfoobar"),
			values.NewInt(3),
		)

		So(out, ShouldEqual, "bar")
	})

	Convey("Right('foobar', 10) should return 'foobar'", t, func() {
		out, _ := strings.Right(
			context.Background(),
			values.NewString("foobar"),
			values.NewInt(10),
		)

		So(out, ShouldEqual, "foobar")
	})
}
