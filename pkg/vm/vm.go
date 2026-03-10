package vm

import (
	"context"
	"errors"
	"io"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

type VM struct {
	cache        *mem.Cache
	program      *bytecode.Program
	instructions []data.ExecInstruction
	state        execState
	options      options
}

func New(program *bytecode.Program) *VM {
	return NewWith(program)
}

func NewWith(program *bytecode.Program, opts ...Option) *VM {
	o := newOptions(opts)

	vm := &VM{
		cache:        mem.NewCache(len(program.Bytecode), o.shapeCacheLimit),
		program:      program,
		options:      o,
		instructions: buildExecInstructions(program.Bytecode),
	}

	vm.state.init(program, buildCatchByPC(len(program.Bytecode), program.CatchTable))

	return vm
}

func (vm *VM) Run(ctx context.Context, env *Environment) (runtime.Value, error) {
	switch vm.options.panicPolicy {
	case PanicPropagate:
		return vm.runUnchecked(ctx, env)
	default:
		return vm.runRecovered(ctx, env)
	}
}

func (vm *VM) runRecovered(ctx context.Context, env *Environment) (result runtime.Value, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = vm.state.runtimeErrorFromPanic(r)
			result = nil

			return
		}

		if err != nil {
			err = vm.state.wrapRuntimeError(err)
		}
	}()

	return vm.runCore(ctx, env)
}

func (vm *VM) runUnchecked(ctx context.Context, env *Environment) (runtime.Value, error) {
	result, err := vm.runCore(ctx, env)

	if err != nil {
		return nil, vm.state.wrapRuntimeError(err)
	}

	return result, nil
}

