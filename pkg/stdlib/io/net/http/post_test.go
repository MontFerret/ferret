package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	h "net/http"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/jarcoal/httpmock"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/io/net/http"
)

func TestPOST(t *testing.T) {
	url := "https://api.montferret.io/users"

	type User struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	Convey("Should successfully make request", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url,
			func(req *h.Request) (*h.Response, error) {
				data, err := io.ReadAll(req.Body)

				if err != nil {
					return nil, err
				}

				user := User{}

				err = json.Unmarshal(data, &user)

				if err != nil {
					return nil, err
				}

				if user.FirstName != "Rob" {
					return nil, fmt.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)
				}

				if user.LastName != "Pike" {
					return nil, fmt.Errorf("Expected LastName to be Pike, but got %s", user.LastName)
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		b, err := json.Marshal(User{
			FirstName: "Rob",
			LastName:  "Pike",
		})

		So(err, ShouldBeNil)

		out, err := http.POST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("body", runtime.NewBinary(b)),
		))

		So(err, ShouldBeNil)
		//So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})

	Convey("Should successfully make request with auto-marshalling to JSON", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		httpmock.RegisterResponder("POST", url,
			func(req *h.Request) (*h.Response, error) {
				data, err := io.ReadAll(req.Body)

				if err != nil {
					return nil, err
				}

				user := User{}

				err = json.Unmarshal(data, &user)

				if err != nil {
					return nil, err
				}

				if user.FirstName != "Rob" {
					return nil, fmt.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)
				}

				if user.LastName != "Pike" {
					return nil, fmt.Errorf("Expected LastName to be Pike, but got %s", user.LastName)
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		j := runtime.NewObjectWith(
			runtime.NewObjectProperty("first_name", runtime.NewString("Rob")),
			runtime.NewObjectProperty("last_name", runtime.NewString("Pike")),
		)

		out, err := http.POST(ctx, runtime.NewObjectWith(
			runtime.NewObjectProperty("url", runtime.NewString(url)),
			runtime.NewObjectProperty("body", j),
		))

		So(err, ShouldBeNil)
		//So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
