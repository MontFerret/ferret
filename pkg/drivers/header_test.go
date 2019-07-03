package drivers_test

import (
	"github.com/MontFerret/ferret/pkg/drivers"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHTTPHeader(t *testing.T) {
	Convey("HTTPHeader", t, func() {
		Convey(".MarshalJSON", func() {
			Convey("Should serialize header values", func() {
				headers := make(drivers.HTTPHeader)

				headers["content-encoding"] = []string{"gzip"}
				headers["content-type"] = []string{"text/html", "charset=utf-8"}

				out, err := headers.MarshalJSON()

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, `{"content-encoding":["gzip"],"content-type":["text/html","charset=utf-8"]}`)
			})
		})
	})
}
