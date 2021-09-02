package dom

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_toCamelCase(t *testing.T) {
	Convey("toCamelCase", t, func() {
		Convey("should format string into camel case", func() {
			inputs := []struct {
				actual   string
				expected string
			}{
				{
					actual:   "foo-bar",
					expected: "fooBar",
				},
				{
					actual:   "foo-1-bar",
					expected: "foo1Bar",
				},
				{
					actual:   "overscroll-behavior-x",
					expected: "overscrollBehaviorX",
				},
				{
					actual:   "x",
					expected: "x",
				},
				{
					actual:   "foo-x",
					expected: "fooX",
				},
				{
					actual:   "foo-$",
					expected: "foo",
				},
				{
					actual:   "color",
					expected: "color",
				},
				{
					actual:   "textDecorationSkipInk",
					expected: "textDecorationSkipInk",
				},
			}

			for _, input := range inputs {
				So(toCamelCase(input.actual), ShouldEqual, input.expected)
			}
		})
	})
}
