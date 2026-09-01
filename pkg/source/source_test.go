package source

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNew(t *testing.T) {
	Convey("New", t, func() {
		Convey("Should create source with name and text", func() {
			name := "test.fql"
			text := "hello\nworld"

			source := New(name, text)

			So(source.Name(), ShouldEqual, name)
			So(source.Content(), ShouldEqual, text)
			So(source.Length(), ShouldEqual, len(text))
			So(source.ID(), ShouldNotResemble, ID{})
			So(source.Empty(), ShouldBeFalse)
		})

		Convey("Should handle empty text", func() {
			name := "empty.fql"
			text := ""

			source := New(name, text)

			So(source.Name(), ShouldEqual, name)
			So(source.Content(), ShouldEqual, text)
			So(source.Length(), ShouldEqual, 0)
			So(source.ID(), ShouldNotResemble, ID{})
			So(source.Empty(), ShouldBeTrue)
		})
	})
}

func TestNewAnonymous(t *testing.T) {
	Convey("NewAnonymous", t, func() {
		Convey("Should create anonymous source", func() {
			text := "test content"

			source := NewAnonymous(text)

			So(source.Name(), ShouldEqual, "anonymous")
			So(source.Content(), ShouldEqual, text)
			So(source.Length(), ShouldEqual, len(text))
			So(source.ID(), ShouldResemble, New("anonymous", text).ID())
		})
	})
}

func TestSourceID(t *testing.T) {
	Convey("Source.ID", t, func() {
		first := New("test.fql", "RETURN 1")
		same := New("test.fql", "RETURN 1")
		differentName := New("other.fql", "RETURN 1")
		differentContent := New("test.fql", "RETURN 2")

		So(first.ID(), ShouldNotResemble, ID{})
		So(first.ID(), ShouldResemble, same.ID())
		So(first.ID(), ShouldNotResemble, differentName.ID())
		So(first.ID(), ShouldNotResemble, differentContent.ID())
	})
}

func TestSourceZeroValue(t *testing.T) {
	Convey("The zero Source value represents no source", t, func() {
		var source Source

		So(source.Name(), ShouldEqual, "")
		So(source.Content(), ShouldEqual, "")
		So(source.Length(), ShouldEqual, 0)
		So(source.ID(), ShouldResemble, ID{})
		So(source.Empty(), ShouldBeTrue)

		position := source.PositionAt(Span{Start: 0, End: 1})
		So(position, ShouldResemble, Position{})
		So(source.LocationAt(Span{Start: 0, End: 1}), ShouldResemble, Location{})
		So(source.RangeAt(Span{Start: 0, End: 1}), ShouldResemble, Range{Span: Span{Start: 0, End: 1}})
		So(source.Snippet(Span{Start: 0, End: 1}), ShouldBeNil)
	})
}

func TestSourceName(t *testing.T) {
	Convey("Source.Name", t, func() {
		Convey("Should return name for valid source", func() {
			source := New("test.fql", "content")

			So(source.Name(), ShouldEqual, "test.fql")
		})
	})
}

func TestSourceEmpty(t *testing.T) {
	Convey("Source.Empty", t, func() {
		Convey("Should return false for non-empty source", func() {
			source := New("test.fql", "content")

			So(source.Empty(), ShouldBeFalse)
		})

		Convey("Should return true for empty text", func() {
			source := New("test.fql", "")

			So(source.Empty(), ShouldBeTrue)
		})
	})
}

func TestSourcePositionAt(t *testing.T) {
	Convey("Source.PositionAt", t, func() {
		Convey("Simple single line text", func() {
			source := New("test.fql", "hello world")

			Convey("Should return correct location at start", func() {
				position := source.PositionAt(Span{Start: 0, End: 1})
				So(position, ShouldResemble, Position{Line: 1, Column: 1})
			})

			Convey("Should return correct location in middle", func() {
				position := source.PositionAt(Span{Start: 6, End: 7})
				So(position, ShouldResemble, Position{Line: 1, Column: 7})
			})

			Convey("Should return correct location at end", func() {
				position := source.PositionAt(Span{Start: 10, End: 11})
				So(position, ShouldResemble, Position{Line: 1, Column: 11})
			})
		})

		Convey("Multi-line text", func() {
			source := New("test.fql", "line1\nline2\nline3")

			Convey("Should return correct location on first line", func() {
				position := source.PositionAt(Span{Start: 2, End: 3})
				So(position, ShouldResemble, Position{Line: 1, Column: 3})
			})

			Convey("Should return correct location on second line", func() {
				position := source.PositionAt(Span{Start: 8, End: 9}) // 'n' in "line2"
				So(position, ShouldResemble, Position{Line: 2, Column: 3})
			})

			Convey("Should return correct location on third line", func() {
				position := source.PositionAt(Span{Start: 14, End: 15}) // 'n' in "line3"
				So(position, ShouldResemble, Position{Line: 3, Column: 3})
			})

			Convey("Should handle location at newline", func() {
				position := source.PositionAt(Span{Start: 5, End: 6}) // First \n
				So(position, ShouldResemble, Position{Line: 1, Column: 6})
			})

			Convey("Should handle location at start of line after newline", func() {
				position := source.PositionAt(Span{Start: 6, End: 7}) // Start of "line2"
				So(position, ShouldResemble, Position{Line: 2, Column: 1})
			})
		})

		Convey("Edge cases", func() {
			source := New("test.fql", "hello\nworld")

			Convey("Should handle negative start", func() {
				position := source.PositionAt(Span{Start: -1, End: 0})
				So(position, ShouldResemble, Position{})
			})

			Convey("Should handle end beyond content", func() {
				position := source.PositionAt(Span{Start: 0, End: 100})
				So(position, ShouldResemble, Position{})
			})

			Convey("Should handle empty source", func() {
				emptySource := New("empty.fql", "")
				position := emptySource.PositionAt(Span{Start: 0, End: 1})
				So(position, ShouldResemble, Position{})
			})
		})
	})
}

