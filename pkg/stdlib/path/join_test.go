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

			So(out, ShouldEqual, "")
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

		So(out, ShouldEqual, "pkg/path")
	})

	Convey("Join('pkg', 'path') should return 'pkg/path'", t, func() {
		out, _ := path.Join(
			context.Background(),
			runtime.NewString("pkg"), runtime.NewString("path"),
		)

		So(out, ShouldEqual, "pkg/path")
	})
}
