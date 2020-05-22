package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Match(context.Background(), values.NewInt(0), values.NewString("/"))

		So(err, ShouldBeError)
	})

	Convey("Second argument is wrong", t, func() {
		var err error
		_, err = path.Match(context.Background(), values.NewString("/"), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Match('http://site.com/*.csv', 'http://site.com/goods.csv') should return true", t, func() {
		out, _ := path.Match(
			context.Background(),
			values.NewString("http://site.com/*.csv"), values.NewString("http://site.com/goods.csv"),
		)

		So(out, ShouldEqual, values.True)
	})

	Convey("Match('ferret*/ferret', 'ferret/bin/ferret') should return false", t, func() {
		out, _ := path.Match(
			context.Background(),
			values.NewString("ferret*/ferret"), values.NewString("ferret/bin/ferret"),
		)

		So(out, ShouldEqual, values.False)
	})

	Convey("Match('[x-]', 'x') should return ad error", t, func() {
		_, err := path.Match(
			context.Background(),
			values.NewString("[x-]"), values.NewString("x"),
		)

		So(err, ShouldBeError)
	})
}
