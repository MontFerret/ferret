package compiler

import (
	"github.com/MontFerret/ferret/pkg/compiler/internal"
	"github.com/MontFerret/ferret/pkg/parser/fql"
)

type Visitor struct {
	*fql.BaseFqlParserVisitor

	Ctx *internal.FuncContext
	Err error
	Src string
}

func NewVisitor(src string) *Visitor {
	v := new(Visitor)
	v.BaseFqlParserVisitor = new(fql.BaseFqlParserVisitor)
	v.Ctx = internal.NewFuncContext()

	v.Src = src

	return v
}

func (v *Visitor) VisitProgram(ctx *fql.ProgramContext) interface{} {
	for _, head := range ctx.AllHead() {
		v.VisitHead(head.(*fql.HeadContext))
	}

	v.Ctx.StmtCompiler.Compile(ctx.Body())

	return nil
}

func (v *Visitor) VisitHead(_ *fql.HeadContext) interface{} {
	return nil
}

//func (v *Visitor) VisitCollectClause(Ctx *fql.CollectClauseContext) interface{} {
//	// TODO: Undefine original loop variables
//	loop := v.Loops.Current()
//
//	// We collect the aggregation keys
//	// And wrap each loop element by a KeyValuePair
//	// Where a key is either a single value or a list of values
//	// These KeyValuePairs are then added to the dataset
//	var kvKeyReg, kvValReg vm.Operand
//	var groupSelectors []fql.ICollectSelectorContext
//	var isGrouping bool
//	grouping := Ctx.CollectGrouping()
//	counter := Ctx.CollectCounter()
//	aggregator := Ctx.CollectAggregator()
//
//	isCollecting := grouping != nil || counter != nil
//
//	if isCollecting {
//		if grouping != nil {
//			isGrouping = true
//			groupSelectors = grouping.AllCollectSelector()
//			kvKeyReg = v.emitCollectGroupKeySelectors(groupSelectors)
//		}
//
//		kvValReg = v.Registers.Allocate(Temp)
//		v.emitIterValue(loop, kvValReg)
//
//		var projectionVariableName string
//		collectorType := CollectorTypeKey
//
//		// If we have a collect group variable, we need to project it
//		if groupVar := Ctx.CollectGroupVariable(); groupVar != nil {
//			// Projection can be either a default projection (identifier) or a custom projection (selector expression)
//			if identifier := groupVar.Identifier(); identifier != nil {
//				projectionVariableName = v.emitCollectDefaultGroupProjection(loop, kvValReg, identifier, groupVar.CollectGroupVariableKeeper())
//			} else if selector := groupVar.CollectSelector(); selector != nil {
//				projectionVariableName = v.emitCollectCustomGroupProjection(loop, kvValReg, selector)
//			}
//
//			collectorType = CollectorTypeKeyGroup
//		} else if counter != nil {
//			projectionVariableName = v.emitCollectCountProjection(loop, kvValReg, counter)
//
//			if isGrouping {
//				collectorType = CollectorTypeKeyCounter
//			} else {
//				collectorType = CollectorTypeCounter
//			}
//		}
//
//		// If we use aggregators, we need to collect group items by key
//		if aggregator != nil && collectorType != CollectorTypeKeyGroup {
//			// We need to patch the loop result to be a collector
//			collectorType = CollectorTypeKeyGroup
//		}
//
//		// We replace DataSet initialization with Collector initialization
//		v.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(collectorType))
//		v.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
//		v.emitIterJumpOrClose(loop)
//
//		// Replace the source with the collector
//		v.emitPatchLoop(loop)
//
//		// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
//		if projectionVariableName != "" {
//			// Now we need to expand group variables from the dataset
//			v.emitIterKey(loop, kvValReg)
//			v.emitIterValue(loop, v.Symbols.DeclareLocal(projectionVariableName))
//		} else {
//			v.emitIterKey(loop, kvKeyReg)
//			v.emitIterValue(loop, kvValReg)
//		}
//	}
//
//	// Aggregation loop
//	if aggregator != nil {
//		v.emitCollectAggregator(aggregator, loop, isCollecting)
//	}
//
//	// TODO: Reuse the Registers
//	v.Registers.Free(loop.Value)
//	v.Registers.Free(loop.Key)
//	loop.Value = vm.NoopOperand
//	loop.Key = vm.NoopOperand
//
//	if isCollecting && isGrouping {
//		// Now we are defining new variables for the group selectors
//		v.emitCollectGroupKeySelectorVariables(groupSelectors, kvKeyReg, kvValReg, aggregator != nil)
//	}
//
//	return nil
//}
//
//func (v *Visitor) emitCollectAggregator(c fql.ICollectAggregatorContext, parentLoop *Loop, isCollected bool) {
//	var accums []vm.Operand
//	var loop *Loop
//	selectors := c.AllCollectAggregateSelector()
//
//	// If data is collected, we need to allocate a temporary accumulators to store aggregation results
//	if isCollected {
//		// First of all, we allocate registers for accumulators
//		accums = make([]vm.Operand, len(selectors))
//
//		// We need to allocate a register for each accumulator
//		for i := 0; i < len(selectors); i++ {
//			reg := v.Registers.Allocate(Temp)
//			accums[i] = reg
//			// TODO: Select persistent List type, we do not know how many items we will have
//			v.Emitter.EmitA(vm.OpList, reg)
//		}
//
//		loop = v.Loops.NewLoop(TemporalLoop, ForLoop, false)
//
//		// Now we iterate over the grouped items
//		v.emitIterValue(parentLoop, loop.Iterator)
//		// We just re-use the same register
//		v.Emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Iterator)
//		// jumpPlaceholder is a placeholder for the exit aggrIterJump position
//		loop.Jump = v.Emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, loop.Iterator)
//		loop.ValueName = parentLoop.ValueName
//	} else {
//		loop = parentLoop
//		// Otherwise, we create a custom collector for aggregators
//		v.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(CollectorTypeKeyGroup))
//	}
//
//	// Store upper scope for aggregators
//	//mainScope := v.Symbols.Scope()
//	// Nested scope for aggregators
//	v.Symbols.EnterScope()
//
//	aggrIterVal := v.Symbols.DeclareLocal(loop.ValueName)
//	v.Emitter.EmitAB(vm.OpIterValue, aggrIterVal, loop.Iterator)
//
//	// Now we add value selectors to the accumulators
//	for i := 0; i < len(selectors); i++ {
//		selector := selectors[i]
//		fcx := selector.FunctionCallExpression()
//		args := fcx.FunctionCall().ArgumentList().AllExpression()
//
//		if len(args) == 0 {
//			// TODO: Better error handling
//			panic("No arguments provided for the function call in the aggregate selector")
//		}
//
//		if len(args) > 1 {
//			// TODO: Better error handling
//			panic("Too many arguments")
//		}
//
//		resultReg := args[0].Accept(v).(vm.Operand)
//
//		if isCollected {
//			v.Emitter.EmitAB(vm.OpPush, accums[i], resultReg)
//		} else {
//			aggrKeyName := selector.Identifier().GetText()
//			aggrKeyReg := v.loadConstant(runtime.String(aggrKeyName))
//			v.Emitter.EmitABC(vm.OpPushKV, loop.Result, aggrKeyReg, resultReg)
//			v.Registers.Free(aggrKeyReg)
//		}
//
//		v.Registers.Free(resultReg)
//	}
//
//	// Now we can iterate over the grouped items
//	v.emitIterJumpOrClose(loop)
//
//	// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
//	// And define variables for each accumulator result
//	if isCollected {
//		for i, selector := range selectors {
//			fcx := selector.FunctionCallExpression()
//			// We won't make any checks here, as we already did it before
//			selectorVarName := selector.Identifier().GetText()
//
//			// We execute the function call with the accumulator as an argument
//			accum := accums[i]
//			result := v.emitFunctionCall(fcx.FunctionCall(), fcx.ErrorOperator() != nil, RegisterSequence{accum})
//
//			// We define the variable for the selector result in the upper scope
//			// Since this temporary scope is only for aggregators and will be closed after the aggregation
//			varReg := v.Symbols.DeclareLocal(selectorVarName)
//			v.Emitter.EmitAB(vm.OpMove, varReg, result)
//			v.Registers.Free(result)
//		}
//
//		v.Loops.Pop()
//		// Now close the aggregators scope
//		v.Symbols.ExitScope()
//	} else {
//		// Now close the aggregators scope
//		v.Symbols.ExitScope()
//
//		parentLoop.ValueName = ""
//		parentLoop.KeyName = ""
//
//		// Since we we in the middle of the loop, we need to patch the loop result
//		// Now we just create a range with 1 item to push the aggregated values to the dataset
//		// Replace source with sorted array
//		zero := v.loadConstant(runtime.Int(0))
//		one := v.loadConstant(runtime.Int(1))
//		aggregator := v.Registers.Allocate(Temp)
//		v.Emitter.EmitAB(vm.OpMove, aggregator, loop.Result)
//		v.Symbols.ExitScope()
//
//		v.Symbols.EnterScope()
//
//		// Create new for loop
//		v.Emitter.EmitABC(vm.OpRange, loop.Src, zero, one)
//		v.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
//
//		// In case of non-collected aggregators, we just iterate over the grouped items
//		// Retrieve the grouped values by key, execute aggregation funcs and assign variable names to the results
//		for _, selector := range selectors {
//			fcx := selector.FunctionCallExpression()
//			// We won't make any checks here, as we already did it before
//			selectorVarName := selector.Identifier().GetText()
//
//			// We execute the function call with the accumulator as an argument
//			key := v.loadConstant(runtime.String(selectorVarName))
//			value := v.Registers.Allocate(Temp)
//			v.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)
//
//			result := v.emitFunctionCall(fcx.FunctionCall(), fcx.ErrorOperator() != nil, RegisterSequence{value})
//
//			// We define the variable for the selector result in the upper scope
//			// Since this temporary scope is only for aggregators and will be closed after the aggregation
//			varReg := v.Symbols.DeclareLocal(selectorVarName)
//			v.Emitter.EmitAB(vm.OpMove, varReg, result)
//			v.Registers.Free(result)
//			v.Registers.Free(value)
//			v.Registers.Free(key)
//		}
//
//		v.Registers.Free(aggregator)
//	}
//
//	// Free the registers for accumulators
//	for _, reg := range accums {
//		v.Registers.Free(reg)
//	}
//
//	// Free the register for the iterator value
//	v.Registers.Free(aggrIterVal)
//}
//
//func (v *Visitor) emitCollectGroupKeySelectors(selectors []fql.ICollectSelectorContext) vm.Operand {
//	if len(selectors) == 0 {
//		return vm.NoopOperand
//	}
//
//	var kvKeyReg vm.Operand
//
//	if len(selectors) > 1 {
//		// We create a sequence of Registers for the clauses
//		// To pack them into an array
//		selectorRegs := v.Registers.AllocateSequence(len(selectors))
//
//		for i, selector := range selectors {
//			reg := selector.Accept(v).(vm.Operand)
//			v.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
//			// Free the register after moving its value to the sequence register
//			v.Registers.Free(reg)
//		}
//
//		kvKeyReg = v.Registers.Allocate(Temp)
//		v.Emitter.EmitAs(vm.OpList, kvKeyReg, selectorRegs)
//		v.Registers.FreeSequence(selectorRegs)
//	} else {
//		kvKeyReg = selectors[0].Accept(v).(vm.Operand)
//	}
//
//	return kvKeyReg
//}
//
//func (v *Visitor) emitCollectGroupKeySelectorVariables(selectors []fql.ICollectSelectorContext, kvKeyReg, kvValReg vm.Operand, isAggregation bool) {
//	if len(selectors) > 1 {
//		variables := make([]vm.Operand, len(selectors))
//
//		for i, selector := range selectors {
//			name := selector.Identifier().GetText()
//
//			if variables[i] == vm.NoopOperand {
//				variables[i] = v.Symbols.DeclareLocal(name)
//			}
//
//			reg := kvValReg
//
//			if isAggregation {
//				reg = kvKeyReg
//			}
//
//			v.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, v.loadConstant(runtime.Int(i)))
//		}
//
//		// Free the register after moving its value to the variable
//		for _, reg := range variables {
//			v.Registers.Free(reg)
//		}
//	} else {
//		// Get the variable name
//		name := selectors[0].Identifier().GetText()
//		// Define a variable for each selector
//		varReg := v.Symbols.DeclareLocal(name)
//
//		reg := kvValReg
//
//		if isAggregation {
//			reg = kvKeyReg
//		}
//
//		// If we have a single selector, we can just move the value
//		v.Emitter.EmitAB(vm.OpMove, varReg, reg)
//	}
//}
//
//func (v *Visitor) emitCollectDefaultGroupProjection(loop *Loop, kvValReg vm.Operand, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
//	if keeper == nil {
//		seq := v.Registers.AllocateSequence(2) // Key and Value for Map
//
//		// TODO: Review this. It's quite a questionable ArrangoDB feature of wrapping group items by a nested object
//		// We will keep it for now for backward compatibility.
//		v.loadConstantTo(runtime.String(loop.ValueName), seq[0]) // Map key
//		v.Emitter.EmitAB(vm.OpMove, seq[1], kvValReg)            // Map value
//		v.Emitter.EmitAs(vm.OpMap, kvValReg, seq)
//
//		v.Registers.FreeSequence(seq)
//	} else {
//		variables := keeper.AllIdentifier()
//		seq := v.Registers.AllocateSequence(len(variables) * 2)
//
//		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
//			varName := variables[i].GetText()
//			v.loadConstantTo(runtime.String(varName), seq[j])
//
//			variable, _, found := v.Symbols.Resolve(varName)
//
//			if !found {
//				panic("variable not found: " + varName)
//			}
//
//			v.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
//		}
//
//		v.Emitter.EmitAs(vm.OpMap, kvValReg, seq)
//		v.Registers.FreeSequence(seq)
//	}
//
//	return identifier.GetText()
//}
//
//func (v *Visitor) emitCollectCustomGroupProjection(_ *Loop, kvValReg vm.Operand, selector fql.ICollectSelectorContext) string {
//	selectorReg := selector.Expression().Accept(v).(vm.Operand)
//	v.Emitter.EmitAB(vm.OpMove, kvValReg, selectorReg)
//	v.Registers.Free(selectorReg)
//
//	return selector.Identifier().GetText()
//}
//
//func (v *Visitor) emitCollectCountProjection(_ *Loop, _ vm.Operand, selector fql.ICollectCounterContext) string {
//	return selector.Identifier().GetText()
//}
//
//func (v *Visitor) VisitCollectSelector(Ctx *fql.CollectSelectorContext) interface{} {
//	if c := Ctx.Expression(); c != nil {
//		return c.Accept(v)
//	}
//
//	panic(runtime.Error(ErrUnexpectedToken, Ctx.GetText()))
//}
//
//func (v *Visitor) VisitForExpressionStatement(Ctx *fql.ForExpressionStatementContext) interface{} {
//	if c := Ctx.VariableDeclaration(); c != nil {
//		return c.Accept(v)
//	}
//
//	if c := Ctx.FunctionCallExpression(); c != nil {
//		return c.Accept(v)
//	}
//
//	panic(runtime.Error(ErrUnexpectedToken, Ctx.GetText()))
//}
//
//func (v *Visitor) VisitExpression(Ctx *fql.ExpressionContext) interface{} {
//	return v.Ctx.ExprCompiler.Compile(Ctx)
//}
//
//// emitIterValue emits an instruction to get the value from the iterator
//func (v *Visitor) emitLoopBegin(loop *Loop) {
//	if loop.Allocate {
//		v.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)
//		loop.ResultPos = v.Emitter.Size() - 1
//	}
//
//	loop.Iterator = v.Registers.Allocate(State)
//
//	if loop.Kind == ForLoop {
//		v.Emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Src)
//		// jumpPlaceholder is a placeholder for the exit jump position
//		loop.Jump = v.Emitter.EmitJumpc(vm.OpIterNext, jumpPlaceholder, loop.Iterator)
//
//		if loop.Value != vm.NoopOperand {
//			v.Emitter.EmitAB(vm.OpIterValue, loop.Value, loop.Iterator)
//		}
//
//		if loop.Key != vm.NoopOperand {
//			v.Emitter.EmitAB(vm.OpIterKey, loop.Key, loop.Iterator)
//		}
//	} else {
//		//counterReg := v.Registers.Allocate(Storage)
//		// TODO: Set JumpOffset here
//	}
//}
//
//// emitIterValue emits an instruction to get the value from the iterator
//func (v *Visitor) emitIterValue(loop *Loop, reg vm.Operand) {
//	v.Emitter.EmitAB(vm.OpIterValue, reg, loop.Iterator)
//}
//
//// emitIterKey emits an instruction to get the key from the iterator
//func (v *Visitor) emitIterKey(loop *Loop, reg vm.Operand) {
//	v.Emitter.EmitAB(vm.OpIterKey, reg, loop.Iterator)
//}
//
//// emitIterJumpOrClose emits an instruction to jump to the end of the loop or close the iterator
//func (v *Visitor) emitIterJumpOrClose(loop *Loop) {
//	v.Emitter.EmitJump(loop.Jump - loop.JumpOffset)
//	v.Emitter.EmitA(vm.OpClose, loop.Iterator)
//
//	if loop.Kind == ForLoop {
//		v.Emitter.PatchJump(loop.Jump)
//	} else {
//		v.Emitter.PatchJumpAB(loop.Jump)
//	}
//}
//
//// emitPatchLoop replaces the source of the loop with a modified dataset
//func (v *Visitor) emitPatchLoop(loop *Loop) {
//	// Replace source with sorted array
//	v.Emitter.EmitAB(vm.OpMove, loop.Src, loop.Result)
//
//	v.Symbols.ExitScope()
//	v.Symbols.EnterScope()
//
//	// Create new for loop
//	v.emitLoopBegin(loop)
//}
//
//func (v *Visitor) emitLoopEnd(loop *Loop) vm.Operand {
//	v.Emitter.EmitJump(loop.Jump - loop.JumpOffset)
//
//	// TODO: Do not allocate for pass-through Loops
//	dst := v.Registers.Allocate(Temp)
//
//	if loop.Allocate {
//		// TODO: Reuse the dsReg register
//		v.Emitter.EmitA(vm.OpClose, loop.Iterator)
//		v.Emitter.EmitAB(vm.OpMove, dst, loop.Result)
//
//		if loop.Kind == ForLoop {
//			v.Emitter.PatchJump(loop.Jump)
//		} else {
//			v.Emitter.PatchJumpAB(loop.Jump)
//		}
//	} else {
//		if loop.Kind == ForLoop {
//			v.Emitter.PatchJumpNext(loop.Jump)
//		} else {
//			v.Emitter.PatchJumpNextAB(loop.Jump)
//		}
//	}
//
//	return dst
//}
//
//func (v *Visitor) loadConstant(constant runtime.Value) vm.Operand {
//	return loadConstant(v.Ctx, constant)
//}
//
//func (v *Visitor) loadConstantTo(constant runtime.Value, reg vm.Operand) {
//	v.Emitter.EmitAB(vm.OpLoadConst, reg, v.Symbols.AddConstant(constant))
//}
