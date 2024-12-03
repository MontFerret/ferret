package compiler

import (
	"regexp"
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
	catchTable []runtime.Catch
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
	v.catchTable = make([]runtime.Catch, 0)

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

	var distinct bool
	var returnRuleCtx antlr.RuleContext
	returnCtx := ctx.ForExpressionReturn()

	if c := returnCtx.ReturnExpression(); c != nil {
		returnRuleCtx = c
		distinct = c.Distinct() != nil
	} else if c := returnCtx.ForExpression(); c != nil {
		returnRuleCtx = c
	}

	loop := v.loops.EnterLoop(v.loopType(ctx), v.loopKind(ctx), distinct)

	if loop.Kind == ForLoop {
		loop.Src = ctx.ForExpressionSource().Accept(v).(runtime.Operand)

		if val := ctx.GetValueVariable(); val != nil {
			if txt := val.GetText(); txt != "" && txt != ignorePseudoVariable {
				loop.ValueName = txt
				loop.Value = v.symbols.DefineVariable(txt)
			}
		}

		if ctr := ctx.GetCounterVariable(); ctr != nil {
			if txt := ctr.GetText(); txt != "" && txt != ignorePseudoVariable {
				loop.KeyName = txt
				loop.Key = v.symbols.DefineVariable(txt)
			}
		}
	} else {
		//srcExpr := ctx.Expression()
		//
		//// Create initial value for the loop counter
		//v.emitter.EmitA(runtime.OpWhileLoopPrep, counterReg)
		//beforeExp := v.emitter.Size()
		//// Loop data source to iterate over
		//cond := srcExpr.Accept(v).(runtime.Operand)
		//jumpOffset = v.emitter.Size() - beforeExp
		//
		//// jumpPlaceholder is a placeholder for the exit jump position
		//loop.Jump = v.emitter.EmitJumpAB(runtime.OpWhileLoopNext, counterReg, cond, jumpPlaceholder)
		//
		//counterVar := ctx.GetCounterVariable().GetText()
		//
		//// declare counter variable
		//valReg := v.symbols.DefineVariable(counterVar)
		//v.emitter.EmitAB(runtime.OpWhileLoopValue, valReg, counterReg)
	}

	v.emitLoopBegin(loop)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			b.Accept(v)
		}
	}

	loop = v.loops.Loop()

	// RETURN
	if loop.Type != PassThroughLoop {
		c := returnRuleCtx.(*fql.ReturnExpressionContext)
		expReg := c.Expression().Accept(v).(runtime.Operand)

		v.emitter.EmitAB(runtime.OpLoopPush, loop.Result, expReg)
	} else if returnRuleCtx != nil {
		returnRuleCtx.Accept(v)
	}

	res := v.emitLoopEnd(loop)

	v.loops.ExitLoop()
	v.symbols.ExitScope()

	return res
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
	src1 := ctx.Expression().Accept(v).(runtime.Operand)
	v.emitter.EmitJumpc(runtime.OpJumpIfFalse, v.loops.Loop().Jump, src1)

	return nil
}

func (v *visitor) VisitLimitClause(ctx *fql.LimitClauseContext) interface{} {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		return v.visitLimit(clauses[0].Accept(v).(runtime.Operand))
	} else {
		v.visitOffset(clauses[0].Accept(v).(runtime.Operand))
		v.visitLimit(clauses[1].Accept(v).(runtime.Operand))
	}

	return nil
}

