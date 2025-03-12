package path_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

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
		_, err = path.Ext(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Ext('dir/main.go') should return '.go'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			core.NewString("dir/main.go"),
		)

		So(out, ShouldEqual, ".go")
	})
}
