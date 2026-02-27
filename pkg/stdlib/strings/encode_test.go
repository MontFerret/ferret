package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/runtime"

	"github.com/MontFerret/ferret/v2/pkg/stdlib/strings"
)

func TestEncodedURIComponent(t *testing.T) {
	Convey("When args are strings", t, func() {
		Convey("EncodeURIComponent('https://github.com/MontFerret/ferret') should return encoded uri", func() {
			out, _ := strings.EncodeURIComponent(
				context.Background(),
				runtime.NewString("https://github.com/MontFerret/ferret"),
			)

			So(out.String(), ShouldEqual, "https%3A%2F%2Fgithub.com%2FMontFerret%2Fferret")
		})
	})
}

func TestMd5(t *testing.T) {
	Convey("Should return hash sum of a string", t, func() {
		str := runtime.NewString("foobar")
		out, _ := strings.Md5(
			context.Background(),
			str,
		)

		So(out, ShouldNotEqual, str)
	})
}

func TestSha1(t *testing.T) {
	Convey("Should return hash sum of a string", t, func() {
		str := runtime.NewString("foobar")
		out, _ := strings.Sha1(
			context.Background(),
			str,
		)

		So(out, ShouldNotEqual, str)
	})
}

func TestSha512(t *testing.T) {
	Convey("Should return hash sum of a string", t, func() {
		str := runtime.NewString("foobar")
		out, _ := strings.Sha512(
			context.Background(),
			str,
		)

		So(out, ShouldNotEqual, str)
	})
}

func TestToBase64(t *testing.T) {
	Convey("Should encode a given value", t, func() {
		out, err := strings.ToBase64(
			context.Background(),
			runtime.NewString("foobar"),
		)

		So(err, ShouldBeNil)
		So(out, ShouldNotEqual, "foobar")
	})
}
