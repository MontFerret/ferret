package file

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewSource(t *testing.T) {
	Convey("NewSource", t, func() {
		Convey("Should create source with name and text", func() {
			name := "test.fql"
			text := "hello\nworld"

			source := NewSource(name, text)

			So(source, ShouldNotBeNil)
			So(source.Name(), ShouldEqual, name)
			So(source.Content(), ShouldEqual, text)
			So(source.Empty(), ShouldBeFalse)
		})

		Convey("Should handle empty text", func() {
			name := "empty.fql"
			text := ""

			source := NewSource(name, text)

			So(source, ShouldNotBeNil)
			So(source.Name(), ShouldEqual, name)
			So(source.Content(), ShouldEqual, text)
			So(source.Empty(), ShouldBeTrue)
		})
	})
}

func TestNewAnonymousSource(t *testing.T) {
	Convey("NewAnonymousSource", t, func() {
		Convey("Should create anonymous source", func() {
			text := "test content"

			source := NewAnonymousSource(text)

			So(source, ShouldNotBeNil)
			So(source.Name(), ShouldEqual, "anonymous")
			So(source.Content(), ShouldEqual, text)
		})
	})
}

func TestSourceName(t *testing.T) {
	Convey("Source.Name", t, func() {
		Convey("Should return name for valid source", func() {
			source := NewSource("test.fql", "content")

			So(source.Name(), ShouldEqual, "test.fql")
		})

		Convey("Should return empty string for nil source", func() {
			var source *Source = nil

			So(source.Name(), ShouldEqual, "")
		})
	})
}

func TestSourceEmpty(t *testing.T) {
	Convey("Source.Empty", t, func() {
		Convey("Should return false for non-empty source", func() {
			source := NewSource("test.fql", "content")

			So(source.Empty(), ShouldBeFalse)
		})

		Convey("Should return true for empty text", func() {
			source := NewSource("test.fql", "")

			So(source.Empty(), ShouldBeTrue)
		})

		Convey("Should return true for nil source", func() {
			var source *Source = nil

			So(source.Empty(), ShouldBeTrue)
		})
	})
}

func TestSourceLocationAt(t *testing.T) {
	Convey("Source.LocationAt", t, func() {
		Convey("Simple single line text", func() {
			source := NewSource("test.fql", "hello world")

			Convey("Should return correct location at start", func() {
				line, col := source.LocationAt(Span{Start: 0, End: 1})
				So(line, ShouldEqual, 1)
				So(col, ShouldEqual, 1)
			})

			Convey("Should return correct location in middle", func() {
				line, col := source.LocationAt(Span{Start: 6, End: 7})
				So(line, ShouldEqual, 1)
				So(col, ShouldEqual, 7)
			})

			Convey("Should return correct location at end", func() {
				line, col := source.LocationAt(Span{Start: 10, End: 11})
				So(line, ShouldEqual, 1)
				So(col, ShouldEqual, 11)
			})
		})

		Convey("Multi-line text", func() {
			source := NewSource("test.fql", "line1\nline2\nline3")

			Convey("Should return correct location on first line", func() {
				line, col := source.LocationAt(Span{Start: 2, End: 3})
				So(line, ShouldEqual, 1)
				So(col, ShouldEqual, 3)
			})

			Convey("Should return correct location on second line", func() {
				line, col := source.LocationAt(Span{Start: 8, End: 9}) // 'n' in "line2"
				So(line, ShouldEqual, 2)
				So(col, ShouldEqual, 3)
			})

			Convey("Should return correct location on third line", func() {
				line, col := source.LocationAt(Span{Start: 14, End: 15}) // 'n' in "line3"
				So(line, ShouldEqual, 3)
				So(col, ShouldEqual, 3)
			})

			Convey("Should handle location at newline", func() {
				line, col := source.LocationAt(Span{Start: 5, End: 6}) // First \n
				So(line, ShouldEqual, 1)
				So(col, ShouldEqual, 6)
			})

			Convey("Should handle location at start of line after newline", func() {
				line, col := source.LocationAt(Span{Start: 6, End: 7}) // Start of "line2"
				So(line, ShouldEqual, 1)                               // Should be treated as end of line1
				So(col, ShouldEqual, 6)                                // Column after last char of "line1" (len("line1") + 1)
			})
		})

		Convey("Edge cases", func() {
			source := NewSource("test.fql", "hello\nworld")

			Convey("Should handle negative start", func() {
				line, col := source.LocationAt(Span{Start: -1, End: 0})
				So(line, ShouldEqual, 0)
				So(col, ShouldEqual, 0)
			})

			Convey("Should handle end beyond content", func() {
				line, col := source.LocationAt(Span{Start: 0, End: 100})
				So(line, ShouldEqual, 0)
				So(col, ShouldEqual, 0)
			})

			Convey("Should handle empty source", func() {
				emptySource := NewSource("empty.fql", "")
				line, col := emptySource.LocationAt(Span{Start: 0, End: 1})
				So(line, ShouldEqual, 0)
				So(col, ShouldEqual, 0)
			})

			Convey("Should handle nil source", func() {
				var nilSource *Source = nil
				line, col := nilSource.LocationAt(Span{Start: 0, End: 1})
				So(line, ShouldEqual, 0)
				So(col, ShouldEqual, 0)
			})
		})
	})
}

