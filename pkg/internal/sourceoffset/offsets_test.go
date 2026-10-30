package sourceoffset

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestOffsets(t *testing.T) {
	for _, text := range []string{"", "ASCII\r\n", "é", "😀", "é中😀\r\nend", "e\u0301", "é\xff😀"} {
		t.Run(text, func(t *testing.T) {
			offsets := New(text)
			var prefix Prefix
			var boundaries []int
			for offset := range text {
				boundaries = append(boundaries, offset)
			}

			boundaries = append(boundaries, len(text))
			if offsets.Len(len(text)) != len(boundaries)-1 {
				t.Fatalf("rune length = %d, want %d", offsets.Len(len(text)), len(boundaries)-1)
			}

			if len(boundaries)-1 == len(text) && offsets != nil {
				t.Fatal("ASCII/empty source allocated a table")
			}

			for start, wantStart := range boundaries {
				if offsets.Byte(start) != wantStart {
					t.Fatalf("Byte(%d) = %d, want %d", start, offsets.Byte(start), wantStart)
				}

				for end := start; end < len(boundaries); end++ {
					got := offsets.Span(source.Span{Start: start, End: end})
					want := source.Span{Start: wantStart, End: boundaries[end]}
					if got != want {
						t.Fatalf("Span(%d,%d) = %+v, want %+v", start, end, got, want)
					}
					if got := prefix.Span(text, source.Span{Start: start, End: end}); got != want {
						t.Fatalf("prefix Span(%d,%d) = %+v, want %+v", start, end, got, want)
					}
					if start == 0 && prefix.bytes != boundaries[end] {
						t.Fatalf("prefix scanned %d bytes for boundary %d, want %d", prefix.bytes, end, boundaries[end])
					}
				}
			}

			for _, span := range []source.Span{
				{Start: -1, End: 1}, {Start: 1, End: -1},
				{Start: len(boundaries), End: len(boundaries) + 1}, {Start: 0, End: len(boundaries)},
			} {
				if got := offsets.Span(span); got != span {
					t.Fatalf("unconvertible span = %+v, want unchanged %+v", got, span)
				}
				if got := prefix.Span(text, span); got != span {
					t.Fatalf("unconvertible prefix span = %+v, want unchanged %+v", got, span)
				}
			}

			for _, offset := range []int{-1, len(boundaries)} {
				if got := offsets.Byte(offset); got != offset {
					t.Fatalf("unconvertible offset = %d, want %d", got, offset)
				}
			}
		})
	}
}
