package runtime

// Functions is a container for functions that organizes them by their argument count.
// It provides separate storage for functions with fixed argument counts (0-4) and
// functions with variable argument counts for optimal performance.
type functionRegistry struct {
	variadic  FunctionCollection[Function]  // Functions with variable number of arguments
	zero      FunctionCollection[Function0] // Functions with no arguments
	single    FunctionCollection[Function1] // Functions with a single argument
	double    FunctionCollection[Function2] // Functions with two arguments
	tripple   FunctionCollection[Function3] // Functions with three arguments
	quadruple FunctionCollection[Function4] // Functions with four arguments
}

// NewFunctions creates and returns a new empty Functions container.
func NewFunctions() Functions {
	return &functionRegistry{}
}

func NewFunctionsFromMap(funcs map[string]Function) Functions {
	return &functionRegistry{
		variadic:  NewFunctionCollectionFromMap(funcs),
		zero:      NewFunctionCollection[Function0](),
		single:    NewFunctionCollection[Function1](),
		double:    NewFunctionCollection[Function2](),
		tripple:   NewFunctionCollection[Function3](),
		quadruple: NewFunctionCollection[Function4](),
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

func (f *functionRegistry) F() FunctionCollection[Function] {
	if f.variadic == nil {
		f.variadic = NewFunctionCollection[Function]()
	}

	return f.variadic
}

func (f *functionRegistry) F0() FunctionCollection[Function0] {
	if f.zero == nil {
		f.zero = NewFunctionCollection[Function0]()
	}

	return f.zero
}

func (f *functionRegistry) F1() FunctionCollection[Function1] {
	if f.single == nil {
		f.single = NewFunctionCollection[Function1]()
	}

	return f.single
}

func (f *functionRegistry) F2() FunctionCollection[Function2] {
	if f.double == nil {
		f.double = NewFunctionCollection[Function2]()
	}

	return f.double
}

func (f *functionRegistry) F3() FunctionCollection[Function3] {
	if f.tripple == nil {
		f.tripple = NewFunctionCollection[Function3]()
	}

	return f.tripple
}

func (f *functionRegistry) F4() FunctionCollection[Function4] {
	if f.quadruple == nil {
		f.quadruple = NewFunctionCollection[Function4]()
	}

	return f.quadruple
}

func (f *functionRegistry) SetAll(otherFns Functions) Functions {
	if otherFns == nil {
		return f
	}

	// Copy functions from each collection
	f.F().SetAll(otherFns.F())
	f.F0().SetAll(otherFns.F0())
	f.F1().SetAll(otherFns.F1())
	f.F2().SetAll(otherFns.F2())
	f.F3().SetAll(otherFns.F3())
	f.F4().SetAll(otherFns.F4())

	return f
}

func (f *functionRegistry) Unset(name string) Functions {
	if f.F().Has(name) {
		f.variadic.Unset(name)
	} else if f.F0().Has(name) {
		f.zero.Unset(name)
	} else if f.F1().Has(name) {
		f.single.Unset(name)
	} else if f.F2().Has(name) {
		f.double.Unset(name)
	} else if f.F3().Has(name) {
		f.tripple.Unset(name)
	} else if f.F4().Has(name) {
		f.quadruple.Unset(name)
	}

	return f
}

func (f *functionRegistry) UnsetAll() Functions {
	if f.variadic != nil {
		f.variadic.UnsetAll()
	}

	if f.zero != nil {
		f.zero.UnsetAll()
	}

	if f.single != nil {
		f.single.UnsetAll()
	}

	if f.double != nil {
		f.double.UnsetAll()
	}

	if f.tripple != nil {
		f.tripple.UnsetAll()
	}

	if f.quadruple != nil {
		f.quadruple.UnsetAll()
	}

	return f
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
	names = append(names, f.variadic.Names()...)
	names = append(names, f.zero.Names()...)
	names = append(names, f.single.Names()...)
	names = append(names, f.double.Names()...)
	names = append(names, f.tripple.Names()...)
	names = append(names, f.quadruple.Names()...)

	return names
}
