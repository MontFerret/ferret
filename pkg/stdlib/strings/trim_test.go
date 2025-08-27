package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestLTrim(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.LTrim(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("LTrim('  foo bar  ') should return 'foo bar  '", t, func() {
		out, _ := strings.LTrim(
			context.Background(),
			runtime.NewString("  foo bar  "),
		)

		So(out.String(), ShouldEqual, "foo bar  ")
	})

	Convey("LTrim('--==[foo-bar]==--', '-=[]') should return 'foo-bar]==--'", t, func() {
		out, _ := strings.LTrim(
			context.Background(),
			runtime.NewString("--==[foo-bar]==--"),
			runtime.NewString("-=[]"),
		)

		So(out.String(), ShouldEqual, "foo-bar]==--")
	})
}

func TestRTrim(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RTrim(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("RTrim('  foo bar  ') should return '  foo bar'", t, func() {
		out, _ := strings.RTrim(
			context.Background(),
			runtime.NewString("  foo bar  "),
		)

		So(out.String(), ShouldEqual, "  foo bar")
	})

	Convey("LTrim('--==[foo-bar]==--', '-=[]') should return '--==[foo-bar'", t, func() {
		out, _ := strings.RTrim(
			context.Background(),
			runtime.NewString("--==[foo-bar]==--"),
			runtime.NewString("-=[]"),
		)

		So(out.String(), ShouldEqual, "--==[foo-bar")
	})
}

func TestTrim(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.Trim(context.Background())

			So(err, ShouldBeError)

		})
	})

	Convey("Trim('  foo bar  ') should return 'foo bar'", t, func() {
		out, _ := strings.Trim(
			context.Background(),
			runtime.NewString("  foo bar  "),
		)

		So(out.String(), ShouldEqual, "foo bar")
	})

	Convey("Trim('--==[foo-bar]==--', '-=[]') should return 'foo-bar'", t, func() {
		out, _ := strings.Trim(
			context.Background(),
			runtime.NewString("--==[foo-bar]==--"),
			runtime.NewString("-=[]"),
		)

		So(out.String(), ShouldEqual, "foo-bar")
	})
}
