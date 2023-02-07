package http_test

import (
	"context"
	"encoding/json"
	"io"
	h "net/http"
	"testing"

	"github.com/jarcoal/httpmock"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/runtime/values/types"
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
					return nil, errors.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)
				}

				if user.LastName != "Pike" {
					return nil, errors.Errorf("Expected LastName to be Pike, but got %s", user.LastName)
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		b, err := json.Marshal(User{
			FirstName: "Rob",
			LastName:  "Pike",
		})

		So(err, ShouldBeNil)

		out, err := http.POST(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString(url)),
			values.NewObjectProperty("body", values.NewBinary(b)),
		))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
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
					return nil, errors.Errorf("Expected FirstName to be Rob, but got %s", user.FirstName)
				}

				if user.LastName != "Pike" {
					return nil, errors.Errorf("Expected LastName to be Pike, but got %s", user.LastName)
				}

				return httpmock.NewStringResponse(200, "OK"), nil
			})

		ctx := context.Background()

		j := values.NewObjectWith(
			values.NewObjectProperty("first_name", values.NewString("Rob")),
			values.NewObjectProperty("last_name", values.NewString("Pike")),
		)

		out, err := http.POST(ctx, values.NewObjectWith(
			values.NewObjectProperty("url", values.NewString(url)),
			values.NewObjectProperty("body", j),
		))

		So(err, ShouldBeNil)
		So(out.Type().ID(), ShouldEqual, types.Binary.ID())
		So(out.String(), ShouldEqual, "OK")
	})
}
