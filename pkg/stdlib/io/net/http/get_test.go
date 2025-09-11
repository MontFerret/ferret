package http_test

import (
	"context"
	"fmt"
	h "net/http"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/io/net/http"
)

func TestGET(t *testing.T) {
	url := "https://api.montferret.io/users"

	Convey("Should successfully make request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", url,
			func(req *h.Request) (*h.Response, error) {
				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		out, err := http.GET(ctx, runtime.NewString(url))

		So(err, ShouldBeNil)
		//So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})

	Convey("Should add headers to a request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", url,
			func(req *h.Request) (*h.Response, error) {
				if req.Header.Get("X-Token") != "Ferret" {
					return nil, fmt.Errorf("Expected X-token to be Ferret, but got %s", req.Header.Get("X-Token"))
				}

				if req.Header.Get("X-From") != "localhost" {
					return nil, fmt.Errorf("Expected X-From to be localhost, but got %s", req.Header.Get("X-From"))
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		out, err := http.GET(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("headers", runtime.NewObjectWith(
				runtime.NewObjectProperty("X-token", runtime.NewString("Ferret")),
				runtime.NewObjectProperty("X-From", runtime.NewString("localhost")),
			)),
		))

		So(err, ShouldBeNil)
		//So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