func (vm *VM) runCore(ctx context.Context, env *Environment) (runtime.Value, error) {
	if env == nil {
		env = noopEnv
	}

	if err := validate(vm.program); err != nil {
		return nil, err
	}

	if err := warmup(vm, env); err != nil {
		return nil, err
	}

	state := &vm.state
	state.reset(env)

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
				if state.applyProtected(err) == errReturn {
					return nil, err
				}
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
				if state.applyProtected(err) == errReturn {
					return nil, err
				}
			}

			if !has {
				state.pc = int(dst)
			}
		case bytecode.OpFail:
			if !dst.IsConstant() {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a constant string message")
				if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
					return nil, callErr
				}

				continue
			}

			idx := dst.Constant()
			if idx < 0 || idx >= len(constants) {
				callErr := runtime.Error(runtime.ErrInvalidOperation, "FAIL expects a valid constant string message")
				if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
					return nil, callErr
				}

				continue
			}

			msg, ok := constants[idx].(runtime.String)
			if !ok {
				callErr := runtime.TypeErrorOf(constants[idx], runtime.TypeString)
				if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
					return nil, callErr
				}

				continue
			}

			callErr := runtime.Error(runtime.ErrInvalidOperation, msg.String())
			if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
				return nil, callErr
			}

			continue
		case bytecode.OpHCall, bytecode.OpProtectedHCall:
			cacheFn := vm.cache.HostFunctions[state.pc-1]
			out, err := callCachedHostFunction(ctx, cacheFn, state.registers.Values, src1, src2)

			if state.setCallResult(op, dst, out, err) == errReturn {
				return nil, err
			}
		case bytecode.OpCall, bytecode.OpProtectedCall:
			if err := state.callUdf(op, dst, src1, src2); err != nil {
				if state.setCallResult(op, dst, runtime.None, err) == errReturn {
					return nil, err
				}
			}
		case bytecode.OpTailCall:
			if err := state.tailCallUdf(dst, src1, src2); err != nil {
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}
		case bytecode.OpDispatch:
			dispatcher, eventName, payload, options, err := vm.castDispatchArgs(ctx, reg[dst], reg[src1], reg[src2])

			if err != nil {
				if state.setOrTryCatch(dst, runtime.None, err) == errReturn {
					return nil, err
				}

				continue
			}

			out, err := dispatcher.Dispatch(ctx, runtime.DispatchEvent{
				Name:    eventName,
				Payload: payload,
				Options: options,
			})

			if out == nil {
				out = runtime.None
			}

			if state.setOrTryCatch(dst, out, err) == errReturn {
				return nil, err
			}
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
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}

			end, err := runtime.ToInt(ctx, reg[src2])

			if err != nil {
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}

			reg[dst] = runtime.NewRange(start, end)
		case bytecode.OpLoadIndex, bytecode.OpLoadIndexOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			if action, err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpLoadKey, bytecode.OpLoadKeyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptional
			arg := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			if action, err := vm.loadKeyAndSet(ctx, dst, state.pc-1, src, arg, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpLoadProperty, bytecode.OpLoadPropertyOptional:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptional
			prop := reg[src2]

			// TODO: inline loadIndexAndSet for better performance
			// I guess the reason it cannot inline is due to a different control flow
			if action, err := vm.loadPropertyAndSet(ctx, dst, state.pc-1, src, prop, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpLoadIndexConst, bytecode.OpLoadIndexOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadIndexOptionalConst
			arg := constants[src2.Constant()]

			if action, err := vm.loadIndexAndSet(ctx, dst, src, arg, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpLoadKeyConst, bytecode.OpLoadKeyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadKeyOptionalConst
			arg := constants[src2.Constant()]

			if action, err := vm.loadKeyConstAndSet(ctx, dst, state.pc-1, inst, src, arg, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpLoadPropertyConst, bytecode.OpLoadPropertyOptionalConst:
			src := reg[src1]
			optional := op == bytecode.OpLoadPropertyOptionalConst
			prop := constants[src2.Constant()]

			if action, err := vm.loadPropertyConstAndSet(ctx, dst, state.pc-1, inst, src, prop, optional); action != errOK {
				if action == errReturn {
					if state.applyProtected(err) == errReturn {
						return nil, err
					}
				}

				continue
			}
		case bytecode.OpPush:
			ds := reg[dst].(runtime.Appendable)

			if err := ds.Append(ctx, reg[src1]); err != nil {
				if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
					return nil, err
				}
			}
		case bytecode.OpPushKV:
			tr := reg[dst].(runtime.KeyWritable)

			if err := tr.Set(ctx, reg[src1], reg[src2]); err != nil {
				if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
					return nil, err
				}
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
					if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
						return nil, err
					}
				}

				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
				return nil, callErr
			}
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
					if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
						return nil, err
					}
				}
				continue
			}

			callErr := runtime.TypeErrorOf(reg[dst], runtime.TypeObject)
			if state.applyCatch(bytecode.NoopOperand, nil, callErr) == errReturn {
				return nil, callErr
			}
		case bytecode.OpIter:
			input := reg[src1]
			iterable, ok := input.(runtime.Iterable)

			if ok {
				iterator, err := iterable.Iterate(ctx)

				if err == nil {
					reg[dst] = data.NewIterator(iterator)
					continue
				}

				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}

			// TODO: replace with inlined version
			callErr := runtime.TypeErrorOf(input, runtime.TypeIterable)
			if state.applyCatch(dst, data.NoopIter, callErr) == errReturn {
				return nil, callErr
			}
		case bytecode.OpIterNext:
			iterator := reg[src1].(*data.Iterator)

			if err := iterator.Next(ctx); err != nil {
				if errors.Is(err, io.EOF) {
					state.pc = int(dst)
					continue
				}

				if state.applyProtected(err) == errReturn {
					return nil, err
				}
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
				if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
					return nil, err
				}

				continue
			}

			stream, err := observable.Subscribe(ctx, runtime.Subscription{
				EventName: eventName,
				Options:   options,
			})

			if err != nil {
				if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
					return nil, err
				}

				continue
			}

			reg[dst] = data.NewStreamValue(stream)
		case bytecode.OpStreamIter:
			stream := reg[src1].(*data.StreamValue)

			var timeout runtime.Int

			if reg[src2] != nil && reg[src2] != runtime.None {
				t, err := runtime.CastInt(reg[src2])

				if err != nil {
					if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
						return nil, err
					}

					t = 0
				}

				timeout = t
			}

			reg[dst] = stream.Iterate(timeout)
		case bytecode.OpQuery:
			src := readOperandValue(reg, constants, src1)
			descriptor := readOperandValue(reg, constants, src2)
			// TODO: unwrap since it's cannot be inlined
			out, err := applyQuery(ctx, src, descriptor)

			if action := state.setOrTryCatch(dst, out, err); action != errOK {
				if action == errReturn {
					return nil, err
				}

				continue
			}
		case bytecode.OpDataSet:
			reg[dst] = data.NewDataSet(src1 == 1)
		case bytecode.OpDataSetCollector:
			collectorType := bytecode.CollectorType(src1)

			if collectorType == bytecode.CollectorTypeAggregate || collectorType == bytecode.CollectorTypeAggregateGroup {
				planIdx := int(src2)

				if planIdx < 0 || planIdx >= len(aggregatePlans) {
					// TODO: is it really recoverable error?
					callErr := runtime.Errorf(runtime.ErrUnexpected, "invalid aggregate plan")
					if state.applyProtected(callErr) == errReturn {
						return nil, callErr
					}

					continue
				}

				plan := aggregatePlans[planIdx]

				if collectorType == bytecode.CollectorTypeAggregate {
					reg[dst] = data.NewAggregateCollector(plan)
				} else {
					reg[dst] = data.NewGroupedAggregateCollector(plan)
				}

				continue
			}

			reg[dst] = data.NewCollector(collectorType)
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

			if state.setOrTryCatch(dst, res, err) == errReturn {
				return nil, err
			}
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
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}
			reg[dst] = runtime.Divide(ctx, reg[src1], reg[src2])
		case bytecode.OpMod:
			if err := state.checkModuloByZero(ctx, reg[src2]); err != nil {
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}
			reg[dst] = runtime.Modulus(ctx, reg[src1], reg[src2])
		case bytecode.OpIncr:
			reg[dst] = runtime.Increment(ctx, reg[dst])
		case bytecode.OpDecr:
			reg[dst] = runtime.Decrement(ctx, reg[dst])
		case bytecode.OpNegate:
			reg[dst] = Negate(reg[src1])
		case bytecode.OpFlipPositive:
			reg[dst] = Positive(reg[src1])
		case bytecode.OpFlipNegative:
			reg[dst] = Negative(reg[src1])
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

			if state.applyProtected(err) == errReturn {
				return nil, err
			}

			continue
		case bytecode.OpRegexp:
			r, err := vm.regexpCached(state.pc-1, reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else {
				if state.applyCatch(dst, runtime.False, err) == errReturn {
					return nil, err
				}

				continue
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
					if state.applyCatch(dst, runtime.False, err) == errReturn {
						return nil, err
					}

					continue
				}

				reg[dst] = runtime.NewBoolean(length != 0)

				continue
			}

			reg[dst] = runtime.True
		case bytecode.OpAllEq, bytecode.OpAllNe, bytecode.OpAllGt, bytecode.OpAllGte, bytecode.OpAllLt, bytecode.OpAllLte, bytecode.OpAllIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAllEq))
			res, err := arrayAll(ctx, cmp, reg[src1], reg[src2])

			if state.setOrTryCatch(dst, res, err) == errReturn {
				return nil, err
			}
		case bytecode.OpAnyEq, bytecode.OpAnyNe, bytecode.OpAnyGt, bytecode.OpAnyGte, bytecode.OpAnyLt, bytecode.OpAnyLte, bytecode.OpAnyIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpAnyEq))
			res, err := arrayAny(ctx, cmp, reg[src1], reg[src2])

			if state.setOrTryCatch(dst, res, err) == errReturn {
				return nil, err
			}
		case bytecode.OpNoneEq, bytecode.OpNoneNe, bytecode.OpNoneGt, bytecode.OpNoneGte, bytecode.OpNoneLt, bytecode.OpNoneLte, bytecode.OpNoneIn:
			cmp := comparatorFromByte(int(op) - int(bytecode.OpNoneEq))
			res, err := arrayNone(ctx, cmp, reg[src1], reg[src2])

			if state.setOrTryCatch(dst, res, err) == errReturn {
				return nil, err
			}
		case bytecode.OpLength:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				length, err := val.Length(ctx)

				if err != nil {
					if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
						return nil, err
					}

					length = 0
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
			if state.applyCatch(dst, runtime.ZeroInt, callErr) == errReturn {
				return runtime.None, callErr
			}
		case bytecode.OpType:
			reg[dst] = runtime.NewString(runtime.TypeName(runtime.TypeOf(reg[src1])))
		case bytecode.OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				closeErr := val.Close()

				if closeErr != nil {
					if state.applyCatch(bytecode.NoopOperand, nil, closeErr) == errReturn {
						return nil, closeErr
					}

					continue
				}
			}
		case bytecode.OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				if state.applyCatch(bytecode.NoopOperand, nil, err) == errReturn {
					return nil, err
				}

				continue
			}

			if err := data.Sleep(ctx, dur); err != nil {
				if state.applyProtected(err) == errReturn {
					return nil, err
				}

				continue
			}
		case bytecode.OpRand:
			reg[dst] = runtime.NewFloat(runtime.RandomDefault())
		default:
			return nil, runtime.Errorf(runtime.ErrUnexpected, "unknown opcode %d at pc %d", op, state.pc-1)
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
