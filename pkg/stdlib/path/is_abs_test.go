package path_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
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
		_, err = path.IsAbs(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("IsAbs('/ferret/bin/ferret') should return true", t, func() {
		out, _ := path.IsAbs(
			context.Background(),
			core.NewString("/ferret/bin/ferret"),
		)

		So(out, ShouldEqual, core.True)
	})

	Convey("IsAbs('..') should return false", t, func() {
		out, _ := path.IsAbs(
			context.Background(),
			core.NewString(".."),
		)

		So(out, ShouldEqual, core.False)
	})
}