func (v *visitor) VisitSortClause(ctx *fql.SortClauseContext) interface{} {
	loop := v.loops.Loop()
	v.emitter.EmitAB(runtime.OpLoopCopy, loop.Result, loop.Iterator)
	v.emitter.EmitJump(runtime.OpJump, loop.Jump-loop.JumpOffset)
	v.emitter.EmitA(runtime.OpClose, loop.Iterator)

	if loop.Kind == ForLoop {
		v.emitter.PatchJump(loop.Jump)
	} else {
		v.emitter.PatchJumpAB(loop.Jump)
	}

	// Allocate registers
	stack := loop.Iterator
	low := v.registers.Allocate(Temp)
	high := v.registers.Allocate(Temp)
	pivot := v.registers.Allocate(Temp)
	n := v.registers.Allocate(Temp)
	i := v.registers.Allocate(Temp)
	j := v.registers.Allocate(Temp)
	comp := v.registers.Allocate(Temp)

	// Init sort
	// Create a new stack
	v.emitter.EmitA(runtime.OpSortPrep, stack)

	// low = 0
	v.emitter.EmitA(runtime.OpLoadZero, low)

	// high = len(arr) - 1
	v.emitter.EmitAB(runtime.OpLength, high, loop.Result)
	v.emitter.EmitA(runtime.OpDecr, high)

	// Push initial range: stack = append(stack, high, low)
	v.emitter.EmitAB(runtime.OpSortPush, stack, high)
	v.emitter.EmitAB(runtime.OpSortPush, stack, low)

	// Main stack loop
	// while len(stack) > 0
	stackLoopStart := v.emitter.Size()
	stackLoopExitJmp := v.emitter.EmitJumpc(runtime.OpJumpIfEmpty, jumpPlaceholder, stack)

	// Pop the range: low = stack.pop(), high = stack.pop()
	v.emitter.EmitAB(runtime.OpSortPop, low, stack)
	v.emitter.EmitAB(runtime.OpSortPop, high, stack)

	// comp = low < high
	v.emitter.EmitABC(runtime.OpLt, comp, low, high)
	v.emitter.EmitJumpc(runtime.OpJumpIfFalse, stackLoopStart, comp)

	// Otherwise start the partition
	// pivot := arr[high]
	v.emitter.EmitABC(runtime.OpLoadProperty, pivot, loop.Result, high)

	// i = low - 1
	v.emitter.EmitAB(runtime.OpMove, i, low)
	v.emitter.EmitA(runtime.OpDecr, i)

	// j = low
	v.emitter.EmitAB(runtime.OpMove, j, low)

	// for j < high
	partitionLoopStart := v.emitter.Size()
	v.emitter.EmitABC(runtime.OpLt, comp, j, high)
	partitionLoopEnd := v.emitter.EmitJumpc(runtime.OpJumpIfFalse, jumpPlaceholder, comp)

	// Load current value into n
	// n = arr[j]
	v.emitter.EmitABC(runtime.OpLoadProperty, n, loop.Result, j)

	// Collect all sort clauses (mostly one or two)
	clauses := ctx.AllSortClauseExpression()
	skipSwapJumps := make([]int, len(clauses))

	// Compare pivot with each clause
	for i, clause := range clauses {
		sort := clause.(*fql.SortClauseExpressionContext)
		exp := sort.Expression()

		// We override the loop variables with the pivot value
		if loop.ValueName != "" {
			varReg := v.symbols.Variable(loop.ValueName)
			v.emitter.EmitAB(runtime.OpSortValue, varReg, pivot)
		}

		if loop.KeyName != "" {
			varReg := v.symbols.Variable(loop.KeyName)
			v.emitter.EmitAB(runtime.OpSortKey, varReg, pivot)
		}

		pivotReg := exp.Accept(v).(runtime.Operand)

		// And now override the loop variables with the n value
		if loop.ValueName != "" {
			varReg := v.symbols.Variable(loop.ValueName)
			v.emitter.EmitAB(runtime.OpSortValue, varReg, n)
		}

		if loop.KeyName != "" {
			varReg := v.symbols.Variable(loop.KeyName)
			v.emitter.EmitAB(runtime.OpSortKey, varReg, n)
		}

		currentReg := exp.Accept(v).(runtime.Operand)

		comparator := runtime.OpLte
		direction := clause.SortDirection()

		if direction != nil && strings.ToLower(direction.GetText()) == "desc" {
			comparator = runtime.OpGte
		}

		// comp = current <= pivot or current >= pivot
		v.emitter.EmitABC(comparator, comp, currentReg, pivotReg)
		// If comp is false, jump to loop end
		skipSwapJumps[i] = v.emitter.EmitJumpc(runtime.OpJumpIfFalse, jumpPlaceholder, comp)
	}

	// i++
	v.emitter.EmitA(runtime.OpIncr, i)
	// swap arr[i], arr[j]
	v.emitter.EmitABC(runtime.OpSortSwap, loop.Result, i, j)

	// Patch all clause skip swap jumps
	for _, jmp := range skipSwapJumps {
		v.emitter.PatchJumpNext(jmp)
	}

	// j++
	v.emitter.EmitA(runtime.OpIncr, j)
	v.emitter.EmitJump(runtime.OpJump, partitionLoopStart)

	// End of partition loop
	v.emitter.PatchJumpNext(partitionLoopEnd)

	// i++
	v.emitter.EmitA(runtime.OpIncr, i)
	// swap arr[i], arr[high]
	v.emitter.EmitABC(runtime.OpSortSwap, loop.Result, i, high)

	// Push left subarray: stack = append(stack, low, i-1)
	v.emitter.EmitAB(runtime.OpMove, n, i)
	v.emitter.EmitA(runtime.OpDecr, n)
	v.emitter.EmitAB(runtime.OpSortPush, stack, n)
	v.emitter.EmitAB(runtime.OpSortPush, stack, low)

	// Push right subarray: stack = append(stack, {i+1, high})
	v.emitter.EmitAB(runtime.OpMove, n, i)
	v.emitter.EmitA(runtime.OpIncr, n)
	v.emitter.EmitAB(runtime.OpSortPush, stack, high)
	v.emitter.EmitAB(runtime.OpSortPush, stack, n)

	// Jump to stack loop beginning
	v.emitter.EmitJump(runtime.OpJump, stackLoopStart)

	// Patch the stack loop exit jump
	v.emitter.PatchJumpNext(stackLoopExitJmp)

	// Replace source with sorted array
	v.emitter.EmitAB(runtime.OpSortCollect, loop.Src, loop.Result)

	// Create new for loop
	// TODO: Reuse existing DataSet instance
	v.emitLoopBegin(loop)

	return nil
}

