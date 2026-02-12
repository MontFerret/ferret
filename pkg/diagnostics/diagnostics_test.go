package diagnostics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/v2/pkg/file"
)

func TestMultiCompilationError(t *testing.T) {
	Convey("MultiCompilationError", t, func() {
		Convey("Diagnostic() should return correct message format", func() {
			tests := []struct {
				name   string
				errors []*Diagnostic
				want   string
			}{
				{
					name:   "no errors",
					errors: []*Diagnostic{},
					want:   "No errors",
				},
				{
					name: "one error",
					errors: []*Diagnostic{
						{Message: "test error"},
					},
					want: "Found 1 errors",
				},
				{
					name: "multiple errors",
					errors: []*Diagnostic{
						{Message: "error 1"},
						{Message: "error 2"},
					},
					want: "Found 2 errors",
				},
			}

			for _, tt := range tests {
				Convey("Should return correct message for "+tt.name, func() {
					e := &Diagnostics[*Diagnostic]{errors: tt.errors}
					So(e.Error(), ShouldEqual, tt.want)
				})
			}
		})

		Convey("Format() should format errors properly", func() {
			src := file.NewSource("test.fql", "LET x = 1")

			tests := []struct {
				name   string
				errors []*Diagnostic
				want   string
			}{
				{
					name:   "no errors",
					errors: []*Diagnostic{},
					want:   "No errors",
				},
				{
					name: "single error",
					errors: []*Diagnostic{
						{
							Kind:    "SyntaxError",
							Message: "test error",
							Source:  src,
						},
					},
				},
				{
					name: "multiple errors",
					errors: []*Diagnostic{
						{
							Kind:    "SyntaxError",
							Message: "error 1",
							Source:  src,
						},
						{
							Kind:    "NameError",
							Message: "error 2",
							Source:  src,
						},
					},
				},
			}

			for _, tt := range tests {
				Convey("Should format correctly for "+tt.name, func() {
					e := NewDiagnosticsOf[*Diagnostic](tt.errors)
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
		errors := []*Diagnostic{
			{Message: "test error 1"},
			{Message: "test error 2"},
		}

		result := NewDiagnosticsOf(errors)

		So(result, ShouldNotBeNil)
		So(len(result.errors), ShouldEqual, 2)
		So(result.errors[0].Message, ShouldEqual, "test error 1")
	})
}
