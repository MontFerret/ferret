package asm_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/MontFerret/ferret/pkg/asm"
)

func TestErrors(t *testing.T) {
	Convey("ErrInvalidProgram", t, func() {
		Convey("Should be defined", func() {
			So(asm.ErrInvalidProgram, ShouldNotBeNil)
			So(asm.ErrInvalidProgram.Error(), ShouldEqual, "invalid program: program cannot be nil or empty")
		})
	})
}