func (v *visitor) visitOffset(src1 runtime.Operand) interface{} {
	state := v.registers.Allocate(State)
	v.emitter.EmitA(runtime.OpIncr, state)

	comp := v.registers.Allocate(Temp)
	v.emitter.EmitABC(runtime.OpGt, comp, state, src1)
	v.emitter.EmitJumpc(runtime.OpJumpIfFalse, v.loops.Loop().Jump, comp)

	return state
}

func (v *visitor) visitLimit(src1 runtime.Operand) interface{} {
	state := v.registers.Allocate(State)
	v.emitter.EmitA(runtime.OpIncr, state)

	comp := v.registers.Allocate(Temp)
	v.emitter.EmitABC(runtime.OpGt, comp, state, src1)
	v.emitter.EmitJumpc(runtime.OpJumpIfTrue, v.loops.Loop().Jump, comp)

	return state
}

func (v *visitor) VisitLimitClauseValue(ctx *fql.LimitClauseValueContext) interface{} {
	if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	}

	panic(core.Error(ErrUnexpectedToken, ctx.GetText()))
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
	op := v.symbols.Variable(ctx.GetText())

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
			jump := -1
			endCatch := v.emitter.Size()

			if c.ForExpression() != nil {
				// We jump back to finalize the loop before exiting
				jump = endCatch - 1
			}

			v.catchTable = append(v.catchTable, [3]int{startCatch, endCatch, jump})
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
		if op.Not() == nil {
			opcode = runtime.OpIn
		} else {
			opcode = runtime.OpNotIn
		}
	} else if op := ctx.LikeOperator(); op != nil {
		if op.(*fql.LikeOperatorContext).Not() == nil {
			opcode = runtime.OpLike
		} else {
			opcode = runtime.OpNotLike
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

		if opcode == runtime.OpRegexpPositive || opcode == runtime.OpRegexpNegative {
			if regRight.IsConstant() {
				val := v.symbols.Constant(regRight)

				// Verify that the expression is a valid regular expression
				regexp.MustCompile(val.String())
			}
		}

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

func (v *visitor) visitFunctionCall(ctx *fql.FunctionCallContext, protected bool) interface{} {
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

	name := v.functionName(ctx)

	switch name {
	case "LENGTH":
		dst := v.registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(core.Error(core.ErrInvalidArgument, "LENGTH: expected 1 argument"))
		}

		v.emitter.EmitAB(runtime.OpLength, dst, seq.Registers[0])

		return dst
	case "TYPENAME":
		dst := v.registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(core.Error(core.ErrInvalidArgument, "TYPENAME: expected 1 argument"))
		}

		v.emitter.EmitAB(runtime.OpType, dst, seq.Registers[0])

		return dst
	default:
		nameAndDest := v.loadConstant(v.functionName(ctx))

		if !protected {
			v.emitter.EmitAs(runtime.OpCall, nameAndDest, seq)
		} else {
			v.emitter.EmitAs(runtime.OpProtectedCall, nameAndDest, seq)
		}

		return nameAndDest
	}
}

