package runtime

// Functions is a container for functions that organizes them by their argument count.
// It provides separate storage for functions with fixed argument counts (0-4) and
// functions with variable argument counts for optimal performance.
type Functions struct {
	av    FunctionCollection[Function]
	a0    FunctionCollection[Function0]
	a1    FunctionCollection[Function1]
	a2    FunctionCollection[Function2]
	a3    FunctionCollection[Function3]
	a4    FunctionCollection[Function4]
	names []string
	hash  uint64
	size  int
}

var (
	// Shared empty collections keep sparse registries read-only on lookup paths.
	emptyFunctionsVar = NewFunctionCollection[Function]()
	emptyFunctions0   = NewFunctionCollection[Function0]()
	emptyFunctions1   = NewFunctionCollection[Function1]()
	emptyFunctions2   = NewFunctionCollection[Function2]()
	emptyFunctions3   = NewFunctionCollection[Function3]()
	emptyFunctions4   = NewFunctionCollection[Function4]()
)

// NewFunctions creates and returns a new empty Functions container.
func NewFunctions() *Functions {
	return &Functions{}
}

// NewFunctionsFrom merges every arity-specific definition from funcs.
// It returns an error only when inputs define the same qualified name and arity.
func NewFunctionsFrom(funcs ...*Functions) (*Functions, error) {
	return NewFunctionsBuilderFrom(funcs...).Build()
}

// NewFunctionsFromMap creates a registry of variadic definitions.
func NewFunctionsFromMap(funcs map[string]Function) (*Functions, error) {
	builder := newRootFunctionsBuilder()

	for name, fn := range funcs {
		builder.Var().Add(name, fn)
	}

	return builder.Build()
}

// Hash returns a deterministic hash of every registered qualified name and arity.
func (f *Functions) Hash() uint64 {
	return f.hash
}

func (f *Functions) Has(name string) bool {
	if f.av != nil && f.av.Has(name) {
		return true
	}

	if f.a0 != nil && f.a0.Has(name) {
		return true
	}

	if f.a1 != nil && f.a1.Has(name) {
		return true
	}

	if f.a2 != nil && f.a2.Has(name) {
		return true
	}

	if f.a3 != nil && f.a3.Has(name) {
		return true
	}

	if f.a4 != nil && f.a4.Has(name) {
		return true
	}

	return false
}

// Size returns the number of registered definitions, including overloads.
func (f *Functions) Size() int {
	return f.size
}

// List returns a sorted defensive copy of the unique canonical function names.
func (f *Functions) List() []string {
	if len(f.names) == 0 {
		return []string{}
	}

	names := make([]string, len(f.names))
	copy(names, f.names)

	return names
}

func (f *Functions) A0() FunctionCollection[Function0] {
	if f.a0 == nil {
		return emptyFunctions0
	}

	return f.a0
}

func (f *Functions) A1() FunctionCollection[Function1] {
	if f.a1 == nil {
		return emptyFunctions1
	}

	return f.a1
}

func (f *Functions) A2() FunctionCollection[Function2] {
	if f.a2 == nil {
		return emptyFunctions2
	}

	return f.a2
}

func (f *Functions) A3() FunctionCollection[Function3] {
	if f.a3 == nil {
		return emptyFunctions3
	}

	return f.a3
}

func (f *Functions) A4() FunctionCollection[Function4] {
	if f.a4 == nil {
		return emptyFunctions4
	}

	return f.a4
}

func (f *Functions) Var() FunctionCollection[Function] {
	if f.av == nil {
		return emptyFunctionsVar
	}

	return f.av
}
