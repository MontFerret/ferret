package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
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
			values.NewString("  foo bar  "),
		)

		So(out, ShouldEqual, "foo bar  ")
	})

	Convey("LTrim('--==[foo-bar]==--', '-=[]') should return 'foo-bar]==--'", t, func() {
		out, _ := strings.LTrim(
			context.Background(),
			values.NewString("--==[foo-bar]==--"),
			values.NewString("-=[]"),
		)

		So(out, ShouldEqual, "foo-bar]==--")
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
			values.NewString("  foo bar  "),
		)

		So(out, ShouldEqual, "  foo bar")
	})

	Convey("LTrim('--==[foo-bar]==--', '-=[]') should return '--==[foo-bar'", t, func() {
		out, _ := strings.RTrim(
			context.Background(),
			values.NewString("--==[foo-bar]==--"),
			values.NewString("-=[]"),
		)

		So(out, ShouldEqual, "--==[foo-bar")
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
			values.NewString("  foo bar  "),
		)

		So(out, ShouldEqual, "foo bar")
	})

	Convey("Trim('--==[foo-bar]==--', '-=[]') should return 'foo-bar'", t, func() {
		out, _ := strings.Trim(
			context.Background(),
			values.NewString("--==[foo-bar]==--"),
			values.NewString("-=[]"),
		)

		So(out, ShouldEqual, "foo-bar")
	})
}
