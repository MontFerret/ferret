package http_test

import (
	stdhttp "net/http"
	"testing"
	"time"

	"github.com/sethgrid/pester"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/http"
)

func TestNewOptions(t *testing.T) {
	Convey("Should create driver options with initial values", t, func() {
		opts := http.NewOptions([]http.Option{})
		So(opts.Options, ShouldNotBeNil)
		So(opts.Name, ShouldEqual, http.DriverName)
		So(opts.Backoff, ShouldEqual, pester.ExponentialBackoff)
		So(opts.Concurrency, ShouldEqual, http.DefaultConcurrency)
		So(opts.MaxRetries, ShouldEqual, http.DefaultMaxRetries)
		So(opts.HTTPCodesFilter, ShouldHaveLength, 0)
	})

	Convey("Should use setters to set values", t, func() {
		expectedName := http.DriverName + "2"
		expectedUA := "Mozilla"
		expectedProxy := "https://proxy.com"
		expectedMaxRetries := 2
		expectedConcurrency := 10
		expectedTransport := &stdhttp.Transport{}
		expectedTimeout := time.Second * 5

		opts := http.NewOptions([]http.Option{
			http.WithCustomName(expectedName),
			http.WithUserAgent(expectedUA),
			http.WithProxy(expectedProxy),
			http.WithCookie(drivers.HTTPCookie{
				Name:     "Session",
				Value:    "fsdfsdfs",
				Path:     "dfsdfsd",
				Domain:   "sfdsfs",
				Expires:  time.Time{},
				MaxAge:   0,
				Secure:   false,
				HTTPOnly: false,
				SameSite: 0,
			}),
			http.WithCookies([]drivers.HTTPCookie{
				{
					Name:     "Use",
					Value:    "Foos",
					Path:     "",
					Domain:   "",
					Expires:  time.Time{},
					MaxAge:   0,
					Secure:   false,
					HTTPOnly: false,
					SameSite: 0,
				},
			}),
			http.WithHeader("Authorization", []string{"Bearer dfsd7f98sd9fsd9fsd"}),
			http.WithHeaders(drivers.NewHTTPHeadersWith(map[string][]string{
				"x-correlation-id": {"232483833833839"},
			})),
			http.WithDefaultBackoff(),
			http.WithMaxRetries(expectedMaxRetries),
			http.WithConcurrency(expectedConcurrency),
			http.WithAllowedHTTPCode(401),
			http.WithAllowedHTTPCodes([]int{403, 404}),
			http.WithCustomTransport(expectedTransport),
			http.WithTimeout(time.Second * 5),
		})
		So(opts.Options, ShouldNotBeNil)
		So(opts.Name, ShouldEqual, expectedName)
		So(opts.UserAgent, ShouldEqual, expectedUA)
		So(opts.Proxy, ShouldEqual, expectedProxy)
		So(opts.Cookies.Length(), ShouldEqual, 2)
		So(opts.Headers.Length(), ShouldEqual, 2)
		So(opts.Backoff, ShouldEqual, pester.DefaultBackoff)
		So(opts.MaxRetries, ShouldEqual, expectedMaxRetries)
		So(opts.Concurrency, ShouldEqual, expectedConcurrency)
		So(opts.HTTPCodesFilter, ShouldHaveLength, 3)
		So(opts.HTTPTransport, ShouldEqual, expectedTransport)
		So(opts.Timeout, ShouldEqual, expectedTimeout)
	})
}
