package runtime

const emptyNS = ""

type (
	// Library represents a collection of functions organized in namespaces.
	// It provides methods to create nested namespaces and register functions within those namespaces,
	// as well as a method to build the final Functions instance.
	Library interface {
		Namespace

		// Size returns the number of function definitions, including arity overloads.
		Size() int

		// Build constructs and returns a finalized Functions instance or an error if the build process fails.
		Build() (*Functions, error)
	}

	// Namespace represents a namespace that can contain functions and nested namespaces.
	// It provides methods to create nested namespaces and register functions within those namespaces.
	Namespace interface {
		// Namespace creates or reuses a nested namespace with the given case-insensitive name.
		Namespace(name string) Namespace
		// Function returns a FunctionDefs interface that allows registering functions within this namespace.
		Function() FunctionDefs
	}

	library struct {
		builder    *FunctionsBuilder
		namespaces *registeredDisplayNames
		name       string
	}
)

func NewLibrary() Library {
	lib := new(library)
	lib.namespaces = newRegisteredDisplayNames()
	lib.builder = newRootFunctionsBuilder()

	return lib
}

func NewNamespace(name string) Namespace {
	lib := new(library)
	lib.namespaces = newRegisteredDisplayNames()
	lib.name = lib.namespaces.Declare(NormalizeRegisteredName(name), name)
	lib.builder = newNamespacedFunctionsBuilder(lib.name)

	return lib
}

func makeFunctionName(namespace, name string) string {
	if namespace == emptyNS {
		return name
	}

	return namespace + NamespaceSeparator + name
}

func (lib *library) Size() int {
	return lib.builder.Size()
}

func (lib *library) Namespace(name string) Namespace {
	newLib := new(library)
	newLib.namespaces = lib.namespaces
	qualifiedName := makeFunctionName(lib.name, name)
	newLib.name = lib.namespaces.Declare(NormalizeRegisteredName(qualifiedName), qualifiedName)
	newLib.builder = newFunctionsBuilderInternalFrom(newLib.name, lib.builder)

	return newLib
}

func (lib *library) Function() FunctionDefs {
	return lib.builder
}

func (lib *library) Build() (*Functions, error) {
	return lib.builder.Build()
}
