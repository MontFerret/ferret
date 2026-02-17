package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func (vm *VM) tryCatch(pos int) (bytecode.Catch, bool) {
	for _, pair := range vm.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return bytecode.Catch{}, false
}

func (vm *VM) callv(ctx context.Context, pc int, src1, src2 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]

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

	return cacheFn.FnV(ctx, args...)
}

func (vm *VM) call0(ctx context.Context, pc int) (runtime.Value, error) {
	cacheFn := vm.cache.Functions[pc]

	if cacheFn.Fn0 != nil {
		return cacheFn.Fn0(ctx)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx)
}

func (vm *VM) call1(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	arg := reg[src1]
	cacheFn := vm.cache.Functions[pc]

	if cacheFn.Fn1 != nil {
		return cacheFn.Fn1(ctx, arg)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg)
}

func (vm *VM) call2(ctx context.Context, pc int, src1, src2 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
	arg1 := reg[src1]
	arg2 := reg[src2]

	if cacheFn.Fn2 != nil {
		return cacheFn.Fn2(ctx, arg1, arg2)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2)
}

func (vm *VM) call3(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
	arg1 := reg[src1]
	arg2 := reg[src1+1]
	arg3 := reg[src1+2]

	if cacheFn.Fn3 != nil {
		return cacheFn.Fn3(ctx, arg1, arg2, arg3)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2, arg3)
}

func (vm *VM) call4(ctx context.Context, pc int, src1 bytecode.Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	cacheFn := vm.cache.Functions[pc]
	arg1 := reg[src1]
	arg2 := reg[src1+1]
	arg3 := reg[src1+2]
	arg4 := reg[src1+3]

	if cacheFn.Fn4 != nil {
		return cacheFn.Fn4(ctx, arg1, arg2, arg3, arg4)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg1, arg2, arg3, arg4)
}

func (vm *VM) applyQuery(ctx context.Context, reg []runtime.Value, src1 bytecode.Operand, constants []runtime.Value, src2 bytecode.Operand, dst bytecode.Operand) error {
	src := reg[src1]

	if src1.IsConstant() {
		src = constants[src1.Constant()]
	}

	var arg runtime.Value

	if src2.IsConstant() {
		arg = constants[src2.Constant()]
	} else {
		arg = reg[src2]
	}

	var query runtime.Query

	switch value := arg.(type) {
	case runtime.Query:
		query = value
	case *runtime.Array:
		length, err := value.Length(ctx)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		if length < 2 {
			if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(arg, runtime.TypeQuery)); err != nil {
				return err
			}

			break
		}

		kindVal, err := value.Get(ctx, runtime.NewInt(0))
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		payloadVal, err := value.Get(ctx, runtime.NewInt(1))
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
				return err
			}

			break
		}

		var paramsVal runtime.Value = runtime.None
		if length > 2 {
			paramsVal, err = value.Get(ctx, runtime.NewInt(2))
			if err != nil {
				if err := vm.setOrTryCatch(dst, runtime.None, err); err != nil {
					return err
				}

				break
			}
		}

		kind, err := runtime.CastString(kindVal)
		if err != nil {
			if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(kindVal, runtime.TypeString)); err != nil {
				return err
			}

			break
		}

		payload := runtime.EmptyString
		if payloadVal != runtime.None {
			payload, err = runtime.CastString(payloadVal)
			if err != nil {
				if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(payloadVal, runtime.TypeString, runtime.TypeNone)); err != nil {
					return err
				}

				break
			}
		}

		query = runtime.NewQuery(kind, payload)
		query.Params = paramsVal
	default:
		if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(arg, runtime.TypeQuery, runtime.TypeArray)); err != nil {
			return err
		}

		return nil
	}

	queryable, ok := src.(runtime.Queryable)

	if !ok {
		if err := vm.setOrTryCatch(dst, runtime.None, runtime.TypeErrorOf(src, runtime.TypeQueryable)); err != nil {
			return err
		}

		return nil
	}

	res, err := queryable.Query(ctx, query)

	if err := vm.setOrTryCatch(dst, res, err); err != nil {
		return err
	}

	return nil
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
		return nil, runtime.TypeErrorOf(value, runtime.TypeString, runtime.TypeRegexp)
	}
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

	shapeID := obj.ShapeID()
	if shapeID != 0 {
		if pc < 0 || pc >= len(vm.cache.LoadKeyICs) {
			return vm.loadKey(ctx, src, arg)
		}

		cache := vm.cache.LoadKeyICs[pc]
		if cache != nil {
			if slot, ok := cache.Lookup(shapeID, key); ok {
				if slot < 0 {
					return nil, runtime.ErrNotFound
				}
				if val, ok := obj.SlotValue(slot); ok {
					return val, nil
				}

				return nil, runtime.ErrNotFound
			}
		}

		slot, ok := obj.LookupSlot(key)

		if !ok {
			if cache == nil {
				cache = mem.NewLoadKeyCache()
				vm.cache.LoadKeyICs[pc] = cache
			}

			cache.Add(shapeID, key, -1)

			return nil, runtime.ErrNotFound
		}

		val, ok := obj.SlotValue(slot)

		if !ok {
			return nil, runtime.ErrNotFound
		}

		if cache == nil {
			cache = mem.NewLoadKeyCache()
			vm.cache.LoadKeyICs[pc] = cache
		}

		cache.Add(shapeID, key, slot)

		return val, nil
	}

	return vm.loadKey(ctx, src, arg)
}

