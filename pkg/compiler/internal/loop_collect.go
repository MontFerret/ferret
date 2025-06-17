package internal

import (
	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/pkg/parser/fql"
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type CollectCompiler struct {
	ctx *CompilerContext
}

func NewCollectCompiler(ctx *CompilerContext) *CollectCompiler {
	return &CollectCompiler{ctx: ctx}
}

func (cc *CollectCompiler) Compile(ctx fql.ICollectClauseContext) {
	// TODO: Undefine original loop variables
	loop := cc.ctx.Loops.Current()

	// We collect the aggregation keys
	// And wrap each loop element by a KeyValuePair
	// Where a key is either a single value or a list of values
	// These KeyValuePairs are then added to the dataset
	var kvKeyReg, kvValReg vm.Operand
	var groupSelectors []fql.ICollectSelectorContext
	var isGrouping bool
	grouping := ctx.CollectGrouping()
	counter := ctx.CollectCounter()
	aggregator := ctx.CollectAggregator()

	isCollecting := grouping != nil || counter != nil

	if isCollecting {
		if grouping != nil {
			isGrouping = true
			groupSelectors = grouping.AllCollectSelector()
			kvKeyReg = cc.compileCollectGroupKeySelectors(groupSelectors)
		}

		kvValReg = cc.ctx.Registers.Allocate(core.Temp)
		loop.EmitValue(kvValReg, cc.ctx.Emitter)

		var projectionVariableName string
		collectorType := core.CollectorTypeKey

		// If we have a collect group variable, we need to project it
		if groupVar := ctx.CollectGroupVariable(); groupVar != nil {
			// Projection can be either a default projection (identifier) or a custom projection (selector expression)
			if identifier := groupVar.Identifier(); identifier != nil {
				projectionVariableName = cc.compileDefaultGroupProjection(loop, kvValReg, identifier, groupVar.CollectGroupVariableKeeper())
			} else if selector := groupVar.CollectSelector(); selector != nil {
				projectionVariableName = cc.compileCustomGroupProjection(loop, kvValReg, selector)
			}

			collectorType = core.CollectorTypeKeyGroup
		} else if counter != nil {
			projectionVariableName = counter.Identifier().GetText()

			if isGrouping {
				collectorType = core.CollectorTypeKeyCounter
			} else {
				collectorType = core.CollectorTypeCounter
			}
		}

		// If we use aggregators, we need to collect group items by key
		if aggregator != nil && collectorType != core.CollectorTypeKeyGroup {
			// We need to patch the loop result to be a collector
			collectorType = core.CollectorTypeKeyGroup
		}

		// We replace DataSet initialization with Collector initialization
		cc.ctx.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(collectorType))
		cc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Result, kvKeyReg, kvValReg)
		loop.EmitFinalization(cc.ctx.Emitter)

		cc.ctx.Emitter.EmitMove(loop.Src, loop.Result)

		cc.ctx.Registers.Free(loop.Value)
		cc.ctx.Registers.Free(loop.Key)
		loop.Value = kvValReg
		loop.Key = vm.NoopOperand

		// If the projection is used, we allocate a new register for the variable and put the iterator's value into it
		if projectionVariableName != "" {
			// Now we need to expand group variables from the dataset
			loop.DeclareValueVar(projectionVariableName, cc.ctx.Symbols)
			cc.ctx.LoopCompiler.EmitLoopBegin(loop)
			loop.EmitKey(kvValReg, cc.ctx.Emitter)
			loop.EmitValue(cc.ctx.Symbols.DeclareLocal(projectionVariableName), cc.ctx.Emitter)
		} else {
			cc.ctx.LoopCompiler.EmitLoopBegin(loop)
			//
			//loop.EmitKey(kvKeyReg, cc.ctx.Emitter)
			//loop.EmitValue(kvValReg, cc.ctx.Emitter)
		}
	}

	// Aggregation loop
	if aggregator != nil {
		cc.compileAggregator(aggregator, loop, isCollecting)
	}

	if isCollecting && isGrouping {
		// Now we are defining new variables for the group selectors
		cc.compileCollectGroupKeySelectorVariables(groupSelectors, kvKeyReg, kvValReg, aggregator != nil)
	}
}

