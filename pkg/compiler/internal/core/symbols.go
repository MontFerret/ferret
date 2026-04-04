package core

import (
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const (
	jumpPlaceholder      = -1
	UndefinedVariable    = -1
	IgnorePseudoVariable = "_"
	PseudoVariable       = "__CURRENT__"
)

type (
	SymbolKind int

	ValueType int

	BindingStorage int

	BindingOptions struct {
		Mutable bool
		Storage BindingStorage
	}
)

const (
	SymbolConst SymbolKind = iota
	SymbolGlobal
	SymbolLocal
	SymbolParam
)

const (
	BindingStorageValue BindingStorage = iota
	BindingStorageCell
)

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

// IsScalar reports whether the type is a known scalar.
func (t ValueType) IsScalar() bool {
	switch t {
	case TypeNone, TypeInt, TypeFloat, TypeString, TypeBool:
		return true
	default:
		return false
	}
}

// IsUntracked reports whether the type is known to never carry direct resource
// / closer ownership at runtime.
func (t ValueType) IsUntracked() bool {
	switch t {
	case TypeNone, TypeInt, TypeFloat, TypeString, TypeBool, TypeArray, TypeObject:
		return true
	default:
		return false
	}
}

type (
	Variable struct {
		Name     string
		Kind     SymbolKind
		Register bytecode.Operand
		Depth    int
		Type     ValueType
		Mutable  bool
		Storage  BindingStorage
	}

	SymbolTable struct {
		registers *RegisterAllocator
		constants *ConstantPool

		globals map[string]*Variable
		locals  []*Variable

		scope int
	}
)

func NewSymbolTable(registers *RegisterAllocator, constants *ConstantPool) *SymbolTable {
	if constants == nil {
		constants = NewConstantPool()
	}

	return &SymbolTable{
		registers: registers,
		constants: constants,
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

func (st *SymbolTable) ProjectionVariables() []Variable {
	vars := make([]Variable, 0)
	seen := make(map[string]struct{})

	add := func(v *Variable) {
		if v == nil || v.Name == IgnorePseudoVariable || v.Name == PseudoVariable {
			return
		}

		if _, ok := seen[v.Name]; ok {
			return
		}

		seen[v.Name] = struct{}{}
		vars = append(vars, *v)
	}

	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]
		if v.Depth == st.scope {
			add(v)
		}
	}

	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]
		if v.Depth >= st.scope || !v.Mutable {
			continue
		}

		add(v)
	}

	if len(st.globals) == 0 {
		return vars
	}

	names := make([]string, 0, len(st.globals))
	for name, v := range st.globals {
		if v == nil || !v.Mutable {
			continue
		}

		if _, ok := seen[name]; ok {
			continue
		}

		names = append(names, name)
	}

	sort.Strings(names)

	for _, name := range names {
		add(st.globals[name])
	}

	return vars
}

func (st *SymbolTable) DeclareLocal(name string, typ ValueType) (bytecode.Operand, bool) {
	return st.DeclareLocalWithOptions(name, typ, BindingOptions{})
}

func (st *SymbolTable) DeclareLocalWithOptions(name string, typ ValueType, opts BindingOptions) (bytecode.Operand, bool) {
	reg := st.registers.Allocate()

	if ok := st.AssignLocalWithOptions(name, typ, reg, opts); !ok {
		return bytecode.NoopOperand, false
	}

	return reg, true
}

func (st *SymbolTable) AssignLocal(name string, typ ValueType, op bytecode.Operand) bool {
	return st.AssignLocalWithOptions(name, typ, op, BindingOptions{})
}

func (st *SymbolTable) AssignLocalWithOptions(name string, typ ValueType, op bytecode.Operand, opts BindingOptions) bool {
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
		Mutable:  opts.Mutable,
		Storage:  opts.Storage,
	})

	return true
}

func (st *SymbolTable) DeclareGlobal(name string, typ ValueType) (bytecode.Operand, bool) {
	return st.DeclareGlobalWithOptions(name, typ, BindingOptions{})
}

func (st *SymbolTable) DeclareGlobalWithOptions(name string, typ ValueType, opts BindingOptions) (bytecode.Operand, bool) {
	op := st.registers.Allocate()

	if ok := st.AssignGlobalWithOptions(name, typ, op, opts); !ok {
		return bytecode.NoopOperand, false
	}

	return op, true
}

func (st *SymbolTable) AssignGlobal(name string, typ ValueType, op bytecode.Operand) bool {
	return st.AssignGlobalWithOptions(name, typ, op, BindingOptions{})
}

func (st *SymbolTable) AssignGlobalWithOptions(name string, typ ValueType, op bytecode.Operand, opts BindingOptions) bool {
	if _, exists := st.globals[name]; exists {
		return false
	}

	st.globals[name] = &Variable{
		Name:     name,
		Kind:     SymbolGlobal,
		Register: op,
		Type:     typ,
		Mutable:  opts.Mutable,
		Storage:  opts.Storage,
	}

	return true
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
	if binding, ok := st.ResolveBinding(name); ok {
		return bytecode.NewRegister(int(binding.Register)), binding.Kind, true
	}

	return bytecode.NoopOperand, SymbolLocal, false
}

func (st *SymbolTable) ResolveBinding(name string) (*Variable, bool) {
	for i := len(st.locals) - 1; i >= 0; i-- {
		v := st.locals[i]
		if v.Name == name {
			return v, true
		}
	}

	if binding, ok := st.globals[name]; ok {
		return binding, true
	}

	return nil, false
}

func (st *SymbolTable) Lookup(name string) (*Variable, bool) {
	return st.ResolveBinding(name)
}
