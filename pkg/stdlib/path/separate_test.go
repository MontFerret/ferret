package path_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/path"
	. "github.com/smartystreets/goconvey/convey"
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
		_, err = path.Separate(context.Background(), values.NewInt(0))

		So(err, ShouldBeError)
	})

	Convey("Separate('http://site.com/logo.png') should return ['http://site.com/', 'logo.png']", t, func() {
		out, _ := path.Separate(
			context.Background(),
			values.NewString("http://site.com/logo.png"),
		)

		expected := values.NewArrayWith(values.NewString("http://site.com/"), values.NewString("logo.png"))
		So(out, ShouldResemble, expected)
	})
}
