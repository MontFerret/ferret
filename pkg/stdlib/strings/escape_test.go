package strings_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"testing"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEscapeHTML(t *testing.T) {
	Convey("EscapeHTML", t, func() {
		Convey("Should escape an HTML string", func() {
			out, err := strings.EscapeHTML(context.Background(), core.NewString(`<body><span>Foobar</span></body>`))

			So(err, ShouldBeNil)
			So(out, ShouldEqual, core.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))
		})
	})
}
