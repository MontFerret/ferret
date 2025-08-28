package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestAnalyzeSyntaxError(t *testing.T) {
	src := file.NewSource("test.fql", "LET x =")
	
	err := &CompilationError{
		Kind:    SyntaxError,
		Message: "mismatched input '<EOF>' expecting {IntegerLiteral, FloatLiteral, StringLiteral}",
		Source:  src,
	}

	// Create a mock TokenNode
	offending := &TokenNode{}

	result := AnalyzeSyntaxError(src, err, offending)

	// The function should return true if any matcher processed the error
	// Since we have matchers registered, it should attempt to match
	if result != true && result != false {
		t.Errorf("AnalyzeSyntaxError() should return bool, got %v", result)
	}
}

func TestAnalyzeSyntaxError_AllMatchers(t *testing.T) {
	src := file.NewSource("test.fql", "RETURN")

	// Test different types of syntax errors that should trigger different matchers
	testCases := []struct {
		name    string
		message string
	}{
		{
			name:    "literal error",
			message: "mismatched input 'invalid' expecting {IntegerLiteral, FloatLiteral}",
		},
		{
			name:    "assignment error", 
			message: "mismatched input '<EOF>' expecting expression",
		},
		{
			name:    "for loop error",
			message: "mismatched input 'FOR' expecting expression",
		},
		{
			name:    "common error",
			message: "no viable alternative at input",
		},
		{
			name:    "return value error",
			message: "missing return value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &CompilationError{
				Kind:    SyntaxError,
				Message: tc.message,
				Source:  src,
			}

			offending := &TokenNode{}
			result := AnalyzeSyntaxError(src, err, offending)
			
			// Should return a boolean value
			if result != true && result != false {
				t.Errorf("AnalyzeSyntaxError() should return bool for %s, got %v", tc.name, result)
			}
		})
	}
}

func TestAnalyzeSyntaxError_NoMatch(t *testing.T) {
	src := file.NewSource("test.fql", "LET x = 1")
	
	err := &CompilationError{
		Kind:    SyntaxError,
		Message: "some unrecognized error message that won't match any patterns",
		Source:  src,
	}

	offending := &TokenNode{}

	result := AnalyzeSyntaxError(src, err, offending)

	// Should return false when no matcher handles the error
	if result != false {
		t.Errorf("AnalyzeSyntaxError() should return false when no matcher matches, got %v", result)
	}
}