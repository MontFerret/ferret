package path_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/path"
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

		So(out.String(), ShouldEqual, ".go")
	})

	Convey("Ext('') should return ''", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString(""),
		)

		So(out.String(), ShouldEqual, "")
	})

	Convey("Ext('file') should return ''", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("file"),
		)

		So(out.String(), ShouldEqual, "")
	})

	Convey("Ext('.hidden') should return '.hidden'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString(".hidden"),
		)

		So(out.String(), ShouldEqual, ".hidden")
	})

	Convey("Ext('archive.tar.gz') should return '.gz'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("archive.tar.gz"),
		)

		So(out.String(), ShouldEqual, ".gz")
	})

	Convey("Ext('/path/to/file.txt') should return '.txt'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			runtime.NewString("/path/to/file.txt"),
		)

		So(out.String(), ShouldEqual, ".txt")
	})
}
