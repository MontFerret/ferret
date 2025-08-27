package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestSubstring(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Substring(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Substring(context.Background(), runtime.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Substring('foobar', 3) should return 'bar'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			runtime.NewString("foobar"),
			runtime.NewInt(3),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "bar")
	})

	Convey("Substring('foobar', 3, 2) should return 'ba'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			runtime.NewString("foobar"),
			runtime.NewInt(3),
			runtime.NewInt(2),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "ba")
	})

	Convey("Substring('foobar', 3, 5) should return 'bar'", t, func() {
		out, err := strings.Substring(
			context.Background(),
			runtime.NewString("foobar"),
			runtime.NewInt(3),
			runtime.NewInt(5),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "bar")
	})

	Convey("Edge cases", t, func() {
		Convey("Substring with negative offset", func() {
			out, err := strings.Substring(
				context.Background(),
				runtime.NewString("foobar"),
				runtime.NewInt(-1),
			)

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "")
		})

		Convey("Substring with offset beyond string length", func() {
			out, err := strings.Substring(
				context.Background(),
				runtime.NewString("foo"),
				runtime.NewInt(10),
			)

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "")
		})

		Convey("Substring with zero length", func() {
			out, err := strings.Substring(
				context.Background(),
				runtime.NewString("foobar"),
				runtime.NewInt(2),
				runtime.NewInt(0),
			)

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "")
		})
	})
}

func TestLeft(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Left(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Left(context.Background(), runtime.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Left('foobarfoobar', 3) should return 'foo'", t, func() {
		out, _ := strings.Left(
			context.Background(),
			runtime.NewString("foobarfoobar"),
			runtime.NewInt(3),
		)

		So(out.String(), ShouldEqual, "foo")
	})

	Convey("Left('foobar', 10) should return 'foobar'", t, func() {
		out, _ := strings.Left(
			context.Background(),
			runtime.NewString("foobar"),
			runtime.NewInt(10),
		)

		So(out.String(), ShouldEqual, "foobar")
	})
}

func TestRight(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Right(context.Background())

			So(err, ShouldBeError)

			_, err = strings.Right(context.Background(), runtime.NewString("foo"))

			So(err, ShouldBeError)
		})
	})

	Convey("Right('foobarfoobar', 3) should return 'bar'", t, func() {
		out, _ := strings.Right(
			context.Background(),
			runtime.NewString("foobarfoobar"),
			runtime.NewInt(3),
		)

		So(out.String(), ShouldEqual, "bar")
	})

	Convey("Right('foobar', 10) should return 'foobar'", t, func() {
		out, _ := strings.Right(
			context.Background(),
			runtime.NewString("foobar"),
			runtime.NewInt(10),
		)

		So(out.String(), ShouldEqual, "foobar")
	})
}
