package diagnostics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestCompilationError(t *testing.T) {
	Convey("CompilationError", t, func() {
		Convey("Diagnostic() should return message", func() {
			err := &diagnostics.Diagnostic{
				Kind:    SyntaxError,
				Message: "test error message",
				Hint:    "test hint",
			}

			So(err.Error(), ShouldEqual, "test error message")
		})

		Convey("Format() should format error with all components", func() {
			src := file.NewSource("test.fql", "LET x = 1")

			err := &diagnostics.Diagnostic{
				Kind:    SyntaxError,
				Message: "test error message",
				Hint:    "test hint",
				Source:  src,
				Spans: []diagnostics.ErrorSpan{
					diagnostics.NewMainErrorSpan(file.Span{Start: 0, End: 5}, "test label"),
				},
			}

			formatted := err.Format()
			So(formatted, ShouldNotBeEmpty)

			// Should contain the error kind and message
			So(formatted, ShouldContainSubstring, "SyntaxError")
			So(formatted, ShouldContainSubstring, "test error message")
			So(formatted, ShouldContainSubstring, "test hint")
		})
	})
}

func TestErrorKind(t *testing.T) {
	Convey("ErrorKind constants", t, func() {
		tests := []struct {
			name string
			kind diagnostics.Kind
			want string
		}{
			{"UnknownError", diagnostics.Unknown, ""},
			{"SyntaxError", SyntaxError, "SyntaxError"},
			{"NameError", NameError, "NameError"},
			{"TypeError", diagnostics.TypeError, "TypeError"},
			{"SemanticError", SemanticError, "SemanticError"},
			{"Unsupported", diagnostics.Unsupported, "Unsupported"},
			{"UnexpectedError", diagnostics.UnexpectedError, "UnexpectedError"},
		}

		for _, tt := range tests {
			Convey("Should have correct string value for "+tt.name, func() {
				So(string(tt.kind), ShouldEqual, tt.want)
			})
		}
	})
}
