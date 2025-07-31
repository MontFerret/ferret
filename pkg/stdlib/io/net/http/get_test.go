package http_test

import (
	"context"
	h "net/http"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/jarcoal/httpmock"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	"github.com/pkg/errors"

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

		out, err := http.GET(ctx, core.NewString(url))

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
					return nil, errors.Errorf("Expected X-token to be Ferret, but got %s", req.Header.Get("X-Token"))
				}

				if req.Header.Get("X-From") != "localhost" {
					return nil, errors.Errorf("Expected X-From to be localhost, but got %s", req.Header.Get("X-From"))
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		out, err := http.GET(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("url", core.NewString(url)),
			runtime.NewObjectProperty("headers", runtime.NewObjectWith(
				runtime.NewObjectProperty("X-token", core.NewString("Ferret")),
				runtime.NewObjectProperty("X-From", core.NewString("localhost")),
			)),
		))

		So(err, ShouldBeNil)
		//So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
