package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMultiCompilationError(t *testing.T) {
	Convey("MultiCompilationError", t, func() {
		Convey("Error() should return correct message format", func() {
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
				Convey("Should return correct message for "+tt.name, func() {
					e := &MultiCompilationError{Errors: tt.errors}
					So(e.Error(), ShouldEqual, tt.want)
				})
			}
		})

		Convey("Format() should format errors properly", func() {
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
				Convey("Should format correctly for "+tt.name, func() {
					e := &MultiCompilationError{Errors: tt.errors}
					formatted := e.Format()
					
					if tt.name == "no errors" {
						So(formatted, ShouldEqual, tt.want)
					} else {
						// For non-empty error cases, just check it's not empty
						So(formatted, ShouldNotBeEmpty)
					}
				})
			}
		})
	})
}

func TestNewMultiCompilationError(t *testing.T) {
	Convey("NewMultiCompilationError should create MultiCompilationError", t, func() {
		errors := []*CompilationError{
			{Message: "test error 1"},
			{Message: "test error 2"},
		}

		result := NewMultiCompilationError(errors)
		
		So(result, ShouldNotBeNil)

		multi, ok := result.(*MultiCompilationError)
		So(ok, ShouldBeTrue)
		So(len(multi.Errors), ShouldEqual, 2)
		So(multi.Errors[0].Message, ShouldEqual, "test error 1")
	})
}