package tests

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	h "net/http"
	"net/http/httptest"
	"testing"

	"github.com/MontFerret/ferret/pkg/drivers"
	"github.com/MontFerret/ferret/pkg/drivers/http"
	"github.com/jarcoal/httpmock"

	"github.com/MontFerret/ferret/pkg/compiler"

	"github.com/MontFerret/ferret/pkg/runtime"
	. "github.com/smartystreets/goconvey/convey"
)

func TestOpenWithSimpleHTTPRequest(t *testing.T) {
	Convey("open with simpleRequest params", t, func() {
		Convey("method GET", func() {
			expected := `{"url":"","status_code":200,"status":"200 OK","headers":{"Set-Cookie":"SampleCookie=sample","Some-Header":"value"},"body":"PCFET0NUWVBFIGh0bWw+PGh0bWw+PGhlYWQ+PC9oZWFkPjxib2R5PnNvbWUgZGF0YTwvYm9keT48L2h0bWw+","response_time":0,"cookies":{"SampleCookie":{"domain":"","expires":"0001-01-01T00:00:00Z","http_only":false,"max_age":0,"name":"SampleCookie","path":"","same_site":"","secure":false,"value":"sample"}}}`

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("GET", "http://localhost:1111", func(request *h.Request) (resp *h.Response, err error) {
				rr := &httptest.ResponseRecorder{}

				rr.Body = bytes.NewBufferString(`<!DOCTYPE html><html><head></head><body>some data</body></html>`)
				h.SetCookie(rr, &h.Cookie{Name: "SampleCookie", Value: "sample", HttpOnly: false})

				resp = rr.Result()
				resp.Header.Add("Some-Header", "value")
				resp.Request = request
				resp.StatusCode = h.StatusOK

				return

			})

			c := compiler.New()

			p, err := c.Compile(`
        let doc = document("http://localhost:1111", {simpleRequest: {method: "GET"}})
        return doc
`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			ctx := drivers.WithContext(context.Background(),
				http.NewDriver(),
				drivers.AsDefault())

			out, err := p.Run(ctx)

			So(err, ShouldBeNil)
			fmt.Println()

			So(string(out), ShouldEqual, expected)
		})

		Convey("method POST. Request should exists body", func() {
			expected := `{"url":"","status_code":200,"status":"200 OK","headers":{"Set-Cookie":"SampleCookie=sample","Some-Header":"value"},"body":"PCFET0NUWVBFIGh0bWw+PGh0bWw+PGhlYWQ+PC9oZWFkPjxib2R5PnNvbWUgZGF0YTwvYm9keT48L2h0bWw+","response_time":0,"cookies":{"SampleCookie":{"domain":"","expires":"0001-01-01T00:00:00Z","http_only":false,"max_age":0,"name":"SampleCookie","path":"","same_site":"","secure":false,"value":"sample"}}}`

			httpmock.Activate()
			defer httpmock.DeactivateAndReset()

			httpmock.RegisterResponder("POST", "http://localhost:1111", func(request *h.Request) (resp *h.Response, err error) {
				if request.Body == nil {
					return nil, errors.New("request Body doesn't exists")
				}

				rr := &httptest.ResponseRecorder{}

				rr.Body = bytes.NewBufferString(`<!DOCTYPE html><html><head></head><body>some data</body></html>`)
				h.SetCookie(rr, &h.Cookie{Name: "SampleCookie", Value: "sample", HttpOnly: false})

				resp = rr.Result()
				resp.Header.Add("Some-Header", "value")
				resp.Request = request
				resp.StatusCode = h.StatusOK

				return

			})

			c := compiler.New()

			p, err := c.Compile(`
        let req = to_binary("some input param")
        let doc = document("http://localhost:1111", {simpleRequest: {method: "POST", body: req}})
        return doc
`)

			So(err, ShouldBeNil)
			So(p, ShouldHaveSameTypeAs, &runtime.Program{})

			ctx := drivers.WithContext(context.Background(),
				http.NewDriver(),
				drivers.AsDefault())

			out, err := p.Run(ctx)

			So(err, ShouldBeNil)
			fmt.Println()

			So(string(out), ShouldEqual, expected)
		})

	})
}
