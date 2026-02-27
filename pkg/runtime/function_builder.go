package runtime

import (
	"fmt"
	"slices"
)

type (
	FnDef[T FunctionConstraint] interface {
		// Add adds a function to the builder.
		// If a function with the same name already exists, it will be ignored and an error will be recorded.
		Add(name string, fn T) FnDef[T]
		// Remove removes a function from the builder.
		// If a function with the same name does not exist, an error will be recorded.
		Remove(name string) FnDef[T]
		// Has checks if a function with the given name exists in the builder.
		Has(name string) bool
		// List retrieves the names of all functions currently registered in the builder.
		List() []string
		// ForEach iterates over all registered functions, calling the provided function with each function and its name.
		// Iteration stops if the provided function returns false.
		ForEach(fn func(fn T, name string) bool)
	}

	FunctionDefs interface {
		// Has checks if an entity with the given name exists in the collection and returns true if found, false otherwise.
		Has(name string) bool
		// A0 returns a function definition interface for managing functions with 0 arguments.
		A0() FnDef[Function0]
		// A1 returns a function definition interface for managing functions with 1 argument.
		A1() FnDef[Function1]
		// A2 returns a function definition interface for managing functions with 2 arguments.
		A2() FnDef[Function2]
		// A3 returns a function definition interface for managing functions with 3 arguments.
		A3() FnDef[Function3]
		// A4 returns a function definition interface for managing functions with 4 arguments.
		A4() FnDef[Function4]
		// Var returns a function definition interface for managing dynamic functions with variable arguments.
		Var() FnDef[Function]
		// From initializes the builder with functions from the given Functions container.
		From(other FunctionDefs) FunctionDefs
	}

	FunctionsBuilder struct {
		FunctionDefs

		namespace string
		av        *defaultFnDef[Function]
		a0        *defaultFnDef[Function0]
		a1        *defaultFnDef[Function1]
		a2        *defaultFnDef[Function2]
		a3        *defaultFnDef[Function3]
		a4        *defaultFnDef[Function4]
	}

	listable interface {
		List() []string
	}

	defaultFnDef[T FunctionConstraint] struct {
		namespace string
		errors    []error
		data      map[string]T
	}
)

func newFnDef[T FunctionConstraint](namespace string) *defaultFnDef[T] {
	return &defaultFnDef[T]{
		namespace: namespace,
		data:      make(map[string]T),
	}
}

func newFnDefFrom[T FunctionConstraint](namespace string, other *defaultFnDef[T]) *defaultFnDef[T] {
	return &defaultFnDef[T]{
		namespace: namespace,
		// We share the same map across all builders to ensure that changes in one builder are reflected in all builders that share the same namespace.
		data: other.data,
	}
}

func (fd *defaultFnDef[T]) Add(name string, fn T) FnDef[T] {
	fname := makeFunctionName(fd.namespace, name)

	if _, exists := fd.data[fname]; exists {
		if fd.errors == nil {
			fd.errors = make([]error, 0)
		}

		fd.errors = append(fd.errors, fmt.Errorf("function with name '%s' already exists in '%s' namespace", name, fd.namespace))

		return fd
	}

	fd.data[fname] = fn

	return fd
}

func (fd *defaultFnDef[T]) Remove(name string) FnDef[T] {
	fname := makeFunctionName(fd.namespace, name)

	if _, exists := fd.data[fname]; !exists {
		if fd.errors == nil {
			fd.errors = make([]error, 0)
		}

		fd.errors = append(fd.errors, fmt.Errorf("function with name '%s' does not exist in '%s' namespace", name, fd.namespace))

		return fd
	}

	delete(fd.data, fname)

	return fd
}

func (fd *defaultFnDef[T]) Has(name string) bool {
	fname := makeFunctionName(fd.namespace, name)
	_, exists := fd.data[fname]
	return exists
}

func (fd *defaultFnDef[T]) ForEach(fn func(fn T, name string) bool) {
	for name, fun := range fd.data {
		if !fn(fun, name) {
			break
		}
	}
}

func (fd *defaultFnDef[T]) List() []string {
	names := make([]string, 0, len(fd.data))

	for name := range fd.data {
		names = append(names, name)
	}

	return names
}

func NewFunctionsBuilder() *FunctionsBuilder {
	return newRootFunctionsBuilder()
}

func NewFunctionsBuilderFrom(funcs ...*Functions) *FunctionsBuilder {
	builder := newRootFunctionsBuilder()

	for _, f := range funcs {
		if f == nil {
			continue
		}

		f.A0().ForEach(func(fun Function0, name string) bool {
			builder.A0().Add(name, fun)

			return true
		})

		f.A1().ForEach(func(fun Function1, name string) bool {
			builder.A1().Add(name, fun)

			return true
		})

		f.A2().ForEach(func(fun Function2, name string) bool {
			builder.A2().Add(name, fun)

			return true
		})

		f.A3().ForEach(func(fun Function3, name string) bool {
			builder.A3().Add(name, fun)

			return true
		})

		f.A4().ForEach(func(fun Function4, name string) bool {
			builder.A4().Add(name, fun)

			return true
		})

		f.Var().ForEach(func(fun Function, name string) bool {
			builder.Var().Add(name, fun)

			return true
		})
	}

	return builder
}

