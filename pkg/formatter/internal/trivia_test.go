package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestTriviaEmitter_PreservesBlockCommentIndent(t *testing.T) {
	input := "/*\n * a\n * b\n */"
	var buf bytes.Buffer
	e := newEngine(file.NewAnonymousSource(input), &buf, DefaultOptions())

	e.trivia.emitTrivia(input, false, false)
	out := buf.String()
	if !strings.Contains(out, "\n * a\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
	if !strings.Contains(out, "\n * b\n") {
		t.Fatalf("expected leading space in block comment line; got:\n%s", out)
	}
}
