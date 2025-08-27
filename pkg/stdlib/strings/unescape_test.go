package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnescapeHTML(t *testing.T) {
	Convey("UnescapeHTML", t, func() {
		Convey("Should unescape an string", func() {
			out, err := strings.UnescapeHTML(context.Background(), runtime.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))

			expected := runtime.NewString("<body><span>Foobar</span></body>")
			So(err, ShouldBeNil)
			So(out, ShouldEqual, expected)
		})
	})
}
