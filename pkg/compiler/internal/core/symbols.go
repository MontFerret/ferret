package core

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

const (
	JumpPlaceholder      = -1
	UndefinedVariable    = -1
	IgnorePseudoVariable = "_"
	PseudoVariable       = "CURRENT"
)

type SymbolKind int

const (
	SymbolVar SymbolKind = iota
	SymbolConst
	SymbolParam
	SymbolGlobal
)

type ValueType int

const (
	TypeUnknown ValueType = iota
	TypeInt
	TypeFloat
	TypeString
	TypeBool
	TypeList
	TypeMap
	TypeAny
)

type Variable struct {
	Name     string
	Kind     SymbolKind
	Register vm.Operand
	Depth    int
	Type     ValueType
}

type SymbolTable struct {
	registers *RegisterAllocator
	constants *ConstantPool

	params  map[string]string
	globals map[string]vm.Operand
	locals  []*Variable

	scope int
}

func NewSymbolTable(registers *RegisterAllocator) *SymbolTable {
	return &SymbolTable{
		registers: registers,
		constants: NewConstantPool(),
		params:    make(map[string]string),
		globals:   make(map[string]vm.Operand),
		locals:    make([]*Variable, 0),
	}
}

func (st *SymbolTable) Scope() int {
	return st.scope
}

func (st *SymbolTable) EnterScope() {
	st.scope++
}

func (st *SymbolTable) ExitScope() {
	st.scope--

	for len(st.locals) > 0 && st.locals[len(st.locals)-1].Depth > st.scope {
		popped := st.locals[len(st.locals)-1]
		st.registers.Free(popped.Register)
		st.locals = st.locals[:len(st.locals)-1]
	}
}

func (st *SymbolTable) DeclareLocal(name string) vm.Operand {
	return st.DeclareLocalTyped(name, TypeUnknown)
}

func (st *SymbolTable) DeclareLocalTyped(name string, typ ValueType) vm.Operand {
	reg := st.registers.Allocate(Var)

	st.locals = append(st.locals, &Variable{
		Name:     name,
		Kind:     SymbolVar,
		Register: reg,
		Depth:    st.scope,
		Type:     typ,
	})

	return reg
}

func (st *SymbolTable) DeclareGlobal(name string) vm.Operand {
	if _, exists := st.globals[name]; exists {
		panic(runtime.Error(ErrVariableNotUnique, name))
	}

	op := st.constants.Add(runtime.NewString(name))
	st.globals[name] = op
	return op
}

func (st *SymbolTable) BindParam(name string) vm.Operand {
	st.params[name] = name
	return st.constants.Add(runtime.NewString(name))
}

func (st *SymbolTable) Constants() []runtime.Value {
	return st.constants.All()
}

func (st *SymbolTable) AddConstant(val runtime.Value) vm.Operand {
	return st.constants.Add(val)
}

func (st *SymbolTable) Constant(addr vm.Operand) runtime.Value {
	return st.constants.Get(addr)
}

func (st *SymbolTable) Resolve(name string) (vm.Operand, SymbolKind, bool) {
	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]
		if v.Name == name {
			return vm.NewRegisterOperand(int(v.Register)), v.Kind, true
		}
	}

	if reg, ok := st.globals[name]; ok {
		return reg, SymbolGlobal, true
	}

	return vm.NoopOperand, SymbolVar, false
}

func (st *SymbolTable) Lookup(name string) (*Variable, bool) {
	for i := len(st.locals) - 1; i >= 0; i-- {
		if st.locals[i].Name == name {
			return st.locals[i], true
		}
	}

	return nil, false
}

func (st *SymbolTable) Params() []string {
	out := make([]string, 0, len(st.params))
	for k := range st.params {
		out = append(out, k)
	}
	return out
}

func (st *SymbolTable) DebugView() []string {
	var out []string

	for _, v := range st.locals {
		out = append(out, fmt.Sprintf("[local] %s -> R%d (%v)", v.Name, v.Register, v.Type))
	}

	for k, r := range st.globals {
		out = append(out, fmt.Sprintf("[global] %s -> R%d", k, r))
	}

	for k, v := range st.params {
		out = append(out, fmt.Sprintf("[param] %s -> %s", k, v))
	}

	for _, c := range st.constants.All() {
		out = append(out, fmt.Sprintf("[constant] %s", c.String()))
	}

	return out
}