func (v *visitor) functionName(ctx *fql.FunctionCallContext) values.String {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	return values.NewString(strings.ToUpper(name))
}

func (v *visitor) emitLoopBegin(loop *Loop) {
	if loop.Allocate {
		v.emitter.EmitAb(runtime.OpLoopBegin, loop.Result, loop.Distinct)
	}

	loop.Iterator = v.registers.Allocate(State)

	if loop.Kind == ForLoop {
		v.emitter.EmitAB(runtime.OpForLoopPrep, loop.Iterator, loop.Src)
		// jumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = v.emitter.EmitJumpc(runtime.OpForLoopNext, jumpPlaceholder, loop.Iterator)

		if loop.Value != runtime.NoopOperand {
			v.emitter.EmitAB(runtime.OpForLoopValue, loop.Value, loop.Iterator)
		}

		if loop.Key != runtime.NoopOperand {
			v.emitter.EmitAB(runtime.OpForLoopKey, loop.Key, loop.Iterator)
		}
	} else {
		//counterReg := v.registers.Allocate(State)
		// TODO: Set JumpOffset here
	}
}

func (v *visitor) emitLoopEnd(loop *Loop) runtime.Operand {
	v.emitter.EmitJump(runtime.OpJump, loop.Jump-loop.JumpOffset)

	// TODO: Do not allocate for pass-through loops
	dst := v.registers.Allocate(Temp)

	if loop.Allocate {
		// TODO: Reuse the dsReg register
		v.emitter.EmitA(runtime.OpClose, loop.Iterator)
		v.emitter.EmitAB(runtime.OpLoopEnd, dst, loop.Result)

		if loop.Kind == ForLoop {
			v.emitter.PatchJump(loop.Jump)
		} else {
			v.emitter.PatchJumpAB(loop.Jump)
		}
	} else {
		if loop.Kind == ForLoop {
			v.emitter.PatchJumpNext(loop.Jump)
		} else {
			v.emitter.PatchJumpNextAB(loop.Jump)
		}
	}

	return dst
}

func (v *visitor) loopType(ctx *fql.ForExpressionContext) LoopType {
	if c := ctx.ForExpressionReturn().ForExpression(); c == nil {
		return NormalLoop
	}

	return PassThroughLoop
}

func (v *visitor) loopKind(ctx *fql.ForExpressionContext) LoopKind {
	if ctx.While() == nil {
		return ForLoop
	}

	if ctx.Do() == nil {
		return WhileLoop
	}

	return DoWhileLoop
}

func (v *visitor) loadConstant(constant core.Value) runtime.Operand {
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(runtime.OpLoadConst, reg, v.symbols.AddConstant(constant))
	return reg
}
