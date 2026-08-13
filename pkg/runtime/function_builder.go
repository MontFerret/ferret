package runtime

import (
	"encoding/binary"
	"fmt"
	"hash/fnv"
	"sort"
)

type (
	FnDef[T FunctionConstraint] interface {
		// Add adds a function using a case-insensitive qualified identity.
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
		// Size returns the total number of function definitions present in the collection.
		Size() int
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
		av        *defaultFnDef[Function]
		a0        *defaultFnDef[Function0]
		a1        *defaultFnDef[Function1]
		a2        *defaultFnDef[Function2]
		a3        *defaultFnDef[Function3]
		a4        *defaultFnDef[Function4]
		names     *registeredDisplayNames
		namespace string
	}

	listable interface {
		List() []string
	}

	arityList struct {
		list   listable
		marker byte
	}

	// fnErrors aggregates build errors shared across nested function definitions.
	fnErrors struct {
		items []error
	}

	// defaultFnDef stores functions for a namespace and shares data/errors with nested builders.
	defaultFnDef[T FunctionConstraint] struct {
		errors    *fnErrors
		data      map[string]registeredFunction[T]
		names     *registeredDisplayNames
		namespace string
	}
)

// Add appends an error to the shared list.
func (e *fnErrors) Add(err error) {
	e.items = append(e.items, err)
}

// All returns all collected errors from the shared list.
func (e *fnErrors) All() []error {
	return e.items
}

// newFnDef creates a function definition with a shared error container.
func newFnDef[T FunctionConstraint](namespace string, errs *fnErrors, names *registeredDisplayNames) *defaultFnDef[T] {
	return &defaultFnDef[T]{
		namespace: namespace,
		errors:    errs,
		names:     names,
		data:      make(map[string]registeredFunction[T]),
	}
}

// newFnDefFrom reuses the parent data and error container for nested builders.
func newFnDefFrom[T FunctionConstraint](namespace string, other *defaultFnDef[T]) *defaultFnDef[T] {
	return &defaultFnDef[T]{
		namespace: namespace,
		errors:    other.errors,
		names:     other.names,
		// We share the same map across all builders to ensure that changes in one builder are reflected in all builders that share the same namespace.
		data: other.data,
	}
}

func (fd *defaultFnDef[T]) addError(err error) {
	if fd.errors == nil {
		fd.errors = &fnErrors{}
	}

	fd.errors.Add(err)
}

func (fd *defaultFnDef[T]) Add(name string, fn T) FnDef[T] {
	registeredName := makeFunctionName(fd.namespace, name)
	if !HasTerminalFunctionName(registeredName) {
		fd.addError(fmt.Errorf("function name cannot be empty in '%s' namespace", fd.namespace))

		return fd
	}

	key := NormalizeRegisteredName(registeredName)
	if _, exists := fd.data[key]; exists {
		fd.addError(fmt.Errorf("function with name '%s' already exists in '%s' namespace", name, fd.namespace))

		return fd
	}

	fd.data[key] = registeredFunction[T]{name: registeredName, function: fn}
	fd.names.Add(key, registeredName)

	return fd
}

func (fd *defaultFnDef[T]) Remove(name string) FnDef[T] {
	key := NormalizeRegisteredName(makeFunctionName(fd.namespace, name))

	if _, exists := fd.data[key]; !exists {
		fd.addError(fmt.Errorf("function with name '%s' does not exist in '%s' namespace", name, fd.namespace))

		return fd
	}

	delete(fd.data, key)
	fd.names.Remove(key)

	return fd
}

func (fd *defaultFnDef[T]) Has(name string) bool {
	key := NormalizeRegisteredName(makeFunctionName(fd.namespace, name))
	_, exists := fd.data[key]

	return exists
}

func (fd *defaultFnDef[T]) ForEach(fn func(fn T, name string) bool) {
	for _, entry := range fd.data {
		if !fn(entry.function, entry.name) {
			break
		}
	}
}

func (fd *defaultFnDef[T]) List() []string {
	names := make([]string, 0, len(fd.data))

	for _, entry := range fd.data {
		names = append(names, entry.name)
	}

	return names
}

// NewFunctionsBuilder creates an empty host function registry builder.
// Host-function qualified names preserve their declared spelling and resolve case-insensitively.
// A logical name may have one definition at each fixed arity and one variadic definition.
func NewFunctionsBuilder() *FunctionsBuilder {
	return newRootFunctionsBuilder()
}

