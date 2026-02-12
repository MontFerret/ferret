package diagnostics

import (
	"errors"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/file"
)

// Simple handler tests without complex ANTLR dependencies
func TestErrorHandler_BasicOperations(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	// Test NewErrorHandler with various thresholds
	handler1 := NewErrorHandler(src, 5)
	if handler1 == nil {
		t.Fatal("NewErrorHandler() returned nil")
	}
	if handler1.threshold != 5 {
		t.Errorf("Expected threshold 5, got %d", handler1.threshold)
	}

	// Test with zero threshold (should default to 10)
	handler2 := NewErrorHandler(src, 0)
	if handler2.threshold != 10 {
		t.Errorf("Expected default threshold 10, got %d", handler2.threshold)
	}

	// Test with negative threshold (should default to 10)
	handler3 := NewErrorHandler(src, -5)
	if handler3.threshold != 10 {
		t.Errorf("Expected default threshold 10, got %d", handler3.threshold)
	}

	// Test initial state
	if handler1.HasErrors() {
		t.Error("New handler should not have errors")
	}

	if handler1.Errors().Size() != 0 {
		t.Error("New handler should have empty errors slice")
	}

	if handler1.Unwrap() != nil {
		t.Error("Unwrap() should return nil for no errors")
	}
}

func TestErrorHandler_AddNilError(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	handler := NewErrorHandler(src, 10)

	// Adding nil error should be ignored
	handler.Add(nil)

	if handler.HasErrors() {
		t.Error("Adding nil error should not create errors")
	}
}

func TestErrorHandler_AddSingleError(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	handler := NewErrorHandler(src, 10)

	err := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "test error",
			Source:  src,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(file.Span{Start: 0, End: 3}, ""),
			},
		},
	}

	handler.Add(err)

	if !handler.HasErrors() {
		t.Error("Handler should have errors after adding one")
	}

	errs := handler.Errors()
	if errs.Size() != 1 {
		t.Errorf("Handler should have 1 error, got %d", errs.Size())
	}

	if !errors.Is(err, errs.Get(0)) {
		t.Error("Added error should be the same as retrieved error")
	}

	// Test Unwrap with single error
	unwrapped := handler.Unwrap()

	if !errors.Is(unwrapped, err) {
		t.Error("Unwrap() should return the single error directly")
	}
}

func TestErrorHandler_AddMultipleErrors(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	handler := NewErrorHandler(src, 10)

	err1 := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "error 1",
		},
	}

	err2 := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    NameError,
			Message: "error 2",
		},
	}

	handler.Add(err1)
	handler.Add(err2)

	if handler.Errors().Size() != 2 {
		t.Errorf("Handler should have 2 errors, got %d", handler.Errors().Size())
	}

	// Test Unwrap with multiple errors
	unwrapped := handler.Unwrap()
	if unwrapped == nil {
		t.Fatal("Unwrap() should not return nil for multiple errors")
	}

	multiErr, ok := unwrapped.(*CompilationErrorSet)
	if !ok {
		t.Error("Unwrap() should return *CompilationErrorSet for multiple errors")
	}

	if multiErr.Size() != 2 {
		t.Errorf("CompilationErrorSet should have 2 errors, got %d", multiErr.Size())
	}
}

func TestErrorHandler_HasErrorOnLine(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1\nRETURN x") // 2 lines
	handler := NewErrorHandler(src, 10)

	// Initially no errors on any line
	if handler.HasErrorOnLine(1) {
		t.Error("Should not have error on line 1 initially")
	}

	// Add error with span that affects line 1
	err := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "test error",
			Source:  src,
			Spans: []diagnostics.ErrorSpan{
				diagnostics.NewMainErrorSpan(file.Span{Start: 0, End: 3}, ""), // Position 0-3 is on line 1
			},
		},
	}

	handler.Add(err)

	if !handler.HasErrorOnLine(1) {
		t.Error("Should have error on line 1 after adding error")
	}
}

func TestErrorHandler_ExceedThreshold(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	handler := NewErrorHandler(src, 2) // Low threshold for testing

	// Add errors up to threshold
	err1 := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "error 1",
		},
	}
	err2 := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "error 2",
		},
	}

	handler.Add(err1)
	handler.Add(err2)

	// At exactly threshold, should trigger "too many errors" message
	errs := handler.Errors()
	if errs.Size() != 3 { // 2 actual errors + 1 "too many" message
		t.Errorf("Handler should have 3 errors (2 + 'too many' message), got %d", errs.Size())
	}

	// Last error should be "Too many errors"
	lastErr := errs.Last()
	if lastErr.Message != "Too many errors" {
		t.Errorf("Last error should be 'Too many errors', got %q", lastErr.Message)
	}

	// Adding more errors should be ignored (since len(errors) > threshold now)
	err3 := &CompilationError{
		Diagnostic: &diagnostics.Diagnostic{
			Kind:    SyntaxError,
			Message: "ignored",
		},
	}
	handler.Add(err3)

	if handler.Errors().Size() != 3 {
		t.Errorf("Handler should still have 3 errors, got %d", handler.Errors().Size())
	}
}
