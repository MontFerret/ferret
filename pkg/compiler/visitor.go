package compiler

import (
	"strconv"
	"strings"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type visitor struct {
	*fql.BaseFqlParserVisitor
	err        error
	src        string
	emitter    *Emitter
	registers  *RegisterAllocator
	symbols    *SymbolTable
	loops      *LoopTable
	catchTable [][2]int
}

const (
	jumpPlaceholder      = -1
	undefinedVariable    = -1
	ignorePseudoVariable = "_"
	pseudoVariable       = "CURRENT"
	waitScope            = "waitfor"
	forScope             = "for"
)

func newVisitor(src string) *visitor {
	v := new(visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.src = src
	v.registers = NewRegisterAllocator()
	v.symbols = NewSymbolTable(v.registers)
	v.loops = NewLoopTable(v.registers)
	v.emitter = NewEmitter()
	v.catchTable = make([][2]int, 0)

	return v
}

func (v *visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	ctx.Body().Accept(v)

	return nil
}

func (v *visitor) VisitBody(ctx *fql.BodyContext) interface{} {
	for _, statement := range ctx.AllBodyStatement() {
		statement.Accept(v)
	}

	ctx.BodyExpression().Accept(v)

	return nil
}

func (v *visitor) VisitBodyStatement(ctx *fql.BodyStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitBodyExpression(ctx *fql.BodyExpressionContext) interface{} {
	if c := ctx.ForExpression(); c != nil {
		out, ok := c.Accept(v).(runtime.Operand)

		if ok && out != runtime.NoopOperand {
			v.emitter.EmitAB(runtime.OpMove, runtime.NoopOperand, out)
		}

		v.emitter.Emit(runtime.OpReturn)

		return out
	} else if c := ctx.ReturnExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func (v *visitor) VisitForExpression(ctx *fql.ForExpressionContext) interface{} {
	v.symbols.EnterScope()

	var passThrough bool
	var distinct bool
	var returnRuleCtx antlr.RuleContext
	// identify whether it's WHILE or FOR loop
	isForLoop := ctx.While() == nil
	returnCtx := ctx.ForExpressionReturn()

	if c := returnCtx.ReturnExpression(); c != nil {
		returnRuleCtx = c
		distinct = c.Distinct() != nil
	} else if c := returnCtx.ForExpression(); c != nil {
		returnRuleCtx = c
		passThrough = true
	}

	loop := v.loops.EnterLoop(passThrough, distinct)
	dsReg := loop.Result

	if loop.Allocate {
		v.emitter.EmitAb(runtime.OpLoopBegin, dsReg, distinct)
	}

	if isForLoop {
		// Loop data source to iterate over
		src1 := ctx.ForExpressionSource().Accept(v).(runtime.Operand)

		iterReg := v.registers.Allocate(State)

		v.emitter.EmitAB(runtime.OpForLoopInit, iterReg, src1)
		// jumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = v.emitter.EmitJumpc(runtime.OpForLoopNext, jumpPlaceholder, iterReg)

		valVar := ctx.GetValueVariable().GetText()
		counterVarCtx := ctx.GetCounterVariable()

		hasValVar := valVar != ignorePseudoVariable
		var hasCounterVar bool
		var counterVar string

		if counterVarCtx != nil {
			counterVar = counterVarCtx.GetText()
			hasCounterVar = true
		}

		var valReg runtime.Operand

		// declare value variable
		if hasValVar {
			valReg = v.symbols.DefineVariable(valVar)
			v.emitter.EmitAB(runtime.OpForLoopValue, valReg, iterReg)
		}

		var keyReg runtime.Operand

		if hasCounterVar {
			// declare counter variable
			keyReg = v.symbols.DefineVariable(counterVar)
			v.emitter.EmitAB(runtime.OpForLoopKey, keyReg, iterReg)
		}
	} else {
		counterReg := v.registers.Allocate(State)

		// Create initial value for the loop counter
		v.emitter.EmitA(runtime.OpWhileLoopInit, counterReg)
		// Loop data source to iterate over
		cond := ctx.Expression().Accept(v).(runtime.Operand)

		// jumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = v.emitter.EmitJumpAB(runtime.OpWhileLoopNext, counterReg, cond, jumpPlaceholder)

		counterVar := ctx.GetCounterVariable().GetText()

		// declare counter variable
		valReg := v.symbols.DefineVariable(counterVar)
		v.emitter.EmitAB(runtime.OpWhileLoopValue, valReg, counterReg)
	}

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			b.Accept(v)
		}
	}

	// RETURN
	if !passThrough {
		c := returnRuleCtx.(*fql.ReturnExpressionContext)
		expReg := c.Expression().Accept(v).(runtime.Operand)

		v.emitter.EmitAB(runtime.OpLoopPush, dsReg, expReg)
	} else if returnRuleCtx != nil {
		returnRuleCtx.Accept(v)
	}

	if isForLoop {
		v.emitter.EmitJump(runtime.OpJump, loop.Jump)
	} else {
		v.emitter.EmitJump(runtime.OpJump, loop.Jump-1)
	}

	// TODO: Do not allocate for pass-through loops
	dst := v.registers.Allocate(Temp)

	if loop.Allocate {
		// TODO: Reuse the dsReg register
		v.emitter.EmitAB(runtime.OpLoopEnd, dst, dsReg)

		if isForLoop {
			v.emitter.PatchJump(loop.Jump)
		} else {
			v.emitter.PatchJumpAB(loop.Jump)
		}
	} else {
		if isForLoop {
			v.emitter.PatchJumpNext(loop.Jump)
		} else {
			v.emitter.PatchJumpNextAB(loop.Jump)
		}
	}

	v.loops.ExitLoop()
	v.symbols.ExitScope()

	return dst
}

func (v *visitor) VisitForExpressionSource(ctx *fql.ForExpressionSourceContext) interface{} {
	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.RangeOperator(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.ArrayLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.ObjectLiteral(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitForExpressionBody(ctx *fql.ForExpressionBodyContext) interface{} {
	if c := ctx.ForExpressionClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.ForExpressionStatement(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitForExpressionClause(ctx *fql.ForExpressionClauseContext) interface{} {
	if c := ctx.LimitClause(); c != nil {
		// TODO: Implement
		return c.Accept(v)
	}

	if c := ctx.FilterClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.SortClause(); c != nil {
		// TODO: Implement
		return c.Accept(v)
	}

	if c := ctx.CollectClause(); c != nil {
		// TODO: Implement
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitFilterClause(ctx *fql.FilterClauseContext) interface{} {
	//ctx.Expression().Accept(v)
	//fwd := v.emitJump(runtime.OpJumpIfTrue)
	//// Pop on false
	//v.emitPop()
	//// And jump back to the beginning of the loop
	//bwd := v.emitJump(runtime.OpJumpBackward)
	//v.patchJumpWith(bwd, len(v.instructions)-v.resolveLoopPosition())
	//
	//// Otherwise, pop on true
	//v.patchJump(fwd)
	//v.emitPop()

	return nil
}

func (v *visitor) VisitForExpressionStatement(ctx *fql.ForExpressionStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) interface{} {
	return v.visitFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext), ctx.ErrorOperator() != nil)
}

func (v *visitor) VisitFunctionCall(ctx *fql.FunctionCallContext) interface{} {
	return v.visitFunctionCall(ctx, false)
}

func (v *visitor) VisitMemberExpression(ctx *fql.MemberExpressionContext) interface{} {
	mes := ctx.MemberExpressionSource().(*fql.MemberExpressionSourceContext)
	segments := ctx.AllMemberExpressionPath()

	var mesOut interface{}

	if c := mes.Variable(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.Param(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.ObjectLiteral(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.ArrayLiteral(); c != nil {
		mesOut = c.Accept(v)
	} else if c := mes.FunctionCall(); c != nil {
		// FOO()?.bar
		segment := segments[0].(*fql.MemberExpressionPathContext)
		mesOut = v.visitFunctionCall(c.(*fql.FunctionCallContext), segment.ErrorOperator() != nil)
	}

	var dst runtime.Operand
	src1 := mesOut.(runtime.Operand)

	for _, segment := range segments {
		var out2 interface{}
		p := segment.(*fql.MemberExpressionPathContext)

		if c := p.PropertyName(); c != nil {
			out2 = c.Accept(v)
		} else if c := p.ComputedPropertyName(); c != nil {
			out2 = c.Accept(v)
		}

		src2 := out2.(runtime.Operand)
		dst = v.registers.Allocate(Temp)

		if p.ErrorOperator() != nil {
			v.emitter.EmitABC(runtime.OpLoadPropertyOptional, dst, src1, src2)
		} else {
			v.emitter.EmitABC(runtime.OpLoadProperty, dst, src1, src2)
		}

		src1 = dst
	}

	return dst
}

func (v *visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	dst := v.registers.Allocate(Temp)
	start := ctx.GetLeft().Accept(v).(runtime.Operand)
	end := ctx.GetRight().Accept(v).(runtime.Operand)

	v.emitter.EmitABC(runtime.OpRange, dst, start, end)

	return dst
}

func (v *visitor) VisitRangeOperand(ctx *fql.RangeOperandContext) interface{} {
	if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	src := ctx.Expression().Accept(v).(runtime.Operand)

	if name != ignorePseudoVariable {
		dest := v.symbols.DefineVariable(name)

		if src.IsConstant() {
			tmp := v.registers.Allocate(Temp)
			v.emitter.EmitAB(runtime.OpLoadConst, tmp, src)
			v.emitter.EmitAB(runtime.OpStoreGlobal, dest, tmp)
		} else if v.symbols.Scope() == 0 {
			v.emitter.EmitAB(runtime.OpStoreGlobal, dest, src)
		} else {
			v.emitter.EmitAB(runtime.OpMove, dest, src)
		}

		return dest
	}

	return runtime.NoopOperand
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	// Just return the register / constant index
	op := v.symbols.LookupVariable(ctx.GetText())

	if op.IsRegister() {
		return op
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadGlobal, reg, op)

	return reg
}

func (v *visitor) VisitArrayLiteral(ctx *fql.ArrayLiteralContext) interface{} {
	// Allocate destination register for the array
	destReg := v.registers.Allocate(Temp)

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size := len(exps)

		if size > 0 {
			// Allocate seq for array elements
			seq := v.registers.AllocateSequence(size)

			// Evaluate each element into seq registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := exp.Accept(v).(runtime.Operand)

				// TODO: Figure out how to remove OpMove and use registers returned from each expression
				v.emitter.EmitAB(runtime.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.registers.Free(srcReg)
				}
			}

			// Initialize an array
			v.emitter.EmitAs(runtime.OpArray, destReg, seq)

			// Free seq registers
			//v.registers.FreeSequence(seq)

			return destReg
		}
	}

	// Empty array
	v.emitter.EmitA(runtime.OpArray, destReg)

	return destReg
}

func (v *visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	dst := v.registers.Allocate(Temp)
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size == 0 {
		v.emitter.EmitA(runtime.OpObject, dst)

		return dst
	}

	seq := v.registers.AllocateSequence(len(assignments) * 2)

	for i := 0; i < size; i++ {
		var propOp runtime.Operand
		var valOp runtime.Operand
		pac := assignments[i].(*fql.PropertyAssignmentContext)

		if prop, ok := pac.PropertyName().(*fql.PropertyNameContext); ok {
			propOp = prop.Accept(v).(runtime.Operand)
			valOp = pac.Expression().Accept(v).(runtime.Operand)
		} else if comProp, ok := pac.ComputedPropertyName().(*fql.ComputedPropertyNameContext); ok {
			propOp = comProp.Accept(v).(runtime.Operand)
			valOp = pac.Expression().Accept(v).(runtime.Operand)
		} else if variable := pac.Variable(); variable != nil {
			propOp = v.loadConstant(values.NewString(variable.GetText()))
			valOp = variable.Accept(v).(runtime.Operand)
		}

		regIndex := i * 2

		v.emitter.EmitAB(runtime.OpMove, seq.Registers[regIndex], propOp)
		v.emitter.EmitAB(runtime.OpMove, seq.Registers[regIndex+1], valOp)

		// Free source register if temporary
		if propOp.IsRegister() {
			//v.registers.Free(propOp)
		}
	}

	v.emitter.EmitAs(runtime.OpObject, dst, seq)

	return dst
}

func (v *visitor) VisitPropertyName(ctx *fql.PropertyNameContext) interface{} {
	if str := ctx.StringLiteral(); str != nil {
		return str.Accept(v)
	}

	var name string

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if word := ctx.SafeReservedWord(); word != nil {
		name = word.GetText()
	} else if word := ctx.UnsafeReservedWord(); word != nil {
		name = word.GetText()
	} else {
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return v.loadConstant(values.NewString(name))
}

func (v *visitor) VisitComputedPropertyName(ctx *fql.ComputedPropertyNameContext) interface{} {
	return ctx.Expression().Accept(v)
}

func (v *visitor) VisitStringLiteral(ctx *fql.StringLiteralContext) interface{} {
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

	return v.loadConstant(values.NewString(b.String()))
}

func (v *visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadConst, reg, v.symbols.AddConstant(values.NewInt(val)))

	return reg
}

func (v *visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadConst, reg, v.symbols.AddConstant(values.NewFloat(val)))

	return reg
}

func (v *visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	reg := v.registers.Allocate(Temp)

	switch strings.ToLower(ctx.GetText()) {
	case "true":
		v.emitter.EmitAB(runtime.OpLoadBool, reg, 1)
	case "false":
		v.emitter.EmitAB(runtime.OpLoadBool, reg, 0)
	default:
		panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return reg
}

func (v *visitor) VisitNoneLiteral(_ *fql.NoneLiteralContext) interface{} {
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitA(runtime.OpLoadNone, reg)

	return reg
}

func (v *visitor) VisitLiteral(ctx *fql.LiteralContext) interface{} {
	if c := ctx.ArrayLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.ObjectLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.StringLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.FloatLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.BooleanLiteral(); c != nil {
		return c.Accept(v)
	} else if c := ctx.NoneLiteral(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitReturnExpression(ctx *fql.ReturnExpressionContext) interface{} {
	valReg := ctx.Expression().Accept(v).(runtime.Operand)

	if valReg.IsConstant() {
		v.emitter.EmitAB(runtime.OpLoadGlobal, runtime.NoopOperand, valReg)
	} else {
		v.emitter.EmitAB(runtime.OpMove, runtime.NoopOperand, valReg)
	}

	v.emitter.Emit(runtime.OpReturn)

	return runtime.NoopOperand
}

func (v *visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if uo := ctx.UnaryOperator(); uo != nil {
		src := ctx.GetRight().Accept(v).(runtime.Operand)
		dst := v.registers.Allocate(Temp)

		uoc := uo.(*fql.UnaryOperatorContext)
		var op runtime.Opcode

		if uoc.Not() != nil {
			op = runtime.OpNot
		} else if uoc.Minus() != nil {
			op = runtime.OpFlipNegative
		} else if uoc.Plus() != nil {
			op = runtime.OpFlipPositive
		} else {
			panic(core.Error(ErrUnexpectedToken, uoc.GetText()))
		}

		// We do not overwrite the source register
		v.emitter.EmitAB(op, dst, src)

		return dst
	}

	if op := ctx.LogicalAndOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(runtime.Operand)
		v.emitter.EmitAB(runtime.OpMove, dst, left)
		// Test if left is false and jump to the end
		end := v.emitter.EmitJumpc(runtime.OpJumpIfFalse, jumpPlaceholder, dst)
		// If left is true, execute right expression
		right := ctx.GetRight().Accept(v).(runtime.Operand)
		// And move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, right)
		v.emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.LogicalOrOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(runtime.Operand)
		// Move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, left)
		// Test if left is true and jump to the end
		end := v.emitter.EmitJumpc(runtime.OpJumpIfTrue, jumpPlaceholder, dst)
		// If left is false, execute right expression
		right := ctx.GetRight().Accept(v).(runtime.Operand)
		// And move the result to the destination register
		v.emitter.EmitAB(runtime.OpMove, dst, right)
		v.emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.GetTernaryOperator(); op != nil {
		dst := v.registers.Allocate(Temp)

		// Compile condition and put result in dst
		condReg := ctx.GetCondition().Accept(v).(runtime.Operand)
		v.emitter.EmitAB(runtime.OpMove, dst, condReg)

		// If condition was temporary, free it
		if condReg.IsRegister() {
			//v.registers.Free(condReg)
		}

		// Jump to 'false' branch if condition is false
		otherwise := v.emitter.EmitJumpc(runtime.OpJumpIfFalse, jumpPlaceholder, dst)

		// True branch
		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			trueReg := onTrue.Accept(v).(runtime.Operand)
			v.emitter.EmitAB(runtime.OpMove, dst, trueReg)

			// Free temporary register if needed
			if trueReg.IsRegister() {
				//v.registers.Free(trueReg)
			}
		}

		// Jump over false branch
		end := v.emitter.EmitJump(runtime.OpJump, jumpPlaceholder)
		v.emitter.PatchJumpNext(otherwise)

		// False branch
		if onFalse := ctx.GetOnFalse(); onFalse != nil {
			falseReg := onFalse.Accept(v).(runtime.Operand)
			v.emitter.EmitAB(runtime.OpMove, dst, falseReg)

			// Free temporary register if needed
			if falseReg.IsRegister() {
				//v.registers.Free(falseReg)
			}
		}

		v.emitter.PatchJumpNext(end)

		return dst
	}

	if c := ctx.Predicate(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitPredicate(ctx *fql.PredicateContext) interface{} {
	if c := ctx.ExpressionAtom(); c != nil {
		startCatch := v.emitter.Size()
		reg := c.Accept(v)

		if c.ErrorOperator() != nil {
			endCatch := v.emitter.Size()
			v.catchTable = append(v.catchTable, [2]int{startCatch, endCatch})
		}

		return reg
	}

	var opcode runtime.Opcode
	dest := v.registers.Allocate(Temp)
	left := ctx.Predicate(0).Accept(v).(runtime.Operand)
	right := ctx.Predicate(1).Accept(v).(runtime.Operand)

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			opcode = runtime.OpEq
		case "!=":
			opcode = runtime.OpNeq
		case ">":
			opcode = runtime.OpGt
		case ">=":
			opcode = runtime.OpGte
		case "<":
			opcode = runtime.OpLt
		case "<=":
			opcode = runtime.OpLte
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
		panic(core.Error(core.ErrNotImplemented, "array operator"))
	} else if op := ctx.InOperator(); op != nil {
		opcode = runtime.OpIn
	} else if op := ctx.LikeOperator(); op != nil {
		if op.(*fql.LikeOperatorContext).Not() != nil {
			opcode = runtime.OpNotLike
		} else {
			opcode = runtime.OpLike
		}
	}

	v.emitter.EmitABC(opcode, dest, left, right)

	return dest
}

func (v *visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
	var opcode runtime.Opcode
	var isSet bool

	if op := ctx.MultiplicativeOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "*":
			opcode = runtime.OpMulti
		case "/":
			opcode = runtime.OpDiv
		case "%":
			opcode = runtime.OpMod
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = runtime.OpAdd
		case "-":
			opcode = runtime.OpSub
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "=~":
			opcode = runtime.OpRegexpPositive
		case "!~":
			opcode = runtime.OpRegexpNegative
		default:
			panic(core.Error(ErrUnexpectedToken, op.GetText()))
		}
	}

	if isSet {
		regLeft := ctx.ExpressionAtom(0).Accept(v).(runtime.Operand)
		regRight := ctx.ExpressionAtom(1).Accept(v).(runtime.Operand)
		dst := v.registers.Allocate(Temp)

		v.emitter.EmitABC(opcode, dst, regLeft, regRight)

		return dst
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.RangeOperator(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Literal(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	} else if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Param(); c != nil {
		return c.Accept(v)
	} else if c := ctx.ForExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.Expression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) visitFunctionCall(ctx *fql.FunctionCallContext, safeCall bool) interface{} {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	var size int
	var seq *RegisterSequence

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size = len(exps)

		if size > 0 {
			// Allocate seq for function arguments
			seq = v.registers.AllocateSequence(size)

			// Evaluate each element into seq registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := exp.Accept(v).(runtime.Operand)

				// TODO: Figure out how to remove OpMove and use registers returned from each expression
				v.emitter.EmitAB(runtime.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.registers.Free(srcReg)
				}
			}
		}
	}

	nameAndDest := v.loadConstant(values.NewString(strings.ToUpper(name)))

	if !safeCall {
		v.emitter.EmitAs(runtime.OpCall, nameAndDest, seq)
	} else {
		v.emitter.EmitAs(runtime.OpCallSafe, nameAndDest, seq)
	}

	return nameAndDest
}

func (v *visitor) loadConstant(constant core.Value) runtime.Operand {
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadConst, reg, v.symbols.AddConstant(constant))
	return reg
}
