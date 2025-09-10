package runtime

type (
	// Functions is a container for functions that organizes them by their argument count.
	// It provides separate storage for functions with fixed argument counts (0-4) and
	// functions with variable argument counts for optimal performance.
	Functions interface {
		Has(name string) bool
		F() FunctionCollection[Function]
		F0() FunctionCollection[Function0]
		F1() FunctionCollection[Function1]
		F2() FunctionCollection[Function2]
		F3() FunctionCollection[Function3]
		F4() FunctionCollection[Function4]
		Names() []string
		Size() int
	}

	functionRegistry struct {
		f  FunctionCollection[Function]  // Functions with variable number of arguments
		f0 FunctionCollection[Function0] // Functions with no arguments
		f1 FunctionCollection[Function1] // Functions with a single argument
		f2 FunctionCollection[Function2] // Functions with two arguments
		f3 FunctionCollection[Function3] // Functions with three arguments
		f4 FunctionCollection[Function4] // Functions with four arguments
	}
)

// NewFunctions creates and returns a new empty Functions container.
func NewFunctions() Functions {
	return &functionRegistry{}
}

func NewFunctionsFrom(funcs ...Functions) Functions {
	builder := NewFunctionsBuilder()

	for _, fn := range funcs {
		builder.SetFrom(fn)
	}

	return builder.Build()
}

func NewFunctionsFromMap(funcs map[string]Function) Functions {
	return &functionRegistry{
		f: NewFunctionCollectionFromMap(funcs),
	}
}

func (f *functionRegistry) Has(name string) bool {
	return f.F().Has(name) ||
		f.F0().Has(name) ||
		f.F1().Has(name) ||
		f.F2().Has(name) ||
		f.F3().Has(name) ||
		f.F4().Has(name)
}

func (f *functionRegistry) Size() int {
	return f.F().Size() +
		f.F0().Size() +
		f.F1().Size() +
		f.F2().Size() +
		f.F3().Size() +
		f.F4().Size()
}

func (f *functionRegistry) F() FunctionCollection[Function] {
	if f.f == nil {
		f.f = NewFunctionCollection[Function]()
	}

	return f.f
}

func (f *functionRegistry) F0() FunctionCollection[Function0] {
	if f.f0 == nil {
		f.f0 = NewFunctionCollection[Function0]()
	}

	return f.f0
}

func (f *functionRegistry) F1() FunctionCollection[Function1] {
	if f.f1 == nil {
		f.f1 = NewFunctionCollection[Function1]()
	}

	return f.f1
}

func (f *functionRegistry) F2() FunctionCollection[Function2] {
	if f.f2 == nil {
		f.f2 = NewFunctionCollection[Function2]()
	}

	return f.f2
}

func (f *functionRegistry) F3() FunctionCollection[Function3] {
	if f.f3 == nil {
		f.f3 = NewFunctionCollection[Function3]()
	}

	return f.f3
}

func (f *functionRegistry) F4() FunctionCollection[Function4] {
	if f.f4 == nil {
		f.f4 = NewFunctionCollection[Function4]()
	}

	return f.f4
}

func (f *functionRegistry) Names() []string {
	// Pre-calculate capacity to avoid reallocations
	capacity := f.F().Size() +
		f.F0().Size() +
		f.F1().Size() +
		f.F2().Size() +
		f.F3().Size() +
		f.F4().Size()

	names := make([]string, 0, capacity)

	// Collect names from all function collections
	names = append(names, f.f.Names()...)
	names = append(names, f.f0.Names()...)
	names = append(names, f.f1.Names()...)
	names = append(names, f.f2.Names()...)
	names = append(names, f.f3.Names()...)
	names = append(names, f.f4.Names()...)

	return names
}
