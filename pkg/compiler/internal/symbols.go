package internal

import (
	"strconv"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

type (
	Variable struct {
		Name     string
		Register vm.Operand
		Depth    int
	}

	SymbolTable struct {
		registers      *RegisterAllocator
		constants      []runtime.Value
		constantsIndex map[uint64]int
		params         map[string]string
		globals        map[string]vm.Operand
		locals         []*Variable
		scope          int
	}
)

func NewSymbolTable(registers *RegisterAllocator) *SymbolTable {
	return &SymbolTable{
		registers:      registers,
		constants:      make([]runtime.Value, 0),
		constantsIndex: make(map[uint64]int),
		params:         make(map[string]string),
		globals:        make(map[string]vm.Operand),
		locals:         make([]*Variable, 0),
	}
}

func (st *SymbolTable) Params() []string {
	params := make([]string, 0, len(st.params))

	for _, name := range st.params {
		params = append(params, name)
	}

	return params
}

func (st *SymbolTable) Constants() []runtime.Value {
	return st.constants
}

func (st *SymbolTable) Scope() int {
	return st.scope
}

func (st *SymbolTable) EnterScope() {
	st.scope++
}

func (st *SymbolTable) AddParam(name string) vm.Operand {
	st.params[name] = name

	return st.AddConstant(runtime.NewString(name))
}

// AddConstant adds a constant to the constants pool and returns its index.
// If the constant is a scalar, it will be deduplicated.
// If the constant is not a scalar, it will be added to the pool without deduplication.
func (st *SymbolTable) AddConstant(constant runtime.Value) vm.Operand {
	var hash uint64
	isNone := constant == runtime.None

	if runtime.IsScalar(constant) {
		hash = constant.Hash()
	}

	if hash > 0 || isNone {
		if p, ok := st.constantsIndex[hash]; ok {
			return vm.NewConstantOperand(p)
		}
	}

	st.constants = append(st.constants, constant)
	p := len(st.constants) - 1

	if hash > 0 || isNone {
		st.constantsIndex[hash] = p
	}

	// We flip the sign to indicate that this is a constant index, not a register.
	return vm.NewConstantOperand(p)
}

// Constant returns a constant by its index.
func (st *SymbolTable) Constant(addr vm.Operand) runtime.Value {
	if !addr.IsConstant() {
		panic(runtime.Error(ErrInvalidOperandType, strconv.Itoa(int(addr))))
	}

	index := addr.Constant()

	if index < 0 || index >= len(st.constants) {
		panic(runtime.Error(ErrConstantNotFound, strconv.Itoa(index)))
	}

	return st.constants[index]
}

func (st *SymbolTable) DefineVariable(name string) vm.Operand {
	if st.scope == 0 {
		// Check for duplicate global variable names.
		_, ok := st.globals[name]

		if ok {
			panic(runtime.Error(ErrVariableNotUnique, name))
		}

		op := st.AddConstant(runtime.NewString(name))
		// Define global variable.
		st.globals[name] = op

		return op
	}

	register := st.registers.Allocate(Var)

	st.DefineScopedVariable(name, register)

	return register
}

func (st *SymbolTable) DefineScopedVariable(name string, register vm.Operand) {
	if st.scope == 0 {
		panic("cannot define scoped variable in global scope")
	}

	st.locals = append(st.locals, &Variable{
		Name:     name,
		Depth:    st.scope,
		Register: register,
	})
}

func (st *SymbolTable) Variable(name string) vm.Operand {
	for i := len(st.locals) - 1; i >= 0; i-- {
		variable := st.locals[i]
		if variable.Name == name {
			return vm.NewRegisterOperand(int(variable.Register))
		}
	}

	op, ok := st.globals[name]

	if !ok {
		panic(runtime.Error(ErrVariableNotFound, name))
	}

	return op
}

// GlobalVariable returns a global variable by its name.
func (st *SymbolTable) GlobalVariable(name string) (vm.Operand, bool) {
	op, ok := st.globals[name]

	return op, ok
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
