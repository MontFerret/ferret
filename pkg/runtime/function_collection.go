package runtime

type functionCollection[T FunctionConstraint] struct {
	values map[string]T
}

// NewFunctionCollection creates a new function collection of the specified type
func NewFunctionCollection[T FunctionConstraint]() FunctionCollection[T] {
	return &functionCollection[T]{
		values: make(map[string]T),
	}
}

// NewFunctionCollectionFromMap creates a new function collection from an existing map
func NewFunctionCollectionFromMap[T FunctionConstraint](values map[string]T) FunctionCollection[T] {
	fc := &functionCollection[T]{
		values: make(map[string]T, len(values)),
	}

	for name, fn := range values {
		fc.values[name] = fn
	}

	return fc
}

func (f *functionCollection[T]) Has(name string) bool {
	_, exists := f.values[name]

	return exists

}

func (f *functionCollection[T]) Set(name string, fn T) FunctionCollection[T] {
	f.values[name] = fn

	return f

}

func (f *functionCollection[T]) SetAll(otherFns FunctionCollection[T]) FunctionCollection[T] {
	if otherFns == nil {
		return f
	}

	for name, fn := range otherFns.GetAll() {
		f.values[name] = fn
	}

	return f
}

func (f *functionCollection[T]) Get(name string) (T, bool) {
	fn, exists := f.values[name]

	return fn, exists
}

func (f *functionCollection[T]) GetAll() map[string]T {
	// Return a copy to prevent external modification
	result := make(map[string]T, len(f.values))

	for name, fn := range f.values {
		result[name] = fn
	}

	return result

}

func (f *functionCollection[T]) Unset(name string) FunctionCollection[T] {
	delete(f.values, name)

	return f

}

func (f *functionCollection[T]) UnsetAll() FunctionCollection[T] {
	f.values = make(map[string]T)

	return f

}

func (f *functionCollection[T]) Names() []string {
	names := make([]string, 0, len(f.values))

	for name := range f.values {
		names = append(names, name)
	}

	return names
}

func (f *functionCollection[T]) Size() int {
	return len(f.values)
}
