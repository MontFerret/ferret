package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

type ScopeProjection struct {
	registers *RegisterAllocator
	emitter   *Emitter
	symbols   *SymbolTable
	types     *TypeTracker
	values    []Variable
}

func NewScopeProjection(
	registers *RegisterAllocator,
	emitter *Emitter,
	symbols *SymbolTable,
	types *TypeTracker,
	scope []Variable,
) *ScopeProjection {
	return &ScopeProjection{
		registers: registers,
		emitter:   emitter,
		symbols:   symbols,
		types:     types,
		values:    scope,
	}
}

func (sp *ScopeProjection) EmitAsArray(dst bytecode.Operand) {
	buildDst := dst

	if sp.usesRegister(dst) {
		// Avoid overwriting a value we're about to project.
		buildDst = sp.registers.Allocate()
	}

	sp.emitter.EmitArray(buildDst, len(sp.values))
	sp.types.Set(buildDst, TypeArray)

	for _, v := range sp.values {
		sp.emitter.EmitArrayPush(buildDst, sp.projectedValue(v))
	}

	if buildDst != dst {
		sp.moveProjectedValue(dst, buildDst)
	}
}

func (sp *ScopeProjection) EmitAsObject(dst bytecode.Operand) {
	size := len(sp.values)
	buildDst := dst

	if sp.usesRegister(dst) {
		// Avoid overwriting a value we're about to project.
		buildDst = sp.registers.Allocate()
	}

	sp.emitter.EmitObject(buildDst, size)
	sp.types.Set(buildDst, TypeObject)

	if size == 0 {
		if buildDst != dst {
			sp.moveProjectedValue(dst, buildDst)
		}

		return
	}

	for _, v := range sp.values {
		// Key (field name)
		keyConst := sp.symbols.AddConstant(runtime.String(v.Name))
		// Set the key-value pair in the object.
		// buildDst may differ from dst when aliasing is detected, so values can be used directly.
		sp.emitter.EmitObjectSetConst(buildDst, keyConst, sp.projectedValue(v))
	}

	if buildDst != dst {
		sp.moveProjectedValue(dst, buildDst)
	}
}

func (sp *ScopeProjection) moveProjectedValue(dst, src bytecode.Operand) {
	if sp.types.Resolve(src).IsUntracked() {
		sp.emitter.EmitPlainMove(dst, src)
	} else {
		sp.emitter.EmitMoveTracked(dst, src)
	}

	sp.types.Set(dst, sp.types.Resolve(src))
}

func (sp *ScopeProjection) RestoreFromArray(src bytecode.Operand) {
	idx := sp.registers.Allocate()

	for i, v := range sp.values {
		sp.emitter.EmitLoadConst(idx, sp.symbols.AddConstant(runtime.Int(i)))

		if v.Storage == BindingStorageCell {
			tmp := sp.registers.Allocate()
			variable, _ := sp.symbols.DeclareLocalWithOptions(v.Name, v.Type, BindingOptions{Mutable: v.Mutable, Storage: v.Storage})
			sp.emitter.EmitABC(bytecode.OpLoadIndex, tmp, src, idx)
			sp.emitter.EmitMakeCell(variable, tmp)
			sp.types.Set(tmp, v.Type)
			sp.types.Set(variable, TypeAny)
			continue
		}

		variable, _ := sp.symbols.DeclareLocalWithOptions(v.Name, v.Type, BindingOptions{Mutable: v.Mutable, Storage: v.Storage})
		sp.emitter.EmitABC(bytecode.OpLoadIndex, variable, src, idx)
	}
}

func (sp *ScopeProjection) RestoreFromObject(src bytecode.Operand) {
	key := sp.registers.Allocate()

	for _, v := range sp.values {
		sp.emitter.EmitLoadConst(key, sp.symbols.AddConstant(runtime.String(v.Name)))

		if v.Storage == BindingStorageCell {
			tmp := sp.registers.Allocate()
			variable, _ := sp.symbols.DeclareLocalWithOptions(v.Name, v.Type, BindingOptions{Mutable: v.Mutable, Storage: v.Storage})
			sp.emitter.EmitABC(bytecode.OpLoadKey, tmp, src, key)
			sp.emitter.EmitMakeCell(variable, tmp)
			sp.types.Set(tmp, v.Type)
			sp.types.Set(variable, TypeAny)
			continue
		}

		variable, _ := sp.symbols.DeclareLocalWithOptions(v.Name, v.Type, BindingOptions{Mutable: v.Mutable, Storage: v.Storage})
		sp.emitter.EmitABC(bytecode.OpLoadKey, variable, src, key)
	}
}

func (sp *ScopeProjection) projectedValue(v Variable) bytecode.Operand {
	if v.Storage != BindingStorageCell {
		return v.Register
	}

	dst := sp.registers.Allocate()
	sp.emitter.EmitLoadCell(dst, v.Register)
	sp.types.Set(dst, v.Type)

	return dst
}

func (sp *ScopeProjection) usesRegister(reg bytecode.Operand) bool {
	for _, v := range sp.values {
		if v.Register == reg {
			return true
		}
	}

	return false
}
