package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestNewErrorSpan(t *testing.T) {
	span := file.Span{Start: 0, End: 10}
	label := "test label"
	main := true

	result := NewErrorSpan(span, label, main)

	if result.Span != span {
		t.Errorf("NewErrorSpan() span = %v, want %v", result.Span, span)
	}
	if result.Label != label {
		t.Errorf("NewErrorSpan() label = %v, want %v", result.Label, label)
	}
	if result.Main != main {
		t.Errorf("NewErrorSpan() main = %v, want %v", result.Main, main)
	}
}

func TestNewMainErrorSpan(t *testing.T) {
	span := file.Span{Start: 0, End: 10}
	label := "main error"

	result := NewMainErrorSpan(span, label)

	if result.Span != span {
		t.Errorf("NewMainErrorSpan() span = %v, want %v", result.Span, span)
	}
	if result.Label != label {
		t.Errorf("NewMainErrorSpan() label = %v, want %v", result.Label, label)
	}
	if !result.Main {
		t.Error("NewMainErrorSpan() should create a main error span")
	}
}

func TestNewSecondaryErrorSpan(t *testing.T) {
	span := file.Span{Start: 5, End: 15}
	label := "secondary error"

	result := NewSecondaryErrorSpan(span, label)

	if result.Span != span {
		t.Errorf("NewSecondaryErrorSpan() span = %v, want %v", result.Span, span)
	}
	if result.Label != label {
		t.Errorf("NewSecondaryErrorSpan() label = %v, want %v", result.Label, label)
	}
	if result.Main {
		t.Error("NewSecondaryErrorSpan() should create a non-main error span")
	}
}