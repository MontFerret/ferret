package compile

import "github.com/MontFerret/ferret/v2/pkg/bytecode"

func CastToProgram(prog any) *bytecode.Program {
	if p, ok := prog.(*bytecode.Program); ok {
		return p
	}

	panic("expected *vm.Program")
}
