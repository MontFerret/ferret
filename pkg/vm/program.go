package vm

import (
	"github.com/MontFerret/ferret/pkg/file"
	"github.com/MontFerret/ferret/pkg/runtime"
)

type (
	Catch [3]int

	Program struct {
		Source     *file.Source
		Locations  []file.Location
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		Functions  map[string]int
		Params     []string
		Registers  int
		Labels     map[int]string
	}
)
