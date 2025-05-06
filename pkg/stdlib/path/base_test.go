package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

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
		_, err = path.Base(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Base('') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			core.NewString(""),
		)

		So(out, ShouldEqual, ".")
	})

	Convey("Base('.') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			core.NewString("."),
		)

		So(out, ShouldEqual, ".")
	})

	Convey("Base('pkg/path/base.go') should return 'base.go'", t, func() {
		out, _ := path.Base(
			context.Background(),
			core.NewString("pkg/path/base.go"),
		)

		So(out, ShouldEqual, "base.go")
	})
}
