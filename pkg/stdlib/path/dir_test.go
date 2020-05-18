package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDir(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.Dir(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Dir(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Dir('pkg/path/dir.go') should return 'pkg/path'", t, func() {
		out, _ := path.Dir(
			context.Background(),
			values.NewString("pkg/path/dir.go"),
		)

		So(out, ShouldEqual, "pkg/path")
	})
}
