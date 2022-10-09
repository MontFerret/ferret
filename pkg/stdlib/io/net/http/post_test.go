package http_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	h "net/http"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
	"github.com/MontFerret/ferret/pkg/stdlib/io/net/http"
)

func TestPOST(t *testing.T) {
	Convey("Should successfully make request", t, func() {
		type User struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}

		port := randPort()

		server := &h.Server{
			Addr: port,
			Handler: h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
				var err error

				defer func() {
					if err != nil {
						w.Write([]byte(err.Error()))
					} else {
						w.Write([]byte("OK"))
					}
				}()

				if r.Method != "POST" {
					err = errors.Errorf("Expected method to be POST, but got %s", r.Method)

					return
				}

				data, err := ioutil.ReadAll(r.Body)

				if err != nil {
					return
				}

				user := User{}

				err = json.Unmarshal(data, &user)

				if err != nil {
					return
				}

				if user.FirstName != "Rob" {
					err = errors.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)

					return
				}

				if user.LastName != "Pike" {
					err = errors.Errorf("Expected LastName to be Pike, but got %s", user.LastName)

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

		b, err := json.Marshal(User{
			FirstName: "Rob",
			LastName:  "Pike",
		})

		So(err, ShouldBeNil)

		out, err := http.POST(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString("http://127.0.0.1"+port)),
			values.NewObjectProperty("body", values.NewBinary(b)),
		))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})

	Convey("Should successfully make request with auto-marshalling to JSON", t, func() {
		type User struct {
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}

		port := randPort()

		server := &h.Server{
			Addr: port,
			Handler: h.HandlerFunc(func(w h.ResponseWriter, r *h.Request) {
				var err error

				defer func() {
					if err != nil {
						w.Write([]byte(err.Error()))
					} else {
						w.Write([]byte("OK"))
					}
				}()

				if r.Method != "POST" {
					err = errors.Errorf("Expected method to be POST, but got %s", r.Method)

					return
				}

				data, err := ioutil.ReadAll(r.Body)

				if err != nil {
					return
				}

				user := User{}

				err = json.Unmarshal(data, &user)

				if err != nil {
					return
				}

				if user.FirstName != "Rob" {
					err = errors.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)

					return
				}

				if user.LastName != "Pike" {
					err = errors.Errorf("Expected LastName to be Pike, but got %s", user.LastName)

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

		j := values.NewObjectWith(
			values.NewObjectProperty("first_name", values.NewString("Rob")),
			values.NewObjectProperty("last_name", values.NewString("Pike")),
		)

		out, err := http.POST(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString("http://127.0.0.1"+port)),
			values.NewObjectProperty("body", j),
		))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
