package vm

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type VM struct {
	cache   *mem.Cache
	program *bytecode.Program
	plan    execPlan
	execState
	options options
}

func New(program *bytecode.Program) (*VM, error) {
	return NewWith(program)
}

func NewWith(program *bytecode.Program, opts ...Option) (*VM, error) {
	if err := validate(program); err != nil {
		return nil, err
	}
	plan, err := buildExecPlan(program)

	if err != nil {
		return nil, err
	}

	o := newOptions(opts)

	vm := &VM{
		cache:   mem.NewCache(len(program.Bytecode), len(plan.hostCallDescriptors), o.shapeCacheLimit),
		program: program,
		plan:    plan,
		options: o,
	}
	state := &vm.execState
	state.init(program, plan.catchByPC)

	return vm, nil
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
	switch vm.options.panicPolicy {
	case PanicPropagate:
		defer vm.execState.end()
		return vm.runUnchecked(ctx, env)
	default:
		result, err := vm.runRecovered(ctx, env)
		vm.execState.end()
		return result, err
	}
}

func (vm *VM) runRecovered(ctx context.Context, env *Environment) (result runtime.Value, err error) {
	defer func() {
		state := &vm.execState

		if r := recover(); r != nil {
			err = state.runtimeErrorFromPanic(r)
			result = nil

			return
		}

		if err != nil {
			err = state.wrapRuntimeError(err)
		}
	}()

	return vm.runCore(ctx, env)
}

func (vm *VM) runUnchecked(ctx context.Context, env *Environment) (runtime.Value, error) {
	result, err := vm.runCore(ctx, env)

	if err != nil {
		var invariantErr *diagnostics.InvariantError
		if errors.As(err, &invariantErr) {
			panic(err)
		}

		return nil, vm.execState.wrapRuntimeError(err)
	}

	return result, nil
}

