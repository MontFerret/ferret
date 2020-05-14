package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/values"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFromBase64(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FromBase64(context.Background())

			So(err, ShouldBeError)
		})
	})

	Convey("When hash is not valid base64", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.FromBase64(
				context.Background(),
				values.NewString("foobar"),
			)

			So(err, ShouldBeError)
		})
	})

	Convey("Should decode a given hash", t, func() {
		out, err := strings.FromBase64(
			context.Background(),
			values.NewString("Zm9vYmFy"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, "Zm9vYmFy")
		So(out, ShouldEqual, "foobar")
	})
}

func TestDecodeURIComponent(t *testing.T) {

	Convey("Decode", t, func() {
		testCases := []struct {
			Name   string
			InURI  string
			OutURI string
		}{
			{
				Name:   "Unicode",
				InURI:  "https://thedomain/alphabet=M\u0026borough=Bronx\u0026a=b",
				OutURI: "https://thedomain/alphabet=M&borough=Bronx&a=b",
			},
			{
				Name:   "Percent-encoding",
				InURI:  "https://ru.wikipedia.org/wiki/%D0%AF%D0%B7%D1%8B%D0%BA_%D0%BF%D1%80%D0%BE%D0%B3%D1%80%D0%B0%D0%BC%D0%BC%D0%B8%D1%80%D0%BE%D0%B2%D0%B0%D0%BD%D0%B8%D1%8F",
				OutURI: "https://ru.wikipedia.org/wiki/Язык_программирования",
			},
		}

		for _, tC := range testCases {
			Convey(tC.Name, func() {
				out, err := strings.DecodeURIComponent(
					context.Background(),
					values.NewString(tC.InURI),
				)
				So(err, ShouldBeNil)

				So(out.String(), ShouldEqual, tC.OutURI)
			})
		}
	})
}
