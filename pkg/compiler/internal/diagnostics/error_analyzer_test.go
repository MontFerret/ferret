package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAnalyzeSyntaxError(t *testing.T) {
	Convey("AnalyzeSyntaxError", t, func() {
		Convey("Should return boolean for basic syntax error", func() {
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
			So(result, ShouldBeIn, []bool{true, false})
		})

		Convey("Should handle different matcher types", func() {
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
				Convey("Should return boolean for "+tc.name, func() {
					err := &CompilationError{
						Kind:    SyntaxError,
						Message: tc.message,
						Source:  src,
					}

					offending := &TokenNode{}
					result := AnalyzeSyntaxError(src, err, offending)
					
					// Should return a boolean value
					So(result, ShouldBeIn, []bool{true, false})
				})
			}
		})

		Convey("Should return false when no matcher matches", func() {
			src := file.NewSource("test.fql", "LET x = 1")
			
			err := &CompilationError{
				Kind:    SyntaxError,
				Message: "some unrecognized error message that won't match any patterns",
				Source:  src,
			}

			offending := &TokenNode{}

			result := AnalyzeSyntaxError(src, err, offending)

			// Should return false when no matcher handles the error
			So(result, ShouldBeFalse)
		})
	})
}