func (vm *VM) loadKeyConstCached(ctx context.Context, pc int, inst *data.ExecInstruction, src, arg runtime.Value) (runtime.Value, error) {
	obj, ok := src.(*data.FastObject)

	if !ok {
		return vm.loadKey(ctx, src, arg)
	}

	shapeID := obj.ShapeID()

	if shapeID != 0 {
		if inst != nil && inst.InlineShapeID == shapeID {
			if inst.InlineSlot < 0 {
				return nil, runtime.ErrNotFound
			}
			if val, ok := obj.SlotValue(inst.InlineSlot); ok {
				return val, nil
			}

			return nil, runtime.ErrNotFound
		}

		if pc < 0 || pc >= len(vm.cache.LoadKeyConstICs) {
			return vm.loadKey(ctx, src, arg)
		}

		cache := vm.cache.LoadKeyConstICs[pc]

		if cache != nil {
			if slot, ok := cache.Lookup(shapeID); ok {
				if inst != nil {
					inst.InlineShapeID = shapeID
					inst.InlineSlot = slot
				}

				if slot < 0 {
					return nil, runtime.ErrNotFound
				}
				if val, ok := obj.SlotValue(slot); ok {
					return val, nil
				}

				return nil, runtime.ErrNotFound
			}
		}

		var key string

		switch v := arg.(type) {
		case runtime.String:
			key = string(v)
		default:
			key = runtime.ToString(v).String()
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

			return nil, runtime.ErrNotFound
		}

		val, ok := obj.SlotValue(slot)

		if !ok {
			return nil, runtime.ErrNotFound
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

	return vm.loadKey(ctx, src, arg)
}

func (vm *VM) objectSetConstCached(inst *data.ExecInstruction, obj *data.FastObject, key runtime.String, value runtime.Value) {
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
		return nil, runtime.TypeErrorOf(src, runtime.TypeIndexReadable)
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
	keyed, ok := src.(runtime.KeyReadable)

	if !ok {
		// Try a more expensive using reflection for structs and maps that don't implement KeyReadable but use ferret tags or map keys.
		return runtime.EncodeField(ctx, src, arg)
	}

	out, err := keyed.Get(ctx, arg)

	if err != nil {
		return nil, err
	}

	return out, nil
}

func (vm *VM) loadIndexAndSet(ctx context.Context, dst bytecode.Operand, src, arg runtime.Value, optional bool) error {
	if optional && src == runtime.None {
		vm.registers.Values[dst] = runtime.None
		return nil
	}

	out, err := vm.loadIndex(ctx, src, arg)

	return vm.setOrOptional(dst, out, err, optional)
}

func (vm *VM) loadKeyAndSet(ctx context.Context, dst bytecode.Operand, pc int, src, arg runtime.Value, optional bool) error {
	if optional && src == runtime.None {
		vm.registers.Values[dst] = runtime.None
		return nil
	}

	out, err := vm.loadKeyCached(ctx, pc, src, arg)

	return vm.setOrOptional(dst, out, err, optional)
}

func (vm *VM) loadKeyConstAndSet(ctx context.Context, dst bytecode.Operand, pc int, inst *data.ExecInstruction, src, arg runtime.Value, optional bool) error {
	if optional && src == runtime.None {
		vm.registers.Values[dst] = runtime.None
		return nil
	}

	out, err := vm.loadKeyConstCached(ctx, pc, inst, src, arg)

	return vm.setOrOptional(dst, out, err, optional)
}

func (vm *VM) loadPropertyAndSet(ctx context.Context, dst bytecode.Operand, pc int, src, prop runtime.Value, optional bool) error {
	if optional && src == runtime.None {
		vm.registers.Values[dst] = runtime.None
		return nil
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

	return vm.setOrOptional(dst, out, err, optional)
}

func (vm *VM) loadPropertyConstAndSet(ctx context.Context, dst bytecode.Operand, pc int, inst *data.ExecInstruction, src, prop runtime.Value, optional bool) error {
	if optional && src == runtime.None {
		vm.registers.Values[dst] = runtime.None
		return nil
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

	return vm.setOrOptional(dst, out, err, optional)
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
		return nil, "", nil, nil, runtime.TypeErrorOf(target, runtime.TypeDispatcher)
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

func (vm *VM) setOrTryCatch(dst bytecode.Operand, val runtime.Value, err error) error {
	reg := vm.registers.Values

	if err == nil {
		reg[dst] = val

		return nil
	}

	if _, catch := vm.tryCatch(vm.pc); catch {
		reg[dst] = runtime.None

		return nil
	}

	return err
}

func (vm *VM) setCallResult(op bytecode.Opcode, dst bytecode.Operand, out runtime.Value, err error) error {
	reg := vm.registers.Values

	if err == nil {
		reg[dst] = out

		return nil
	}

	if isProtectedCall(op) {
		reg[dst] = runtime.None

		return nil
	}

	if catch, ok := vm.tryCatch(vm.pc); ok {
		reg[dst] = runtime.None

		if catch[2] > 0 {
			vm.pc = catch[2]
		}

		return nil
	}

	return err
}

func (vm *VM) setOrOptional(dst bytecode.Operand, val runtime.Value, err error, optional bool) error {
	if err == nil {
		vm.registers.Values[dst] = val

		return nil
	}

	if optional {
		vm.registers.Values[dst] = runtime.None

		return nil
	}

	return err
}

func isProtectedCall(op bytecode.Opcode) bool {
	switch op {
	case bytecode.OpProtectedCall, bytecode.OpProtectedCall0, bytecode.OpProtectedCall1, bytecode.OpProtectedCall2, bytecode.OpProtectedCall3, bytecode.OpProtectedCall4:
		return true
	default:
		return false
	}
}
