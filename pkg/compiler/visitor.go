package compiler

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"

	"github.com/antlr4-go/antlr/v4"
)

type visitor struct {
	*fql.BaseFqlParserVisitor
	err        error
	src        string
	emitter    *Emitter
	registers  *RegisterAllocator
	symbols    *SymbolTable
	loops      *LoopTable
	catchTable []vm.Catch
}

const (
	jumpPlaceholder      = -1
	undefinedVariable    = -1
	ignorePseudoVariable = "_"
	pseudoVariable       = "CURRENT"
	forScope             = "for"
)

// Runtime functions
const (
	runtimeTypename = "TYPENAME"
	runtimeLength   = "LENGTH"
	runtimeWait     = "WAIT"
)

func newVisitor(src string) *visitor {
	v := new(visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.src = src
	v.registers = NewRegisterAllocator()
	v.symbols = NewSymbolTable(v.registers)
	v.loops = NewLoopTable(v.registers)
	v.emitter = NewEmitter()
	v.catchTable = make([]vm.Catch, 0)

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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitBodyExpression(ctx *fql.BodyExpressionContext) interface{} {
	if c := ctx.ForExpression(); c != nil {
		out, ok := c.Accept(v).(vm.Operand)

		if ok && out != vm.NoopOperand {
			v.emitter.EmitAB(vm.OpMove, vm.NoopOperand, out)
		}

		v.emitter.Emit(vm.OpReturn)

		return out
	} else if c := ctx.ReturnExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func (v *visitor) VisitWaitForExpression(ctx *fql.WaitForExpressionContext) interface{} {
	if ctx.Event() != nil {
		return v.visitWaitForEventExpression(ctx)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) visitWaitForEventExpression(ctx *fql.WaitForExpressionContext) interface{} {
	v.symbols.EnterScope()

	srcReg := ctx.WaitForEventSource().Accept(v).(vm.Operand)
	eventReg := ctx.WaitForEventName().Accept(v).(vm.Operand)

	var optsReg vm.Operand

	if opts := ctx.OptionsClause(); opts != nil {
		optsReg = opts.Accept(v).(vm.Operand)
	}

	var timeoutReg vm.Operand

	if timeout := ctx.TimeoutClause(); timeout != nil {
		timeoutReg = timeout.Accept(v).(vm.Operand)
	}

	streamReg := v.registers.Allocate(Temp)

	// We move the source object to the stream register in order to re-use it in OpStream
	v.emitter.EmitAB(vm.OpMove, streamReg, srcReg)
	v.emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
	v.emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)

	var valReg vm.Operand

	// Now we start iterating over the stream
	jumpToNext := v.emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, streamReg)

	if filter := ctx.FilterClause(); filter != nil {
		valReg = v.symbols.DefineVariable(pseudoVariable)
		v.emitter.EmitAB(vm.OpIterValue, valReg, streamReg)

		filter.Expression().Accept(v)

		v.emitter.EmitJumpc(vm.OpJumpIfFalse, jumpToNext, valReg)

		// TODO: Do we need to use timeout here too? We can really get stuck in the loop if no event satisfies the filter
	}

	// Clean up the stream
	v.emitter.EmitA(vm.OpClose, streamReg)

	v.symbols.ExitScope()

	return nil
}

func (v *visitor) VisitWaitForEventName(ctx *fql.WaitForEventNameContext) interface{} {
	if c := ctx.StringLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitWaitForEventSource(ctx *fql.WaitForEventSourceContext) interface{} {
	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitTimeoutClauseContext(ctx *fql.TimeoutClauseContext) interface{} {
	if c := ctx.IntegerLiteral(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Variable(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.Param(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.MemberExpression(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCall(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitOptionsClause(ctx *fql.OptionsClauseContext) interface{} {
	if c := ctx.ObjectLiteral(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
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
		loop.Src = ctx.ForExpressionSource().Accept(v).(vm.Operand)

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
		expReg := c.Expression().Accept(v).(vm.Operand)

		v.emitter.EmitAB(vm.OpDataSetAdd, loop.Result, expReg)
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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitForExpressionBody(ctx *fql.ForExpressionBodyContext) interface{} {
	if c := ctx.ForExpressionClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.ForExpressionStatement(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitForExpressionClause(ctx *fql.ForExpressionClauseContext) interface{} {
	if c := ctx.LimitClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FilterClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.SortClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.CollectClause(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitFilterClause(ctx *fql.FilterClauseContext) interface{} {
	src1 := ctx.Expression().Accept(v).(vm.Operand)
	v.emitter.EmitJumpc(vm.OpJumpIfFalse, v.loops.Loop().Jump, src1)

	return nil
}

func (v *visitor) VisitLimitClause(ctx *fql.LimitClauseContext) interface{} {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		return v.visitLimit(clauses[0].Accept(v).(vm.Operand))
	} else {
		v.visitOffset(clauses[0].Accept(v).(vm.Operand))
		v.visitLimit(clauses[1].Accept(v).(vm.Operand))
	}

	return nil
}

func (v *visitor) VisitCollectClause(ctx *fql.CollectClauseContext) interface{} {
	if c := ctx.CollectGrouping(); c != nil {
		// Collect by grouping
		return c.Accept(v)
	}

	return nil
}

func (v *visitor) VisitCollectGrouping(ctx *fql.CollectGroupingContext) interface{} {
	// TODO: Undefine original loop variables
	loop := v.loops.Loop()

	// We collect the aggregation keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	var kvKeyReg vm.Operand
	selectors := ctx.AllCollectSelector()
	isMultiSelector := len(selectors) > 1

	if isMultiSelector {
		// We create a sequence of registers for the clauses
		// To pack them into an array
		selectorRegs := v.registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := selector.Accept(v).(vm.Operand)
			v.emitter.EmitAB(vm.OpMove, selectorRegs.Registers[i], reg)
			// Free the register after moving its value to the sequence register
			v.registers.Free(reg)
		}

		kvKeyReg = v.registers.Allocate(Temp)
		v.emitter.EmitAs(vm.OpList, kvKeyReg, selectorRegs)
		v.registers.FreeSequence(selectorRegs)
	} else {
		kvKeyReg = selectors[0].Accept(v).(vm.Operand)
	}

	kvValReg := v.registers.Allocate(Temp)

	if loop.Kind == ForLoop {
		v.emitter.EmitAB(vm.OpIterValue, kvValReg, loop.Iterator)
	} else {
		v.emitter.EmitAB(vm.OpWhileLoopValue, kvValReg, loop.Iterator)
	}

	v.emitter.EmitABC(vm.OpDataSetAddKV, loop.Result, kvKeyReg, kvValReg)
	v.emitter.EmitJump(vm.OpJump, loop.Jump-loop.JumpOffset)
	v.emitter.EmitA(vm.OpClose, loop.Iterator)

	if loop.Kind == ForLoop {
		v.emitter.PatchJump(loop.Jump)
	} else {
		v.emitter.PatchJumpAB(loop.Jump)
	}

	v.emitter.EmitA(vm.OpCollectGrouping, loop.Result)

	// Replace source with sorted array
	v.emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	v.symbols.ExitScope()
	v.symbols.EnterScope()

	loop.ValueName = ""
	loop.KeyName = ""
	v.registers.Free(loop.Value)
	v.registers.Free(loop.Key)
	loop.Value = vm.NoopOperand
	loop.Key = vm.NoopOperand

	// Create new for loop
	v.emitLoopBegin(loop)

	// Now we need to expand group variables from the dataset
	v.emitter.EmitAB(vm.OpIterKey, kvValReg, loop.Iterator)

	if isMultiSelector {
		// Define a variable for each selector
		for i, selector := range selectors {
			// Get the variable name
			name := selector.Identifier().GetText()

			v.emitter.EmitABC(vm.OpLoadIndex, v.symbols.DefineVariable(name), kvValReg, v.loadConstant(runtime.Int(i)))
		}
	} else {
		// Get the variable name
		name := selectors[0].Identifier().GetText()
		// Define a variable for each selector
		varReg := v.symbols.DefineVariable(name)
		// If we have a single selector, we can just move the value
		v.emitter.EmitAB(vm.OpMove, varReg, kvValReg)
	}

	return nil
}

func (v *visitor) VisitCollectSelector(ctx *fql.CollectSelectorContext) interface{} {
	if c := ctx.Expression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitSortClause(ctx *fql.SortClauseContext) interface{} {
	loop := v.loops.Loop()

	// We collect the sorting conditions (keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	kvKeyReg := v.registers.Allocate(Temp)
	clauses := ctx.AllSortClauseExpression()
	isSortMany := len(clauses) > 1

	// For multi-sort
	var directionRegs *RegisterSequence

	if isSortMany {
		clausesRegs := make([]vm.Operand, len(clauses))
		// We create a sequence of registers for the clauses
		// To pack them into an array
		keyRegs := v.registers.AllocateSequence(len(clauses))

		// We create a sequence of registers for the directions
		directionRegs = v.registers.AllocateSequence(len(clauses))

		for i, clause := range clauses {
			clauseReg := clause.Accept(v).(vm.Operand)
			v.emitter.EmitAB(vm.OpMove, keyRegs.Registers[i], clauseReg)
			clausesRegs[i] = keyRegs.Registers[i]
			v.visitSortDirection(clause.SortDirection(), directionRegs.Registers[i])

			// TODO: Free registers
		}

		arrReg := v.registers.Allocate(Temp)
		v.emitter.EmitAs(vm.OpList, arrReg, keyRegs)
		v.emitter.EmitAB(vm.OpMove, kvKeyReg, arrReg) // TODO: Free registers
	} else {
		clausesReg := clauses[0].Accept(v).(vm.Operand)
		v.emitter.EmitAB(vm.OpMove, kvKeyReg, clausesReg)
	}

	var kvValReg vm.Operand

	// In case the value is not used in the loop body, and only key is used
	if loop.ValueName != "" {
		kvValReg = loop.Value
	} else {
		// If so, we need to load it from the iterator
		kvValReg = v.registers.Allocate(Temp)

		if loop.Kind == ForLoop {
			v.emitter.EmitAB(vm.OpIterValue, kvValReg, loop.Iterator)
		} else {
			v.emitter.EmitAB(vm.OpWhileLoopValue, kvValReg, loop.Iterator)
		}
	}

	v.emitter.EmitABC(vm.OpDataSetAddKV, loop.Result, kvKeyReg, kvValReg)
	v.emitter.EmitJump(vm.OpJump, loop.Jump-loop.JumpOffset)
	v.emitter.EmitA(vm.OpClose, loop.Iterator)

	if loop.Kind == ForLoop {
		v.emitter.PatchJump(loop.Jump)
	} else {
		v.emitter.PatchJumpAB(loop.Jump)
	}

	if isSortMany {
		v.emitter.EmitAs(vm.OpSortMany, loop.Result, directionRegs)
	} else {
		directionReg := v.registers.Allocate(Temp)
		v.visitSortDirection(clauses[0].SortDirection(), directionReg)
		v.emitter.EmitAB(vm.OpSort, loop.Result, directionReg)
	}

	// Replace source with sorted array
	v.emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	// Create new for loop
	// TODO: Reuse existing DataSet instance
	v.emitLoopBegin(loop)

	return nil
}

func (v *visitor) visitSortDirection(dir antlr.TerminalNode, dest vm.Operand) {
	var val runtime.Int = vm.SortAsc

	if dir != nil {
		if strings.ToLower(dir.GetText()) == "desc" {
			val = vm.SortDesc
		}
	}

	// TODO: Free constant registers
	v.emitter.EmitAB(vm.OpMove, dest, v.loadConstant(val))
}

func (v *visitor) VisitSortClauseExpression(ctx *fql.SortClauseExpressionContext) interface{} {
	return ctx.Expression().Accept(v).(vm.Operand)
}

func (v *visitor) visitOffset(src1 vm.Operand) interface{} {
	state := v.registers.Allocate(State)
	v.emitter.EmitABx(vm.OpSkip, state, src1, v.loops.Loop().Jump)

	return state
}

func (v *visitor) visitLimit(src1 vm.Operand) interface{} {
	state := v.registers.Allocate(State)
	v.emitter.EmitABx(vm.OpLimit, state, src1, v.loops.Loop().Jump)

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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitForExpressionStatement(ctx *fql.ForExpressionStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
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

	var dst vm.Operand
	src1 := mesOut.(vm.Operand)

	for _, segment := range segments {
		var out2 interface{}
		p := segment.(*fql.MemberExpressionPathContext)

		if c := p.PropertyName(); c != nil {
			out2 = c.Accept(v)
		} else if c := p.ComputedPropertyName(); c != nil {
			out2 = c.Accept(v)
		}

		src2 := out2.(vm.Operand)
		dst = v.registers.Allocate(Temp)

		if p.ErrorOperator() != nil {
			v.emitter.EmitABC(vm.OpLoadPropertyOptional, dst, src1, src2)
		} else {
			v.emitter.EmitABC(vm.OpLoadProperty, dst, src1, src2)
		}

		src1 = dst
	}

	return dst
}

func (v *visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	dst := v.registers.Allocate(Temp)
	start := ctx.GetLeft().Accept(v).(vm.Operand)
	end := ctx.GetRight().Accept(v).(vm.Operand)

	v.emitter.EmitABC(vm.OpRange, dst, start, end)

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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitParam(ctx *fql.ParamContext) interface{} {
	name := ctx.Identifier().GetText()
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(vm.OpLoadParam, reg, v.symbols.AddParam(name))

	return reg
}

func (v *visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	src := ctx.Expression().Accept(v).(vm.Operand)

	if name != ignorePseudoVariable {
		dest := v.symbols.DefineVariable(name)

		if src.IsConstant() {
			tmp := v.registers.Allocate(Temp)
			v.emitter.EmitAB(vm.OpLoadConst, tmp, src)
			v.emitter.EmitAB(vm.OpStoreGlobal, dest, tmp)
		} else if v.symbols.Scope() == 0 {
			v.emitter.EmitAB(vm.OpStoreGlobal, dest, src)
		} else {
			v.emitter.EmitAB(vm.OpMove, dest, src)
		}

		return dest
	}

	return vm.NoopOperand
}

func (v *visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	// Just return the register / constant index
	op := v.symbols.Variable(ctx.GetText())

	if op.IsRegister() {
		return op
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(vm.OpLoadGlobal, reg, op)

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
				srcReg := exp.Accept(v).(vm.Operand)

				// TODO: Figure out how to remove OpMove and use registers returned from each expression
				v.emitter.EmitAB(vm.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.registers.Free(srcReg)
				}
			}

			// Initialize an array
			v.emitter.EmitAs(vm.OpList, destReg, seq)

			// Free seq registers
			//v.registers.FreeSequence(seq)

			return destReg
		}
	}

	// Empty array
	v.emitter.EmitA(vm.OpList, destReg)

	return destReg
}

func (v *visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	dst := v.registers.Allocate(Temp)
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size == 0 {
		v.emitter.EmitA(vm.OpMap, dst)

		return dst
	}

	seq := v.registers.AllocateSequence(len(assignments) * 2)

	for i := 0; i < size; i++ {
		var propOp vm.Operand
		var valOp vm.Operand
		pac := assignments[i].(*fql.PropertyAssignmentContext)

		if prop, ok := pac.PropertyName().(*fql.PropertyNameContext); ok {
			propOp = prop.Accept(v).(vm.Operand)
			valOp = pac.Expression().Accept(v).(vm.Operand)
		} else if comProp, ok := pac.ComputedPropertyName().(*fql.ComputedPropertyNameContext); ok {
			propOp = comProp.Accept(v).(vm.Operand)
			valOp = pac.Expression().Accept(v).(vm.Operand)
		} else if variable := pac.Variable(); variable != nil {
			propOp = v.loadConstant(runtime.NewString(variable.GetText()))
			valOp = variable.Accept(v).(vm.Operand)
		}

		regIndex := i * 2

		v.emitter.EmitAB(vm.OpMove, seq.Registers[regIndex], propOp)
		v.emitter.EmitAB(vm.OpMove, seq.Registers[regIndex+1], valOp)

		// Free source register if temporary
		if propOp.IsRegister() {
			//v.registers.Free(propOp)
		}
	}

	v.emitter.EmitAs(vm.OpMap, dst, seq)

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
		panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return v.loadConstant(runtime.NewString(name))
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

	return v.loadConstant(runtime.NewString(b.String()))
}

func (v *visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(vm.OpLoadConst, reg, v.symbols.AddConstant(runtime.NewInt(val)))

	return reg
}

func (v *visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(vm.OpLoadConst, reg, v.symbols.AddConstant(runtime.NewFloat(val)))

	return reg
}

func (v *visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	reg := v.registers.Allocate(Temp)

	switch strings.ToLower(ctx.GetText()) {
	case "true":
		v.emitter.EmitAB(vm.OpLoadBool, reg, 1)
	case "false":
		v.emitter.EmitAB(vm.OpLoadBool, reg, 0)
	default:
		panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return reg
}

func (v *visitor) VisitNoneLiteral(_ *fql.NoneLiteralContext) interface{} {
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitA(vm.OpLoadNone, reg)

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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *visitor) VisitReturnExpression(ctx *fql.ReturnExpressionContext) interface{} {
	valReg := ctx.Expression().Accept(v).(vm.Operand)

	if valReg.IsConstant() {
		v.emitter.EmitAB(vm.OpLoadGlobal, vm.NoopOperand, valReg)
	} else {
		v.emitter.EmitAB(vm.OpMove, vm.NoopOperand, valReg)
	}

	v.emitter.Emit(vm.OpReturn)

	return vm.NoopOperand
}

func (v *visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if uo := ctx.UnaryOperator(); uo != nil {
		src := ctx.GetRight().Accept(v).(vm.Operand)
		dst := v.registers.Allocate(Temp)

		uoc := uo.(*fql.UnaryOperatorContext)
		var op vm.Opcode

		if uoc.Not() != nil {
			op = vm.OpNot
		} else if uoc.Minus() != nil {
			op = vm.OpFlipNegative
		} else if uoc.Plus() != nil {
			op = vm.OpFlipPositive
		} else {
			panic(runtime.Error(ErrUnexpectedToken, uoc.GetText()))
		}

		// We do not overwrite the source register
		v.emitter.EmitAB(op, dst, src)

		return dst
	}

	if op := ctx.LogicalAndOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(vm.Operand)
		v.emitter.EmitAB(vm.OpMove, dst, left)
		// Test if left is false and jump to the end
		end := v.emitter.EmitJumpc(vm.OpJumpIfFalse, jumpPlaceholder, dst)
		// If left is true, execute right expression
		right := ctx.GetRight().Accept(v).(vm.Operand)
		// And move the result to the destination register
		v.emitter.EmitAB(vm.OpMove, dst, right)
		v.emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.LogicalOrOperator(); op != nil {
		dst := v.registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(vm.Operand)
		// Move the result to the destination register
		v.emitter.EmitAB(vm.OpMove, dst, left)
		// Test if left is true and jump to the end
		end := v.emitter.EmitJumpc(vm.OpJumpIfTrue, jumpPlaceholder, dst)
		// If left is false, execute right expression
		right := ctx.GetRight().Accept(v).(vm.Operand)
		// And move the result to the destination register
		v.emitter.EmitAB(vm.OpMove, dst, right)
		v.emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.GetTernaryOperator(); op != nil {
		dst := v.registers.Allocate(Temp)

		// Compile condition and put result in dst
		condReg := ctx.GetCondition().Accept(v).(vm.Operand)
		v.emitter.EmitAB(vm.OpMove, dst, condReg)

		// If condition was temporary, free it
		if condReg.IsRegister() {
			//v.registers.Free(condReg)
		}

		// Jump to 'false' branch if condition is false
		otherwise := v.emitter.EmitJumpc(vm.OpJumpIfFalse, jumpPlaceholder, dst)

		// True branch
		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			trueReg := onTrue.Accept(v).(vm.Operand)
			v.emitter.EmitAB(vm.OpMove, dst, trueReg)

			// Free temporary register if needed
			if trueReg.IsRegister() {
				//v.registers.Free(trueReg)
			}
		}

		// Jump over false branch
		end := v.emitter.EmitJump(vm.OpJump, jumpPlaceholder)
		v.emitter.PatchJumpNext(otherwise)

		// False branch
		if onFalse := ctx.GetOnFalse(); onFalse != nil {
			falseReg := onFalse.Accept(v).(vm.Operand)
			v.emitter.EmitAB(vm.OpMove, dst, falseReg)

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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
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

	var opcode vm.Opcode
	dest := v.registers.Allocate(Temp)
	left := ctx.Predicate(0).Accept(v).(vm.Operand)
	right := ctx.Predicate(1).Accept(v).(vm.Operand)

	if op := ctx.EqualityOperator(); op != nil {
		switch op.GetText() {
		case "==":
			opcode = vm.OpEq
		case "!=":
			opcode = vm.OpNeq
		case ">":
			opcode = vm.OpGt
		case ">=":
			opcode = vm.OpGte
		case "<":
			opcode = vm.OpLt
		case "<=":
			opcode = vm.OpLte
		default:
			panic(runtime.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.ArrayOperator(); op != nil {
		// TODO: Implement me
		panic(runtime.Error(runtime.ErrNotImplemented, "array operator"))
	} else if op := ctx.InOperator(); op != nil {
		if op.Not() == nil {
			opcode = vm.OpIn
		} else {
			opcode = vm.OpNotIn
		}
	} else if op := ctx.LikeOperator(); op != nil {
		if op.(*fql.LikeOperatorContext).Not() == nil {
			opcode = vm.OpLike
		} else {
			opcode = vm.OpNotLike
		}
	}

	v.emitter.EmitABC(opcode, dest, left, right)

	return dest
}

func (v *visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
	var opcode vm.Opcode
	var isSet bool

	if op := ctx.MultiplicativeOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "*":
			opcode = vm.OpMulti
		case "/":
			opcode = vm.OpDiv
		case "%":
			opcode = vm.OpMod
		default:
			panic(runtime.Error(ErrUnexpectedToken, op.GetText()))
		}
	} else if op := ctx.AdditiveOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "+":
			opcode = vm.OpAdd
		case "-":
			opcode = vm.OpSub
		default:
			panic(runtime.Error(ErrUnexpectedToken, op.GetText()))
		}

	} else if op := ctx.RegexpOperator(); op != nil {
		isSet = true

		switch op.GetText() {
		case "=~":
			opcode = vm.OpRegexpPositive
		case "!~":
			opcode = vm.OpRegexpNegative
		default:
			panic(runtime.Error(ErrUnexpectedToken, op.GetText()))
		}
	}

	if isSet {
		regLeft := ctx.ExpressionAtom(0).Accept(v).(vm.Operand)
		regRight := ctx.ExpressionAtom(1).Accept(v).(vm.Operand)
		dst := v.registers.Allocate(Temp)

		if opcode == vm.OpRegexpPositive || opcode == vm.OpRegexpNegative {
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

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
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
				srcReg := exp.Accept(v).(vm.Operand)

				// TODO: Figure out how to remove OpMove and use registers returned from each expression
				v.emitter.EmitAB(vm.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.registers.Free(srcReg)
				}
			}
		}
	}

	name := v.functionName(ctx)

	switch name {
	case runtimeLength:
		dst := v.registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
		}

		v.emitter.EmitAB(vm.OpLength, dst, seq.Registers[0])

		return dst
	case runtimeTypename:
		dst := v.registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
		}

		v.emitter.EmitAB(vm.OpType, dst, seq.Registers[0])

		return dst
	case runtimeWait:
		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
		}

		v.emitter.EmitA(vm.OpSleep, seq.Registers[0])

		return seq.Registers[0]
	default:
		nameAndDest := v.loadConstant(v.functionName(ctx))

		if !protected {
			v.emitter.EmitAs(vm.OpCall, nameAndDest, seq)
		} else {
			v.emitter.EmitAs(vm.OpProtectedCall, nameAndDest, seq)
		}

		return nameAndDest
	}
}

func (v *visitor) functionName(ctx *fql.FunctionCallContext) runtime.String {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	return runtime.NewString(strings.ToUpper(name))
}

func (v *visitor) emitLoopBegin(loop *Loop) {
	if loop.Allocate {
		v.emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
	}

	loop.Iterator = v.registers.Allocate(State)

	if loop.Kind == ForLoop {
		v.emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Src)
		// jumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = v.emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, loop.Iterator)

		if loop.Value != vm.NoopOperand {
			v.emitter.EmitAB(vm.OpIterValue, loop.Value, loop.Iterator)
		}

		if loop.Key != vm.NoopOperand {
			v.emitter.EmitAB(vm.OpIterKey, loop.Key, loop.Iterator)
		}
	} else {
		//counterReg := v.registers.Allocate(Storage)
		// TODO: Set JumpOffset here
	}
}

func (v *visitor) emitLoopEnd(loop *Loop) vm.Operand {
	v.emitter.EmitJump(vm.OpJump, loop.Jump-loop.JumpOffset)

	// TODO: Do not allocate for pass-through loops
	dst := v.registers.Allocate(Temp)

	if loop.Allocate {
		// TODO: Reuse the dsReg register
		v.emitter.EmitA(vm.OpClose, loop.Iterator)
		v.emitter.EmitAB(vm.OpMove, dst, loop.Result)

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

func (v *visitor) loadConstant(constant runtime.Value) vm.Operand {
	reg := v.registers.Allocate(Temp)
	v.emitter.EmitAB(vm.OpLoadConst, reg, v.symbols.AddConstant(constant))
	return reg
}
