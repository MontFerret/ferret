package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSnippet(t *testing.T) {
	Convey("NewSnippet", t, func() {
		Convey("Should create snippet from valid line", func() {
			src := []string{"line1", "line2", "line3"}

			snippet := NewSnippet(src, 2)

			So(snippet.Line, ShouldEqual, 2)
			So(snippet.Text, ShouldEqual, "line2")
			So(snippet.Caret, ShouldEqual, "")
		})

		Convey("Should create snippet for first line", func() {
			src := []string{"first line", "second line"}

			snippet := NewSnippet(src, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "first line")
		})

		Convey("Should create snippet for last line", func() {
			src := []string{"first line", "last line"}

			snippet := NewSnippet(src, 2)

			So(snippet.Line, ShouldEqual, 2)
			So(snippet.Text, ShouldEqual, "last line")
		})

		Convey("Should handle bounds issues gracefully", func() {
			src := []string{"line1", "line2"}

			Convey("Line number 0 should return empty snippet with line number", func() {
				snippet := NewSnippet(src, 0)
				So(snippet.Line, ShouldEqual, 0)
				So(snippet.Text, ShouldEqual, "")
			})

			Convey("Line number beyond bounds should return empty snippet with line number", func() {
				snippet := NewSnippet(src, 10)
				So(snippet.Line, ShouldEqual, 10)
				So(snippet.Text, ShouldEqual, "")
			})
		})
	})
}

func TestNewSnippetWithCaret(t *testing.T) {
	Convey("NewSnippetWithCaret", t, func() {
		Convey("Should create snippet with single character caret", func() {
			lines := []string{"hello world"}
			span := Span{Start: 6, End: 7} // Single character 'w'

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "hello world")
			So(snippet.Caret, ShouldEqual, "      ^")
		})

		Convey("Should create snippet with multi-character caret", func() {
			lines := []string{"hello world"}
			span := Span{Start: 6, End: 11} // "world"

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "hello world")
			So(snippet.Caret, ShouldEqual, "      ~~~~~")
		})

		Convey("Should handle multi-line text", func() {
			lines := []string{"first line", "second line", "third line"}
			span := Span{Start: 11, End: 17} // "second" in "second line" (starts at position 11)

			snippet := NewSnippetWithCaret(lines, span, 2)

			So(snippet.Line, ShouldEqual, 2)
			So(snippet.Text, ShouldEqual, "second line")
			So(snippet.Caret, ShouldEqual, "~~~~~~")
		})

		Convey("Should handle tabs in visual offset", func() {
			lines := []string{"\thello world"}
			span := Span{Start: 5, End: 10} // "world" after tab

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "\thello world")
			// Tab should expand to 4 spaces, so caret starts at position 4+5=9
			So(snippet.Caret, ShouldEqual, "        ~~~~~")
		})

		Convey("Should handle line number out of bounds", func() {
			lines := []string{"line1", "line2"}
			span := Span{Start: 0, End: 5}

			Convey("Line number too small", func() {
				snippet := NewSnippetWithCaret(lines, span, 0)

				So(snippet, ShouldResemble, Snippet{})
			})

			Convey("Line number too large", func() {
				snippet := NewSnippetWithCaret(lines, span, 3)

				So(snippet, ShouldResemble, Snippet{})
			})
		})

		Convey("Should handle empty lines", func() {
			lines := []string{""}
			span := Span{Start: 0, End: 0}

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "")
			So(snippet.Caret, ShouldEqual, "^")
		})

		Convey("Should handle span at line start", func() {
			lines := []string{"hello world"}
			span := Span{Start: 0, End: 1} // First character

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "hello world")
			So(snippet.Caret, ShouldEqual, "^")
		})

		Convey("Should handle span at line end", func() {
			lines := []string{"hello world"}
			span := Span{Start: 10, End: 11} // Last character

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, "hello world")
			So(snippet.Caret, ShouldEqual, "          ^")
		})

		Convey("Should handle mixed whitespace", func() {
			lines := []string{" \t hello"}
			span := Span{Start: 5, End: 6} // First character of "hello"

			snippet := NewSnippetWithCaret(lines, span, 1)

			So(snippet.Line, ShouldEqual, 1)
			So(snippet.Text, ShouldEqual, " \t hello")
			// space(1) + tab(3 more to align to 4) + space(1) = 5 spaces before caret
			So(snippet.Caret, ShouldEqual, "       ^")
		})
	})
}
