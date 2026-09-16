package compiler

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestUnicodeDebugMetadataAddressesOriginalBytes(t *testing.T) {
	for _, prefix := range []string{"abc", "é", "中", "😀", "😀😀", "é中😀", "e\u0301", "\t", "\xff", "\xf0\x9f", "é\x80😀\xff"} {
		for _, separator := range []string{" ", "\n", "\r\n"} {
			query := "LET text = \"" + prefix + "\"" + separator + "RETURN TEST(text)"
			program, err := mustNewCompiler(t, WithDebugInfo()).Compile(t.Context(), source.New("unicode.fql", query))
			if err != nil {
				t.Fatalf("%q: %v", query, err)
			}

			want := source.Span{Start: strings.Index(query, "RETURN"), End: len(query)}
			pointFound, instructionFound, argumentFound := false, false, false
			for _, point := range program.Metadata.DebugPoints {
				pointFound = pointFound || point.Span == want
			}

			for _, span := range program.Metadata.DebugSpans {
				instructionFound = instructionFound || span == want
			}

			argumentStart := strings.LastIndex(query, "text")
			argument := source.Span{Start: argumentStart, End: argumentStart + len("text")}
			for _, spans := range program.Metadata.CallArgumentSpans {
				for _, span := range spans {
					argumentFound = argumentFound || span == argument
				}
			}

			if !pointFound || !instructionFound || !argumentFound {
				t.Fatalf("%q: missing byte spans: point=%t instruction=%t argument=%t; want %+v / %+v", query, pointFound, instructionFound, argumentFound, want, argument)
			}
		}
	}
}
