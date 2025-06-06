package internal

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"

	"github.com/antlr4-go/antlr/v4"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor
	Err        error
	Src        string
	Emitter    *Emitter
	Registers  *RegisterAllocator
	Symbols    *SymbolTable
	Loops      *LoopTable
	CatchTable []vm.Catch
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

func NewVisitor(src string) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Src = src
	v.Registers = NewRegisterAllocator()
	v.Symbols = NewSymbolTable(v.Registers)
	v.Loops = NewLoopTable(v.Registers)
	v.Emitter = NewEmitter()
	v.CatchTable = make([]vm.Catch, 0)

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	ctx.Body().Accept(v)

	return nil
}

func (v *Visitor) VisitBody(ctx *fql.BodyContext) interface{} {
	for _, statement := range ctx.AllBodyStatement() {
		statement.Accept(v)
	}

	ctx.BodyExpression().Accept(v)

	return nil
}

func (v *Visitor) VisitBodyStatement(ctx *fql.BodyStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	} else if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	} else if c := ctx.WaitForExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitBodyExpression(ctx *fql.BodyExpressionContext) interface{} {
	if c := ctx.ForExpression(); c != nil {
		out, ok := c.Accept(v).(vm.Operand)

		if ok && out != vm.NoopOperand {
			v.Emitter.EmitAB(vm.OpMove, vm.NoopOperand, out)
		}

		v.Emitter.Emit(vm.OpReturn)

		return out
	} else if c := ctx.ReturnExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

func (v *Visitor) VisitWaitForExpression(ctx *fql.WaitForExpressionContext) interface{} {
	if ctx.Event() != nil {
		return v.visitWaitForEventExpression(ctx)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) visitWaitForEventExpression(ctx *fql.WaitForExpressionContext) interface{} {
	v.Symbols.EnterScope()

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

	streamReg := v.Registers.Allocate(Temp)

	// We move the source object to the stream register in order to re-use it in OpStream
	v.Emitter.EmitAB(vm.OpMove, streamReg, srcReg)
	v.Emitter.EmitABC(vm.OpStream, streamReg, eventReg, optsReg)
	v.Emitter.EmitAB(vm.OpStreamIter, streamReg, timeoutReg)

	var valReg vm.Operand

	// Now we start iterating over the stream
	jumpToNext := v.Emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, streamReg)

	if filter := ctx.FilterClause(); filter != nil {
		valReg = v.Symbols.DefineVariable(pseudoVariable)
		v.Emitter.EmitAB(vm.OpIterValue, valReg, streamReg)

		filter.Expression().Accept(v)

		v.Emitter.EmitJumpc(vm.OpJumpIfFalse, jumpToNext, valReg)

		// TODO: Do we need to use timeout here too? We can really get stuck in the loop if no event satisfies the filter
	}

	// Clean up the stream
	v.Emitter.EmitA(vm.OpClose, streamReg)

	v.Symbols.ExitScope()

	return nil
}

func (v *Visitor) VisitWaitForEventName(ctx *fql.WaitForEventNameContext) interface{} {
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

func (v *Visitor) VisitWaitForEventSource(ctx *fql.WaitForEventSourceContext) interface{} {
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

func (v *Visitor) VisitTimeoutClauseContext(ctx *fql.TimeoutClauseContext) interface{} {
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

func (v *Visitor) VisitOptionsClause(ctx *fql.OptionsClauseContext) interface{} {
	if c := ctx.ObjectLiteral(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitForExpression(ctx *fql.ForExpressionContext) interface{} {
	v.Symbols.EnterScope()

	var distinct bool
	var returnRuleCtx antlr.RuleContext
	returnCtx := ctx.ForExpressionReturn()

	if c := returnCtx.ReturnExpression(); c != nil {
		returnRuleCtx = c
		distinct = c.Distinct() != nil
	} else if c := returnCtx.ForExpression(); c != nil {
		returnRuleCtx = c
	}

	loop := v.Loops.EnterLoop(v.loopType(ctx), v.loopKind(ctx), distinct)

	if loop.Kind == ForLoop {
		loop.Src = ctx.ForExpressionSource().Accept(v).(vm.Operand)

		if val := ctx.GetValueVariable(); val != nil {
			if txt := val.GetText(); txt != "" && txt != ignorePseudoVariable {
				loop.ValueName = txt
				loop.Value = v.Symbols.DefineVariable(txt)
			}
		}

		if ctr := ctx.GetCounterVariable(); ctr != nil {
			if txt := ctr.GetText(); txt != "" && txt != ignorePseudoVariable {
				loop.KeyName = txt
				loop.Key = v.Symbols.DefineVariable(txt)
			}
		}
	} else {
		//srcExpr := ctx.Expression()
		//
		//// Create initial value for the loop counter
		//v.Emitter.EmitA(runtime.OpWhileLoopPrep, counterReg)
		//beforeExp := v.Emitter.Size()
		//// Loop data source to iterate over
		//cond := srcExpr.Accept(v).(runtime.Operand)
		//jumpOffset = v.Emitter.Size() - beforeExp
		//
		//// jumpPlaceholder is a placeholder for the exit jump position
		//loop.Jump = v.Emitter.EmitJumpAB(runtime.OpWhileLoopNext, counterReg, cond, jumpPlaceholder)
		//
		//counterVar := ctx.GetCounterVariable().GetText()
		//
		//// declare counter variable
		//valReg := v.Symbols.DefineVariable(counterVar)
		//v.Emitter.EmitAB(runtime.OpWhileLoopValue, valReg, counterReg)
	}

	v.emitLoopBegin(loop)

	// body
	if body := ctx.AllForExpressionBody(); body != nil && len(body) > 0 {
		for _, b := range body {
			b.Accept(v)
		}
	}

	loop = v.Loops.Loop()

	// RETURN
	if loop.Type != PassThroughLoop {
		c := returnRuleCtx.(*fql.ReturnExpressionContext)
		expReg := c.Expression().Accept(v).(vm.Operand)

		v.Emitter.EmitAB(vm.OpPush, loop.Result, expReg)
	} else if returnRuleCtx != nil {
		returnRuleCtx.Accept(v)
	}

	res := v.emitLoopEnd(loop)

	v.Loops.ExitLoop()
	v.Symbols.ExitScope()

	return res
}

func (v *Visitor) VisitForExpressionSource(ctx *fql.ForExpressionSourceContext) interface{} {
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

func (v *Visitor) VisitForExpressionBody(ctx *fql.ForExpressionBodyContext) interface{} {
	if c := ctx.ForExpressionClause(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.ForExpressionStatement(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitForExpressionClause(ctx *fql.ForExpressionClauseContext) interface{} {
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

func (v *Visitor) VisitFilterClause(ctx *fql.FilterClauseContext) interface{} {
	src1 := ctx.Expression().Accept(v).(vm.Operand)
	v.Emitter.EmitJumpc(vm.OpJumpIfFalse, v.Loops.Loop().Jump, src1)

	return nil
}

func (v *Visitor) VisitLimitClause(ctx *fql.LimitClauseContext) interface{} {
	clauses := ctx.AllLimitClauseValue()

	if len(clauses) == 1 {
		return v.visitLimit(clauses[0].Accept(v).(vm.Operand))
	} else {
		v.visitOffset(clauses[0].Accept(v).(vm.Operand))
		v.visitLimit(clauses[1].Accept(v).(vm.Operand))
	}

	return nil
}

func (v *Visitor) VisitCollectClause(ctx *fql.CollectClauseContext) interface{} {
	// TODO: Undefine original loop variables
	loop := v.Loops.Loop()

	// We collect the aggregation keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	var kvKeyReg vm.Operand
	kvValReg := v.Registers.Allocate(Temp)
	var groupSelectors []fql.ICollectSelectorContext
	var isGrouping bool
	grouping := ctx.CollectGrouping()

	if grouping != nil {
		isGrouping = true
		groupSelectors = grouping.AllCollectSelector()
		kvKeyReg = v.emitCollectGroupKeySelectors(groupSelectors)
	}

	v.emitIterValue(loop, kvValReg)

	var projectionVariableName string
	// TODO: Create enum for better readability
	collectorType := 1

	// If we have a collect group variable, we need to project it
	if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
		// Projection can be either a default projection (identifier) or a custom projection (selector expression)
		if identifier := groupVar.Identifier(); identifier != nil {
			projectionVariableName = v.emitCollectDefaultGroupProjection(loop, kvValReg, identifier, groupVar.CollectGroupVariableKeeper())
		} else if selector := groupVar.CollectSelector(); selector != nil {
			projectionVariableName = v.emitCollectCustomGroupProjection(loop, kvValReg, selector)
		}

		collectorType = 3
	} else if countVar := ctx.CollectCounter(); countVar != nil {
		projectionVariableName = v.emitCollectCountProjection(loop, kvValReg, countVar)

		if isGrouping {
			collectorType = 2
		} else {
			collectorType = 0
		}
	}

	// We replace DataSet initialization with Collector initialization
	v.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, collectorType)
	v.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
	v.emitIterJumpOrClose(loop)

	// Replace source with sorted array
	v.patchJoinLoop(loop)

	// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
	if projectionVariableName != "" {
		// Now we need to expand group variables from the dataset
		v.emitIterKey(loop, kvValReg)
		v.emitIterValue(loop, v.Symbols.DefineVariable(projectionVariableName))
	} else {
		v.emitIterValue(loop, kvValReg)
	}

	//loop.ValueName = ""
	//loop.KeyName = ""
	// TODO: Reuse the Registers
	v.Registers.Free(loop.Value)
	v.Registers.Free(loop.Key)
	loop.Value = vm.NoopOperand
	loop.Key = vm.NoopOperand

	if isGrouping {
		v.emitCollectGroupKeySelectorVariables(groupSelectors, kvValReg)
	}

	return nil
}

func (v *Visitor) emitCollectGroupKeySelectors(selectors []fql.ICollectSelectorContext) vm.Operand {
	var kvKeyReg vm.Operand

	if len(selectors) > 1 {
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		selectorRegs := v.Registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := selector.Accept(v).(vm.Operand)
			v.Emitter.EmitAB(vm.OpMove, selectorRegs.Registers[i], reg)
			// Free the register after moving its value to the sequence register
			v.Registers.Free(reg)
		}

		kvKeyReg = v.Registers.Allocate(Temp)
		v.Emitter.EmitAs(vm.OpList, kvKeyReg, selectorRegs)
		v.Registers.FreeSequence(selectorRegs)
	} else {
		kvKeyReg = selectors[0].Accept(v).(vm.Operand)
	}

	return kvKeyReg
}

func (v *Visitor) emitCollectGroupKeySelectorVariables(selectors []fql.ICollectSelectorContext, kvValReg vm.Operand) {
	if len(selectors) > 1 {
		variables := make([]vm.Operand, len(selectors))

		for i, selector := range selectors {
			name := selector.Identifier().GetText()

			if variables[i] == vm.NoopOperand {
				variables[i] = v.Symbols.DefineVariable(name)
			}

			v.Emitter.EmitABC(vm.OpLoadIndex, variables[i], kvValReg, v.loadConstant(runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			v.Registers.Free(reg)
		}
	} else {
		// Get the variable name
		name := selectors[0].Identifier().GetText()
		// Define a variable for each selector
		varReg := v.Symbols.DefineVariable(name)
		// If we have a single selector, we can just move the value
		v.Emitter.EmitAB(vm.OpMove, varReg, kvValReg)
	}
}

func (v *Visitor) emitCollectDefaultGroupProjection(loop *Loop, kvValReg vm.Operand, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		seq := v.Registers.AllocateSequence(2) // Key and Value for Map

		// TODO: Review this. It's quite a questionable ArrangoDB feature of wrapping group items by a nested object
		// We will keep it for now for backward compatibility.
		v.loadConstantTo(runtime.String(loop.ValueName), seq.Registers[0]) // Map key
		v.Emitter.EmitAB(vm.OpMove, seq.Registers[1], kvValReg)            // Map value
		v.Emitter.EmitAs(vm.OpMap, kvValReg, seq)

		v.Registers.FreeSequence(seq)
	} else {
		variables := keeper.AllIdentifier()
		seq := v.Registers.AllocateSequence(len(variables) * 2)

		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			v.loadConstantTo(runtime.String(varName), seq.Registers[j])
			v.Emitter.EmitAB(vm.OpMove, seq.Registers[j+1], v.Symbols.Variable(varName))
		}

		v.Emitter.EmitAs(vm.OpMap, kvValReg, seq)
		v.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (v *Visitor) emitCollectCustomGroupProjection(_ *Loop, kvValReg vm.Operand, selector fql.ICollectSelectorContext) string {
	selectorReg := selector.Expression().Accept(v).(vm.Operand)
	v.Emitter.EmitAB(vm.OpMove, kvValReg, selectorReg)
	v.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}

func (v *Visitor) emitCollectCountProjection(_ *Loop, _ vm.Operand, selector fql.ICollectCounterContext) string {
	return selector.Identifier().GetText()
}

func (v *Visitor) VisitCollectSelector(ctx *fql.CollectSelectorContext) interface{} {
	if c := ctx.Expression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitSortClause(ctx *fql.SortClauseContext) interface{} {
	loop := v.Loops.Loop()

	// We collect the sorting conditions (keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	kvKeyReg := v.Registers.Allocate(Temp)
	clauses := ctx.AllSortClauseExpression()
	var directions []runtime.SortDirection
	isSortMany := len(clauses) > 1

	if isSortMany {
		clausesRegs := make([]vm.Operand, len(clauses))
		directions = make([]runtime.SortDirection, len(clauses))
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		keyRegs := v.Registers.AllocateSequence(len(clauses))

		for i, clause := range clauses {
			clauseReg := clause.Accept(v).(vm.Operand)
			v.Emitter.EmitAB(vm.OpMove, keyRegs.Registers[i], clauseReg)
			clausesRegs[i] = keyRegs.Registers[i]
			directions[i] = v.sortDirection(clause.SortDirection())
			// TODO: Free Registers
		}

		arrReg := v.Registers.Allocate(Temp)
		v.Emitter.EmitAs(vm.OpList, arrReg, keyRegs)
		v.Emitter.EmitAB(vm.OpMove, kvKeyReg, arrReg) // TODO: Free Registers
	} else {
		clausesReg := clauses[0].Accept(v).(vm.Operand)
		v.Emitter.EmitAB(vm.OpMove, kvKeyReg, clausesReg)
	}

	var kvValReg vm.Operand

	// In case the value is not used in the loop body, and only key is used
	if loop.ValueName != "" {
		kvValReg = loop.Value
	} else {
		// If so, we need to load it from the iterator
		kvValReg = v.Registers.Allocate(Temp)
		v.emitIterValue(loop, kvValReg)
	}

	if isSortMany {
		encoded := runtime.EncodeSortDirections(directions)
		count := len(clauses)

		v.Emitter.PatchSwapAxy(loop.ResultPos, vm.OpDataSetMultiSorter, loop.Result, encoded, count)
	} else {
		dir := v.sortDirection(clauses[0].SortDirection())
		v.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetSorter, loop.Result, int(dir))
	}

	v.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
	v.emitIterJumpOrClose(loop)

	// Replace source with the Sorter
	v.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	// Create a new loop
	v.emitLoopBegin(loop)

	return nil
}

func (v *Visitor) sortDirection(dir antlr.TerminalNode) runtime.SortDirection {
	if dir == nil {
		return runtime.SortDirectionAsc
	}

	if strings.ToLower(dir.GetText()) == "desc" {
		return runtime.SortDirectionDesc
	}

	return runtime.SortDirectionAsc
}

func (v *Visitor) VisitSortClauseExpression(ctx *fql.SortClauseExpressionContext) interface{} {
	return ctx.Expression().Accept(v).(vm.Operand)
}

func (v *Visitor) visitOffset(src1 vm.Operand) interface{} {
	state := v.Registers.Allocate(State)
	v.Emitter.EmitABx(vm.OpIterSkip, state, src1, v.Loops.Loop().Jump)

	return state
}

func (v *Visitor) visitLimit(src1 vm.Operand) interface{} {
	state := v.Registers.Allocate(State)
	v.Emitter.EmitABx(vm.OpIterLimit, state, src1, v.Loops.Loop().Jump)

	return state
}

func (v *Visitor) VisitLimitClauseValue(ctx *fql.LimitClauseValueContext) interface{} {
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

func (v *Visitor) VisitForExpressionStatement(ctx *fql.ForExpressionStatementContext) interface{} {
	if c := ctx.VariableDeclaration(); c != nil {
		return c.Accept(v)
	}

	if c := ctx.FunctionCallExpression(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitFunctionCallExpression(ctx *fql.FunctionCallExpressionContext) interface{} {
	return v.visitFunctionCall(ctx.FunctionCall().(*fql.FunctionCallContext), ctx.ErrorOperator() != nil)
}

func (v *Visitor) VisitFunctionCall(ctx *fql.FunctionCallContext) interface{} {
	return v.visitFunctionCall(ctx, false)
}

func (v *Visitor) VisitMemberExpression(ctx *fql.MemberExpressionContext) interface{} {
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
		dst = v.Registers.Allocate(Temp)

		if p.ErrorOperator() != nil {
			v.Emitter.EmitABC(vm.OpLoadPropertyOptional, dst, src1, src2)
		} else {
			v.Emitter.EmitABC(vm.OpLoadProperty, dst, src1, src2)
		}

		src1 = dst
	}

	return dst
}

func (v *Visitor) VisitRangeOperator(ctx *fql.RangeOperatorContext) interface{} {
	dst := v.Registers.Allocate(Temp)
	start := ctx.GetLeft().Accept(v).(vm.Operand)
	end := ctx.GetRight().Accept(v).(vm.Operand)

	v.Emitter.EmitABC(vm.OpRange, dst, start, end)

	return dst
}

func (v *Visitor) VisitRangeOperand(ctx *fql.RangeOperandContext) interface{} {
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

func (v *Visitor) VisitParam(ctx *fql.ParamContext) interface{} {
	name := ctx.Identifier().GetText()
	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitAB(vm.OpLoadParam, reg, v.Symbols.AddParam(name))

	return reg
}

func (v *Visitor) VisitVariableDeclaration(ctx *fql.VariableDeclarationContext) interface{} {
	name := ignorePseudoVariable

	if id := ctx.Identifier(); id != nil {
		name = id.GetText()
	} else if reserved := ctx.SafeReservedWord(); reserved != nil {
		name = reserved.GetText()
	}

	src := ctx.Expression().Accept(v).(vm.Operand)

	if name != ignorePseudoVariable {
		dest := v.Symbols.DefineVariable(name)

		if src.IsConstant() {
			tmp := v.Registers.Allocate(Temp)
			v.Emitter.EmitAB(vm.OpLoadConst, tmp, src)
			v.Emitter.EmitAB(vm.OpStoreGlobal, dest, tmp)
		} else if v.Symbols.Scope() == 0 {
			v.Emitter.EmitAB(vm.OpStoreGlobal, dest, src)
		} else {
			v.Emitter.EmitAB(vm.OpMove, dest, src)
		}

		return dest
	}

	return vm.NoopOperand
}

func (v *Visitor) VisitVariable(ctx *fql.VariableContext) interface{} {
	// Just return the register / constant index
	op := v.Symbols.Variable(ctx.GetText())

	if op.IsRegister() {
		return op
	}

	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitAB(vm.OpLoadGlobal, reg, op)

	return reg
}

func (v *Visitor) VisitArrayLiteral(ctx *fql.ArrayLiteralContext) interface{} {
	// Allocate destination register for the array
	destReg := v.Registers.Allocate(Temp)

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size := len(exps)

		if size > 0 {
			// Allocate seq for array elements
			seq := v.Registers.AllocateSequence(size)

			// Evaluate each element into seq Registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := exp.Accept(v).(vm.Operand)

				// TODO: Figure out how to remove OpMove and use Registers returned from each expression
				v.Emitter.EmitAB(vm.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.Registers.Free(srcReg)
				}
			}

			// Initialize an array
			v.Emitter.EmitAs(vm.OpList, destReg, seq)

			// Free seq Registers
			//v.Registers.FreeSequence(seq)

			return destReg
		}
	}

	// Empty array
	v.Emitter.EmitA(vm.OpList, destReg)

	return destReg
}

func (v *Visitor) VisitObjectLiteral(ctx *fql.ObjectLiteralContext) interface{} {
	dst := v.Registers.Allocate(Temp)
	assignments := ctx.AllPropertyAssignment()
	size := len(assignments)

	if size == 0 {
		v.Emitter.EmitA(vm.OpMap, dst)

		return dst
	}

	seq := v.Registers.AllocateSequence(len(assignments) * 2)

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

		v.Emitter.EmitAB(vm.OpMove, seq.Registers[regIndex], propOp)
		v.Emitter.EmitAB(vm.OpMove, seq.Registers[regIndex+1], valOp)

		// Free source register if temporary
		if propOp.IsRegister() {
			//v.Registers.Free(propOp)
		}
	}

	v.Emitter.EmitAs(vm.OpMap, dst, seq)

	return dst
}

func (v *Visitor) VisitPropertyName(ctx *fql.PropertyNameContext) interface{} {
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

func (v *Visitor) VisitComputedPropertyName(ctx *fql.ComputedPropertyNameContext) interface{} {
	return ctx.Expression().Accept(v)
}

func (v *Visitor) VisitStringLiteral(ctx *fql.StringLiteralContext) interface{} {
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

func (v *Visitor) VisitIntegerLiteral(ctx *fql.IntegerLiteralContext) interface{} {
	val, err := strconv.Atoi(ctx.GetText())

	if err != nil {
		panic(err)
	}

	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitAB(vm.OpLoadConst, reg, v.Symbols.AddConstant(runtime.NewInt(val)))

	return reg
}

func (v *Visitor) VisitFloatLiteral(ctx *fql.FloatLiteralContext) interface{} {
	val, err := strconv.ParseFloat(ctx.GetText(), 64)

	if err != nil {
		panic(err)
	}

	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitAB(vm.OpLoadConst, reg, v.Symbols.AddConstant(runtime.NewFloat(val)))

	return reg
}

func (v *Visitor) VisitBooleanLiteral(ctx *fql.BooleanLiteralContext) interface{} {
	reg := v.Registers.Allocate(Temp)

	switch strings.ToLower(ctx.GetText()) {
	case "true":
		v.Emitter.EmitAB(vm.OpLoadBool, reg, 1)
	case "false":
		v.Emitter.EmitAB(vm.OpLoadBool, reg, 0)
	default:
		panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
	}

	return reg
}

func (v *Visitor) VisitNoneLiteral(_ *fql.NoneLiteralContext) interface{} {
	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitA(vm.OpLoadNone, reg)

	return reg
}

func (v *Visitor) VisitLiteral(ctx *fql.LiteralContext) interface{} {
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

func (v *Visitor) VisitReturnExpression(ctx *fql.ReturnExpressionContext) interface{} {
	valReg := ctx.Expression().Accept(v).(vm.Operand)

	if valReg.IsConstant() {
		v.Emitter.EmitAB(vm.OpLoadGlobal, vm.NoopOperand, valReg)
	} else {
		v.Emitter.EmitAB(vm.OpMove, vm.NoopOperand, valReg)
	}

	v.Emitter.Emit(vm.OpReturn)

	return vm.NoopOperand
}

func (v *Visitor) VisitExpression(ctx *fql.ExpressionContext) interface{} {
	if uo := ctx.UnaryOperator(); uo != nil {
		src := ctx.GetRight().Accept(v).(vm.Operand)
		dst := v.Registers.Allocate(Temp)

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
		v.Emitter.EmitAB(op, dst, src)

		return dst
	}

	if op := ctx.LogicalAndOperator(); op != nil {
		dst := v.Registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(vm.Operand)
		v.Emitter.EmitAB(vm.OpMove, dst, left)
		// Test if left is false and jump to the end
		end := v.Emitter.EmitJumpc(vm.OpJumpIfFalse, jumpPlaceholder, dst)
		// If left is true, execute right expression
		right := ctx.GetRight().Accept(v).(vm.Operand)
		// And move the result to the destination register
		v.Emitter.EmitAB(vm.OpMove, dst, right)
		v.Emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.LogicalOrOperator(); op != nil {
		dst := v.Registers.Allocate(Temp)
		// Execute left expression
		left := ctx.GetLeft().Accept(v).(vm.Operand)
		// Move the result to the destination register
		v.Emitter.EmitAB(vm.OpMove, dst, left)
		// Test if left is true and jump to the end
		end := v.Emitter.EmitJumpc(vm.OpJumpIfTrue, jumpPlaceholder, dst)
		// If left is false, execute right expression
		right := ctx.GetRight().Accept(v).(vm.Operand)
		// And move the result to the destination register
		v.Emitter.EmitAB(vm.OpMove, dst, right)
		v.Emitter.PatchJumpNext(end)

		return dst
	}

	if op := ctx.GetTernaryOperator(); op != nil {
		dst := v.Registers.Allocate(Temp)

		// Compile condition and put result in dst
		condReg := ctx.GetCondition().Accept(v).(vm.Operand)
		v.Emitter.EmitAB(vm.OpMove, dst, condReg)

		// If condition was temporary, free it
		if condReg.IsRegister() {
			//v.Registers.Free(condReg)
		}

		// Jump to 'false' branch if condition is false
		otherwise := v.Emitter.EmitJumpc(vm.OpJumpIfFalse, jumpPlaceholder, dst)

		// True branch
		if onTrue := ctx.GetOnTrue(); onTrue != nil {
			trueReg := onTrue.Accept(v).(vm.Operand)
			v.Emitter.EmitAB(vm.OpMove, dst, trueReg)

			// Free temporary register if needed
			if trueReg.IsRegister() {
				//v.Registers.Free(trueReg)
			}
		}

		// Jump over false branch
		end := v.Emitter.EmitJump(vm.OpJump, jumpPlaceholder)
		v.Emitter.PatchJumpNext(otherwise)

		// False branch
		if onFalse := ctx.GetOnFalse(); onFalse != nil {
			falseReg := onFalse.Accept(v).(vm.Operand)
			v.Emitter.EmitAB(vm.OpMove, dst, falseReg)

			// Free temporary register if needed
			if falseReg.IsRegister() {
				//v.Registers.Free(falseReg)
			}
		}

		v.Emitter.PatchJumpNext(end)

		return dst
	}

	if c := ctx.Predicate(); c != nil {
		return c.Accept(v)
	}

	panic(runtime.Error(ErrUnexpectedToken, ctx.GetText()))
}

func (v *Visitor) VisitPredicate(ctx *fql.PredicateContext) interface{} {
	if c := ctx.ExpressionAtom(); c != nil {
		startCatch := v.Emitter.Size()
		reg := c.Accept(v)

		if c.ErrorOperator() != nil {
			jump := -1
			endCatch := v.Emitter.Size()

			if c.ForExpression() != nil {
				// We jump back to finalize the loop before exiting
				jump = endCatch - 1
			}

			v.CatchTable = append(v.CatchTable, [3]int{startCatch, endCatch, jump})
		}

		return reg
	}

	var opcode vm.Opcode
	dest := v.Registers.Allocate(Temp)
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

	v.Emitter.EmitABC(opcode, dest, left, right)

	return dest
}

func (v *Visitor) VisitExpressionAtom(ctx *fql.ExpressionAtomContext) interface{} {
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
		dst := v.Registers.Allocate(Temp)

		if opcode == vm.OpRegexpPositive || opcode == vm.OpRegexpNegative {
			if regRight.IsConstant() {
				val := v.Symbols.Constant(regRight)

				// Verify that the expression is a valid regular expression
				regexp.MustCompile(val.String())
			}
		}

		v.Emitter.EmitABC(opcode, dst, regLeft, regRight)

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

func (v *Visitor) visitFunctionCall(ctx *fql.FunctionCallContext, protected bool) interface{} {
	var size int
	var seq *RegisterSequence

	if list := ctx.ArgumentList(); list != nil {
		// Get all array element expressions
		exps := list.(*fql.ArgumentListContext).AllExpression()
		size = len(exps)

		if size > 0 {
			// Allocate seq for function arguments
			seq = v.Registers.AllocateSequence(size)

			// Evaluate each element into seq Registers
			for i, exp := range exps {
				// Compile expression and move to seq register
				srcReg := exp.Accept(v).(vm.Operand)

				// TODO: Figure out how to remove OpMove and use Registers returned from each expression
				v.Emitter.EmitAB(vm.OpMove, seq.Registers[i], srcReg)

				// Free source register if temporary
				if srcReg.IsRegister() {
					//v.Registers.Free(srcReg)
				}
			}
		}
	}

	name := v.functionName(ctx)

	switch name {
	case runtimeLength:
		dst := v.Registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeLength+": expected 1 argument"))
		}

		v.Emitter.EmitAB(vm.OpLength, dst, seq.Registers[0])

		return dst
	case runtimeTypename:
		dst := v.Registers.Allocate(Temp)

		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeTypename+": expected 1 argument"))
		}

		v.Emitter.EmitAB(vm.OpType, dst, seq.Registers[0])

		return dst
	case runtimeWait:
		if seq == nil || len(seq.Registers) != 1 {
			panic(runtime.Error(runtime.ErrInvalidArgument, runtimeWait+": expected 1 argument"))
		}

		v.Emitter.EmitA(vm.OpSleep, seq.Registers[0])

		return seq.Registers[0]
	default:
		nameAndDest := v.loadConstant(v.functionName(ctx))

		if !protected {
			v.Emitter.EmitAs(vm.OpCall, nameAndDest, seq)
		} else {
			v.Emitter.EmitAs(vm.OpProtectedCall, nameAndDest, seq)
		}

		return nameAndDest
	}
}

func (v *Visitor) functionName(ctx *fql.FunctionCallContext) runtime.String {
	var name string
	funcNS := ctx.Namespace()

	if funcNS != nil {
		name += funcNS.GetText()
	}

	name += ctx.FunctionName().GetText()

	return runtime.NewString(strings.ToUpper(name))
}

// emitIterValue emits an instruction to get the value from the iterator
func (v *Visitor) emitLoopBegin(loop *Loop) {
	if loop.Allocate {
		v.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
		loop.ResultPos = v.Emitter.Size() - 1
	}

	loop.Iterator = v.Registers.Allocate(State)

	if loop.Kind == ForLoop {
		v.Emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Src)
		// jumpPlaceholder is a placeholder for the exit jump position
		loop.Jump = v.Emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, loop.Iterator)

		if loop.Value != vm.NoopOperand {
			v.Emitter.EmitAB(vm.OpIterValue, loop.Value, loop.Iterator)
		}

		if loop.Key != vm.NoopOperand {
			v.Emitter.EmitAB(vm.OpIterKey, loop.Key, loop.Iterator)
		}
	} else {
		//counterReg := v.Registers.Allocate(Storage)
		// TODO: Set JumpOffset here
	}
}

// emitIterValue emits an instruction to get the value from the iterator
func (v *Visitor) emitIterValue(loop *Loop, reg vm.Operand) {
	v.Emitter.EmitAB(vm.OpIterValue, reg, loop.Iterator)
}

// emitIterKey emits an instruction to get the key from the iterator
func (v *Visitor) emitIterKey(loop *Loop, reg vm.Operand) {
	v.Emitter.EmitAB(vm.OpIterKey, reg, loop.Iterator)
}

// emitIterJumpOrClose emits an instruction to jump to the end of the loop or close the iterator
func (v *Visitor) emitIterJumpOrClose(loop *Loop) {
	v.Emitter.EmitJump(vm.OpJump, loop.Jump-loop.JumpOffset)
	v.Emitter.EmitA(vm.OpClose, loop.Iterator)

	if loop.Kind == ForLoop {
		v.Emitter.PatchJump(loop.Jump)
	} else {
		v.Emitter.PatchJumpAB(loop.Jump)
	}
}

// patchJoinLoop replaces the source of the loop with a modified dataset
func (v *Visitor) patchJoinLoop(loop *Loop) {
	// Replace source with sorted array
	v.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)

	v.Symbols.ExitScope()
	v.Symbols.EnterScope()

	// Create new for loop
	v.emitLoopBegin(loop)
}

func (v *Visitor) emitLoopEnd(loop *Loop) vm.Operand {
	v.Emitter.EmitJump(vm.OpJump, loop.Jump-loop.JumpOffset)

	// TODO: Do not allocate for pass-through Loops
	dst := v.Registers.Allocate(Temp)

	if loop.Allocate {
		// TODO: Reuse the dsReg register
		v.Emitter.EmitA(vm.OpClose, loop.Iterator)
		v.Emitter.EmitAB(vm.OpMove, dst, loop.Result)

		if loop.Kind == ForLoop {
			v.Emitter.PatchJump(loop.Jump)
		} else {
			v.Emitter.PatchJumpAB(loop.Jump)
		}
	} else {
		if loop.Kind == ForLoop {
			v.Emitter.PatchJumpNext(loop.Jump)
		} else {
			v.Emitter.PatchJumpNextAB(loop.Jump)
		}
	}

	return dst
}

func (v *Visitor) loopType(ctx *fql.ForExpressionContext) LoopType {
	if c := ctx.ForExpressionReturn().ForExpression(); c == nil {
		return NormalLoop
	}

	return PassThroughLoop
}

func (v *Visitor) loopKind(ctx *fql.ForExpressionContext) LoopKind {
	if ctx.While() == nil {
		return ForLoop
	}

	if ctx.Do() == nil {
		return WhileLoop
	}

	return DoWhileLoop
}

func (v *Visitor) loadConstant(constant runtime.Value) vm.Operand {
	reg := v.Registers.Allocate(Temp)
	v.Emitter.EmitAB(vm.OpLoadConst, reg, v.Symbols.AddConstant(constant))
	return reg
}

func (v *Visitor) loadConstantTo(constant runtime.Value, reg vm.Operand) {
	v.Emitter.EmitAB(vm.OpLoadConst, reg, v.Symbols.AddConstant(constant))
}
