package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEscapeHTML(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			_, err := strings.EscapeHTML(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("EscapeHTML", t, func() {
		Convey("Should escape an HTML string", func() {
			out, err := strings.EscapeHTML(context.Background(), core.NewString(`<body><span>Foobar</span></body>`))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))
		})

		Convey("Should escape special HTML characters", func() {
			out, err := strings.EscapeHTML(context.Background(), core.NewString(`<>&"'`))

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "&lt;&gt;&amp;&#34;&#39;")
		})

		Convey("Should handle empty string", func() {
			out, err := strings.EscapeHTML(context.Background(), core.NewString(""))

			So(err, ShouldBeNil)
			So(out.String(), ShouldEqual, "")
		})
	})
}
