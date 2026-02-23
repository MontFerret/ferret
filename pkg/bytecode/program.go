package bytecode

import (
	"github.com/MontFerret/ferret/v2/pkg/file"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type (
	Catch [3]int

	Metadata struct {
		AggregatePlans []AggregatePlan `json:"aggregatePlans"`
		DebugSpans     []file.Span     `json:"debugSpans"`
		Functions      map[string]int  `json:"functions"`
		Labels         map[int]string  `json:"labels"`
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
