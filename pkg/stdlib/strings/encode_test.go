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

		So(out, ShouldEqual, runtime.NewString("3858f62230ac3c915f300c664312c63f"))
	})
}

func TestSha1(t *testing.T) {
	Convey("Should return hash sum of a string", t, func() {
		str := runtime.NewString("foobar")
		out, _ := strings.Sha1(
			context.Background(),
			str,
		)

		So(out, ShouldEqual, runtime.NewString("8843d7f92416211de9ebb963ff4ce28125932878"))
	})
}

func TestSha512(t *testing.T) {
	Convey("Should return hash sum of a string", t, func() {
		str := runtime.NewString("foobar")
		out, _ := strings.Sha512(
			context.Background(),
			str,
		)

		So(out, ShouldEqual, runtime.NewString("0a50261ebd1a390fed2bf326f2673c145582a6342d523204973d0219337f81616a8069b012587cf5635f6925f1b56c360230c19b273500ee013e030601bf2425"))
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
