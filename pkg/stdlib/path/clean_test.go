package path_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/path"
)

func TestClean(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.Clean(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Clean(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Clean('pkg//path//clean.go') should return 'pkg/path/clean.go'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString("pkg//path//clean.go"),
		)

		So(out.String(), ShouldEqual, "pkg/path/clean.go")
	})

	Convey("Clean('/cmd/main/../../..') should return '/'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString("/cmd/main/../../.."),
		)

		So(out.String(), ShouldEqual, "/")
	})

	Convey("Clean('') should return '.'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString(""),
		)

		So(out.String(), ShouldEqual, ".")
	})

	Convey("Clean('.') should return '.'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString("."),
		)

		So(out.String(), ShouldEqual, ".")
	})

	Convey("Clean('./a') should return 'a'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString("./a"),
		)

		So(out.String(), ShouldEqual, "a")
	})

	Convey("Clean('a//b') should return 'a/b'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			runtime.NewString("a//b"),
		)

		So(out.String(), ShouldEqual, "a/b")
	})
}
