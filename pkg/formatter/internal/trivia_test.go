package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestTriviaEmitter_PreservesBlockCommentIndent(t *testing.T) {
	input := "/*\n * a\n * b\n */"
	var buf bytes.Buffer
	e := newTestEngine(source.NewAnonymous(input), &buf, defaultTestConfig())

	e.trivia.emitTrivia(input, false, false)
	out := buf.String()
	if !strings.Contains(out, "\n * a\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
	if !strings.Contains(out, "\n * b\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
}

func TestTriviaRuneBoundaries(t *testing.T) {
	var out bytes.Buffer
	e := newTestEngine(source.NewAnonymous("a😀é// 尾\r\n"), &out, defaultTestConfig())
	for _, test := range []struct {
		want  string
		start int
		end   int
	}{
		{start: 0, end: 1, want: "a"}, {start: 1, end: 2, want: "😀"},
		{start: 2, end: 3, want: "é"}, {start: 3, end: 9, want: "// 尾\r\n"},
		{start: 9, end: 9}, {start: -1, end: 2}, {start: 2, end: 1},
		{start: 0, end: 10}, {start: 10, end: 11},
	} {
		if got := e.trivia.sliceBetween(test.start, test.end); got != test.want {
			t.Fatalf("sliceBetween(%d,%d) = %q, want %q", test.start, test.end, got, test.want)
		}
	}

	for _, test := range []struct {
		want  string
		start int
		stop  int
	}{
		{start: 1, stop: 1, want: "😀"}, {start: 2, stop: 2, want: "é"},
		{start: 9, stop: 8}, {start: -1, stop: -1},
	} {
		tok := antlr.NewCommonToken(&antlr.TokenSourceCharStreamPair{}, antlr.TokenEOF, 0, test.start, test.stop)
		node := antlr.NewTerminalNodeImpl(tok)
		ctx := antlr.NewBaseParserRuleContext(nil, -1)
		ctx.SetStart(tok)
		ctx.SetStop(tok)
		if got := e.trivia.sliceBetween(e.trivia.tokenStart(node), e.trivia.tokenStop(node)+1); got != test.want {
			t.Fatalf("token slice = %q, want %q", got, test.want)
		}

		if got := e.trivia.sliceBetween(e.trivia.startIndex(ctx), e.trivia.stopIndex(ctx)+1); got != test.want {
			t.Fatalf("rule slice = %q, want %q", got, test.want)
		}
	}

	if e.trivia.sliceBetween(e.trivia.tokenStop(nil)+1, 2) != "" || e.trivia.sliceBetween(e.trivia.stopIndex(nil)+1, 2) != "" {
		t.Fatal("absent stop produced source trivia")
	}
}
