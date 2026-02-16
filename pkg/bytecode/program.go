package bytecode

import (
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	Catch [3]int

	Program struct {
		Source     *file.Source
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		DebugSpans []file.Span
		Functions  map[string]int
		Params     []string
		Registers  int
		Labels     map[int]string
	}
)
