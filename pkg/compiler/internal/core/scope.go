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
	reg := sp.registers.Allocate()
	args := sp.registers.AllocateSequence(len(sp.values))

	for i, v := range sp.values {
		sp.emitter.EmitAB(vm.OpMove, args[i], v.Register)
	}

	sp.emitter.EmitArray(reg, args)
	sp.emitter.EmitMove(dst, reg)
}

func (sp *ScopeProjection) EmitAsObject(dst vm.Operand) {
	if len(sp.values) == 0 {
		sp.emitter.EmitA(vm.OpLoadObject, dst)
		return
	}

	pairs := sp.registers.AllocateSequence(len(sp.values) * 2)

	for i, v := range sp.values {
		// Key (field name)
		keyReg := pairs[i*2]
		sp.emitter.EmitLoadConst(keyReg, sp.symbols.AddConstant(runtime.String(v.Name)))

		// Value (actual variable value)
		valReg := pairs[i*2+1]
		sp.emitter.EmitAB(vm.OpMove, valReg, v.Register)
	}

	sp.emitter.EmitObject(dst, pairs)
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
