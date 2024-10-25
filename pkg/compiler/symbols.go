package compiler

import (
	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/runtime/core"
	"github.com/MontFerret/ferret/pkg/runtime/values"
)

type SymbolTable struct {
	constants      []core.Value
	constantsIndex map[uint64]int
	globals        map[string]int
	scope          int
	locals         []Variable
	registers      *RegisterAllocator
}

func NewSymbolTable(registers *RegisterAllocator) *SymbolTable {
	return &SymbolTable{
		constants:      make([]core.Value, 0),
		constantsIndex: make(map[uint64]int),
		globals:        make(map[string]int),
		locals:         make([]Variable, 0),
		registers:      registers,
	}
}

func (st *SymbolTable) Scope() int {
	return st.scope
}

func (st *SymbolTable) EnterScope() {
	st.scope++
}

// AddConstant adds a constant to the constants pool and returns its index.
// If the constant is a scalar, it will be deduplicated.
// If the constant is not a scalar, it will be added to the pool without deduplication.
func (st *SymbolTable) AddConstant(constant core.Value) runtime.Operand {
	var hash uint64
	isNone := constant == values.None

	if values.IsScalar(constant) {
		hash = constant.Hash()
	}

	if hash > 0 || isNone {
		if p, ok := st.constantsIndex[hash]; ok {
			return runtime.NewConstantOperand(p)
		}
	}

	st.constants = append(st.constants, constant)
	p := len(st.constants) - 1

	if hash > 0 || isNone {
		st.constantsIndex[hash] = p
	}

	// We flip the sign to indicate that this is a constant index, not a register.
	return runtime.NewConstantOperand(p)
}

func (st *SymbolTable) DefineVariable(name string) runtime.Operand {
	var index int

	if st.scope == 0 {
		// Check for duplicate global variable names.
		_, ok := st.globals[name]

		if ok {
			panic(core.Error(ErrVariableNotUnique, name))
		}

		index = len(st.globals)
		// Define global variable.
		st.globals[name] = index

		// Return a constant operand to indicate that this is a global variable and use its index.
		return st.AddConstant(values.NewString(name))
	}

	register := st.registers.AllocateLocalVar(name)

	st.locals = append(st.locals, Variable{
		Name:     name,
		Depth:    st.scope,
		Register: register,
	})

	return runtime.NewRegisterOperand(index)
}

func (st *SymbolTable) LookupVariable(name string) runtime.Operand {
	for i := len(st.locals) - 1; i >= 0; i-- {
		variable := st.locals[i]
		if variable.Name == name {
			return runtime.NewRegisterOperand(int(variable.Register))
		}
	}

	index, ok := st.globals[name]

	if !ok {
		panic(core.Error(ErrVariableNotFound, name))
	}

	return runtime.NewConstantOperand(index)
}

func (st *SymbolTable) ExitScope() {
	st.scope--

	// Pop all local variables from the stack within the closed scope.
	for len(st.locals) > 0 && st.locals[len(st.locals)-1].Depth > st.scope {
		popped := st.locals[len(st.locals)-1:]

		// Free the register.
		for _, v := range popped {
			st.registers.Free(v.Register)
		}

		st.locals = st.locals[:len(st.locals)-1]
	}
}
