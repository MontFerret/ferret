package vm

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type VM struct {
	cache        *mem.Cache
	program      *bytecode.Program
	catchByPC    []int
	hostWarmups  []hostCallWarmupDescriptor
	instructions []data.ExecInstruction
	statePool    []*execState
	options      options
}

func New(program *bytecode.Program) (*VM, error) {
	return NewWith(program)
}

func NewWith(program *bytecode.Program, opts ...Option) (*VM, error) {
	if err := validate(program); err != nil {
		return nil, err
	}

	o := newOptions(opts)
	catchByPC := buildCatchByPC(len(program.Bytecode), program.CatchTable)

	vm := &VM{
		cache:        mem.NewCache(len(program.Bytecode), o.shapeCacheLimit),
		program:      program,
		catchByPC:    catchByPC,
		hostWarmups:  buildHostWarmupDescriptors(program),
		options:      o,
		instructions: buildExecInstructions(program.Bytecode),
	}

	return vm, nil
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
	state := vm.acquireRunState()
	defer vm.releaseRunState(state)

	switch vm.options.panicPolicy {
	case PanicPropagate:
		return vm.runUnchecked(ctx, env, state)
	default:
		return vm.runRecovered(ctx, env, state)
	}
}

func (vm *VM) acquireRunState() *execState {
	n := len(vm.statePool)
	if n > 0 {
		state := vm.statePool[n-1]
		vm.statePool = vm.statePool[:n-1]
		return state
	}

	state := &execState{}
	state.init(vm.program, vm.catchByPC)

	return state
}

func (vm *VM) releaseRunState(state *execState) {
	if state == nil {
		return
	}

	state.cleanupForPool()
	vm.statePool = append(vm.statePool, state)
}

func (vm *VM) runRecovered(ctx context.Context, env *Environment, state *execState) (result runtime.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = state.runtimeErrorFromPanic(r)
			result = nil

			return
		}

		if err != nil {
			err = state.wrapRuntimeError(err)
		}
	}()

	return vm.runCore(ctx, env, state)
}

func (vm *VM) runUnchecked(ctx context.Context, env *Environment, state *execState) (runtime.Value, error) {
	result, err := vm.runCore(ctx, env, state)

	if err != nil {
		var invariantErr *diagnostic.InvariantError
		if errors.As(err, &invariantErr) {
			panic(err)
		}

		return nil, state.wrapRuntimeError(err)
	}

	return result, nil
}