func (vm *VM) runCore(ctx context.Context, env *Environment) (runtime.Value, error) {
	if env == nil {
		env = noopEnv
	}

	state := &vm.execState
	state.start(env)
	if err := state.bindParams(env); err != nil {
		return nil, err
	}

	ensureRegexpsWarmed(vm)
	if err := ensureHostFunctionsBound(vm, env); err != nil {
		return nil, err
	}

	instructions := vm.plan.instructions
	constants := vm.program.Constants
	aggregatePlans := vm.program.Metadata.AggregatePlans
	shapeCache := vm.cache.ShapeCache
	hostFunctions := vm.cache.HostFunctions
	udfs := vm.program.Functions.UserDefined
	hostCallDescriptors := vm.plan.hostCallDescriptors
	udfCallDescriptors := vm.plan.udfCallDescriptors
	udfTailCallDescriptors := vm.plan.udfTailCallDescriptors
	paramSlots := state.scratch.Params
loop:
	for state.pc < len(instructions) {
		pc := state.pc
		inst := &instructions[pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := state.registers.Values
		state.pc = pc + 1

		switch op {
		case bytecode.OpReturn:
			retVal := reg[dst]

			if state.returnToCaller(retVal) {
				continue
			}

			reg[bytecode.NoopOperand] = retVal

			break loop
		case bytecode.OpJump:
			state.pc = int(dst)
		case bytecode.OpJumpIfFalse:
			if !coerceBool(reg[src1]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfTrue:
			if coerceBool(reg[src1]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfNone:
			if reg[src1] == runtime.None {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfNe:
			if ne(ctx, reg[src1], reg[src2]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfNeConst:
			if ne(ctx, reg[src1], constants[src2.Constant()]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfEq:
			if eq(ctx, reg[src1], reg[src2]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfEqConst:
			if eq(ctx, reg[src1], constants[src2.Constant()]) {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingProperty:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				state.pc = int(dst)
				continue
			}

			key, ok := reg[src2].(runtime.String)
			if !ok {
				state.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			if !has {
				state.pc = int(dst)
			}
		case bytecode.OpJumpIfMissingPropertyConst:
			obj, ok := reg[src1].(runtime.Map)
			if !ok {
				state.pc = int(dst)

				continue
			}

			key, ok := constants[src2.Constant()].(runtime.String)
			if !ok {
				state.pc = int(dst)
				continue
			}

			has, err := obj.ContainsKey(ctx, key)
			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			if !has {
				state.pc = int(dst)
			}
		case bytecode.OpFail:
			if !dst.IsConstant() {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message")
				state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			idx := dst.Constant()
			if idx < 0 || idx >= len(constants) {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message")
				state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			msg, ok := constants[idx].(runtime.String)
			if !ok {
				callErr := runtime.TypeErrorOf(constants[idx], runtime.TypeString)
				state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			callErr := runtime.Error(runtime.ErrInvalidOperation, msg.String())
			state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			hostID := inst.InlineSlot
			if hostID < 0 || hostID >= len(hostFunctions) {
				invariantErr := diagnostics.NewInvariantError(
					"invalid host call slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid host call slot %d at pc %d", hostID, pc),
				)
				state.raiseInvariantAt(pc, invariantErr)
				break
			}

			call := &hostCallDescriptors[hostID]
			cacheFn := &hostFunctions[call.ID]
			out, err := callCachedHostFunction(ctx, call, cacheFn, reg, &state.scratch)

			state.setCallResult(pc, op, dst, out, err)
		case bytecode.OpCall, bytecode.OpProtectedCall:
			callID := inst.InlineSlot
			call := &udfCallDescriptors[callID]

			if call.ID < 0 || call.ID >= len(udfs) {
				invariantErr := diagnostics.NewInvariantError(
					"invalid udf call slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid udf call slot %d at pc %d", callID, pc),
				)
				state.raiseInvariantAt(pc, invariantErr)
				break
			}

			udf := &udfs[call.ID]

			if err := state.callUdf(call, udf); err != nil {
				state.setCallResult(pc, op, dst, runtime.None, err)
			}
		case bytecode.OpTailCall:
			callID := inst.InlineSlot
			call := &udfTailCallDescriptors[callID]

			if call.ID < 0 || call.ID >= len(udfs) {
				invariantErr := diagnostics.NewInvariantError(
					"invalid udf call slot",
					runtime.Errorf(runtime.ErrUnexpected, "invalid udf call slot %d at pc %d", callID, pc),
				)
				state.raiseInvariantAt(pc, invariantErr)
				break
			}

			udf := &udfs[call.ID]

			if err := state.tailCallUdf(call, udf); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
			}
		case bytecode.OpDispatch:
			dispatcher, eventName, payload, opts, err := coerceDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])

			if err != nil {
				state.setOrRaiseDefault(pc, dst, runtime.None, err)
				break
			}

			out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
				Name:    eventName,
				Payload: payload,
				Options: opts,
			})

			if out == nil {
				out = runtime.None
			}

			state.setOrRaiseDefault(pc, dst, out, err)
		case bytecode.OpMove:
			reg[dst] = reg[src1]
		case bytecode.OpLoadNone:
			reg[dst] = runtime.None
		case bytecode.OpLoadBool:
			reg[dst] = runtime.Boolean(src1 == 1)
		case bytecode.OpLoadZero:
			reg[dst] = runtime.ZeroInt
		case bytecode.OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case bytecode.OpLoadParam:
			reg[dst] = paramSlots[int(src1)-1]
		case bytecode.OpLoadArray:
			reg[dst] = runtime.NewArray(int(src1))
		case bytecode.OpLoadObject:
			reg[dst] = data.NewFastObjectOf(shapeCache, vm.options.fastObjectDictThreshold, int(src1))
		case bytecode.OpLoadRange:
			start, err := runtime.ToInt(ctx, reg[src1])

			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			end, err := runtime.ToInt(ctx, reg[src2])

			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = runtime.NewRange(start, end)
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			vm.loadIndexAndSet(state, ctx, dst, pc, src, arg, optional)
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			vm.loadKeyAndSet(state, ctx, dst, pc, src, arg, optional)
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			// I guess the reason it cannot inline is due to a different control flow
			vm.loadPropertyAndSet(state, ctx, dst, pc, src, prop, optional)
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			vm.loadIndexAndSet(state, ctx, dst, pc, src, arg, optional)
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			vm.loadKeyConstAndSet(state, ctx, dst, pc, inst, src, arg, optional)
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			vm.loadPropertyConstAndSet(state, ctx, dst, pc, inst, src, prop, optional)
		case bytecode.OpPush:
			ds := reg[dst].(runtime.Appendable)

			if err := ds.Append(ctx, reg[src1]); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
			}
		case bytecode.OpPushKV:
			tr := reg[dst].(runtime.KeyWritable)

			if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
			}
		case bytecode.OpArrayPush:
			ds := reg[dst].(*runtime.Array)

			_ = ds.Append(ctx, reg[src1])
		case bytecode.OpObjectSet:
			key := runtime.ToString(reg[src1])
			value := reg[src2]
			obj, ok := reg[dst].(*data.FastObject)

			if ok {
				_ = obj.Set(ctx, key, value)
				continue
			}

			writable, ok := reg[dst].(runtime.KeyWritable)

			if ok {
				if err := writable.Set(ctx, key, value); err != nil {
					state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
		case bytecode.OpObjectSetConst:
			objVal := reg[dst]
			key := runtime.ToString(constants[src1.Constant()])
			value := reg[src2]

			if obj, ok := objVal.(*data.FastObject); ok {
				vm.objectSetConstCached(inst, obj, key, value)
				continue
			}

			writable, ok := reg[dst].(runtime.KeyWritable)

			if ok {
				if err := writable.Set(ctx, key, value); err != nil {
					state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			state.raiseRuntimeAt(pc, callErr, recoverDefault, bytecode.NoopOperand, nil, false)
		case bytecode.OpIter:
			input := reg[src1]
			iterable, ok := input.(runtime.Iterable)

			if ok {
				iterator, err := iterable.Iterate(ctx)

				if err == nil {
					reg[dst] = data.NewIterator(iterator)
					continue
				}

				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			callErr := runtime.TypeErrorOf(input, runtime.TypeIterable)
			state.raiseRuntimeAt(pc, callErr, recoverDefault, dst, data.NoopIter, true)
		case bytecode.OpIterNext:
			iterator := reg[src1].(*data.Iterator)

			if err := iterator.Next(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					state.pc = int(dst)
					continue
				}

				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
			}
		case bytecode.OpIterValue:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Value()
		case bytecode.OpIterKey:
			iterator := reg[src1].(*data.Iterator)
			reg[dst] = iterator.Key()
		case bytecode.OpIterSkip:
			iterState := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if iterState < threshold {
				iterState++
				reg[src1] = iterState
				state.pc = int(dst)
			}
		case bytecode.OpIterLimit:
			iterState := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if iterState < threshold {
				iterState++
				reg[src1] = iterState
			} else {
				state.pc = int(dst)
			}
		case bytecode.OpStream:
			observable, eventName, opts, err := coerceSubscribeArgs(reg[dst], reg[src1], reg[src2])

			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			stream, err := observable.Subscribe(ctx, runtime.Subscription{
				EventName: eventName,
				Options:   opts,
			})

			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = data.NewStreamValue(stream)
		case bytecode.OpStreamIter:
			stream := reg[src1].(*data.StreamValue)

			var timeout runtime.Int

			if reg[src2] != runtime.None {
				t, err := runtime.CastInt(reg[src2])

				if err != nil {
					state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				timeout = t
			}

			reg[dst] = stream.Iterate(timeout)
		case bytecode.OpQuery:
			src := readOperandValue(reg, constants, src1)
			descriptor := readOperandValue(reg, constants, src2)
			out, err := applyQuery(ctx, src, descriptor)

			state.setOrRaiseDefault(pc, dst, out, err)
		case bytecode.OpDataSet:
			reg[dst] = data.NewDataSet(src1 == 1)
		case bytecode.OpDataSetCollector:
			collectorType := bytecode.CollectorType(src1)

			if collectorType == bytecode.CollectorTypeAggregate || collectorType == bytecode.CollectorTypeAggregateGroup {
				planIdx := int(src2)

				if planIdx < 0 || planIdx >= len(aggregatePlans) {
					invariantErr := diagnostics.NewInvariantError(
						"invalid aggregate plan index",
						runtime.Errorf(runtime.ErrUnexpected, "invalid aggregate plan index %d", planIdx),
					)
					state.raiseInvariantAt(pc, invariantErr)
					break
				}

				plan := aggregatePlans[planIdx]

				if collectorType == bytecode.CollectorTypeAggregate {
					reg[dst] = data.NewAggregateCollector(plan)
				} else {
					reg[dst] = data.NewGroupedAggregateCollector(plan)
				}

				continue
			}

			collector, err := data.NewCollectorSafe(collectorType)
			if err != nil {
				invariantErr := diagnostics.NewInvariantError("invalid collector configuration", err)
				state.raiseInvariantAt(pc, invariantErr)
				break
			}

			reg[dst] = collector
		case bytecode.OpDataSetSorter:
			reg[dst] = data.NewSorter(runtime.SortDirection(src1))
		case bytecode.OpDataSetMultiSorter:
			encoded := src1.Register()
			count := src2.Register()

			reg[dst] = data.NewMultiSorter(runtime.DecodeSortDirections(encoded, count))
		case bytecode.OpFlatten:
			depth := src2.Register()

			if depth < 1 {
				depth = 1
			}

			res, err := arrayFlatten(ctx, reg[src1], depth)

			state.setOrRaiseDefault(pc, dst, res, err)
		case bytecode.OpAdd:
			reg[dst] = runtime.Add(ctx, reg[src1], reg[src2])
		case bytecode.OpAddConst:
			reg[dst] = runtime.Add(ctx, reg[src1], constants[src2.Constant()])
		case bytecode.OpConcat:
			concatStrings(reg, dst, src1, src2)
		case bytecode.OpSub:
			reg[dst] = runtime.Subtract(ctx, reg[src1], reg[src2])
		case bytecode.OpMul:
			reg[dst] = runtime.Multiply(ctx, reg[src1], reg[src2])
		case bytecode.OpDiv:
			if err := state.checkDivisionByZeroAt(ctx, pc, reg[src1], reg[src2]); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := state.checkModuloByZeroAt(ctx, pc, reg[src2]); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = runtime.Modulus(ctx, reg[src1], reg[src2])
		case bytecode.OpIncr:
			reg[dst] = runtime.Increment(ctx, reg[dst])
		case bytecode.OpDecr:
			reg[dst] = runtime.Decrement(ctx, reg[dst])
		case bytecode.OpNegate:
			reg[dst] = negate(reg[src1])
		case bytecode.OpFlipPositive:
			reg[dst] = positive(reg[src1])
		case bytecode.OpFlipNegative:
			reg[dst] = negative(reg[src1])
		case bytecode.OpCastBool:
			reg[dst] = coerceBool(reg[src1])
		case bytecode.OpCmp:
			reg[dst] = cmp(ctx, reg[src1], reg[src2])
		case bytecode.OpNot:
			reg[dst] = !coerceBool(reg[src1])
		case bytecode.OpEq:
			reg[dst] = eq(ctx, reg[src1], reg[src2])
		case bytecode.OpNe:
			reg[dst] = ne(ctx, reg[src1], reg[src2])
		case bytecode.OpGt:
			reg[dst] = gt(ctx, reg[src1], reg[src2])
		case bytecode.OpLt:
			reg[dst] = lt(ctx, reg[src1], reg[src2])
		case bytecode.OpGte:
			reg[dst] = gte(ctx, reg[src1], reg[src2])
		case bytecode.OpLte:
			reg[dst] = lte(ctx, reg[src1], reg[src2])
		case bytecode.OpIn:
			reg[dst] = contains(ctx, reg[src2], reg[src1])
		case bytecode.OpLike:
			res, err := Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
				break
			}

			state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(pc, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
				continue
			}

			state.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.False, true)
		case bytecode.OpExists:
			val := reg[src1]

			if val == runtime.None {
				reg[dst] = runtime.False
				continue
			}

			if measurable, ok := val.(runtime.Measurable); ok {
				length, err := measurable.Length(ctx)

				if err != nil {
					state.raiseRuntimeAt(pc, err, recoverDefault, dst, runtime.False, true)
					break
				}

				reg[dst] = runtime.NewBoolean(length != 0)

				continue
			}

			reg[dst] = runtime.True
		case bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := arrayAll(ctx, cmp, reg[src1], reg[src2])

			state.setOrRaiseDefault(pc, dst, res, err)
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := arrayAny(ctx, cmp, reg[src1], reg[src2])

			state.setOrRaiseDefault(pc, dst, res, err)
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := arrayNone(ctx, cmp, reg[src1], reg[src2])

			state.setOrRaiseDefault(pc, dst, res, err)
		case bytecode.OpLength:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				length, err := val.Length(ctx)

				if err != nil {
					state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				reg[dst] = length
				continue
			}

			callErr := runtime.TypeErrorOf(reg[src1],
				runtime.TypeString,
				runtime.TypeList,
				runtime.TypeMap,
				runtime.TypeBinary,
				runtime.TypeMeasurable,
			)
			state.raiseRuntimeAt(pc, callErr, recoverDefault, dst, runtime.ZeroInt, true)
		case bytecode.OpType:
			reg[dst] = runtime.NewString(runtime.TypeName(runtime.TypeOf(reg[src1])))
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					state.raiseRuntimeAt(pc, closeErr, recoverDefault, bytecode.NoopOperand, nil, false)
				}
			}
		case bytecode.OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			if err := data.Sleep(ctx, dur); err != nil {
				state.raiseRuntimeAt(pc, err, recoverDefault, bytecode.NoopOperand, nil, false)
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		default:
			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, pc)
		}

		// Sticky checkpoint: opcode branches only raise failures; resolution happens here.
		if state.hasFail {
			if state.resolveFailure() == errReturn {
				return nil, state.failure.err
			}

			continue
		}
	}

	return state.registers.Values[bytecode.NoopOperand], nil
}

func (vm *VM) regexpCached(pc int, value runtime.Value) (*data.Regexp, error) {
	// We compare patterns to ensure that the cached regexp is the same as the one we're trying to use.
	// This is necessary because the same compiled function can be used in different places with different regexps,
	// and we want to avoid caching a regexp that doesn't match the current pattern.
	switch v := value.(type) {
	case *data.Regexp:
		pattern := v.String()

		if cached := vm.cache.Regexps[pc]; cached == nil || cached.Pattern != pattern {
			vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: v}
		}

		return v, nil
	case runtime.String:
		pattern := v.String()

		if cached := vm.cache.Regexps[pc]; cached != nil && cached.Pattern == pattern {
			return cached.Regexp, nil
		}

		r, err := data.NewRegexp(v)
		if err != nil {
			return nil, err
		}

		vm.cache.Regexps[pc] = &mem.CachedRegexp{Pattern: pattern, Regexp: r}

		return r, nil
	default:
		return nil, runtime.TypeErrorOf(value, runtime.TypeString, data.TypeRegexp)
	}
}

func (vm *VM) loadFastKeyCached(
	ctx context.Context,
	pc int,
	inst *execInstruction,
	obj *data.FastObject,
	arg runtime.Value,
	key string,
	constKey bool,
) (runtime.Value, error) {
	shapeID := obj.ShapeID()
	if shapeID == 0 {
		return vm.loadKey(ctx, obj, arg)
	}

	if constKey {
		if inst != nil && inst.InlineShapeID == shapeID {
			if inst.InlineSlot < 0 {
				return runtime.None, nil
			}
			if val, ok := obj.SlotValue(inst.InlineSlot); ok {
				return val, nil
			}

			return runtime.None, nil
		}

		if pc < 0 || pc >= len(vm.cache.LoadKeyConstICs) {
			return vm.loadKey(ctx, obj, arg)
		}

		cache := vm.cache.LoadKeyConstICs[pc]
		if cache != nil {
			if slot, ok := cache.Lookup(shapeID); ok {
				if inst != nil {
					inst.InlineShapeID = shapeID
					inst.InlineSlot = slot
				}

				if slot < 0 {
					return runtime.None, nil
				}
				if val, ok := obj.SlotValue(slot); ok {
					return val, nil
				}

				return runtime.None, nil
			}
		}

		slot, ok := obj.LookupSlot(key)
		if !ok {
			if cache == nil {
				cache = mem.NewLoadKeyConstCache()
				vm.cache.LoadKeyConstICs[pc] = cache
			}

			cache.Add(shapeID, -1)

			if inst != nil {
				inst.InlineShapeID = shapeID
				inst.InlineSlot = -1
			}

			return runtime.None, nil
		}

		val, ok := obj.SlotValue(slot)
		if !ok {
			return runtime.None, nil
		}

		if cache == nil {
			cache = mem.NewLoadKeyConstCache()
			vm.cache.LoadKeyConstICs[pc] = cache
		}

		cache.Add(shapeID, slot)

		if inst != nil {
			inst.InlineShapeID = shapeID
			inst.InlineSlot = slot
		}

		return val, nil
	}

	if pc < 0 || pc >= len(vm.cache.LoadKeyICs) {
		return vm.loadKey(ctx, obj, arg)
	}

	cache := vm.cache.LoadKeyICs[pc]
	if cache != nil {
		if slot, ok := cache.Lookup(shapeID, key); ok {
			if slot < 0 {
				return runtime.None, nil
			}
			if val, ok := obj.SlotValue(slot); ok {
				return val, nil
			}

			return runtime.None, nil
		}
	}

	slot, ok := obj.LookupSlot(key)
	if !ok {
		if cache == nil {
			cache = mem.NewLoadKeyCache()
			vm.cache.LoadKeyICs[pc] = cache
		}

		cache.Add(shapeID, key, -1)

		return runtime.None, nil
	}

	val, ok := obj.SlotValue(slot)
	if !ok {
		return runtime.None, nil
	}

	if cache == nil {
		cache = mem.NewLoadKeyCache()
		vm.cache.LoadKeyICs[pc] = cache
	}

	cache.Add(shapeID, key, slot)

	return val, nil
}

func (vm *VM) loadKeyCached(ctx context.Context, pc int, src, arg runtime.Value) (runtime.Value, error) {
	obj, ok := src.(*data.FastObject)

	if !ok {
		return vm.loadKey(ctx, src, arg)
	}

	var key string

	switch v := arg.(type) {
	case runtime.String:
		key = string(v)
	default:
		key = runtime.ToString(v).String()
	}

	return vm.loadFastKeyCached(ctx, pc, nil, obj, arg, key, false)
}

func (vm *VM) loadKeyConstCached(ctx context.Context, pc int, inst *execInstruction, src, arg runtime.Value) (runtime.Value, error) {
	obj, ok := src.(*data.FastObject)

	if !ok {
		return vm.loadKey(ctx, src, arg)
	}

	var key string

	switch v := arg.(type) {
	case runtime.String:
		key = string(v)
	default:
		key = runtime.ToString(v).String()
	}

	return vm.loadFastKeyCached(ctx, pc, inst, obj, arg, key, true)
}

func (vm *VM) objectSetConstCached(inst *execInstruction, obj *data.FastObject, key runtime.String, value runtime.Value) {
	if obj == nil {
		return
	}

	if inst != nil {
		shape := obj.Shape()

		if shape != nil && inst.InlineSetShape == shape {
			if obj.SetSlotWithShape(inst.InlineSetNextShape, inst.InlineSlot, value) {
				return
			}
		}
	}

	prev, next, slot, ok := obj.SetStringCached(string(key), value)

	if ok && inst != nil {
		inst.InlineSetShape = prev
		inst.InlineSetNextShape = next
		inst.InlineSlot = slot
	}
}

func (vm *VM) loadIndex(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	indexed, ok := src.(runtime.IndexReadable)

	if !ok {
		return nil, diagnostics.MemberAccessErrorOf(src, diagnostics.MemberAccessIndex, arg)
	}

	var idx runtime.Int
	var err error

	switch v := arg.(type) {
	case runtime.Int:
		idx = v
	case runtime.Float:
		// Convert float to int, rounding down
		idx = runtime.Int(v)
	default:
		err = runtime.TypeErrorOf(arg, runtime.TypeInt, runtime.TypeFloat)
	}

	if err != nil {
		return nil, err
	}

	return indexed.At(ctx, idx)
}

func (vm *VM) loadKey(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	keyed, ok := src.(runtime.KeyReadable)

	if !ok {
		return nil, diagnostics.MemberAccessErrorOf(src, diagnostics.MemberAccessProperty, arg)
	}

	out, err := keyed.Get(ctx, arg)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (vm *VM) loadIndexAndSet(state *execState, ctx context.Context, dst bytecode.Operand, pc int, src, arg runtime.Value, optional bool) {
	if optional && src == runtime.None {
		state.registers.Values[dst] = runtime.None
		return
	}

	out, err := vm.loadIndex(ctx, src, arg)
	state.setOrOptional(pc, dst, out, err, optional)
}

func (vm *VM) loadKeyAndSet(state *execState, ctx context.Context, dst bytecode.Operand, pc int, src, arg runtime.Value, optional bool) {
	if optional && src == runtime.None {
		state.registers.Values[dst] = runtime.None
		return
	}

	out, err := vm.loadKeyCached(ctx, pc, src, arg)
	state.setOrOptional(pc, dst, out, err, optional)
}

func (vm *VM) loadKeyConstAndSet(state *execState, ctx context.Context, dst bytecode.Operand, pc int, inst *execInstruction, src, arg runtime.Value, optional bool) {
	if optional && src == runtime.None {
		state.registers.Values[dst] = runtime.None
		return
	}

	out, err := vm.loadKeyConstCached(ctx, pc, inst, src, arg)
	state.setOrOptional(pc, dst, out, err, optional)
}

func (vm *VM) loadPropertyAndSet(state *execState, ctx context.Context, dst bytecode.Operand, pc int, src, prop runtime.Value, optional bool) {
	if optional && src == runtime.None {
		state.registers.Values[dst] = runtime.None
		return
	}

	var out runtime.Value
	var err error

	switch getter := prop.(type) {
	case runtime.String:
		out, err = vm.loadKeyCached(ctx, pc, src, getter)
	case runtime.Float, runtime.Int:
		out, err = vm.loadIndex(ctx, src, getter)
	default:
		out, err = vm.loadKeyCached(ctx, pc, src, runtime.ToString(prop))
	}

	state.setOrOptional(pc, dst, out, err, optional)
}

func (vm *VM) loadPropertyConstAndSet(state *execState, ctx context.Context, dst bytecode.Operand, pc int, inst *execInstruction, src, prop runtime.Value, optional bool) {
	if optional && src == runtime.None {
		state.registers.Values[dst] = runtime.None
		return
	}

	var out runtime.Value
	var err error

	switch getter := prop.(type) {
	case runtime.String:
		out, err = vm.loadKeyConstCached(ctx, pc, inst, src, getter)
	case runtime.Float, runtime.Int:
		out, err = vm.loadIndex(ctx, src, getter)
	default:
		out, err = vm.loadKeyConstCached(ctx, pc, inst, src, runtime.ToString(prop))
	}

	state.setOrOptional(pc, dst, out, err, optional)
}
