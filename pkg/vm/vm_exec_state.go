package vm

import "github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"

type execState struct {
	vm        *VM
	env       *Environment
	registers *mem.RegisterFile
	scratch   *mem.Scratch
	pc        int
	errors    errorHandler
}
