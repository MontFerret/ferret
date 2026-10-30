package internal

import (
	"errors"
	"testing"
)

type indentFailureWriter struct {
	first  error
	later  error
	failed bool
}

func (w *indentFailureWriter) Write([]byte) (int, error) {
	if w.failed {
		return 0, w.later
	}

	w.failed = true

	return 0, w.first
}

func TestPrinterStopsOnIndentFailure(t *testing.T) {
	first := errors.New("indent failure")
	w := &indentFailureWriter{first: first, later: errors.New("write after failure")}
	p := newPrinter(w, defaultTestConfig())
	p.indent = 1
	p.write("return")
	p.space()
	p.newline()
	p.writeRaw("ignored\n")

	if !errors.Is(p.Err(), first) {
		t.Fatalf("first writer error was replaced: %v", p.Err())
	}
}
