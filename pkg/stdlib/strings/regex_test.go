package strings_test

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/runtime/values"
	"github.com/MontFerret/ferret/pkg/stdlib/strings"
)

func TestRegexMatch(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexMatch(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexMatch(context.Background(), values.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should match with case insensitive regexp", t, func() {
		out, err := strings.RegexMatch(
			context.Background(),
			values.NewString("My-us3r_n4m3"),
			values.NewString("[a-z0-9_-]{3,16}$"),
			values.True,
		)

		So(err, ShouldBeNil)
		So(out.String(), ShouldEqual, `["My-us3r_n4m3"]`)
	})

	Convey("Should match with case sensitive regexp", t, func() {
		out, err := strings.RegexMatch(
			context.Background(),
			values.NewString("john@doe.com"),
			values.NewString(`([a-z0-9_\.-]+)@([\da-z-]+)\.([a-z\.]{2,6})$`),
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

			_, err = strings.RegexSplit(context.Background(), values.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should split with regexp", t, func() {
		out, err := strings.RegexSplit(
			context.Background(),
			values.NewString("This is a line.\n This is yet another line\r\n This again is a line.\r Mac line "),
			values.NewString(`\.?(\n|\r)`),
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

			_, err = strings.RegexTest(context.Background(), values.NewString(""))

			So(err, ShouldBeError)

		})
	})

	Convey("Should return true when matches", t, func() {
		out, _ := strings.RegexTest(
			context.Background(),
			values.NewString("the quick brown fox"),
			values.NewString("the.*fox"),
		)

		So(out, ShouldEqual, true)
	})
}

func TestRegexReplace(t *testing.T) {
	Convey("When args are not passed", t, func() {
		Convey("It should return an error", func() {
			var err error
			_, err = strings.RegexReplace(context.Background())

			So(err, ShouldBeError)

			_, err = strings.RegexReplace(context.Background(), values.NewString(""))

			So(err, ShouldBeError)

			_, err = strings.RegexReplace(context.Background(), values.NewString(""), values.NewString(""))

			So(err, ShouldBeError)
		})
	})

	Convey("Should replace with regexp", t, func() {
		out, _ := strings.RegexReplace(
			context.Background(),
			values.NewString("the quick brown fox"),
			values.NewString("the.*fox"),
			values.NewString("jumped over"),
		)

		So(out, ShouldEqual, "jumped over")

		out, _ = strings.RegexReplace(
			context.Background(),
			values.NewString("the quick brown fox"),
			values.NewString("o"),
			values.NewString("i"),
		)

		So(out, ShouldEqual, "the quick briwn fix")
	})
}
