package http_test

import (
	"context"
	h "net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"

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

		out, err := http.GET(ctx, values.NewString(url))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})

	Convey("Should add headers to a request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", url,
			func(req *h.Request) (*h.Response, error) {
				if req.Header.Get("X-Token") != "Ferret" {
					return nil, errors.Errorf("Expected X-Token to be Ferret, but got %s", req.Header.Get("X-Token"))
				}

				if req.Header.Get("X-From") != "localhost" {
					return nil, errors.Errorf("Expected X-From to be localhost, but got %s", req.Header.Get("X-From"))
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		out, err := http.GET(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString(url)),
			values.NewObjectProperty("headers", values.NewObjectWith(
				values.NewObjectProperty("X-Token", values.NewString("Ferret")),
				values.NewObjectProperty("X-From", values.NewString("localhost")),
			)),
		))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