func (cc *CollectCompiler) compileAggregator(c fql.ICollectAggregatorContext, parentLoop *core.Loop, isCollected bool) {
	var accums []vm.Operand
	var loop *core.Loop
	selectors := c.AllCollectAggregateSelector()

	// If data is collected, we need to allocate a temporary accumulators to store aggregation results
	if isCollected {
		// First of all, we allocate registers for accumulators
		accums = make([]vm.Operand, len(selectors))

		// We need to allocate a register for each accumulator
		for i := 0; i < len(selectors); i++ {
			reg := cc.ctx.Registers.Allocate(core.Temp)
			accums[i] = reg
			// TODO: Select persistent List type, we do not know how many items we will have
			cc.ctx.Emitter.EmitA(vm.OpList, reg)
		}

		loop = cc.ctx.Loops.NewLoop(core.TemporalLoop, core.ForLoop, false)

		// Now we iterate over the grouped items
		parentLoop.EmitValue(loop.Iterator, cc.ctx.Emitter)
		// We just re-use the same register
		cc.ctx.Emitter.EmitAB(vm.OpIter, loop.Iterator, loop.Iterator)
		// jumpPlaceholder is a placeholder for the exit aggrIterJump position
		loop.Jump = cc.ctx.Emitter.EmitJumpc(vm.OpIterNext, core.JumpPlaceholder, loop.Iterator)
		loop.ValueName = parentLoop.ValueName
	} else {
		loop = parentLoop
		// Otherwise, we create a custom collector for aggregators
		cc.ctx.Emitter.PatchSwapAx(loop.ResultPos, vm.OpDataSetCollector, loop.Result, int(core.CollectorTypeKeyGroup))
	}

	// Store upper scope for aggregators
	//mainScope := cc.ctx.Symbols.Scope()
	// Nested scope for aggregators
	cc.ctx.Symbols.EnterScope()

	aggrIterVal := cc.ctx.Symbols.DeclareLocal(loop.ValueName)
	cc.ctx.Emitter.EmitAB(vm.OpIterValue, aggrIterVal, loop.Iterator)

	// Now we add value selectors to the accumulators
	for i := 0; i < len(selectors); i++ {
		selector := selectors[i]
		fcx := selector.FunctionCallExpression()
		args := cc.ctx.ExprCompiler.CompileArgumentList(fcx.FunctionCall().ArgumentList())

		if len(args) == 0 {
			// TODO: Better error handling
			panic("No arguments provided for the function call in the aggregate selector")
		}

		if len(args) > 1 {
			// TODO: Better error handling
			panic("Too many arguments")
		}

		resultReg := args[0]

		if isCollected {
			cc.ctx.Emitter.EmitAB(vm.OpPush, accums[i], resultReg)
		} else {
			aggrKeyName := selector.Identifier().GetText()
			aggrKeyReg := loadConstant(cc.ctx, runtime.String(aggrKeyName))
			cc.ctx.Emitter.EmitABC(vm.OpPushKV, loop.Result, aggrKeyReg, resultReg)
			cc.ctx.Registers.Free(aggrKeyReg)
		}

		cc.ctx.Registers.Free(resultReg)
	}

	// Now we can iterate over the grouped items
	loop.EmitFinalization(cc.ctx.Emitter)

	// Now we can iterate over the selectors and execute the aggregation functions by passing the accumulators
	// And define variables for each accumulator result
	if isCollected {
		for i, selector := range selectors {
			fcx := selector.FunctionCallExpression()
			// We won't make any checks here, as we already did it before
			selectorVarName := selector.Identifier().GetText()

			// We execute the function call with the accumulator as an argument
			accum := accums[i]
			result := cc.ctx.ExprCompiler.CompileFunctionCallWith(fcx.FunctionCall(), fcx.ErrorOperator() != nil, core.RegisterSequence{accum})

			// We define the variable for the selector result in the upper scope
			// Since this temporary scope is only for aggregators and will be closed after the aggregation
			varReg := cc.ctx.Symbols.DeclareLocal(selectorVarName)
			cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
			cc.ctx.Registers.Free(result)
		}

		cc.ctx.Loops.Pop()
		// Now close the aggregators scope
		cc.ctx.Symbols.ExitScope()
	} else {
		// Now close the aggregators scope
		cc.ctx.Symbols.ExitScope()

		parentLoop.ValueName = ""
		parentLoop.KeyName = ""

		// Since we we in the middle of the loop, we need to patch the loop result
		// Now we just create a range with 1 item to push the aggregated values to the dataset
		// Replace source with sorted array
		zero := loadConstant(cc.ctx, runtime.Int(0))
		one := loadConstant(cc.ctx, runtime.Int(1))
		aggregator := cc.ctx.Registers.Allocate(core.Temp)
		cc.ctx.Emitter.EmitAB(vm.OpMove, aggregator, loop.Result)
		cc.ctx.Symbols.ExitScope()

		cc.ctx.Symbols.EnterScope()

		// Create new for loop
		cc.ctx.Emitter.EmitABC(vm.OpRange, loop.Src, zero, one)
		cc.ctx.Emitter.EmitAb(vm.OpDataSet, loop.Result, loop.Distinct)

		// In case of non-collected aggregators, we just iterate over the grouped items
		// Retrieve the grouped values by key, execute aggregation funcs and assign variable names to the results
		for _, selector := range selectors {
			fcx := selector.FunctionCallExpression()
			// We won't make any checks here, as we already did it before
			selectorVarName := selector.Identifier().GetText()

			// We execute the function call with the accumulator as an argument
			key := loadConstant(cc.ctx, runtime.String(selectorVarName))
			value := cc.ctx.Registers.Allocate(core.Temp)
			cc.ctx.Emitter.EmitABC(vm.OpLoadKey, value, aggregator, key)

			result := cc.ctx.ExprCompiler.CompileFunctionCallWith(fcx.FunctionCall(), fcx.ErrorOperator() != nil, core.RegisterSequence{value})

			// We define the variable for the selector result in the upper scope
			// Since this temporary scope is only for aggregators and will be closed after the aggregation
			varReg := cc.ctx.Symbols.DeclareLocal(selectorVarName)
			cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, result)
			cc.ctx.Registers.Free(result)
			cc.ctx.Registers.Free(value)
			cc.ctx.Registers.Free(key)
		}

		cc.ctx.Registers.Free(aggregator)
	}

	// Free the registers for accumulators
	for _, reg := range accums {
		cc.ctx.Registers.Free(reg)
	}

	// Free the register for the iterator value
	cc.ctx.Registers.Free(aggrIterVal)
}

