package source

import (
	"strings"
	"testing"
)

func TestUnicodePositionsAddressOriginalBytes(t *testing.T) {
	for _, prefix := range []string{"abc", "é", "中", "😀", "😀😀", "é中😀", "e\u0301", "\t", "\xff", "\xf0\x9f", "é\x80😀\xff"} {
		for _, separator := range []string{" ", "\n", "\r\n"} {
			text := "LET text = \"" + prefix + "\"" + separator + "RETURN text"
			src := New("unicode.fql", text)
			start := strings.Index(text, "RETURN")
			span := Span{Start: start, End: len(text)}
			line := strings.Count(text[:start], "\n") + 1
			column := start - strings.LastIndex(text[:start], "\n")
			want := Position{Line: line, Column: column}
			if got := src.PositionAt(span); got != want {
				t.Fatalf("%q: position = %+v, want %+v", text, got, want)
			}

			got := src.RangeAt(span)
			if got.Span != span || text[got.Span.Start:got.Span.End] != "RETURN text" || got.Position != want {
				t.Fatalf("%q: range = %+v, want original byte span %+v", text, got, span)
			}
		}
	}
}
