package vm

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

func validate(_ *Environment, program *bytecode.Program) error {
	if program == nil {
		return runtime.Error(runtime.ErrInvalidOperation, "unsupported bytecode version; recompile query")
	}

	if program.ISAVersion != bytecode.Version {
		return runtime.Error(runtime.ErrInvalidOperation, "unsupported bytecode version; recompile query")
	}

	return nil
}