func (vm *VM) runCore(ctx context.Context, env *Environment, state *execState) (runtime.Value, error) {
	if env == nil {
		env = noopEnv
	}

	state.prepareRun(env)

	if err := warmup(vm, state, env); err != nil {
		return nil, err
	}

	instructions := vm.instructions
	constants := vm.program.Constants
	aggregatePlans := vm.program.Metadata.AggregatePlans
	shapeCache := vm.cache.ShapeCache
	paramSlots := state.scratch.Params
loop:
	for state.pc < len(instructions) {
		inst := &instructions[state.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		reg := state.registers.Values
		state.pc++

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
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
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
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			if !has {
				state.pc = int(dst)
			}
		case bytecode.OpFail:
			if !dst.IsConstant() {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message")
				state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			idx := dst.Constant()
			if idx < 0 || idx >= len(constants) {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message")
				state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			msg, ok := constants[idx].(runtime.String)
			if !ok {
				callErr := runtime.TypeErrorOf(constants[idx], runtime.TypeString)
				state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			callErr := runtime.Error(runtime.ErrInvalidOperation, msg.String())
			state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
			break
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			cacheFn := vm.cache.HostFunctions[state.pc-1]
			out, err := callCachedHostFunction(ctx, cacheFn, state.registers.Values, state.scratch, reg[dst], src1, src2)

			state.setCallResult(op, dst, out, err)
		case bytecode.OpCall, bytecode.OpProtectedCall:
			if err := state.callUdf(op, dst, src1, src2); err != nil {
				state.setCallResult(op, dst, runtime.None, err)
			}
		case bytecode.OpTailCall:
			if err := state.tailCallUdf(dst, src1, src2); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}
		case bytecode.OpDispatch:
			dispatcher, eventName, payload, options, err := vm.castDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])

			if err != nil {
				state.setOrTryCatch(dst, runtime.None, err)
				break
			}

			out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
				Name:    eventName,
				Payload: payload,
				Options: options,
			})

			if out == nil {
				out = runtime.None
			}

			state.setOrTryCatch(dst, out, err)
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
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			end, err := runtime.ToInt(ctx, reg[src2])

			if err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = runtime.NewRange(start, end)
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			vm.loadIndexAndSet(state, ctx, dst, src, arg, optional)
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			vm.loadKeyAndSet(state, ctx, dst, state.pc-1, src, arg, optional)
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			// I guess the reason it cannot inline is due to a different control flow
			vm.loadPropertyAndSet(state, ctx, dst, state.pc-1, src, prop, optional)
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			vm.loadIndexAndSet(state, ctx, dst, src, arg, optional)
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			vm.loadKeyConstAndSet(state, ctx, dst, state.pc-1, inst, src, arg, optional)
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			vm.loadPropertyConstAndSet(state, ctx, dst, state.pc-1, inst, src, prop, optional)
		case bytecode.OpPush:
			ds := reg[dst].(runtime.Appendable)

			if err := ds.Append(ctx, reg[src1]); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}
		case bytecode.OpPushKV:
			tr := reg[dst].(runtime.KeyWritable)

			if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
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
					state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
			break
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
					state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}
				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			state.raiseRuntime(callErr, recoverDefault, bytecode.NoopOperand, nil, false)
			break
		case bytecode.OpIter:
			input := reg[src1]
			iterable, ok := input.(runtime.Iterable)

			if ok {
				iterator, err := iterable.Iterate(ctx)

				if err == nil {
					reg[dst] = data.NewIterator(iterator)
					continue
				}

				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			// TODO: replace with inlined version
			callErr := runtime.TypeErrorOf(input, runtime.TypeIterable)
			state.raiseRuntime(callErr, recoverDefault, dst, data.NoopIter, true)
			break
		case bytecode.OpIterNext:
			iterator := reg[src1].(*data.Iterator)

			if err := iterator.Next(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					state.pc = int(dst)
					continue
				}

				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
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
			observable, eventName, options, err := vm.castSubscribeArgs(reg[dst], reg[src1], reg[src2])

			if err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			stream, err := observable.Subscribe(ctx, runtime.Subscription{
				EventName: eventName,
				Options:   options,
			})

			if err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			reg[dst] = data.NewStreamValue(stream)
		case bytecode.OpStreamIter:
			stream := reg[src1].(*data.StreamValue)

			var timeout runtime.Int

			if reg[src2] != runtime.None {
				t, err := runtime.CastInt(reg[src2])

				if err != nil {
					state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}

				timeout = t
			}

			reg[dst] = stream.Iterate(timeout)
		case bytecode.OpQuery:
			src := readOperandValue(reg, constants, src1)
			descriptor := readOperandValue(reg, constants, src2)
			// TODO: unwrap since it's cannot be inlined
			out, err := applyQuery(ctx, src, descriptor)

			state.setOrTryCatch(dst, out, err)
		case bytecode.OpDataSet:
			reg[dst] = data.NewDataSet(src1 == 1)
		case bytecode.OpDataSetCollector:
			collectorType := bytecode.CollectorType(src1)

			if collectorType == bytecode.CollectorTypeAggregate || collectorType == bytecode.CollectorTypeAggregateGroup {
				planIdx := int(src2)

				if planIdx < 0 || planIdx >= len(aggregatePlans) {
					invariantErr := diagnostic.NewInvariantError(
						"invalid aggregate plan index",
						runtime.Errorf(runtime.ErrUnexpected, "invalid aggregate plan index %d", planIdx),
					)
					state.raiseInvariant(invariantErr)
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
				invariantErr := diagnostic.NewInvariantError("invalid collector configuration", err)
				state.raiseInvariant(invariantErr)
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

			state.setOrTryCatch(dst, res, err)
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
			if err := state.checkDivisionByZero(ctx, reg[src1], reg[src2]); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := state.checkModuloByZero(ctx, reg[src2]); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
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

			state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
			break
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(state.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else {
				state.raiseRuntime(err, recoverDefault, dst, runtime.False, true)
				break
			}
		case bytecode.OpExists:
			val := reg[src1]

			if val == runtime.None {
				reg[dst] = runtime.False
				continue
			}

			if measurable, ok := val.(runtime.Measurable); ok {
				length, err := measurable.Length(ctx)

				if err != nil {
					state.raiseRuntime(err, recoverDefault, dst, runtime.False, true)
					break
				}

				reg[dst] = runtime.NewBoolean(length != 0)

				continue
			}

			reg[dst] = runtime.True
		case bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := arrayAll(ctx, cmp, reg[src1], reg[src2])

			state.setOrTryCatch(dst, res, err)
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := arrayAny(ctx, cmp, reg[src1], reg[src2])

			state.setOrTryCatch(dst, res, err)
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := arrayNone(ctx, cmp, reg[src1], reg[src2])

			state.setOrTryCatch(dst, res, err)
		case bytecode.OpLength:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				length, err := val.Length(ctx)

				if err != nil {
					state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
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
			state.raiseRuntime(callErr, recoverDefault, dst, runtime.ZeroInt, true)
			break
		case bytecode.OpType:
			reg[dst] = runtime.NewString(runtime.TypeName(runtime.TypeOf(reg[src1])))
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					state.raiseRuntime(closeErr, recoverDefault, bytecode.NoopOperand, nil, false)
					break
				}
			}
		case bytecode.OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}

			if err := data.Sleep(ctx, dur); err != nil {
				state.raiseRuntime(err, recoverDefault, bytecode.NoopOperand, nil, false)
				break
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		default:
			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, state.pc-1)
		}

		// Sticky checkpoint: opcode branches only raise failures; resolution happens here.
		if state.hasFailure() {
			if state.resolveFailure() == errReturn {
				return nil, state.failureError()
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

func (vm *VM) castSubscribeArgs(dst, eventName, opts runtime.Value) (runtime.Observable, runtime.String, runtime.Map, error) {
	observable, ok := dst.(runtime.Observable)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(dst, runtime.TypeObservable)
	}

	eventNameStr, ok := eventName.(runtime.String)

	if !ok {
		return nil, "", nil, runtime.TypeErrorOf(eventName, runtime.TypeString)
	}

	var options runtime.Map

	if opts != nil && opts != runtime.None {
		m, ok := opts.(runtime.Map)

		if !ok {
			return nil, "", nil, runtime.TypeErrorOf(opts, runtime.TypeMap)
		}

		options = m
	}

	return observable, eventNameStr, options, nil
}

func (vm *VM) castDispatchArgs(
	ctx context.Context,
	target, eventName, args runtime.Value,
) (runtime.Dispatchable, runtime.String, runtime.Value, runtime.Value, error) {
	dispatcher, ok := target.(runtime.Dispatchable)

	if !ok {
		return nil, "", nil, nil, runtime.TypeErrorOf(target, runtime.TypeDispatchable)
	}

	eventNameStr, err := runtime.CastString(eventName)

	if err != nil {
		return nil, "", nil, nil, err
	}

	var payload runtime.Value = runtime.None
	var options runtime.Value = runtime.None

	if args == nil || args == runtime.None {
		return dispatcher, eventNameStr, payload, options, nil
	}

	argMap, err := runtime.CastMap(args)

	if err != nil {
		return nil, "", nil, nil, err
	}

	if val, err := argMap.Get(ctx, runtime.NewString("payload")); err == nil {
		payload = val
	}

	if val, err := argMap.Get(ctx, runtime.NewString("options")); err == nil {
		options = val
	}

	return dispatcher, eventNameStr, payload, options, nil
}
