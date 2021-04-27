package drivers_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func TestHTTPHeaders(t *testing.T) {
	Convey("HTTPHeaders", t, func() {
		Convey(".MarshalJSON", func() {
			Convey("Should serialize header values", func() {
				headers := drivers.NewHTTPHeadersWith(map[string][]string{
					"Content-Encoding": []string{"gzip"},
					"Content-Type":     []string{"text/html", "charset=utf-8"},
				})

				out, err := headers.MarshalJSON()

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, `{"Content-Encoding":"gzip","Content-Type":"text/html, charset=utf-8"}`)
			})

			Convey("Should set proper values", func() {
				headers := drivers.NewHTTPHeaders()

				headers.Set("Authorization", `["Basic e40b7d5eff464a4fb51efed2d1a19a24"]`)

				_, err := jettison.MarshalOpts(headers, jettison.NoHTMLEscaping())

				So(err, ShouldBeNil)
			})
		})
	})
}