func TestSourceLocationAndRangeAt(t *testing.T) {
	Convey("Source structured locations", t, func() {
		source := New("test.fql", "hello\nworld")
		span := Span{Start: 6, End: 11}

		So(source.LocationAt(span), ShouldResemble, Location{
			File:     "test.fql",
			Position: Position{Line: 2, Column: 1},
		})
		So(source.RangeAt(span), ShouldResemble, Range{
			Location: Location{
				File:     "test.fql",
				Position: Position{Line: 2, Column: 1},
			},
			Span: span,
		})

		invalid := Span{Start: -1, End: 0}
		So(source.LocationAt(invalid), ShouldResemble, Location{File: "test.fql"})
		So(source.RangeAt(invalid), ShouldResemble, Range{
			Location: Location{File: "test.fql"},
			Span:     invalid,
		})
	})
}

func TestSourceSnippet(t *testing.T) {
	Convey("Source.Snippet", t, func() {
		Convey("Single line source", func() {
			source := New("test.fql", "hello world")
			span := Span{Start: 6, End: 11}

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 1) // Only one line, no previous/next
			So(snippets[0].Line, ShouldEqual, 1)
			So(snippets[0].Text, ShouldEqual, "hello world")
			So(snippets[0].Caret, ShouldNotBeEmpty)
		})

		Convey("Multi-line source", func() {
			source := New("test.fql", "line1\nline2\nline3")
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
			source := New("test.fql", "line1\nline2\nline3")
			span := Span{Start: 2, End: 4} // "ne" in "line1"

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 2) // No previous line
			So(snippets[0].Line, ShouldEqual, 1)
			So(snippets[0].Text, ShouldEqual, "line1")
			So(snippets[1].Line, ShouldEqual, 2)
			So(snippets[1].Text, ShouldEqual, "line2")
		})

		Convey("Last line span", func() {
			source := New("test.fql", "line1\nline2\nline3")
			span := Span{Start: 14, End: 16} // "ne" in "line3"

			snippets := source.Snippet(span)

			So(len(snippets), ShouldEqual, 2) // No next line
			So(snippets[0].Line, ShouldEqual, 2)
			So(snippets[0].Text, ShouldEqual, "line2")
			So(snippets[1].Line, ShouldEqual, 3)
			So(snippets[1].Text, ShouldEqual, "line3")
		})

		Convey("Empty source", func() {
			source := New("empty.fql", "")
			span := Span{Start: 0, End: 1}

			snippets := source.Snippet(span)

			So(snippets, ShouldBeNil)
		})
	})
}

func TestSourceJSONRoundTrip(t *testing.T) {
	Convey("Source JSON round-trip", t, func() {
		original := New("roundtrip.fql", "first\nsecond")

		encoded, err := original.MarshalJSON()
		So(err, ShouldBeNil)
		So(string(encoded), ShouldEqual, `{"name":"roundtrip.fql","text":"first\nsecond"}`)

		var decoded Source

		err = decoded.UnmarshalJSON(encoded)
		So(err, ShouldBeNil)
		So(decoded.Name(), ShouldEqual, original.Name())
		So(decoded.Content(), ShouldEqual, original.Content())
		So(decoded.Length(), ShouldEqual, original.Length())
		So(decoded.ID(), ShouldResemble, original.ID())

		position := decoded.PositionAt(Span{Start: 6, End: 7})
		So(position, ShouldResemble, Position{Line: 2, Column: 1})
		So(decoded.Snippet(Span{Start: 6, End: 7}), ShouldHaveLength, 2)
	})

	Convey("UnmarshalJSON rejects a nil receiver", t, func() {
		var source *Source

		err := source.UnmarshalJSON([]byte(`{"name":"test.fql","text":"RETURN 1"}`))
		So(err, ShouldNotBeNil)
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
