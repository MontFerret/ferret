package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Base(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Base('') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			values.NewString(""),
		)

		So(out, ShouldEqual, ".")
	})

	Convey("Base('.') should return '.'", t, func() {
		out, _ := path.Base(
			context.Background(),
			values.NewString("."),
		)

		So(out, ShouldEqual, ".")
	})

	Convey("Base('pkg/path/base.go') should return 'base.go'", t, func() {
		out, _ := path.Base(
			context.Background(),
			values.NewString("pkg/path/base.go"),
		)

		So(out, ShouldEqual, "base.go")
	})
}