func TestSourceSnippet(t *testing.T) {
	Convey("Source.Snippet", t, func() {
		Convey("Single line source", func() {
			source := NewSource("test.fql", "hello world")
			span := Span{Start: 6, End: 11}

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 1) // Only one line, no previous/next
			So(snippets[0].Line, ShouldEqual, 1)
			So(snippets[0].Text, ShouldEqual, "hello world")
			So(snippets[0].Caret, ShouldNotBeEmpty)
		})

		Convey("Multi-line source", func() {
			source := NewSource("test.fql", "line1\nline2\nline3")
			span := Span{Start: 8, End: 10} // "in" in "line2"

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 3) // Previous, current, and next line
			So(snippets[0].Line, ShouldEqual, 1)
			So(snippets[0].Text, ShouldEqual, "line1")
			So(snippets[1].Line, ShouldEqual, 2)
			So(snippets[1].Text, ShouldEqual, "line2")
			So(snippets[2].Line, ShouldEqual, 3)
			So(snippets[2].Text, ShouldEqual, "line3")
		})

		Convey("First line span", func() {
			source := NewSource("test.fql", "line1\nline2\nline3")
			span := Span{Start: 2, End: 4} // "ne" in "line1"

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 2) // No previous line
			So(snippets[0].Line, ShouldEqual, 1)
			So(snippets[0].Text, ShouldEqual, "line1")
			So(snippets[1].Line, ShouldEqual, 2)
			So(snippets[1].Text, ShouldEqual, "line2")
		})

		Convey("Last line span", func() {
			source := NewSource("test.fql", "line1\nline2\nline3")
			span := Span{Start: 14, End: 16} // "ne" in "line3"

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 2) // No next line
			So(snippets[0].Line, ShouldEqual, 2)
			So(snippets[0].Text, ShouldEqual, "line2")
			So(snippets[1].Line, ShouldEqual, 3)
			So(snippets[1].Text, ShouldEqual, "line3")
		})

		Convey("Empty source", func() {
			source := NewSource("empty.fql", "")
			span := Span{Start: 0, End: 1}

			snippets := source.Snippet(span)

			So(snippets, ShouldBeNil)
		})

		Convey("Nil source", func() {
			var source *Source = nil
			span := Span{Start: 0, End: 1}

			snippets := source.Snippet(span)

			So(snippets, ShouldBeNil)
		})
	})
}

func TestComputeVisualOffset(t *testing.T) {
	Convey("computeVisualOffset", t, func() {
		Convey("Should handle regular characters", func() {
			line := "hello world"
			offset := computeVisualOffset(line, 6)
			So(offset, ShouldEqual, 5) // Column 6 = position 5 (0-based)
		})

		Convey("Should handle tabs with default width 4", func() {
			line := "\thello"
			offset := computeVisualOffset(line, 2)
			So(offset, ShouldEqual, 4) // Tab expands to 4 spaces
		})

		Convey("Should handle multiple tabs", func() {
			line := "\t\thello"
			offset := computeVisualOffset(line, 3)
			So(offset, ShouldEqual, 8) // Two tabs = 8 spaces
		})

		Convey("Should handle mixed tabs and spaces", func() {
			line := " \t hello"
			offset := computeVisualOffset(line, 4)
			So(offset, ShouldEqual, 5) // space(1) + tab(3 more to 4) + space(1) = 5
		})

		Convey("Should handle column beyond line length", func() {
			line := "hello"
			offset := computeVisualOffset(line, 10)
			So(offset, ShouldEqual, 5) // Should not go beyond line length
		})

		Convey("Should handle empty line", func() {
			line := ""
			offset := computeVisualOffset(line, 1)
			So(offset, ShouldEqual, 0)
		})

		Convey("Should handle unicode characters", func() {
			line := "hello 世界"
			offset := computeVisualOffset(line, 8)
			So(offset, ShouldEqual, 7) // Unicode characters count as 1 visual position each
		})
	})
}
