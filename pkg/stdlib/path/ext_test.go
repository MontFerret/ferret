package path_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestExt(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.Ext(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Ext(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Ext('dir/main.go') should return '.go'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("dir/main.go"),
		)

		So(out.Unwrap(), ShouldEqual, ".go")
	})

	Convey("Ext('') should return ''", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString(""),
		)

		So(out.Unwrap(), ShouldEqual, "")
	})

	Convey("Ext('file') should return ''", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("file"),
		)

		So(out.Unwrap(), ShouldEqual, "")
	})

	Convey("Ext('.hidden') should return '.hidden'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString(".hidden"),
		)

		So(out.Unwrap(), ShouldEqual, ".hidden")
	})

	Convey("Ext('archive.tar.gz') should return '.gz'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("archive.tar.gz"),
		)

		So(out.Unwrap(), ShouldEqual, ".gz")
	})

	Convey("Ext('/path/to/file.txt') should return '.txt'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("/path/to/file.txt"),
		)

		So(out.Unwrap(), ShouldEqual, ".txt")
	})
}
