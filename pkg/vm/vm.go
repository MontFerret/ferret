package vm

import (
	"context"
	"io"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal"
)

type VM struct {
	env       *Environment
	program   *Program
	registers []runtime.Value
	pc        int
}

func New(program *Program) *VM {
	vm := new(VM)
	vm.program = program

	return vm
}

func (vm *VM) Run(ctx context.Context, opts []EnvironmentOption) (runtime.Value, error) {
	env := newEnvironment(opts)

	if err := validate(env, vm.program); err != nil {
		return nil, err
	}

	vm.env = env
	vm.registers = make([]runtime.Value, vm.program.Registers)
	vm.pc = 0

	bytecode := vm.program.Bytecode
	constants := vm.program.Constants
	reg := vm.registers

loop:
	for vm.pc < len(bytecode) {
		inst := bytecode[vm.pc]
		op := inst.Opcode
		dst, src1, src2 := inst.Operands[0], inst.Operands[1], inst.Operands[2]
		vm.pc++

		switch op {
		case OpLoadNone:
			reg[dst] = runtime.None
		case OpLoadZero:
			reg[dst] = runtime.ZeroInt
		case OpLoadBool:
			reg[dst] = runtime.Boolean(src1 == 1)
		case OpMove:
			reg[dst] = reg[src1]
		case OpLoadConst:
			reg[dst] = constants[src1.Constant()]
		case OpLoadParam:
			name := constants[src1.Constant()]
			reg[dst] = vm.env.params[name.String()]
		case OpJump:
			vm.pc = int(dst)
		case OpJumpIfFalse:
			if !runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpJumpIfTrue:
			if runtime.ToBoolean(reg[src1]) {
				vm.pc = int(dst)
			}
		case OpAdd:
			reg[dst] = internal.Add(ctx, reg[src1], reg[src2])
		case OpSub:
			reg[dst] = internal.Subtract(ctx, reg[src1], reg[src2])
		case OpMulti:
			reg[dst] = internal.Multiply(ctx, reg[src1], reg[src2])
		case OpDiv:
			reg[dst] = internal.Divide(ctx, reg[src1], reg[src2])
		case OpMod:
			reg[dst] = internal.Modulus(ctx, reg[src1], reg[src2])
		case OpIncr:
			reg[dst] = internal.Increment(ctx, reg[dst])
		case OpDecr:
			reg[dst] = internal.Decrement(ctx, reg[dst])
		case OpCastBool:
			reg[dst] = runtime.ToBoolean(reg[src1])
		case OpNegate:
			reg[dst] = runtime.Negate(reg[src1])
		case OpFlipPositive:
			reg[dst] = runtime.Positive(reg[src1])
		case OpFlipNegative:
			reg[dst] = runtime.Negative(reg[src1])
		case OpCmp:
			reg[dst] = runtime.Int(runtime.CompareValues(reg[src1], reg[src2]))
		case OpNot:
			reg[dst] = !runtime.ToBoolean(reg[src1])
		case OpEq:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) == 0)
		case OpNeq:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) != 0)
		case OpGt:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) > 0)
		case OpLt:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) < 0)
		case OpGte:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) >= 0)
		case OpLte:
			reg[dst] = runtime.Boolean(runtime.CompareValues(reg[src1], reg[src2]) <= 0)
		case OpIn:
			reg[dst] = internal.Contains(ctx, reg[src2], reg[src1])
		case OpLike:
			res, err := internal.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpRegexpPositive:
			// TODO: Add caching to avoid recompilation
			r, err := internal.ToRegexp(reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
				return nil, err
			}
		case OpLoadArray:
			var size int

			if src1 > 0 {
				size = src2.Register() - src1.Register() + 1
			}

			arr := runtime.NewArray(size)
			start := int(src1)
			end := int(src1) + size

			for i := start; i < end; i++ {
				_ = arr.Add(ctx, reg[i])
			}

			reg[dst] = arr
		case OpLoadObject:
			obj := runtime.NewObject()
			var args int

			if src1 > 0 {
				args = src2.Register() - src1.Register() + 1
			}

			start := int(src1)
			end := int(src1) + args

			for i := start; i < end; i += 2 {
				key := reg[i]
				value := reg[i+1]

				_ = obj.Set(ctx, runtime.ToString(key), value)
			}

			reg[dst] = obj
		case OpLoadIndex, OpLoadIndexOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadIndex(ctx, src, arg)

			if err == nil {
				reg[dst] = out
			} else if op == OpLoadIndexOptional {
				reg[dst] = runtime.None
			} else {
				return nil, err
			}

		case OpLoadKey, OpLoadKeyOptional:
			src := reg[src1]
			arg := reg[src2]
			out, err := vm.loadKey(ctx, src, arg)

			if err == nil {
				reg[dst] = out
			} else if op == OpLoadKeyOptional {
				reg[dst] = runtime.None
			} else {
				return nil, err
			}

		case OpLoadProperty, OpLoadPropertyOptional:
			src := reg[src1]
			prop := reg[src2]

			var out runtime.Value
			var err error

			switch getter := prop.(type) {
			case runtime.String:
				out, err = vm.loadKey(ctx, src, getter)
			case runtime.Float, runtime.Int:
				out, err = vm.loadIndex(ctx, src, getter)
			default:
				out, err = vm.loadKey(ctx, src, runtime.ToString(prop))
			}

			if err == nil {
				reg[dst] = out
			} else if op == OpLoadPropertyOptional {
				reg[dst] = runtime.None
			} else {
				return nil, err
			}
		case OpCall, OpProtectedCall:
			out, err := vm.call(ctx, dst, src1, src2)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}
		case OpCall0, OpProtectedCall0:
			out, err := vm.call0(ctx, dst)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall0 {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}

		case OpCall1, OpProtectedCall1:
			out, err := vm.call1(ctx, dst, src1)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall1 {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}

		case OpCall2, OpProtectedCall2:
			out, err := vm.call2(ctx, dst, src1, src2)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall2 {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}

		case OpCall3, OpProtectedCall3:
			out, err := vm.call3(ctx, dst, src1, src2)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall3 {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}

		case OpCall4, OpProtectedCall4:
			out, err := vm.call4(ctx, dst, src1, src2)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall4 {
				reg[dst] = runtime.None
			} else if catch, ok := vm.tryCatch(vm.pc); ok {
				reg[dst] = runtime.None

				if catch[2] > 0 {
					vm.pc = catch[2]
				}
			} else {
				return nil, err
			}

		case OpLength:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				length, err := val.Length(ctx)

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); catch {
						length = 0
					} else {
						return nil, err
					}
				}

				reg[dst] = length
			} else if _, catch := vm.tryCatch(vm.pc); catch {
				reg[dst] = runtime.ZeroInt
			} else {
				return runtime.None, runtime.TypeErrorOf(reg[src1],
					runtime.TypeString,
					runtime.TypeList,
					runtime.TypeMap,
					runtime.TypeBinary,
					runtime.TypeMeasurable,
				)
			}
		case OpType:
			reg[dst] = runtime.String(runtime.Reflect(reg[src1]))
		case OpClose:
			val, ok := reg[dst].(io.Closer)
			reg[dst] = runtime.None

			if ok {
				err := val.Close()

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); !catch {
						return nil, err
					}
				}
			}
		case OpLoadRange:
			res, err := internal.ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpDataSet:
			reg[dst] = internal.NewDataSet(src1 == 1)
		case OpDataSetSorter:
			reg[dst] = internal.NewSorter(runtime.SortDirection(src1))
		case OpDataSetMultiSorter:
			encoded := src1.Register()
			count := src2.Register()

			reg[dst] = internal.NewMultiSorter(runtime.DecodeSortDirections(encoded, count))
		case OpDataSetCollector:
			reg[dst] = internal.NewCollector(internal.CollectorType(src1))
		case OpPush:
			ds := reg[dst].(runtime.List)

			if err := ds.Add(ctx, reg[src1]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}
		case OpPushKV:
			tr := reg[dst].(internal.Transformer)

			if err := tr.Add(ctx, reg[src1], reg[src2]); err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				}

				return nil, err
			}
		case OpIter:
			input := reg[src1]

			switch src := input.(type) {
			case runtime.Iterable:
				iterator, err := src.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				reg[dst] = internal.NewIterator(iterator)
			default:
				if _, catch := vm.tryCatch(vm.pc); catch {
					// Fall back to an empty iterator
					reg[dst] = internal.NoopIter
				} else {
					return nil, runtime.TypeErrorOf(src, runtime.TypeIterable)
				}
			}
		case OpIterNext:
			iterator := reg[src1].(*internal.Iterator)
			hasNext, err := iterator.HasNext(ctx)

			if err != nil {
				return nil, err
			}

			if hasNext {
				if err := iterator.Next(ctx); err != nil {
					return nil, err
				}
			} else {
				vm.pc = int(dst)
			}
		case OpIterValue:
			iterator := reg[src1].(*internal.Iterator)
			reg[dst] = iterator.Value()
		case OpIterKey:
			iterator := reg[src1].(*internal.Iterator)
			reg[dst] = iterator.Key()
		case OpIterSkip:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
				vm.pc = int(dst)
			}
		case OpIterLimit:
			state := runtime.ToIntSafe(ctx, reg[src1])
			threshold := runtime.ToIntSafe(ctx, reg[src2])

			if state < threshold {
				state++
				reg[src1] = state
			} else {
				vm.pc = int(dst)
			}
		case OpStream:
			observable, eventName, options, err := vm.castSubscribeArgs(reg[dst], reg[src1], reg[src2])

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			stream, err := observable.Subscribe(ctx, runtime.Subscription{
				EventName: eventName,
				Options:   options,
			})

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			reg[dst] = internal.NewStreamValue(stream)
		case OpStreamIter:
			stream := reg[src1].(*internal.StreamValue)

			var timeout runtime.Int

			if reg[src2] != nil && reg[src2] != runtime.None {
				t, err := runtime.CastInt(reg[src1])

				if err != nil {
					if _, catch := vm.tryCatch(vm.pc); catch {
						continue
					} else {
						return nil, err
					}
				}

				timeout = t
			}

			reg[dst] = stream.Iterate(timeout)
		case OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				if _, catch := vm.tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			if err := internal.Sleep(ctx, dur); err != nil {
				return nil, err
			}
		case OpReturn:
			reg[NoopOperand] = reg[dst]

			break loop
		}
	}

	return vm.registers[NoopOperand], nil
}

