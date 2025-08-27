package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSkipWhitespaceForward(t *testing.T) {
	Convey("SkipWhitespaceForward", t, func() {
		Convey("Should skip spaces and tabs", func() {
			content := "   \t  hello"
			result := SkipWhitespaceForward(content, 0)
			So(result, ShouldEqual, 6)
		})

		Convey("Should skip newlines and carriage returns", func() {
			content := " \n\r\t  world"
			result := SkipWhitespaceForward(content, 0)
			So(result, ShouldEqual, 6)
		})

		Convey("Should handle empty string", func() {
			content := ""
			result := SkipWhitespaceForward(content, 0)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle offset at end of string", func() {
			content := "hello"
			result := SkipWhitespaceForward(content, 5)
			So(result, ShouldEqual, 5)
		})

		Convey("Should handle offset beyond end of string", func() {
			content := "hello"
			result := SkipWhitespaceForward(content, 10)
			So(result, ShouldEqual, 10) // Should not access out of bounds
		})

		Convey("Should handle string with no whitespace", func() {
			content := "hello"
			result := SkipWhitespaceForward(content, 0)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle string with only whitespace", func() {
			content := "   \t\n\r  "
			result := SkipWhitespaceForward(content, 0)
			So(result, ShouldEqual, len(content))
		})

		Convey("Should start from given offset", func() {
			content := "abc   def"
			result := SkipWhitespaceForward(content, 3)
			So(result, ShouldEqual, 6)
		})
	})
}

func TestSkipHorizontalWhitespaceForward(t *testing.T) {
	Convey("SkipHorizontalWhitespaceForward", t, func() {
		Convey("Should skip spaces and tabs only", func() {
			content := "   \t  hello"
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, 6)
		})

		Convey("Should stop at newlines", func() {
			content := "  \nworld"
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, 2) // Should stop at \n
		})

		Convey("Should stop at carriage returns", func() {
			content := "  \rworld"
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, 2) // Should stop at \r
		})

		Convey("Should handle empty string", func() {
			content := ""
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle offset at end of string", func() {
			content := "hello"
			result := SkipHorizontalWhitespaceForward(content, 5)
			So(result, ShouldEqual, 5)
		})

		Convey("Should handle offset beyond end of string", func() {
			content := "hello"
			result := SkipHorizontalWhitespaceForward(content, 10)
			So(result, ShouldEqual, 10) // Should not access out of bounds
		})

		Convey("Should handle string with no whitespace", func() {
			content := "hello"
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, 0)
		})

		Convey("Should handle string with only horizontal whitespace", func() {
			content := "   \t  "
			result := SkipHorizontalWhitespaceForward(content, 0)
			So(result, ShouldEqual, len(content))
		})
	})
}
