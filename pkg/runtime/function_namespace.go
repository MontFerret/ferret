package runtime

import (
	"strings"
)

const NamespaceSeparator = "::"
const emptyNS = ""

type (
	// RootNamespace is the top-level namespace that can contain multiple nested namespaces and functions.
	// It provides methods to create nested namespaces and register functions within those namespaces.
	RootNamespace interface {
		Namespace

		// Build constructs and returns a finalized Functions instance or an error if the build process fails.
		Build() (*Functions, error)
	}

	// Namespace represents a namespace that can contain functions and nested namespaces.
	// It provides methods to create nested namespaces and register functions within those namespaces.
	Namespace interface {
		// Namespace creates a new nested namespace with the given name and returns it.
		Namespace(name string) Namespace
		// Function returns a FunctionDefs interface that allows registering functions within this namespace.
		Function() FunctionDefs
	}

	defaultNamespace struct {
		builder *FunctionsBuilder
		name    string
	}
)

func NewRootNamespace() RootNamespace {
	ns := new(defaultNamespace)
	ns.builder = newRootFunctionsBuilder()

	return ns
}

func NewNamespace(name string) Namespace {
	ns := new(defaultNamespace)
	ns.name = strings.ToUpper(name)
	ns.builder = newNamespaceFunctionsBuilder(ns.name)

	return ns
}

func makeFunctionName(namespace, name string) string {
	name = strings.ToUpper(name)

	if namespace == emptyNS {
		return name
	}

	return namespace + NamespaceSeparator + name
}

func (nc *defaultNamespace) Namespace(name string) Namespace {
	ns := new(defaultNamespace)
	ns.name = makeFunctionName(nc.name, name)
	ns.builder = newFunctionsBuilderInternalFrom(ns.name, nc.builder)

	return ns
}

func (nc *defaultNamespace) Function() FunctionDefs {
	return nc.builder
}

func (nc *defaultNamespace) Build() (*Functions, error) {
	return nc.builder.Build()
}
