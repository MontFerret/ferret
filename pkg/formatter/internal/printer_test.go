package internal

import (
	"bytes"
	"testing"
)

func TestPrinter_WriteRawTracksByteWidth(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf, DefaultOptions())

	p.writeRaw("πa")

	if got := buf.String(); got != "πa" {
		t.Fatalf("unexpected raw output: %q", got)
	}
	if p.lineColumn != len("πa") {
		t.Fatalf("expected byte width %d, got %d", len("πa"), p.lineColumn)
	}
	if p.atLineStart {
		t.Fatalf("expected printer to leave line start after raw write")
	}
	if p.lastWasSpace {
		t.Fatalf("expected lastWasSpace to remain false")
	}
}

func TestPrinter_WriteRawPreservesNewlineState(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf, DefaultOptions())

	p.writeRaw("π\nβ")

	if got := buf.String(); got != "π\nβ" {
		t.Fatalf("unexpected raw output: %q", got)
	}
	if !p.sawHardNewline {
		t.Fatalf("expected sawHardNewline to be set")
	}
	if p.lineColumn != len("β") {
		t.Fatalf("expected line column %d, got %d", len("β"), p.lineColumn)
	}
	if p.atLineStart {
		t.Fatalf("expected printer to end off line start after trailing segment")
	}
	if p.lastWasSpace {
		t.Fatalf("expected lastWasSpace to be false after trailing segment")
	}
}

func TestPrinter_WriteRawForceSingleLineCollapsesNewline(t *testing.T) {
	var buf bytes.Buffer
	p := newPrinter(&buf, DefaultOptions())
	p.forceSingleLine = true

	p.writeRaw("π\nβ")

	if got := buf.String(); got != "π β" {
		t.Fatalf("unexpected raw output: %q", got)
	}
	if !p.sawHardNewline {
		t.Fatalf("expected sawHardNewline to be set")
	}
	if p.lineColumn != len("π β") {
		t.Fatalf("expected line column %d, got %d", len("π β"), p.lineColumn)
	}
	if p.atLineStart {
		t.Fatalf("expected printer to end off line start after trailing segment")
	}
	if p.lastWasSpace {
		t.Fatalf("expected lastWasSpace to be false after trailing segment")
	}
}
