package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestMultiCompilationError_Error(t *testing.T) {
	tests := []struct {
		name   string
		errors []*CompilationError
		want   string
	}{
		{
			name:   "no errors",
			errors: []*CompilationError{},
			want:   "No errors",
		},
		{
			name: "one error",
			errors: []*CompilationError{
				{Message: "test error"},
			},
			want: "Found 1 errors",
		},
		{
			name: "multiple errors",
			errors: []*CompilationError{
				{Message: "error 1"},
				{Message: "error 2"},
			},
			want: "Found 2 errors",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MultiCompilationError{Errors: tt.errors}
			if got := e.Error(); got != tt.want {
				t.Errorf("MultiCompilationError.Error() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMultiCompilationError_Format(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")

	tests := []struct {
		name   string
		errors []*CompilationError
		want   string
	}{
		{
			name:   "no errors",
			errors: []*CompilationError{},
			want:   "No errors",
		},
		{
			name: "single error",
			errors: []*CompilationError{
				{
					Kind:    SyntaxError,
					Message: "test error",
					Source:  src,
				},
			},
		},
		{
			name: "multiple errors",
			errors: []*CompilationError{
				{
					Kind:    SyntaxError,
					Message: "error 1",
					Source:  src,
				},
				{
					Kind:    NameError,
					Message: "error 2",
					Source:  src,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &MultiCompilationError{Errors: tt.errors}
			formatted := e.Format()
			
			if tt.name == "no errors" {
				if formatted != tt.want {
					t.Errorf("MultiCompilationError.Format() = %v, want %v", formatted, tt.want)
				}
			} else {
				// For non-empty error cases, just check it's not empty
				if formatted == "" {
					t.Error("MultiCompilationError.Format() returned empty string for non-empty errors")
				}
			}
		})
	}
}

func TestNewMultiCompilationError(t *testing.T) {
	errors := []*CompilationError{
		{Message: "test error 1"},
		{Message: "test error 2"},
	}

	result := NewMultiCompilationError(errors)
	
	if result == nil {
		t.Fatal("NewMultiCompilationError() returned nil")
	}

	multi, ok := result.(*MultiCompilationError)
	if !ok {
		t.Fatal("NewMultiCompilationError() did not return *MultiCompilationError")
	}

	if len(multi.Errors) != 2 {
		t.Errorf("NewMultiCompilationError() errors length = %v, want %v", len(multi.Errors), 2)
	}

	if multi.Errors[0].Message != "test error 1" {
		t.Errorf("NewMultiCompilationError() first error = %v, want %v", multi.Errors[0].Message, "test error 1")
	}
}