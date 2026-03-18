package core

import (
	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	jumpPlaceholder      = -1
	UndefinedVariable    = -1
	IgnorePseudoVariable = "_"
	PseudoVariable       = "__CURRENT__"
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
	TypeNone
	TypeInt
	TypeFloat
	TypeString
	TypeBool
	TypeArray
	TypeObject
	TypeList
	TypeMap
	TypeAny
)

// IsScalar reports whether the type is a known scalar that can never
// carry resource / closer semantics. Mirrors runtime mem.CanTrackValue.
func (t ValueType) IsScalar() bool {
	switch t {
	case TypeNone, TypeInt, TypeFloat, TypeString, TypeBool:
		return true
	default:
		return false
	}
}

type Variable struct {
	Name     string
	Kind     SymbolKind
	Register bytecode.Operand
	Depth    int
	Type     ValueType
}

type SymbolTable struct {
	registers *RegisterAllocator
	constants *ConstantPool

	params    map[string]int
	paramList []string
	functions map[string]int
	globals   map[string]bytecode.Operand
	locals    []*Variable

	scope int
}

func NewSymbolTable(registers *RegisterAllocator, constants *ConstantPool) *SymbolTable {
	if constants == nil {
		constants = NewConstantPool()
	}

	return &SymbolTable{
		registers: registers,
		constants: constants,
		params:    make(map[string]int),
		paramList: make([]string, 0),
		globals:   make(map[string]bytecode.Operand),
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

func (st *SymbolTable) DeclareLocal(name string, typ ValueType) (bytecode.Operand, bool) {
	reg := st.registers.Allocate()

	if ok := st.AssignLocal(name, typ, reg); !ok {
		return bytecode.NoopOperand, false
	}

	return reg, true
}

func (st *SymbolTable) AssignLocal(name string, typ ValueType, op bytecode.Operand) bool {
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

func (st *SymbolTable) DeclareGlobal(name string, typ ValueType) (bytecode.Operand, bool) {
	op := st.registers.Allocate()

	if ok := st.AssignGlobal(name, typ, op); !ok {
		return bytecode.NoopOperand, false
	}

	return op, true
}

func (st *SymbolTable) AssignGlobal(name string, typ ValueType, op bytecode.Operand) bool {
	if _, exists := st.globals[name]; exists {
		return false
	}

	st.globals[name] = op

	return true
}

func (st *SymbolTable) BindParam(name string) bytecode.Operand {
	if slot, exists := st.params[name]; exists {
		return bytecode.Operand(slot + 1)
	}

	slot := len(st.paramList)
	st.params[name] = slot
	st.paramList = append(st.paramList, name)

	return bytecode.Operand(slot + 1)
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

func (st *SymbolTable) AddConstant(val runtime.Value) bytecode.Operand {
	return st.constants.Add(val)
}

func (st *SymbolTable) Constant(addr bytecode.Operand) runtime.Value {
	return st.constants.Get(addr)
}

func (st *SymbolTable) Resolve(name string) (bytecode.Operand, SymbolKind, bool) {
	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]
		if v.Name == name {
			return bytecode.NewRegister(int(v.Register)), v.Kind, true
		}
	}

	if reg, ok := st.globals[name]; ok {
		return reg, SymbolGlobal, true
	}

	return bytecode.NoopOperand, SymbolLocal, false
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
	out := make([]string, len(st.paramList))
	copy(out, st.paramList)

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
