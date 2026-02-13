package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

func (vm *VM) tryCatch(pos int) (Catch, bool) {
	for _, pair := range vm.program.CatchTable {
		if pos >= pair[0] && pos <= pair[1] {
			return pair, true
		}
	}

	return Catch{}, false
}

func (vm *VM) callv(ctx context.Context, pc int, src1, src2 Operand) (runtime.Value, error) {
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

func (vm *VM) call1(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
	reg := vm.registers.Values
	arg := reg[src1]
	cacheFn := vm.cache.Functions[pc]

	if cacheFn.Fn1 != nil {
		return cacheFn.Fn1(ctx, arg)
	}

	// Fall back to a variadic function call
	return cacheFn.FnV(ctx, arg)
}

func (vm *VM) call2(ctx context.Context, pc int, src1, src2 Operand) (runtime.Value, error) {
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

func (vm *VM) call3(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
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

func (vm *VM) call4(ctx context.Context, pc int, src1 Operand) (runtime.Value, error) {
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

func (vm *VM) loadKeyConstCached(ctx context.Context, pc int, inst *Instruction, src, arg runtime.Value) (runtime.Value, error) {
	obj, ok := src.(*data.FastObject)

	if !ok {
		return vm.loadKey(ctx, src, arg)
	}

	shapeID := obj.ShapeID()

	if shapeID != 0 {
		if inst != nil && inst.inlineShapeID == shapeID {
			if inst.inlineSlot < 0 {
				return nil, runtime.ErrNotFound
			}
			if val, ok := obj.SlotValue(inst.inlineSlot); ok {
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
					inst.inlineShapeID = shapeID
					inst.inlineSlot = slot
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
				inst.inlineShapeID = shapeID
				inst.inlineSlot = -1
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
			inst.inlineShapeID = shapeID
			inst.inlineSlot = slot
		}

		return val, nil
	}

	return vm.loadKey(ctx, src, arg)
}

func (vm *VM) objectSetConstCached(inst *Instruction, obj *data.FastObject, key runtime.String, value runtime.Value) {
	if obj == nil {
		return
	}

	if inst != nil {
		shape := obj.Shape()

		if shape != nil && inst.inlineSetShape == shape {
			if obj.SetSlotWithShape(inst.inlineSetNextShape, inst.inlineSlot, value) {
				return
			}
		}
	}

	prev, next, slot, ok := obj.SetStringCached(string(key), value)

	if ok && inst != nil {
		inst.inlineSetShape = prev
		inst.inlineSetNextShape = next
		inst.inlineSlot = slot
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
		return nil, runtime.TypeErrorOf(src, runtime.TypeKeyReadable)
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

func (vm *VM) setOrTryCatch(dst Operand, val runtime.Value, err error) error {
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

func (vm *VM) setCallResult(op Opcode, dst Operand, out runtime.Value, err error) error {
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

func (vm *VM) setOrOptional(dst Operand, val runtime.Value, err error, optional bool) error {
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

func isProtectedCall(op Opcode) bool {
	switch op {
	case OpProtectedCall, OpProtectedCall0, OpProtectedCall1, OpProtectedCall2, OpProtectedCall3, OpProtectedCall4:
		return true
	default:
		return false
	}
}
