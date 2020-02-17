package http_test

import (
	"context"
	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/pkg/errors"
	h "net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/io/net/http"
)

func TestGET(t *testing.T) {
	Convey("Should successfully make request", t, func() {
		server := &h.Server{
			Addr: ":9999",
			Handler: h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
				if r.Method == "GET" {
					w.Write([]byte("OK"))
				} else {
					w.Write([]byte("Expected method to be GET"))
				}
			}),
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			server.ListenAndServe()
		}()

		defer func() {
			cancel()
			server.Shutdown(ctx)
		}()

		out, err := http.GET(ctx, values.NewString("http://localhost:9999"))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})

	Convey("Should add headers to a request", t, func() {
		server := &h.Server{
			Addr: ":9999",
			Handler: h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
				var err error

				defer func() {
					if err != nil {
						w.Write([]byte(err.Error()))
					} else {
						w.Write([]byte("OK"))
					}
				}()

				if r.Method != "GET" {
					err = errors.Errorf("Expected method to be GET, but got %s", r.Method)

					return
				}

				token := r.Header.Get("X-Token")

				if token != "Ferret" {
					err = errors.Errorf("Expected X-Token header to equal to Ferret, but got %s", token)

					return
				}

				from := r.Header.Get("X-From")

				if from != "localhost" {
					err = errors.Errorf("Expected X-From header to equal to localhost, but got %s", from)

					return
				}
			}),
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			server.ListenAndServe()
		}()

		defer func() {
			cancel()
			server.Shutdown(ctx)
		}()

		out, err := http.GET(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString("http://localhost:9999")),
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