func (cc *CollectCompiler) compileCollectGroupKeySelectors(selectors []fql.ICollectSelectorContext) vm.Operand {
	if len(selectors) == 0 {
		return vm.NoopOperand
	}

	var kvKeyReg vm.Operand

	if len(selectors) > 1 {
		// We create a sequence of Registers for the clauses
		// To pack them into an array
		selectorRegs := cc.ctx.Registers.AllocateSequence(len(selectors))

		for i, selector := range selectors {
			reg := cc.ctx.ExprCompiler.Compile(selector.Expression())
			cc.ctx.Emitter.EmitAB(vm.OpMove, selectorRegs[i], reg)
			// Free the register after moving its value to the sequence register
			cc.ctx.Registers.Free(reg)
		}

		kvKeyReg = cc.ctx.Registers.Allocate(core.Temp)
		cc.ctx.Emitter.EmitAs(vm.OpList, kvKeyReg, selectorRegs)
		cc.ctx.Registers.FreeSequence(selectorRegs)
	} else {
		kvKeyReg = cc.ctx.ExprCompiler.Compile(selectors[0].Expression())
	}

	return kvKeyReg
}

func (cc *CollectCompiler) compileCollectGroupKeySelectorVariables(selectors []fql.ICollectSelectorContext, kvKeyReg, kvValReg vm.Operand, isAggregation bool) {
	if len(selectors) > 1 {
		variables := make([]vm.Operand, len(selectors))

		for i, selector := range selectors {
			name := selector.Identifier().GetText()

			if variables[i] == vm.NoopOperand {
				variables[i] = cc.ctx.Symbols.DeclareLocal(name)
			}

			reg := kvValReg

			if isAggregation {
				reg = kvKeyReg
			}

			cc.ctx.Emitter.EmitABC(vm.OpLoadIndex, variables[i], reg, loadConstant(cc.ctx, runtime.Int(i)))
		}

		// Free the register after moving its value to the variable
		for _, reg := range variables {
			cc.ctx.Registers.Free(reg)
		}
	} else {
		// Get the variable name
		name := selectors[0].Identifier().GetText()
		// Define a variable for each selector
		varReg := cc.ctx.Symbols.DeclareLocal(name)

		reg := kvValReg

		if isAggregation {
			reg = kvKeyReg
		}

		// If we have a single selector, we can just move the value
		cc.ctx.Emitter.EmitAB(vm.OpMove, varReg, reg)
	}
}

func (cc *CollectCompiler) compileDefaultGroupProjection(loop *core.Loop, kvValReg vm.Operand, identifier antlr.TerminalNode, keeper fql.ICollectGroupVariableKeeperContext) string {
	if keeper == nil {
		seq := cc.ctx.Registers.AllocateSequence(2) // Key and Value for Map

		// TODO: Review this. It's quite a questionable ArrangoDB feature of wrapping group items by a nested object
		// We will keep it for now for backward compatibility.
		loadConstantTo(cc.ctx, runtime.String(loop.ValueName), seq[0]) // Map key
		cc.ctx.Emitter.EmitAB(vm.OpMove, seq[1], kvValReg)             // Map value
		cc.ctx.Emitter.EmitAs(vm.OpMap, kvValReg, seq)

		cc.ctx.Registers.FreeSequence(seq)
	} else {
		variables := keeper.AllIdentifier()
		seq := cc.ctx.Registers.AllocateSequence(len(variables) * 2)

		for i, j := 0, 0; i < len(variables); i, j = i+1, j+2 {
			varName := variables[i].GetText()
			loadConstantTo(cc.ctx, runtime.String(varName), seq[j])

			variable, _, found := cc.ctx.Symbols.Resolve(varName)

			if !found {
				panic("variable not found: " + varName)
			}

			cc.ctx.Emitter.EmitAB(vm.OpMove, seq[j+1], variable)
		}

		cc.ctx.Emitter.EmitAs(vm.OpMap, kvValReg, seq)
		cc.ctx.Registers.FreeSequence(seq)
	}

	return identifier.GetText()
}

func (cc *CollectCompiler) compileCustomGroupProjection(_ *core.Loop, kvValReg vm.Operand, selector fql.ICollectSelectorContext) string {
	selectorReg := cc.ctx.ExprCompiler.Compile(selector.Expression())
	cc.ctx.Emitter.EmitMove(kvValReg, selectorReg)
	cc.ctx.Registers.Free(selectorReg)

	return selector.Identifier().GetText()
}
