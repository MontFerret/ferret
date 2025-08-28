package diagnostics

import (
	"testing"

	"github.com/MontFerret/ferret/pkg/file"
	. "github.com/smartystreets/goconvey/convey"
)

func TestErrorConstructors(t *testing.T) {
	Convey("Error constructors", t, func() {
		Convey("NewEmptyQueryErr should create syntax error for empty query", func() {
			src := file.NewSource("test.fql", "")

			err := NewEmptyQueryErr(src)

			So(err, ShouldNotBeNil)
			So(err.Kind, ShouldEqual, SyntaxError)
			So(err.Message, ShouldEqual, "Query is empty")
			So(err.Source, ShouldEqual, src)
		})

		Convey("NewInternalErr should create internal error without cause", func() {
			src := file.NewSource("test.fql", "LET x = 1")
			msg := "internal error message"

			err := NewInternalErr(src, msg)

			So(err, ShouldNotBeNil)
			So(err.Kind, ShouldEqual, InternalError)
			So(err.Message, ShouldEqual, msg)
			So(err.Source, ShouldEqual, src)
			So(err.Cause, ShouldBeNil)
		})

		Convey("NewInternalErrWith should create internal error with cause", func() {
			src := file.NewSource("test.fql", "LET x = 1")
			msg := "internal error with cause"
			cause := &CompilationError{Message: "original error"}

			err := NewInternalErrWith(src, msg, cause)

			So(err, ShouldNotBeNil)
			So(err.Kind, ShouldEqual, InternalError)
			So(err.Message, ShouldEqual, msg)
			So(err.Source, ShouldEqual, src)
			So(err.Cause, ShouldEqual, cause)
		})
	})
}

func TestErrorConstants(t *testing.T) {
	Convey("Error constants should have correct values", t, func() {
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