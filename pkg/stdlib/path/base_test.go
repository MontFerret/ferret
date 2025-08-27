package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestBase(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.Base(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Base(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Base('') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString(""),
		)

		So(out.Unwrap(), ShouldEqual, ".")
	})

	Convey("Base('.') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString("."),
		)

		So(out.Unwrap(), ShouldEqual, ".")
	})

	Convey("Base('pkg/path/base.go') should return 'base.go'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString("pkg/path/base.go"),
		)

		So(out.Unwrap(), ShouldEqual, "base.go")
	})

	Convey("Base('/') should return '/'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString("/"),
		)

		So(out.Unwrap(), ShouldEqual, "/")
	})

	Convey("Base('/usr/bin/') should return 'bin'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString("/usr/bin/"),
		)

		So(out.Unwrap(), ShouldEqual, "bin")
	})

	Convey("Base('a/b/c') should return 'c'", t, func() {
		out, _ := path.Base(
			context.Background(),
			runtime.NewString("a/b/c"),
		)

		So(out.Unwrap(), ShouldEqual, "c")
	})
}
