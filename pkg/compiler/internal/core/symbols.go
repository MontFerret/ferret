package core

import (
	"fmt"

	"github.com/MontFerret/ferret/pkg/runtime"
	"github.com/MontFerret/ferret/pkg/vm"
)

const (
	jumpPlaceholder      = -1
	UndefinedVariable    = -1
	IgnorePseudoVariable = "_"
	PseudoVariable       = "CURRENT"
)

type SymbolKind int

const (
	SymbolConst SymbolKind = iota
	SymbolGlobal
	SymbolLocal
	SymbolParam
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

	params    map[string]string
	functions map[string]int
	globals   map[string]*Variable
	locals    []*Variable

	scope int
}

func NewSymbolTable(registers *RegisterAllocator) *SymbolTable {
	return &SymbolTable{
		registers: registers,
		constants: NewConstantPool(),
		params:    make(map[string]string),
		globals:   make(map[string]*Variable),
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

func (st *SymbolTable) LocalVariables() []Variable {
	// collect all local variables in the current scope
	locals := make([]Variable, 0)

	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]

		if v.Depth == st.scope {
			locals = append(locals, *v)
		}
	}

	return locals
}

func (st *SymbolTable) DeclareLocal(name string, typ ValueType) (vm.Operand, bool) {
	reg := st.registers.Allocate(Var)

	if ok := st.AssignLocal(name, typ, reg); !ok {
		st.registers.Free(reg)
		return vm.NoopOperand, false
	}

	return reg, true
}

func (st *SymbolTable) AssignLocal(name string, typ ValueType, op vm.Operand) bool {
	if name != IgnorePseudoVariable && name != PseudoVariable {
		// Check if the variable already exists in the current scope
		for i := len(st.locals) - 1; i >= 0; i-- {
			v := st.locals[i]

			if v.Name == name && v.Depth == st.scope {
				return false
			}
		}
	}

	st.locals = append(st.locals, &Variable{
		Name:     name,
		Kind:     SymbolLocal,
		Register: op,
		Depth:    st.scope,
		Type:     typ,
	})

	return true
}

func (st *SymbolTable) DeclareGlobal(name string, typ ValueType) (vm.Operand, bool) {
	op := st.registers.Allocate(Var)

	if ok := st.AssignGlobal(name, typ, op); !ok {
		st.registers.Free(op)

		return vm.NoopOperand, false
	}

	return op, true
}

func (st *SymbolTable) AssignGlobal(name string, typ ValueType, op vm.Operand) bool {
	if _, exists := st.globals[name]; exists {
		return false
	}

	st.globals[name] = &Variable{
		Name:     name,
		Kind:     SymbolGlobal,
		Register: op,
		Depth:    0,
		Type:     typ,
	}

	return true
}

func (st *SymbolTable) BindParam(name string) vm.Operand {
	st.params[name] = name
	return st.constants.Add(runtime.NewString(name))
}

func (st *SymbolTable) BindFunction(name string, args int) {
	if st.functions == nil {
		st.functions = make(map[string]int)
	}

	if currArgs, exists := st.functions[name]; exists {
		// we need to ensure that the number of arguments is not greater than the current one
		// if it is not, we will not override the current one
		if currArgs > args {
			return
		}
	}

	st.functions[name] = args
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
			return vm.NewRegister(int(v.Register)), v.Kind, true
		}
	}

	if variable, ok := st.globals[name]; ok {
		return variable.Register, SymbolGlobal, true
	}

	return vm.NoopOperand, SymbolLocal, false
}

func (st *SymbolTable) Lookup(name string) (*Variable, bool) {
	for i := len(st.locals) - 1; i >= 0; i-- {
		if st.locals[i].Name == name {
			return st.locals[i], true
		}
	}

	if variable, ok := st.globals[name]; ok {
		return variable, true
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

func (st *SymbolTable) Functions() map[string]int {
	// Returns a copy of the functions map to avoid external modifications
	funcs := make(map[string]int, len(st.functions))

	for k, v := range st.functions {
		funcs[k] = v
	}

	return funcs
}

// RuntimeTypeToValueType converts a runtime type to compiler ValueType
func RuntimeTypeToValueType(runtimeType runtime.Type) ValueType {
	switch runtimeType {
	case runtime.TypeInt:
		return TypeInt
	case runtime.TypeFloat:
		return TypeFloat
	case runtime.TypeString:
		return TypeString
	case runtime.TypeBoolean:
		return TypeBool
	case runtime.TypeArray, runtime.TypeList:
		return TypeList
	case runtime.TypeObject, runtime.TypeMap:
		return TypeMap
	default:
		return TypeUnknown
	}
}

func (st *SymbolTable) DebugView() []string {
	var out []string

	for _, v := range st.locals {
		out = append(out, fmt.Sprintf("[local] %s -> R%d (%v)", v.Name, v.Register, v.Type))
	}

	for k, v := range st.globals {
		out = append(out, fmt.Sprintf("[global] %s -> R%d (%v)", k, v.Register, v.Type))
	}

	for k, v := range st.params {
		out = append(out, fmt.Sprintf("[param] %s -> %s", k, v))
	}

	for _, c := range st.constants.All() {
		out = append(out, fmt.Sprintf("[constant] %s", c.String()))
	}

	return out
}
