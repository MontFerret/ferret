package http_test

import (
	"context"
	h "net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/pkg/errors"

	"github.com/MontFerret/ferret/pkg/runtime"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/io/net/http"
)

func TestRegisterLib(t *testing.T) {
	Convey("Should register HTTP namespace functions", t, func() {
		ns := runtime.NewRootNamespace()
		
		err := http.RegisterLib(ns)
		
		So(err, ShouldBeNil)
		
		// Verify that functions were registered by checking registered function names
		functions := ns.RegisteredFunctions()
		So(len(functions), ShouldBeGreaterThan, 0)
		
		// Check that HTTP functions are registered
		hasGet := false
		hasPost := false
		hasPut := false
		hasDelete := false
		hasDo := false
		
		for _, fn := range functions {
			if fn == "HTTP::GET" {
				hasGet = true
			}
			if fn == "HTTP::POST" {
				hasPost = true
			}
			if fn == "HTTP::PUT" {
				hasPut = true
			}
			if fn == "HTTP::DELETE" {
				hasDelete = true
			}
			if fn == "HTTP::DO" {
				hasDo = true
			}
		}
		
		So(hasGet, ShouldBeTrue)
		So(hasPost, ShouldBeTrue)
		So(hasPut, ShouldBeTrue)
		So(hasDelete, ShouldBeTrue)
		So(hasDo, ShouldBeTrue)
	})
}

func TestREQUEST(t *testing.T) {
	url := "https://api.montferret.io/test"

	Convey("Should successfully make GET request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("GET", url,
			func(req *h.Request) (*h.Response, error) {
				return httpmock.NewStringResponse(200, "GET OK"), nil
			})

		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("method", runtime.NewString("GET")),
			runtime.NewObjectProperty("url", runtime.NewString(url)),
		))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "GET OK")
	})

	Convey("Should successfully make POST request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url,
			func(req *h.Request) (*h.Response, error) {
				return httpmock.NewStringResponse(200, "POST OK"), nil
			})

		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("method", runtime.NewString("POST")),
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("body", runtime.NewBinary([]byte("test data"))),
		))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "POST OK")
	})

	Convey("Should add headers to request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url,
			func(req *h.Request) (*h.Response, error) {
				if req.Header.Get("X-Token") != "test-token" {
					return nil, errors.Errorf("Expected X-Token to be test-token, but got %s", req.Header.Get("X-Token"))
				}
				return httpmock.NewStringResponse(200, "Headers OK"), nil
			})

		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("method", runtime.NewString("POST")),
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("headers", runtime.NewObjectWith(
				runtime.NewObjectProperty("X-Token", runtime.NewString("test-token")),
			)),
		))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "Headers OK")
	})

	Convey("Should handle JSON body", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url,
			func(req *h.Request) (*h.Response, error) {
				if req.Header.Get("Content-Type") != "application/json" {
					return nil, errors.Errorf("Expected Content-Type to be application/json, but got %s", req.Header.Get("Content-Type"))
				}
				return httpmock.NewStringResponse(200, "JSON OK"), nil
			})

		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("method", runtime.NewString("POST")),
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("body", runtime.NewObjectWith(
				runtime.NewObjectProperty("test", runtime.NewString("data")),
			)),
		))

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, "JSON OK")
	})

	Convey("Should return error when url is missing", t, func() {
		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("method", runtime.NewString("GET")),
		))

		So(out, ShouldEqual, runtime.None)
		So(err, ShouldBeError)
		So(err.Error(), ShouldContainSubstring, ".url")
	})

	Convey("Should return error with invalid arguments", t, func() {
		ctx := context.Background()

		out, err := http.REQUEST(ctx)

		So(out, ShouldEqual, runtime.None)
		So(err, ShouldBeError)
	})

	Convey("Should return error with invalid argument type", t, func() {
		ctx := context.Background()

		out, err := http.REQUEST(ctx, runtime.NewString("invalid"))

		So(out, ShouldEqual, runtime.None)
		So(err, ShouldBeError)
	})
}