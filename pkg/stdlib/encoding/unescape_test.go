package encoding_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/stdlib/encoding"
)

func TestHTMLUnescape(t *testing.T) {
	Convey("HTMLUnescape", t, func() {
		Convey("Should unescape an string", func() {
			out, err := encoding.HTMLUnescape(context.Background(), runtime.NewString("&lt;body&gt;&lt;span&gt;Foobar&lt;/span&gt;&lt;/body&gt;"))

			expected := runtime.NewString("<body><span>Foobar</span></body>")
			So(err, ShouldBeNil)
			So(out, ShouldEqual, expected)
		})
	})
}
