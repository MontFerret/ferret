package strings_test

import (
	"context"
	"testing"

	"github.com/MontFerret/ferret/pkg/runtime/core"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestRegexMatch(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexMatch(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexMatch(context.Background(), core.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should match with case insensitive regexp", t, func() {
		out, err := strings.RegexMatch(
			context.Background(),
			core.NewString("My-us3r_n4m3"),
			core.NewString("[a-z0-9_-]{3,16}$"),
			core.True,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["My-us3r_n4m3"]`)
	})

	Convey("Should match with case sensitive regexp", t, func() {
		out, err := strings.RegexMatch(
			context.Background(),
			core.NewString("john@doe.com"),
			core.NewString(`([a-z0-9_\.-]+)@([\da-z-]+)\.([a-z\.]{2,6})$`),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["john@doe.com","john","doe","com"]`)
	})
}

func TestRegexSplit(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexSplit(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexSplit(context.Background(), core.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should split with regexp", t, func() {
		out, err := strings.RegexSplit(
			context.Background(),
			core.NewString("This is a line.\n This is yet another line\r\n This again is a line.\r Mac line "),
			core.NewString(`\.?(\n|\r)`),
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["This is a line"," This is yet another line",""," This again is a line"," Mac line "]`)
	})
}

func TestRegexTest(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexTest(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexTest(context.Background(), core.NewString(""))

			So(err, ShouldBeError)

		})
	})

	Convey("Should return true when matches", t, func() {
		out, _ := strings.RegexTest(
			context.Background(),
			core.NewString("the quick brown fox"),
			core.NewString("the.*fox"),
		)

		So(out, ShouldEqual, core.True)
	})
}

func TestRegexReplace(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexReplace(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexReplace(context.Background(), core.NewString(""))

			So(err, ShouldBeError)

			_, err = strings.RegexReplace(context.Background(), core.NewString(""), core.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should replace with regexp", t, func() {
		out, _ := strings.RegexReplace(
			context.Background(),
			core.NewString("the quick brown fox"),
			core.NewString("the.*fox"),
			core.NewString("jumped over"),
		)

		So(out.String(), ShouldEqual, "jumped over")

		out, _ = strings.RegexReplace(
			context.Background(),
			core.NewString("the quick brown fox"),
			core.NewString("o"),
			core.NewString("i"),
		)

		So(out.String(), ShouldEqual, "the quick briwn fix")
	})
}
