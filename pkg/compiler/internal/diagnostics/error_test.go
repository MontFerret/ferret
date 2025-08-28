package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestCompilationError(t *testing.T) {
	Convey("CompilationError", t, func() {
		Convey("Error() should return message", func() {
			err := &CompilationError{
				Kind:    SyntaxError,
				Message: "test error message",
				Hint:    "test hint",
			}

			So(err.Error(), ShouldEqual, "test error message")
		})

		Convey("Format() should format error with all components", func() {
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
			Convey("Should have correct string value for "+tt.name, func() {
				So(string(tt.kind), ShouldEqual, tt.want)
			})
		}
	})
}

