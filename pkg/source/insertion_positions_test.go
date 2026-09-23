package source

import "testing"

func TestEOFInsertionPositions(t *testing.T) {
	for _, test := range []struct {
		text     string
		position Position
	}{
		{"RETURN", Position{Line: 1, Column: 7}},
		{"let arr = []\n\nreturn", Position{Line: 3, Column: 7}},
		{"é😀 RETURN", Position{Line: 1, Column: 14}},
	} {
		t.Run(test.text, func(t *testing.T) {
			src := New("test.fql", test.text)
			span := Span{Start: len(test.text), End: len(test.text)}
			if got := src.PositionAt(span); got != test.position {
				t.Fatalf("position = %+v, want %+v", got, test.position)
			}

			got := src.RangeAt(span)
			if got.Span != span || got.Position != test.position || got.SourceName != src.Name() || got.Location != src.LocationAt(span) {
				t.Fatalf("range = %+v, want unchanged span %+v at %+v", got, span, test.position)
			}

			invalid := Span{Start: span.Start, End: span.End + 1}
			if src.PositionAt(invalid) != (Position{}) || src.RangeAt(invalid).Span != invalid || src.LocationAt(invalid).Position != (Position{}) {
				t.Fatal("out-of-bounds span was accepted or repaired")
			}
		})
	}
}
