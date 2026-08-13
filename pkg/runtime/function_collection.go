package runtime

type (
	// FunctionConstraint is a type constraint that includes all function types
	FunctionConstraint interface {
		Function | Function0 | Function1 | Function2 | Function3 | Function4
	}

	// FunctionCollection is an immutable collection of functions of a specific type (e.g., Function, Function0, etc.)
	FunctionCollection[T FunctionConstraint] interface {
		Has(name string) bool
		Get(name string) (T, bool)
		GetAll() map[string]T
		Names() []string
		Size() int
		ForEach(fn func(T, string) bool)
	}

	defaultFunctionCollection[T FunctionConstraint] struct {
		values     map[string]registeredFunction[T]
		normalized bool
	}
)

// NewFunctionCollection creates a new function collection of the specified type
func NewFunctionCollection[T FunctionConstraint]() FunctionCollection[T] {
	return NewFunctionCollectionFromMap[T](nil)
}

// NewFunctionCollectionFromMap creates a new function collection from an existing map
// It makes a copy of the provided map to ensure that the original map remains unmodified
func NewFunctionCollectionFromMap[T FunctionConstraint](values map[string]T) FunctionCollection[T] {
	fc := &defaultFunctionCollection[T]{
		values: make(map[string]registeredFunction[T], len(values)),
	}

	for name, fn := range values {
		fc.values[name] = registeredFunction[T]{name: name, function: fn}
	}

	return fc
}

func newNormalizedFunctionCollectionFromMap[T FunctionConstraint](values map[string]registeredFunction[T]) FunctionCollection[T] {
	fc := &defaultFunctionCollection[T]{
		values:     make(map[string]registeredFunction[T], len(values)),
		normalized: true,
	}

	for key, entry := range values {
		fc.values[key] = entry
	}

	return fc
}

func (f *defaultFunctionCollection[T]) Has(name string) bool {
	if f.normalized {
		name = NormalizeRegisteredName(name)
	}

	_, exists := f.values[name]

	return exists
}

func (f *defaultFunctionCollection[T]) Get(name string) (T, bool) {
	if f.normalized {
		name = NormalizeRegisteredName(name)
	}

	entry, exists := f.values[name]

	return entry.function, exists
}

func (f *defaultFunctionCollection[T]) GetAll() map[string]T {
	// Return a copy to prevent external modification
	result := make(map[string]T, len(f.values))

	for _, entry := range f.values {
		result[entry.name] = entry.function
	}

	return result
}

func (f *defaultFunctionCollection[T]) Names() []string {
	names := make([]string, 0, len(f.values))

	for _, entry := range f.values {
		names = append(names, entry.name)
	}

	return names
}

func (f *defaultFunctionCollection[T]) Size() int {
	return len(f.values)
}

func (f *defaultFunctionCollection[T]) ForEach(fn func(T, string) bool) {
	for _, entry := range f.values {
		if !fn(entry.function, entry.name) {
			break
		}
	}
}
