package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnescapeHTML(t *testing.T) {
	Convey("UnescapeHTML", t, func() {
		Convey("Should unescape an string", func() {
			out, err := strings.UnescapeHTML(context.Background(), core.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))

			expected := core.NewString("<body><span>Foobar</span></body>")
			So(err, ShouldBeNil)
			So(out, ShouldEqual, expected)
		})
	})
}
