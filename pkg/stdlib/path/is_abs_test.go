package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
)

func TestIsAbs(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = path.IsAbs(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.IsAbs(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("IsAbs('/ferret/bin/ferret') should return true", t, func() {
		out, _ := path.IsAbs(
			context.Background(),
			values.NewString("/ferret/bin/ferret"),
		)

		So(out, ShouldEqual, values.True)
	})

	Convey("IsAbs('..') should return false", t, func() {
		out, _ := path.IsAbs(
			context.Background(),
			values.NewString(".."),
		)

		So(out, ShouldEqual, values.False)
	})
}
