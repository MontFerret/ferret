package path_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestDir(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.Dir(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Dir(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Dir('pkg/path/dir.go') should return 'pkg/path'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			runtime.NewString("pkg/path/dir.go"),
		)

		So(out.Unwrap(), ShouldEqual, "pkg/path")
	})

	Convey("Dir('/') should return '/'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			runtime.NewString("/"),
		)

		So(out.Unwrap(), ShouldEqual, "/")
	})

	Convey("Dir('') should return '.'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			runtime.NewString(""),
		)

		So(out.Unwrap(), ShouldEqual, ".")
	})

	Convey("Dir('file') should return '.'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			runtime.NewString("file"),
		)

		So(out.Unwrap(), ShouldEqual, ".")
	})

	Convey("Dir('/a/b/c') should return '/a/b'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			runtime.NewString("/a/b/c"),
		)

		So(out.Unwrap(), ShouldEqual, "/a/b")
	})
}
