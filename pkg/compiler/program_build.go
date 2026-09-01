package compiler

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func buildProgram(visitor *Visitor, src source.Source, level optimization.Level) (*bytecode.Program, error) {
	var udfs []bytecode.UDF

	if visitor.Session.Program.UDFs != nil {
		udfs = visitor.Session.Program.UDFs.Metadata()
	}

	registers := visitor.Session.Function.Registers.Size()

	for _, udf := range udfs {
		if udf.Registers > registers {
			registers = udf.Registers
		}
	}

	program := &bytecode.Program{
		ISAVersion: bytecode.Version,
		Functions: bytecode.Functions{
			Host:        visitor.Session.Program.HostFunctions.All(),
			UserDefined: udfs,
		},
		Metadata: bytecode.Metadata{
			CompilerVersion:        Version,
			OptimizationLevel:      int(level),
			AggregatePlans:         visitor.Session.Program.AggregatePlans(),
			AggregateSelectorSlots: visitor.Session.Program.Emitter.AggregateSelectorSlots(),
			CallArgumentSpans:      visitor.Session.Program.Emitter.CallArgumentSpans(),
			MatchFailTargets:       visitor.Session.Program.Emitter.MatchFailTargets(),
			DebugSpans:             visitor.Session.Program.Emitter.Spans(),
			DebugPoints:            visitor.Session.Program.DebugPoints,
			Labels:                 visitor.Session.Program.Emitter.Labels(),
		},
		Source:     src,
		Bytecode:   visitor.Session.Program.Emitter.Bytecode(),
		Constants:  visitor.Session.Function.Symbols.Constants(),
		CatchTable: visitor.Session.Program.CatchTable.All(),
		Registers:  registers,
		Params:     visitor.Session.Program.HostParams.Names(),
	}

	if err := optimization.Run(program, level); err != nil {
		return nil, err
	}

	return program, nil
}
