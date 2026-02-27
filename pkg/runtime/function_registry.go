package runtime

import (
	"hash/fnv"
	"sort"
)

// Functions is a container for functions that organizes them by their argument count.
// It provides separate storage for functions with fixed argument counts (0-4) and
// functions with variable argument counts for optimal performance.
type Functions struct {
	hash  uint64
	size  int
	names []string
	av    FunctionCollection[Function]  // Functions with variable number of arguments
	a0    FunctionCollection[Function0] // Functions with no arguments
	a1    FunctionCollection[Function1] // Functions with a single argument
	a2    FunctionCollection[Function2] // Functions with two arguments
	a3    FunctionCollection[Function3] // Functions with three arguments
	a4    FunctionCollection[Function4] // Functions with four arguments
}

// NewFunctions creates and returns a new empty Functions container.
func NewFunctions() *Functions {
	return &Functions{}
}

func NewFunctionsFrom(funcs ...*Functions) (*Functions, error) {
	return NewFunctionsBuilderFrom(funcs...).Build()
}

func NewFunctionsFromMap(funcs map[string]Function) (*Functions, error) {
	builder := newRootFunctionsBuilder()

	for name, fn := range funcs {
		builder.Var().Add(name, fn)
	}

	return builder.Build()
}

func functionsHash(f *Functions) uint64 {
	if f == nil {
		return 0
	}

	names := f.List()
	sort.Strings(names)

	hasher := fnv.New64a()

	for _, name := range names {
		_, _ = hasher.Write([]byte(name))
	}

	return hasher.Sum64()
}

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

func (f *Functions) Size() int {
	return f.size
}

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
		f.a0 = NewFunctionCollection[Function0]()
	}

	return f.a0
}

func (f *Functions) A1() FunctionCollection[Function1] {
	if f.a1 == nil {
		f.a1 = NewFunctionCollection[Function1]()
	}

	return f.a1
}

func (f *Functions) A2() FunctionCollection[Function2] {
	if f.a2 == nil {
		f.a2 = NewFunctionCollection[Function2]()
	}

	return f.a2
}

func (f *Functions) A3() FunctionCollection[Function3] {
	if f.a3 == nil {
		f.a3 = NewFunctionCollection[Function3]()
	}

	return f.a3
}

func (f *Functions) A4() FunctionCollection[Function4] {
	if f.a4 == nil {
		f.a4 = NewFunctionCollection[Function4]()
	}

	return f.a4
}

func (f *Functions) Var() FunctionCollection[Function] {
	if f.av == nil {
		f.av = NewFunctionCollection[Function]()
	}

	return f.av
}
