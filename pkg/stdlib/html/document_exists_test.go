package html_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/html"
)

func TestDocumentExists(t *testing.T) {
	Convey("DOCUMENT_EXISTS", t, func() {
		Convey("Should return error if url is not a string", func() {
			_, err := html.DocumentExists(context.Background(), values.None)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return error if options is not an object", func() {
			_, err := html.DocumentExists(context.Background(), values.NewString("http://fsdfsdfdsdsf.fdf"), values.None)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return error if headers is not an object", func() {
			opts := values.NewObjectWith(values.NewObjectProperty("headers", values.None))
			_, err := html.DocumentExists(context.Background(), values.NewString("http://fsdfsdfdsdsf.fdf"), opts)

			So(err, ShouldNotBeNil)
		})

		Convey("Should return 'false' when a website does not exist by a given url", func() {
			out, err := html.DocumentExists(context.Background(), values.NewString("http://fsdfsdfdsdsf.fdf"))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.False)
		})

		Convey("Should return 'true' when a website exists by a given url", func() {
			out, err := html.DocumentExists(context.Background(), values.NewString("https://www.google.com/"))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, values.True)
		})
	})
}