// NewFunctionsBuilderFrom creates a builder containing every definition from funcs.
// Merging rejects only definitions with the same qualified name and arity.
func NewFunctionsBuilderFrom(funcs ...*Functions) *FunctionsBuilder {
	builder := newRootFunctionsBuilder()

	for _, f := range funcs {
		if f == nil {
			continue
		}

		for _, name := range f.List() {
			builder.names.Declare(NormalizeRegisteredName(name), name)
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
	return newNamespacedFunctionsBuilder("")
}

// newNamespacedFunctionsBuilder creates a builder with shared errors across its FnDefs.
func newNamespacedFunctionsBuilder(namespace string) *FunctionsBuilder {
	errs := &fnErrors{}
	names := newRegisteredDisplayNames()

	return &FunctionsBuilder{
		namespace: namespace,
		names:     names,
		av:        newFnDef[Function](namespace, errs, names),
		a0:        newFnDef[Function0](namespace, errs, names),
		a1:        newFnDef[Function1](namespace, errs, names),
		a2:        newFnDef[Function2](namespace, errs, names),
		a3:        newFnDef[Function3](namespace, errs, names),
		a4:        newFnDef[Function4](namespace, errs, names),
	}
}

// newFunctionsBuilderInternalFrom creates a nested builder sharing parent data and errors.
func newFunctionsBuilderInternalFrom(namespace string, other *FunctionsBuilder) *FunctionsBuilder {
	return &FunctionsBuilder{
		namespace: namespace,
		names:     other.names,
		av:        newFnDefFrom[Function](namespace, other.av),
		a0:        newFnDefFrom[Function0](namespace, other.a0),
		a1:        newFnDefFrom[Function1](namespace, other.a1),
		a2:        newFnDefFrom[Function2](namespace, other.a2),
		a3:        newFnDefFrom[Function3](namespace, other.a3),
		a4:        newFnDefFrom[Function4](namespace, other.a4),
	}
}

func (b *FunctionsBuilder) Size() int {
	var size int

	size += len(b.av.data)
	size += len(b.a0.data)
	size += len(b.a1.data)
	size += len(b.a2.data)
	size += len(b.a3.data)
	size += len(b.a4.data)

	return size
}

func (b *FunctionsBuilder) Has(name string) bool {
	fname := NormalizeRegisteredName(makeFunctionName(b.namespace, name))

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

	if source, ok := other.(*FunctionsBuilder); ok {
		b.names.Seed(source.names)
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
	errs := collectFnErrors(b.av)

	if len(errs) > 0 {
		return nil, fmt.Errorf(
			"failed to build functions: %d error(s) occurred: %v",
			len(errs),
			errs,
		)
	}

	collections := []arityList{
		{marker: 0, list: b.a0},
		{marker: 1, list: b.a1},
		{marker: 2, list: b.a2},
		{marker: 3, list: b.a3},
		{marker: 4, list: b.a4},
		{marker: 0xff, list: b.av},
	}

	flookup := make(map[string]struct{})
	fnames := make([]string, 0, b.Size())
	hasher := fnv.New64a()
	var nameLen [8]byte

	for _, collection := range collections {
		names := collection.list.List()
		sort.Slice(names, func(i, j int) bool {
			return NormalizeRegisteredName(names[i]) < NormalizeRegisteredName(names[j])
		})

		for _, name := range names {
			key := NormalizeRegisteredName(name)
			binary.LittleEndian.PutUint64(nameLen[:], uint64(len(key)))
			_, _ = hasher.Write([]byte{collection.marker})
			_, _ = hasher.Write(nameLen[:])
			_, _ = hasher.Write([]byte(key))

			if _, exists := flookup[key]; exists {
				continue
			}

			flookup[key] = struct{}{}
			displayName, exists := b.names.Name(key)
			if !exists {
				displayName = name
			}

			fnames = append(fnames, displayName)
		}
	}

	sort.Slice(fnames, func(i, j int) bool {
		return NormalizeRegisteredName(fnames[i]) < NormalizeRegisteredName(fnames[j])
	})

	registry := new(Functions)

	if len(b.av.data) > 0 {
		registry.av = newNormalizedFunctionCollectionFromMap(b.av.data)
	}

	if len(b.a0.data) > 0 {
		registry.a0 = newNormalizedFunctionCollectionFromMap(b.a0.data)
	}

	if len(b.a1.data) > 0 {
		registry.a1 = newNormalizedFunctionCollectionFromMap(b.a1.data)
	}

	if len(b.a2.data) > 0 {
		registry.a2 = newNormalizedFunctionCollectionFromMap(b.a2.data)
	}

	if len(b.a3.data) > 0 {
		registry.a3 = newNormalizedFunctionCollectionFromMap(b.a3.data)
	}

	if len(b.a4.data) > 0 {
		registry.a4 = newNormalizedFunctionCollectionFromMap(b.a4.data)
	}

	registry.names = fnames
	registry.size = b.Size()
	registry.hash = hasher.Sum64()

	return registry, nil
}

// collectFnErrors returns errors collected for a definition, if any.
func collectFnErrors[T FunctionConstraint](fd *defaultFnDef[T]) []error {
	if fd == nil || fd.errors == nil {
		return nil
	}

	return fd.errors.All()
}
