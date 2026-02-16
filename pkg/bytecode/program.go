package bytecode

import (
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	Catch [3]int

	Metadata struct {
		AggregatePlans []AggregatePlan
		DebugSpans     []file.Span
		Functions      map[string]int
		Labels         map[int]string
	}

	Program struct {
		Source     *file.Source
		Registers  int
		Bytecode   []Instruction
		Constants  []runtime.Value
		CatchTable []Catch
		Params     []string
		Metadata   Metadata
	}
)
