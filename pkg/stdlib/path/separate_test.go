package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/path"
)

func TestSeparate(t *testing.T) {
	Convey("When arg is not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := path.Separate(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("Wrong argument", t, func() {
		var err error
		_, err = path.Separate(context.Background(), core.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Separate('http://site.com/logo.png') should return ['http://site.com/', 'logo.png']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			core.NewString("http://site.com/logo.png"),
		)

		expected := internal.NewArrayWith(core.NewString("http://site.com/"), core.NewString("logo.png"))
		So(out, ShouldResemble, expected)
	})
}
