package drivers_test

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/wI2L/jettison"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func TestHTTPCookies(t *testing.T) {
	Convey("HTTPCookies", t, func() {
		Convey(".MarshalJSON", func() {
			Convey("Should serialize cookies", func() {
				expires := time.Now()
				headers := drivers.NewHTTPCookiesWith(map[string]drivers.HTTPCookie{
					"Session": {
						Name:     "Session",
						Value:    "asdfg",
						Path:     "/",
						Domain:   "www.google.com",
						Expires:  expires,
						MaxAge:   0,
						Secure:   true,
						HTTPOnly: true,
						SameSite: drivers.SameSiteLaxMode,
					},
				})

				out, err := headers.MarshalJSON()

				t, e := expires.MarshalJSON()
				So(e, ShouldBeNil)

				expected := fmt.Sprintf(`{"Session":{"domain":"www.google.com","expires":%s,"http_only":true,"max_age":0,"name":"Session","path":"/","same_site":"Lax","secure":true,"value":"asdfg"}}`, string(t))

				So(err, ShouldBeNil)
				So(string(out), ShouldEqual, expected)
			})

			Convey("Should set proper values", func() {
				headers := drivers.NewHTTPCookies()

				headers.Set(drivers.HTTPCookie{
					Name:     "Authorization",
					Value:    "e40b7d5eff464a4fb51efed2d1a19a24",
					Path:     "/",
					Domain:   "www.google.com",
					Expires:  time.Now(),
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				})

				_, err := jettison.MarshalOpts(headers, jettison.NoHTMLEscaping())

				So(err, ShouldBeNil)
			})
		})
	})
}
