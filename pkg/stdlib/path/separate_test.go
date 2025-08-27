package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestSeparate(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := path.Separate(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Separate(context.Background(), runtime.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Separate('http://site.com/logo.png') should return ['http://site.com/', 'logo.png']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			runtime.NewString("http://site.com/logo.png"),
		)

		expected := runtime.NewArrayWith(runtime.NewString("http://site.com/"), runtime.NewString("logo.png"))
		So(out, ShouldResemble, expected)
	})

	Convey("Separate('/') should return ['/', '']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			runtime.NewString("/"),
		)

		expected := runtime.NewArrayWith(runtime.NewString("/"), runtime.NewString(""))
		So(out, ShouldResemble, expected)
	})

	Convey("Separate('file') should return ['', 'file']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			runtime.NewString("file"),
		)

		expected := runtime.NewArrayWith(runtime.NewString(""), runtime.NewString("file"))
		So(out, ShouldResemble, expected)
	})

	Convey("Separate('') should return ['', '']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			runtime.NewString(""),
		)

		expected := runtime.NewArrayWith(runtime.NewString(""), runtime.NewString(""))
		So(out, ShouldResemble, expected)
	})

	Convey("Separate('/a/b/c') should return ['/a/b/', 'c']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			runtime.NewString("/a/b/c"),
		)

		expected := runtime.NewArrayWith(runtime.NewString("/a/b/"), runtime.NewString("c"))
		So(out, ShouldResemble, expected)
	})
}
