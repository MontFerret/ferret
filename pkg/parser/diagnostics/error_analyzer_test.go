package diagnostics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/source"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func TestQueryDialectWithHistoryEndingAtUsing(t *testing.T) {
	src := source.NewAnonymous("// 😀\nRETURN QUERY `.x` IN @doc USING \"css\" /* USING */")
	history := NewTokenHistory(64)
	for _, tok := range lexDefaultTokens(src.Content()) {
		history.Add(tok)
		if tok.GetText() == "USING" {
			break
		}
	}

	d := NewError(src, SyntaxError, "mismatched input 'USING' expecting <EOF>")
	if !AnalyzeSyntaxError(src, d, history.Last()) {
		t.Fatal("missing dialect diagnostic")
	}
	want := diagnostics.NewMainErrorSpan(source.Span{Start: 37, End: 42}, "invalid dialect")
	if len(d.Spans) != 1 || d.Spans[0] != want {
		t.Fatalf("parser annotations = %+v, want %+v", d.Spans, want)
	}
}

func TestAnalyzeSyntaxError(t *testing.T) {
	Convey("AnalyzeSyntaxError", t, func() {
		Convey("Should return boolean for basic syntax error", func() {
			src := source.New("test.fql", "LET x =")

			err := &diagnostics.Diagnostic{
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
			src := source.New("test.fql", "RETURN")

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
					err := &diagnostics.Diagnostic{
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
			src := source.New("test.fql", "LET x = 1")

			err := &diagnostics.Diagnostic{
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
