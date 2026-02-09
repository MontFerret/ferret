package diagnostics

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/file"
)

func TestErrorConstructors(t *testing.T) {
	Convey("Diagnostic constructors", t, func() {
		Convey("NewEmptyQueryError should create syntax error for empty query", func() {
			src := file.NewSource("test.fql", "")

			err := NewEmptyQueryError(src)

			So(err, ShouldNotBeNil)
			So(err.Kind, ShouldEqual, SyntaxError)
			So(err.Message, ShouldEqual, "Query is empty")
			So(err.Source, ShouldEqual, src)
		})
	})
}

func TestErrorConstants(t *testing.T) {
	Convey("Diagnostic constants should have correct values", t, func() {
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
			Convey("Should have correct value for "+tt.name, func() {
				So(tt.constant, ShouldEqual, tt.expected)
			})
		}
	})
}
