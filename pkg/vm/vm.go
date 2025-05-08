package vm

import (
	"context"
	"io"
	"strings"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm/internal"
)

type VM struct {
	env       *Environment
	program   *Program
	globals   map[string]runtime.Value
	registers []runtime.Value
	pc        int
}

func NewVM(program *Program) *VM {
	vm := new(VM)
	vm.program = program

	return vm
}

func (vm *VM) Run(ctx context.Context, opts []EnvironmentOption) (runtime.Value, error) {
	tryCatch := func(pos int) (Catch, bool) {
		for _, pair := range vm.program.CatchTable {
			if pos >= pair[0] && pos <= pair[1] {
				return pair, true
			}
		}

		return Catch{}, false
	}

	env := newEnvironment(opts)

	if err := vm.validateParams(env); err != nil {
		return nil, err
	}

	vm.env = env
	vm.registers = make([]runtime.Value, vm.program.Registers)
	vm.globals = make(map[string]runtime.Value)
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
		case OpStoreGlobal:
			vm.globals[constants[dst.Constant()].String()] = reg[src1]
		case OpLoadGlobal:
			reg[dst] = vm.globals[constants[src1.Constant()].String()]
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
		case OpJumpIfEmpty:
			val, ok := reg[src1].(runtime.Measurable)

			if ok {
				size, err := val.Length(ctx)

				if err != nil {
					return nil, err
				}

				if size == 0 {
					vm.pc = int(dst)
				}
			} else {
				// If the value is not measurable, we consider it empty
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
		case OpComp:
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
		case OpNotIn:
			reg[dst] = !internal.Contains(ctx, reg[src2], reg[src1])
		case OpLike:
			res, err := internal.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpNotLike:
			res, err := internal.Like(reg[src1], reg[src2])

			if err == nil {
				reg[dst] = !res
			} else {
				return nil, err
			}
		case OpRegexpPositive:
			// TODO: Add caching to avoid recompilation
			r, err := internal.ToRegexp(reg[src2])

			if err == nil {
				reg[dst] = r.Match(reg[src1])
			} else if _, catch := tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
				return nil, err
			}
		case OpRegexpNegative:
			// TODO: Add caching to avoid recompilation
			r, err := internal.ToRegexp(reg[src2])

			if err == nil {
				reg[dst] = !r.Match(reg[src1])
			} else if _, catch := tryCatch(vm.pc); catch {
				reg[dst] = runtime.False
			} else {
				return nil, err
			}
		case OpArray:
			var size int

			if src1 > 0 {
				size = src2.Register() - src1.Register() + 1
			}

			arr := runtime.NewArray(size)
			start := int(src1)
			end := int(src1) + size

			// Iterate over registers starting from src1 and up to the src2
			for i := start; i < end; i++ {
				_ = arr.Add(ctx, reg[i])
			}

			reg[dst] = arr
		case OpObject:
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
		case OpLoadProperty, OpLoadPropertyOptional:
			val := reg[src1]
			prop := reg[src2]

			switch getter := prop.(type) {
			case runtime.String:
				switch src := val.(type) {
				case runtime.Keyed:
					out, err := src.Get(ctx, getter)

					if err == nil {
						reg[dst] = out
					} else if op == OpLoadPropertyOptional {
						reg[dst] = runtime.None
					} else {
						return nil, err
					}
				default:
					if op != OpLoadPropertyOptional {
						return nil, runtime.TypeError(src, runtime.TypeMap)
					}

					reg[dst] = runtime.None
				}
			case runtime.Float, runtime.Int:
				// TODO: Optimize this. Avoid extra type conversion
				idx, _ := runtime.ToInt(ctx, getter)
				switch src := val.(type) {
				case runtime.Indexed:
					out, err := src.Get(ctx, idx)

					if err == nil {
						reg[dst] = out
					} else if op == OpLoadPropertyOptional {
						reg[dst] = runtime.None
					} else {
						return nil, err
					}
				case *internal.DataSet:
					reg[dst] = src.Get(ctx, idx)
				default:
					if op != OpLoadPropertyOptional {
						return nil, runtime.TypeError(src, runtime.TypeList)
					}

					reg[dst] = runtime.None
				}
			}
		case OpCall, OpProtectedCall:
			fnName := reg[dst].String()
			fn := vm.env.GetFunction(fnName)

			if fn == nil {
				if op == OpProtectedCall {
					reg[dst] = runtime.None
					continue
				} else {
					return nil, runtime.Error(ErrFunctionNotFound, fnName)
				}
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
				args[i-start] = reg[i]
			}

			out, err := fn(ctx, args...)

			if err == nil {
				reg[dst] = out
			} else if op == OpProtectedCall {
				reg[dst] = runtime.None
			} else if catch, ok := tryCatch(vm.pc); ok {
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
					if _, catch := tryCatch(vm.pc); catch {
						length = 0
					} else {
						return nil, err
					}
				}

				reg[dst] = length
			} else if _, catch := tryCatch(vm.pc); catch {
				reg[dst] = runtime.ZeroInt
			} else {
				return runtime.None, runtime.TypeError(reg[src1],
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
					if _, catch := tryCatch(vm.pc); !catch {
						return nil, err
					}
				}
			}
		case OpRange:
			res, err := internal.ToRange(ctx, reg[src1], reg[src2])

			if err == nil {
				reg[dst] = res
			} else {
				return nil, err
			}
		case OpLoopBegin:
			reg[dst] = internal.NewDataSet(src1 == 1)
		case OpLoopEnd:
			// TODO: Optimize this. Avoid extra type conversion
			ds, ok := reg[src1].(*internal.DataSet)

			if ok {
				reg[dst] = ds.ToList()
			} else {
				// Recover from an error
				reg[dst] = runtime.None
			}
		case OpLoopSkip:
			state := runtime.ToIntSafe(ctx, reg[dst])
			threshold := runtime.ToIntSafe(ctx, reg[src1])
			jump := int(src2)

			if state < threshold {
				state++
				reg[dst] = state
				vm.pc = jump
			}
		case OpLoopLimit:
			state := runtime.ToIntSafe(ctx, reg[dst])
			threshold := runtime.ToIntSafe(ctx, reg[src1])
			jump := int(src2)

			if state < threshold {
				state++
				reg[dst] = state
			} else {
				vm.pc = jump
			}
		case OpForLoopPrep:
			input := reg[src1]

			switch src := input.(type) {
			case runtime.Iterable:
				iterator, err := src.Iterate(ctx)

				if err != nil {
					return nil, err
				}

				reg[dst] = internal.NewIterator(iterator)
			default:
				if _, catch := tryCatch(vm.pc); catch {
					// Fall back to an empty iterator
					reg[dst] = internal.NoopIter
				} else {
					return nil, runtime.TypeError(src, runtime.TypeIterable)
				}
			}
		case OpForLoopNext:
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
		case OpForLoopValue:
			iterator := reg[src1].(*internal.Iterator)
			reg[dst] = iterator.Value()
		case OpForLoopKey:
			iterator := reg[src1].(*internal.Iterator)
			reg[dst] = iterator.Key()
		case OpWhileLoopPrep:
			reg[dst] = runtime.Int(-1)
		case OpWhileLoopNext:
			cond := runtime.ToBoolean(reg[src1])

			if cond {
				reg[dst] = internal.Increment(ctx, reg[dst])
			} else {
				vm.pc = int(src2)
			}
		case OpWhileLoopValue:
			reg[dst] = reg[src1]
		case OpLoopPush:
			ds := reg[dst].(*internal.DataSet)
			ds.Push(ctx, reg[src1])
		case OpLoopPushIter:
			ds := reg[dst].(*internal.DataSet)
			iterator := reg[src1].(*internal.Iterator)
			ds.Push(ctx, &internal.KeyValuePair{
				Key:   iterator.Key(),
				Value: iterator.Value(),
			})
		case OpLoopSequence:
			ds := reg[src1].(*internal.DataSet)
			reg[dst] = internal.NewSequence(ds.ToList())
		case OpSortPrep:
			reg[dst] = internal.NewStack(3)
		case OpSortPush:
			stack := reg[dst].(*internal.Stack)
			stack.Push(reg[src1])
		case OpSortPop:
			stack := reg[src1].(*internal.Stack)
			reg[dst] = stack.Pop()
		case OpSortValue:
			pair := reg[src1].(*internal.KeyValuePair)
			reg[dst] = pair.Value
		case OpSortKey:
			pair := reg[src1].(*internal.KeyValuePair)
			reg[dst] = pair.Key
		case OpSortSwap:
			ds := reg[dst].(*internal.DataSet)
			i, _ := runtime.ToInt(ctx, reg[src1])
			j, _ := runtime.ToInt(ctx, reg[src2])
			ds.Swap(ctx, i, j)
		case OpGroupPrep:
			reg[dst] = internal.NewCollector()
		case OpGroupAdd:
			collector := reg[dst].(*internal.Collector)
			key := reg[src1]
			value := reg[src2]
			collector.Add(key, value)
		case OpSleep:
			dur, err := runtime.ToInt(ctx, reg[dst])

			if err != nil {
				if _, catch := tryCatch(vm.pc); catch {
					continue
				} else {
					return nil, err
				}
			}

			if err := internal.Sleep(ctx, dur); err != nil {
				return nil, err
			}
		case OpReturn:
			break loop
		}
	}

	return vm.registers[NoopOperand], nil
}

func (vm *VM) validateParams(env *Environment) error {
	if len(vm.program.Params) == 0 {
		return nil
	}

	// There might be no errors.
	// Thus, we allocate this slice lazily, on a first error.
	var missedParams []string

	for _, n := range vm.program.Params {
		_, exists := env.params[n]

		if !exists {
			if missedParams == nil {
				missedParams = make([]string, 0, len(vm.program.Params))
			}

			missedParams = append(missedParams, "@"+n)
		}
	}

	if len(missedParams) > 0 {
		return runtime.Error(ErrMissedParam, strings.Join(missedParams, ", "))
	}

	return nil
}
