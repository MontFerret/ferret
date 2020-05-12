package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Join(context.Background(), values.NewString("/"), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Wrong argument within an array", t, func() {
		var err error
		_, err = path.Join(
			context.Background(),
			values.NewArrayWith(values.NewString("/"), values.NewInt(0)),
		)

		So(err, ShouldBeError)
	})

	Convey("Join(['pkg', 'path']) should return 'pkg/path'", t, func() {
		out, _ := path.Join(
			context.Background(),
			values.NewArrayWith(values.NewString("pkg"), values.NewString("path")),
		)

		So(out, ShouldEqual, "pkg/path")
	})

	Convey("Join('pkg', 'path') should return 'pkg/path'", t, func() {
		out, _ := path.Join(
			context.Background(),
			values.NewString("pkg"), values.NewString("path"),
		)

		So(out, ShouldEqual, "pkg/path")
	})
}