func (vm *VM) tryCatch(pos int) (Catch, bool) {
	for _, pair := range vm.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return Catch{}, false
}

func (vm *VM) call(ctx context.Context, dst, src1, src2 Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F().Get(fnName)

	if !found {
		return nil, runtime.Error(ErrFunctionNotFound, fnName)
	}

	var size int

	if src1 > 0 {
		size = src2.Register() - src1.Register() + 1
	}

	start := int(src1)
	end := int(src1) + size
	args := make([]runtime.Value, size)

	// Iterate over registers starting from src1 and up to the src2
	for i := start; i < end; i++ {
		args[i-start] = vm.registers[i]
	}

	return fn(ctx, args...)
}

func (vm *VM) call0(ctx context.Context, dst Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F0().Get(fnName)

	if found {
		return fn(ctx)
	}

	// Fall back to a variadic function call
	return vm.call(ctx, dst, NoopOperand, NoopOperand)
}

func (vm *VM) call1(ctx context.Context, dst, src1 Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F1().Get(fnName)

	if found {
		return fn(ctx, vm.registers[src1])
	}

	// Fall back to a variadic function call
	return vm.call(ctx, dst, src1, src1)
}

func (vm *VM) call2(ctx context.Context, dst, src1, src2 Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F2().Get(fnName)

	if found {
		return fn(ctx, vm.registers[src1], vm.registers[src2])
	}

	// Fall back to a variadic function call
	return vm.call(ctx, dst, src1, src2)
}

