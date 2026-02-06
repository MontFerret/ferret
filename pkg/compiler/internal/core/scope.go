package core

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type ScopeProjection struct {
	registers *RegisterAllocator
	emitter   *Emitter
	symbols   *SymbolTable
	values    []Variable
}

func NewScopeProjection(
	registers *RegisterAllocator,
	emitter *Emitter,
	symbols *SymbolTable,
	scope []Variable,
) *ScopeProjection {
	return &ScopeProjection{
		registers: registers,
		emitter:   emitter,
		symbols:   symbols,
		values:    scope,
	}
}

func (sp *ScopeProjection) EmitAsArray(dst vm.Operand) {
	buildDst := dst

	if sp.usesRegister(dst) {
		// Avoid overwriting a value we're about to project.
		buildDst = sp.registers.Allocate()
	}

	sp.emitter.EmitArray(buildDst, len(sp.values))

	for _, v := range sp.values {
		sp.emitter.EmitArrayPush(buildDst, v.Register)
	}

	if buildDst != dst {
		sp.emitter.EmitMove(dst, buildDst)
	}
}

func (sp *ScopeProjection) EmitAsObject(dst vm.Operand) {
	size := len(sp.values)
	buildDst := dst

	if sp.usesRegister(dst) {
		// Avoid overwriting a value we're about to project.
		buildDst = sp.registers.Allocate()
	}

	sp.emitter.EmitObject(buildDst, size)

	if size == 0 {
		if buildDst != dst {
			sp.emitter.EmitMove(dst, buildDst)
		}

		return
	}

	for _, v := range sp.values {
		// Key (field name)
		keyReg := sp.registers.Allocate()
		sp.emitter.EmitLoadConst(keyReg, sp.symbols.AddConstant(runtime.String(v.Name)))

		// Value (actual variable value)
		valReg := sp.registers.Allocate()
		sp.emitter.EmitAB(vm.OpMove, valReg, v.Register)

		// Set the key-value pair in the object
		sp.emitter.EmitObjectSet(buildDst, keyReg, valReg)
	}

	if buildDst != dst {
		sp.emitter.EmitMove(dst, buildDst)
	}
}

func (sp *ScopeProjection) RestoreFromArray(src vm.Operand) {
	idx := sp.registers.Allocate()

	for i, v := range sp.values {
		sp.emitter.EmitLoadConst(idx, sp.symbols.AddConstant(runtime.Int(i)))
		variable, _ := sp.symbols.DeclareLocal(v.Name, v.Type)
		sp.emitter.EmitABC(vm.OpLoadIndex, variable, src, idx)
	}
}

func (sp *ScopeProjection) RestoreFromObject(src vm.Operand) {
	key := sp.registers.Allocate()

	for _, v := range sp.values {
		sp.emitter.EmitLoadConst(key, sp.symbols.AddConstant(runtime.String(v.Name)))
		variable, _ := sp.symbols.DeclareLocal(v.Name, v.Type)
		sp.emitter.EmitABC(vm.OpLoadKey, variable, src, key)
	}
}

func (sp *ScopeProjection) usesRegister(reg vm.Operand) bool {
	for _, v := range sp.values {
		if v.Register == reg {
			return true
		}
	}

	return false
}
