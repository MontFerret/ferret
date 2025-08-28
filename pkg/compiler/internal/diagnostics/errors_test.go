package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestNewEmptyQueryErr(t *testing.T) {
	src := file.NewSource("test.fql", "")

	err := NewEmptyQueryErr(src)

	if err == nil {
		t.Fatal("NewEmptyQueryErr() returned nil")
	}

	if err.Kind != SyntaxError {
		t.Errorf("NewEmptyQueryErr() Kind = %v, want %v", err.Kind, SyntaxError)
	}

	if err.Message != "Query is empty" {
		t.Errorf("NewEmptyQueryErr() Message = %v, want %v", err.Message, "Query is empty")
	}

	if err.Source != src {
		t.Errorf("NewEmptyQueryErr() Source = %v, want %v", err.Source, src)
	}
}

func TestNewInternalErr(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	msg := "internal error message"

	err := NewInternalErr(src, msg)

	if err == nil {
		t.Fatal("NewInternalErr() returned nil")
	}

	if err.Kind != InternalError {
		t.Errorf("NewInternalErr() Kind = %v, want %v", err.Kind, InternalError)
	}

	if err.Message != msg {
		t.Errorf("NewInternalErr() Message = %v, want %v", err.Message, msg)
	}

	if err.Source != src {
		t.Errorf("NewInternalErr() Source = %v, want %v", err.Source, src)
	}

	if err.Cause != nil {
		t.Errorf("NewInternalErr() Cause = %v, want nil", err.Cause)
	}
}

func TestNewInternalErrWith(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	msg := "internal error with cause"
	cause := &CompilationError{Message: "original error"}

	err := NewInternalErrWith(src, msg, cause)

	if err == nil {
		t.Fatal("NewInternalErrWith() returned nil")
	}

	if err.Kind != InternalError {
		t.Errorf("NewInternalErrWith() Kind = %v, want %v", err.Kind, InternalError)
	}

	if err.Message != msg {
		t.Errorf("NewInternalErrWith() Message = %v, want %v", err.Message, msg)
	}

	if err.Source != src {
		t.Errorf("NewInternalErrWith() Source = %v, want %v", err.Source, src)
	}

	if err.Cause != cause {
		t.Errorf("NewInternalErrWith() Cause = %v, want %v", err.Cause, cause)
	}
}

func TestErrorConstants(t *testing.T) {
	tests := []struct {
		name     string
		constant string
		expected string
	}{
		{"ErrNotImplemented", ErrNotImplemented, "not implemented"},
		{"ErrInvalidToken", ErrInvalidToken, "invalid token"},
		{"ErrConstantNotFound", ErrConstantNotFound, "constant not found"},
		{"ErrInvalidDataSource", ErrInvalidDataSource, "invalid data source"},
		{"ErrUnknownOpcode", ErrUnknownOpcode, "unknown opcode"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.constant != tt.expected {
				t.Errorf("%s = %v, want %v", tt.name, tt.constant, tt.expected)
			}
		})
	}
}