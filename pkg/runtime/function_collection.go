package runtime

import "strings"

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
		ForEach(fn func(T, string) error) error
	}

	defaultFunctionCollection[T FunctionConstraint] struct {
		values map[string]T
	}
)

// NewFunctionCollection creates a new function collection of the specified type
func NewFunctionCollection[T FunctionConstraint]() FunctionCollection[T] {
	return &defaultFunctionCollection[T]{
		values: make(map[string]T),
	}
}

// NewFunctionCollectionFromMap creates a new function collection from an existing map
// It makes a copy of the provided map to ensure that the original map remains unmodified
func NewFunctionCollectionFromMap[T FunctionConstraint](values map[string]T) FunctionCollection[T] {
	fc := &defaultFunctionCollection[T]{
		values: make(map[string]T, len(values)),
	}

	for name, fn := range values {
		fc.values[name] = fn
	}

	return fc
}

func (f *defaultFunctionCollection[T]) Has(name string) bool {
	_, exists := f.values[name]

	return exists

}

func (f *defaultFunctionCollection[T]) Get(name string) (T, bool) {
	fn, exists := f.values[strings.ToUpper(name)]

	return fn, exists
}

func (f *defaultFunctionCollection[T]) GetAll() map[string]T {
	// Return a copy to prevent external modification
	result := make(map[string]T, len(f.values))

	for name, fn := range f.values {
		result[name] = fn
	}

	return result

}

func (f *defaultFunctionCollection[T]) Names() []string {
	names := make([]string, 0, len(f.values))

	for name := range f.values {
		names = append(names, name)
	}

	return names
}

func (f *defaultFunctionCollection[T]) Size() int {
	return len(f.values)
}

func (f *defaultFunctionCollection[T]) ForEach(fn func(T, string) error) error {
	for name, value := range f.values {
		if err := fn(value, name); err != nil {
			return err
		}
	}

	return nil
}
