package path_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestMatch(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := path.Match(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("First argument is wrong", t, func() {
		var err error
		_, err = path.Match(context.Background(), core.NewInt(0), core.NewString("/"))

		So(err, ShouldBeError)
	})

	Convey("Second argument is wrong", t, func() {
		var err error
		_, err = path.Match(context.Background(), core.NewString("/"), core.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Match('http://site.com/*.csv', 'http://site.com/goods.csv') should return true", t, func() {
		out, _ := path.Match(
			context.Background(),
			core.NewString("http://site.com/*.csv"), core.NewString("http://site.com/goods.csv"),
		)

		So(out, ShouldEqual, core.True)
	})

	Convey("Match('ferret*/ferret', 'ferret/bin/ferret') should return false", t, func() {
		out, _ := path.Match(
			context.Background(),
			core.NewString("ferret*/ferret"), core.NewString("ferret/bin/ferret"),
		)

		So(out, ShouldEqual, core.False)
	})

	Convey("Match('[x-]', 'x') should return ad error", t, func() {
		_, err := path.Match(
			context.Background(),
			core.NewString("[x-]"), core.NewString("x"),
		)

		So(err, ShouldBeError)
	})
}
