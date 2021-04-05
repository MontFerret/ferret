package drivers_test

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/drivers"
)

func TestSetDefaultParams(t *testing.T) {
	Convey("Should take values from Options if not present in Params", t, func() {
		opts := &drivers.Options{
			Name:      "Test",
			UserAgent: "Mozilla",
			Headers: drivers.NewHTTPHeadersWith(map[string][]string{
				"Accept": {"application/json"},
			}),
			Cookies: drivers.NewHTTPCookiesWith(map[string]drivers.HTTPCookie{
				"Session": drivers.HTTPCookie{
					Name:     "Session",
					Value:    "fsfsdfsd",
					Path:     "",
					Domain:   "",
					Expires:  time.Time{},
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
			}),
		}

		params := drivers.SetDefaultParams(opts, drivers.Params{})

		So(params.UserAgent, ShouldEqual, opts.UserAgent)
		So(params.Headers, ShouldNotBeNil)
		So(params.Cookies, ShouldNotBeNil)
	})
}
