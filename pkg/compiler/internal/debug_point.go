package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

// WithDebugPoint records a logical source stop only when compiling the location
// emits executable bytecode. Bindings are captured before compilation.
func (c *CompilationSession) WithDebugPoint(ctx antlr.ParserRuleContext, compile func()) {
	if compile == nil {
		return
	}

	if c == nil || c.Program == nil || c.Function == nil || !c.Program.DebugInfo || ctx == nil {
		compile()
		return
	}

	pc := c.Program.Emitter.Size()
	vars := c.Function.Symbols.VisibleVariables()
	bindings := make([]bytecode.DebugBinding, 0, len(vars))

	for _, variable := range vars {
		bindings = append(bindings, bytecode.DebugBinding{
			Name:     variable.Name,
			Register: variable.Register,
			Mutable:  variable.Mutable,
			Cell:     variable.Storage == core.BindingStorageCell,
		})
	}

	span := diagnostics.SpanFromRuleContext(ctx)
	pointID := len(c.Program.DebugPoints)
	c.Program.DebugPoints = append(c.Program.DebugPoints, bytecode.DebugPoint{
		PC:         pc,
		FunctionID: c.Function.FunctionID,
		Span:       span,
		Bindings:   bindings,
	})

	c.Program.Emitter.WithSpan(span, func() {
		c.Program.Emitter.EmitA(bytecode.OpSourcePoint, bytecode.Operand(pointID))
	})

	compile()

	if c.Program.Emitter.Size() == pc+1 {
		c.Program.Emitter.Truncate(pc)
		c.Program.DebugPoints = c.Program.DebugPoints[:len(c.Program.DebugPoints)-1]
	}
}
