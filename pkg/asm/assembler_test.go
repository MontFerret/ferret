package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
	"github.com/MontFerret/ferret/pkg/vm"
)

func TestAssemble(t *testing.T) {
	Convey("Should return a VM Program", t, func() {
		result, err := asm.Assemble("")

		So(err, ShouldBeNil)
		So(result, ShouldHaveSameTypeAs, &vm.Program{})
		So(result, ShouldNotBeNil)
	})

	Convey("Should handle empty input", t, func() {
		result, err := asm.Assemble("")

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})

	Convey("Should handle non-empty input", t, func() {
		result, err := asm.Assemble("some fasm code")

		So(err, ShouldBeNil)
		So(result, ShouldNotBeNil)
	})
}