func newRootFunctionsBuilder() *FunctionsBuilder {
	return newNamespaceFunctionsBuilder("")
}

func newNamespaceFunctionsBuilder(namespace string) *FunctionsBuilder {
	return &FunctionsBuilder{
		namespace: namespace,
		av:        newFnDef[Function](namespace),
		a0:        newFnDef[Function0](namespace),
		a1:        newFnDef[Function1](namespace),
		a2:        newFnDef[Function2](namespace),
		a3:        newFnDef[Function3](namespace),
		a4:        newFnDef[Function4](namespace),
	}
}

func newFunctionsBuilderInternalFrom(namespace string, other *FunctionsBuilder) *FunctionsBuilder {
	return &FunctionsBuilder{
		namespace: namespace,
		av:        newFnDefFrom[Function](namespace, other.av),
		a0:        newFnDefFrom[Function0](namespace, other.a0),
		a1:        newFnDefFrom[Function1](namespace, other.a1),
		a2:        newFnDefFrom[Function2](namespace, other.a2),
		a3:        newFnDefFrom[Function3](namespace, other.a3),
		a4:        newFnDefFrom[Function4](namespace, other.a4),
	}
}

func (b *FunctionsBuilder) Has(name string) bool {
	fname := makeFunctionName(b.namespace, name)

	if _, ok := b.av.data[fname]; ok {
		return true
	}

	if _, ok := b.a0.data[fname]; ok {
		return true
	}

	if _, ok := b.a1.data[fname]; ok {
		return true
	}

	if _, ok := b.a2.data[fname]; ok {
		return true
	}

	if _, ok := b.a3.data[fname]; ok {
		return true
	}

	if _, ok := b.a4.data[fname]; ok {
		return true
	}

	return false
}

func (b *FunctionsBuilder) Var() FnDef[Function] {
	return b.av
}

func (b *FunctionsBuilder) A0() FnDef[Function0] {
	return b.a0
}

func (b *FunctionsBuilder) A1() FnDef[Function1] {
	return b.a1
}

func (b *FunctionsBuilder) A2() FnDef[Function2] {
	return b.a2
}

func (b *FunctionsBuilder) A3() FnDef[Function3] {
	return b.a3
}

func (b *FunctionsBuilder) A4() FnDef[Function4] {
	return b.a4
}

func (b *FunctionsBuilder) From(other FunctionDefs) FunctionDefs {
	if other == nil {
		return b
	}

	other.A0().ForEach(func(fun Function0, name string) bool {
		b.a0.Add(name, fun)

		return true
	})

	other.A1().ForEach(func(fun Function1, name string) bool {
		b.a1.Add(name, fun)

		return true
	})

	other.A2().ForEach(func(fun Function2, name string) bool {
		b.a2.Add(name, fun)

		return true
	})

	other.A3().ForEach(func(fun Function3, name string) bool {
		b.a3.Add(name, fun)

		return true
	})

	other.A4().ForEach(func(fun Function4, name string) bool {
		b.a4.Add(name, fun)

		return true
	})

	other.Var().ForEach(func(fun Function, name string) bool {
		b.av.Add(name, fun)

		return true
	})

	return b
}

func (b *FunctionsBuilder) Build() (*Functions, error) {
	errs := slices.Concat(
		b.av.errors,
		b.a0.errors,
		b.a1.errors,
		b.a2.errors,
		b.a3.errors,
		b.a4.errors,
	)

	flookup := make(map[string]struct{})
	fnames := make([]string, 0, len(b.av.data)+len(b.a0.data)+len(b.a1.data)+len(b.a2.data)+len(b.a3.data)+len(b.a4.data))
	listables := []listable{b.av, b.a0, b.a1, b.a2, b.a3, b.a4}

	var exit bool

	for _, l := range listables {
		names := l.List()
		for _, name := range names {
			if _, exists := flookup[name]; exists {
				errs = append(errs, fmt.Errorf("function with name '%s' already exists in '%s' namespace", name, b.namespace))
				exit = true
				break
			}

			flookup[name] = struct{}{}
			fnames = append(fnames, name)
		}

		if exit {
			break
		}
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to build functions: %d error(s) occurred: %v", len(errs), errs)
	}

	registry := new(Functions)

	if len(b.av.data) > 0 {
		registry.av = NewFunctionCollectionFromMap(b.av.data)
	}

	if len(b.a0.data) > 0 {
		registry.a0 = NewFunctionCollectionFromMap(b.a0.data)
	}

	if len(b.a1.data) > 0 {
		registry.a1 = NewFunctionCollectionFromMap(b.a1.data)
	}

	if len(b.a2.data) > 0 {
		registry.a2 = NewFunctionCollectionFromMap(b.a2.data)
	}

	if len(b.a3.data) > 0 {
		registry.a3 = NewFunctionCollectionFromMap(b.a3.data)
	}

	if len(b.a4.data) > 0 {
		registry.a4 = NewFunctionCollectionFromMap(b.a4.data)
	}

	registry.names = fnames
	registry.size = len(fnames)
	registry.hash = functionsHash(registry)

	return registry, nil
}