func (vm *VM) call3(ctx context.Context, dst, src1, src2 Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F3().Get(fnName)

	if found {
		arg1 := vm.registers[src1]
		arg2 := vm.registers[src1+1]
		arg3 := vm.registers[src1+2]

		return fn(ctx, arg1, arg2, arg3)
	}

	// Fall back to a variadic function call
	return vm.call(ctx, dst, src1, src2)
}

func (vm *VM) call4(ctx context.Context, dst, src1, src2 Operand) (runtime.Value, error) {
	fnName := vm.registers[dst].String()
	fn, found := vm.env.functions.F4().Get(fnName)

	if found {
		arg1 := vm.registers[src1]
		arg2 := vm.registers[src1+1]
		arg3 := vm.registers[src1+2]
		arg4 := vm.registers[src1+3]

		return fn(ctx, arg1, arg2, arg3, arg4)
	}

	// Fall back to a variadic function call
	return vm.call(ctx, dst, src1, src2)
}

func (vm *VM) loadIndex(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	indexed, ok := src.(runtime.Indexed)

	if !ok {
		return nil, runtime.TypeErrorOf(src, runtime.TypeIndexed)
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

	return indexed.Get(ctx, idx)
}

func (vm *VM) loadKey(ctx context.Context, src, arg runtime.Value) (runtime.Value, error) {
	keyed, ok := src.(runtime.Keyed)

	if !ok {
		return nil, runtime.TypeErrorOf(src, runtime.TypeKeyed)
	}

	out, err := keyed.Get(ctx, arg)

	if err != nil {
		return nil, err
	}

	return out, nil
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
