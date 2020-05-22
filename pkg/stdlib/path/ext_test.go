package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Ext(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Ext('dir/main.go') should return '.go'", t, func() {
		out, _ := path.Ext(
			context.Background(),
			values.NewString("dir/main.go"),
		)

		So(out, ShouldEqual, ".go")
	})
}
