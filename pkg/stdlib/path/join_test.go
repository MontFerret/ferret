package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestJoin(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an empty string without error", func() {
			out, err := path.Join(context.Background())

			So(out.Unwrap(), ShouldEqual, "")
			So(err, ShouldBeNil)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Join(context.Background(), runtime.NewString("/"), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Wrong argument within an array", t, func() {
		var err error
		_, err = path.Join(
			context.Background(),
			runtime.NewArrayWith(runtime.NewString("/"), runtime.NewInt(0)),
		)

		So(err, ShouldBeError)
	})

	Convey("Join(['pkg', 'path']) should return 'pkg/path'", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewArrayWith(runtime.NewString("pkg"), runtime.NewString("path")),
		)

		So(out.Unwrap(), ShouldEqual, "pkg/path")
	})

	Convey("Join('pkg', 'path') should return 'pkg/path'", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewString("pkg"), runtime.NewString("path"),
		)

		So(out.Unwrap(), ShouldEqual, "pkg/path")
	})

	Convey("Join with empty array should return empty string", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewArray(0),
		)

		So(out.Unwrap(), ShouldEqual, "")
	})

	Convey("Join with single element should return that element", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewArrayWith(runtime.NewString("single")),
		)

		So(out.Unwrap(), ShouldEqual, "single")
	})

	Convey("Join('/', 'home', 'user') should return '/home/user'", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewString("/"), runtime.NewString("home"), runtime.NewString("user"),
		)

		So(out.Unwrap(), ShouldEqual, "/home/user")
	})

	Convey("Join with empty strings should handle correctly", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewString(""), runtime.NewString("path"),
		)

		So(out.Unwrap(), ShouldEqual, "path")
	})
}
