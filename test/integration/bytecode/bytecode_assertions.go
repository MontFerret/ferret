package bytecode_test

import (
	"github.com/MontFerret/ferret/pkg/vm"
	"github.com/smartystreets/goconvey/convey"
)

func CastToProgram(prog any) *vm.Program {
	if p, ok := prog.(*vm.Program); ok {
		return p
	}

	panic("expected *vm.Program")
}

func ShouldEqualBytecode(e any, a ...any) string {
	expected := CastToProgram(e).Bytecode
	actual := CastToProgram(a[0]).Bytecode

	for i := 0; i < len(expected); i++ {
		if err := convey.ShouldEqual(actual[i].String(), expected[i].String()); err != "" {
			return err
		}
	}

	return ""
}
