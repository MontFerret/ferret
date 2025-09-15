package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
)

func TestDisassemblerOptions(t *testing.T) {
	Convey("WithDebug", t, func() {
		Convey("Should create debug option", func() {
			opt := asm.WithDebug()
			So(opt, ShouldNotBeNil)
		})
	})

	// Note: The internal disassemblerOptions struct and newDisassemblerOptions are not exported,
	// so we can only test them indirectly through the public API (Disassemble function).
	// The actual functionality will be tested in the disassembler tests.
}