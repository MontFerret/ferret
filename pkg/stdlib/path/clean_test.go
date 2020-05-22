package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Clean(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Clean('pkg//path//clean.go') should return 'pkg/path/clean.go'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			values.NewString("pkg//path//clean.go"),
		)

		So(out, ShouldEqual, "pkg/path/clean.go")
	})

	Convey("Clean('/cmd/main/../../..') should return '/'", t, func() {
		out, _ := path.Clean(
			context.Background(),
			values.NewString("/cmd/main/../../.."),
		)

		So(out, ShouldEqual, "/")
	})
}
