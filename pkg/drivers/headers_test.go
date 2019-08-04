package drivers_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func TestHTTPHeader(t *testing.T) {
	Convey("HTTPHeaders", t, func() {
		Convey(".MarshalJSON", func() {
			Convey("Should serialize header values", func() {
				headers := make(drivers.HTTPHeaders)

				headers["Content-Encoding"] = []string{"gzip"}
				headers["Content-Type"] = []string{"text/html", "charset=utf-8"}

				out, err := headers.MarshalJSON()

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, `{"Content-Encoding":"gzip","Content-Type":"text/html, charset=utf-8"}`)
			})
		})
	})
}
