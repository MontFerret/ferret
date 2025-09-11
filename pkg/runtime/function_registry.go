package runtime

type (
	// Functions is a container for functions that organizes them by their argument count.
	// It provides separate storage for functions with fixed argument counts (0-4) and
	// functions with variable argument counts for optimal performance.
	Functions interface {
		Hash() uint64
		Has(name string) bool
		F0() FunctionCollection[Function0]
		F1() FunctionCollection[Function1]
		F2() FunctionCollection[Function2]
		F3() FunctionCollection[Function3]
		F4() FunctionCollection[Function4]
		FV() FunctionCollection[Function]
		Names() []string
		Size() int
	}

	functionRegistry struct {
		hash uint64
		fv   FunctionCollection[Function]  // Functions with variable number of arguments
		f0   FunctionCollection[Function0] // Functions with no arguments
		f1   FunctionCollection[Function1] // Functions with a single argument
		f2   FunctionCollection[Function2] // Functions with two arguments
		f3   FunctionCollection[Function3] // Functions with three arguments
		f4   FunctionCollection[Function4] // Functions with four arguments
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
	f := &functionRegistry{
		fv: NewFunctionCollectionFromMap(funcs),
	}

	f.hash = functionsHash(f)

	return f
}

func (f *functionRegistry) Hash() uint64 {
	return f.hash
}

func (f *functionRegistry) Has(name string) bool {
	if f.fv != nil && f.fv.Has(name) {
		return true
	}

	if f.f0 != nil && f.f0.Has(name) {
		return true
	}

	if f.f1 != nil && f.f1.Has(name) {
		return true
	}

	if f.f2 != nil && f.f2.Has(name) {
		return true
	}

	if f.f3 != nil && f.f3.Has(name) {
		return true
	}

	if f.f4 != nil && f.f4.Has(name) {
		return true
	}

	return false
}

func (f *functionRegistry) Size() int {
	var size int

	if f.f0 != nil {
		size += f.f0.Size()
	}

	if f.f1 != nil {
		size += f.f1.Size()
	}

	if f.f2 != nil {
		size += f.f2.Size()
	}

	if f.f3 != nil {
		size += f.f3.Size()
	}

	if f.f4 != nil {
		size += f.f4.Size()
	}

	if f.fv != nil {
		size += f.fv.Size()
	}

	return size
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

func (f *functionRegistry) FV() FunctionCollection[Function] {
	if f.fv == nil {
		f.fv = NewFunctionCollection[Function]()
	}

	return f.fv
}

func (f *functionRegistry) Names() []string {
	// Pre-calculate capacity to avoid reallocations
	names := make([]string, 0, f.Size())

	if f.f0 != nil {
		names = append(names, f.f0.Names()...)
	}

	if f.f1 != nil {
		names = append(names, f.f1.Names()...)
	}

	if f.f2 != nil {
		names = append(names, f.f2.Names()...)
	}

	if f.f3 != nil {
		names = append(names, f.f3.Names()...)
	}

	if f.f4 != nil {
		names = append(names, f.f4.Names()...)
	}

	if f.fv != nil {
		names = append(names, f.fv.Names()...)
	}

	return names
}
