package internal

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type LiteralCompiler struct {
	ctx *CompilerContext
}

func NewLiteralCompiler(ctx *CompilerContext) *LiteralCompiler {
	return &LiteralCompiler{
		ctx: ctx,
	}
}

func (lc *LiteralCompiler) Compile(ctx fql.ILiteralContext) vm.Operand {
	if c := ctx.StringLiteral(); c != nil {
		return lc.CompileStringLiteral(c)
	} else if c := ctx.IntegerLiteral(); c != nil {
		return lc.CompileIntegerLiteral(c)
	} else if c := ctx.FloatLiteral(); c != nil {
		return lc.CompileFloatLiteral(c)
	} else if c := ctx.BooleanLiteral(); c != nil {
		return lc.CompileBooleanLiteral(c)
	} else if c := ctx.ArrayLiteral(); c != nil {
		return lc.CompileArrayLiteral(c)
	} else if c := ctx.ObjectLiteral(); c != nil {
		return lc.CompileObjectLiteral(c)
	} else if c := ctx.NoneLiteral(); c != nil {
		return lc.CompileNoneLiteral(c)
	}

	panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
}

func (lc *LiteralCompiler) CompileStringLiteral(ctx fql.IStringLiteralContext) vm.Operand {
	var b strings.Builder

	for _, child := range ctx.GetChildren() {
		tree := child.(antlr.TerminalNode)
		sym := tree.GetSymbol()
		input := sym.GetInputStream()

		if input == nil {
			continue
		}

		size := input.Size()
		// skip quotes
		start := sym.GetStart() + 1
		stop := sym.GetStop() - 1

		if stop >= size {
			stop = size - 1
		}

		if start < size && stop < size {
			for i := start; i <= stop; i++ {
				c := input.GetText(i, i)

				switch c {
				case "\\":
					c2 := input.GetText(i, i+1)

					switch c2 {
					case "\\n":
						b.WriteString("\n")
					case "\\t":
						b.WriteString("\t")
					default:
						b.WriteString(c2)
					}

					i++
				default:
					b.WriteString(c)
				}
			}
		}
	}

	return loadConstant(lc.ctx, runtime.NewString(b.String()))
}

func (lc *LiteralCompiler) CompileIntegerLiteral(ctx fql.IIntegerLiteralContext) vm.Operand {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	return loadConstant(lc.ctx, runtime.NewInt(val))
}

func (lc *LiteralCompiler) CompileFloatLiteral(ctx fql.IFloatLiteralContext) vm.Operand {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	return loadConstant(lc.ctx, runtime.NewFloat(val))
}

func (lc *LiteralCompiler) CompileBooleanLiteral(ctx fql.IBooleanLiteralContext) vm.Operand {
	reg := lc.ctx.Registers.Allocate(core.Temp)

	switch strings.ToLower(ctx.GetText()) {
	case "true":
		lc.ctx.Emitter.EmitBoolean(reg, true)
	case "false":
		lc.ctx.Emitter.EmitBoolean(reg, false)
	default:
		panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
	}

	return reg
}

func (lc *LiteralCompiler) CompileNoneLiteral(_ fql.INoneLiteralContext) vm.Operand {
	reg := lc.ctx.Registers.Allocate(core.Temp)
	lc.ctx.Emitter.EmitA(vm.OpLoadNone, reg)

	return reg
}

func (lc *LiteralCompiler) CompileArrayLiteral(ctx fql.IArrayLiteralContext) vm.Operand {
	// Allocate destination register for the array
	destReg := lc.ctx.Registers.Allocate(core.Temp)

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(fql.IArgumentListContext).AllExpression()
		size := len(exps)

		if size > 0 {
			// Allocate seq for array elements
			seq := lc.ctx.Registers.AllocateSequence(size)

			// Evaluate each element into seq Registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := lc.ctx.ExprCompiler.Compile(exp)

				// TODO: Figure out how to remove OpMove and use Registers returned from each expression
				lc.ctx.Emitter.EmitMove(seq[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//lc.ctx.Registers.Free(srcReg)
				}
			}

			// Initialize an array
			lc.ctx.Emitter.EmitList(destReg, seq)

			// Free seq Registers
			//lc.ctx.Registers.FreeSequence(seq)

			return destReg
		}
	}

	// Empty array
	lc.ctx.Emitter.EmitEmptyList(destReg)

	return destReg
}

func (lc *LiteralCompiler) CompileObjectLiteral(ctx fql.IObjectLiteralContext) vm.Operand {
	dst := lc.ctx.Registers.Allocate(core.Temp)
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size == 0 {
		lc.ctx.Emitter.EmitEmptyMap(dst)

		return dst
	}

	seq := lc.ctx.Registers.AllocateSequence(len(assignments) * 2)

	for i := 0; i < size; i++ {
		var propOp vm.Operand
		var valOp vm.Operand
		pac := assignments[i]

		if prop := pac.PropertyName(); prop != nil {
			propOp = lc.CompilePropertyName(prop)
			valOp = lc.ctx.ExprCompiler.Compile(pac.Expression())
		} else if comProp := pac.ComputedPropertyName(); comProp != nil {
			propOp = lc.CompileComputedPropertyName(comProp)
			valOp = lc.ctx.ExprCompiler.Compile(pac.Expression())
		} else if variable := pac.Variable(); variable != nil {
			propOp = loadConstant(lc.ctx, runtime.NewString(variable.GetText()))
			valOp = lc.ctx.ExprCompiler.CompileVariable(variable)
		}

		regIndex := i * 2

		lc.ctx.Emitter.EmitMove(seq[regIndex], propOp)
		lc.ctx.Emitter.EmitMove(seq[regIndex+1], valOp)

		// Free source register if temporary
		if propOp.IsRegister() {
			//lc.ctx.Registers.Free(propOp)
		}
	}

	lc.ctx.Emitter.EmitMap(dst, seq)

	return dst
}

func (lc *LiteralCompiler) CompilePropertyName(ctx fql.IPropertyNameContext) vm.Operand {
	if str := ctx.StringLiteral(); str != nil {
		return lc.CompileStringLiteral(str)
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if word := ctx.SafeReservedWord(); word != nil {
		name = word.GetText()
	} else if word := ctx.UnsafeReservedWord(); word != nil {
		name = word.GetText()
	} else {
		panic(runtime.Error(core.ErrUnexpectedToken, ctx.GetText()))
	}

	return loadConstant(lc.ctx, runtime.NewString(name))
}

func (lc *LiteralCompiler) CompileComputedPropertyName(ctx fql.IComputedPropertyNameContext) vm.Operand {
	return lc.ctx.ExprCompiler.Compile(ctx.Expression())
}
