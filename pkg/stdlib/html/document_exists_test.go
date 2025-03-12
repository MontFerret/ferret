package html_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/internal"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/html"
)

func TestDocumentExists(t *testing.T) {
	Convey("DOCUMENT_EXISTS", t, func() {
		Convey("Should return error if url is not a string", func() {
			_, err := html.DocumentExists(context.Background(), core.None)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return error if options is not an object", func() {
			_, err := html.DocumentExists(context.Background(), core.NewString("http://fsdfsdfdsdsf.fdf"), core.None)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return error if headers is not an object", func() {
			opts := internal.NewObjectWith(internal.NewObjectProperty("headers", core.None))
			_, err := html.DocumentExists(context.Background(), core.NewString("http://fsdfsdfdsdsf.fdf"), opts)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return 'false' when a website does not exist by a given url", func() {
			out, err := html.DocumentExists(context.Background(), core.NewString("http://fsdfsdfdsdsf.fdf"))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.False)
		})

		Convey("Should return 'true' when a website exists by a given url", func() {
			out, err := html.DocumentExists(context.Background(), core.NewString("https://www.google.com/"))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.True)
		})
	})
}
