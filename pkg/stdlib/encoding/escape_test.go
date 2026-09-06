package encoding_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/encoding"
)

func TestHTMLEscape(t *testing.T) {
	Convey("HTMLEscape", t, func() {
		Convey("Should escape an HTML string", func() {
			out, err := encoding.HTMLEscape(context.Background(), runtime.NewString(`<body><span>Foobar</span></body>`))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, runtime.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))
		})

		Convey("Should escape special HTML characters", func() {
			out, err := encoding.HTMLEscape(context.Background(), runtime.NewString(`<>&"'`))

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "&lt;&gt;&amp;&#34;&#39;")
		})

		Convey("Should handle empty string", func() {
			out, err := encoding.HTMLEscape(context.Background(), runtime.NewString(""))

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "")
		})
	})
}
