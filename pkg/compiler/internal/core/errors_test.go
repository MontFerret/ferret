package core

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

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
