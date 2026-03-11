package vm

import (
	"context"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/data"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/diagnostic"
	"github.com/MontFerret/ferret/v2/pkg/vm/internal/mem"
)

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
				return nil, runtime.ErrNotFound
			}
			if val, ok := obj.SlotValue(inst.InlineSlot); ok {
				return val, nil
			}

			return nil, runtime.ErrNotFound
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

	if pc < 0 || pc >= len(vm.cache.LoadKeyICs) {
		return vm.loadKey(ctx, obj, arg)
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
		return nil, diagnostic.MemberAccessErrorOf(src, diagnostic.MemberAccessIndex, arg)
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
		return nil, diagnostic.MemberAccessErrorOf(src, diagnostic.MemberAccessProperty, arg)
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
