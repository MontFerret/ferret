package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestCompilationError_Error(t *testing.T) {
	err := &CompilationError{
		Kind:    SyntaxError,
		Message: "test error message",
		Hint:    "test hint",
	}

	if err.Error() != "test error message" {
		t.Errorf("Error() = %v, want %v", err.Error(), "test error message")
	}
}

func TestCompilationError_Format(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	
	err := &CompilationError{
		Kind:    SyntaxError,
		Message: "test error message",
		Hint:    "test hint",
		Source:  src,
		Spans: []ErrorSpan{
			NewMainErrorSpan(file.Span{Start: 0, End: 5}, "test label"),
		},
	}

	formatted := err.Format()
	if formatted == "" {
		t.Error("Format() returned empty string")
	}

	// Should contain the error kind and message
	if !contains(formatted, "SyntaxError") {
		t.Error("Format() should contain error kind")
	}
	
	if !contains(formatted, "test error message") {
		t.Error("Format() should contain error message")
	}

	if !contains(formatted, "test hint") {
		t.Error("Format() should contain hint")
	}
}

func TestErrorKindConstants(t *testing.T) {
	tests := []struct {
		name string
		kind ErrorKind
		want string
	}{
		{"UnknownError", UnknownError, ""},
		{"SyntaxError", SyntaxError, "SyntaxError"},
		{"NameError", NameError, "NameError"},
		{"TypeError", TypeError, "TypeError"},
		{"SemanticError", SemanticError, "SemanticError"},
		{"UnsupportedError", UnsupportedError, "UnsupportedError"},
		{"InternalError", InternalError, "InternalError"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.kind) != tt.want {
				t.Errorf("ErrorKind = %v, want %v", string(tt.kind), tt.want)
			}
		})
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || s[0:len(substr)] == substr || s[len(s)-len(substr):] == substr || findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}