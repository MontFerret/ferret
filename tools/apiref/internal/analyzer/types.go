package analyzer

import (
	"go/ast"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const moduleID = "montferret/core"

// Options identifies the source tree, API version, and authoritative runtime registry to analyze.
type Options struct {
	Functions *runtime.Functions
	Root      string
	Version   string
}

type registeredSignature struct {
	QualifiedName string
	Namespace     string
	Name          string
	Symbol        string
	PackagePath   string
	File          string
	Line          int
	Arity         int
	Variadic      bool
}

type sourceDeclaration struct {
	Expression    ast.Expr
	Key           string
	Name          string
	PackagePath   string
	File          string
	Documentation string
	Factory       string
	Line          int
}

type assertionDescriptor struct {
	Declaration *sourceDeclaration
	Min         int
	Max         int
}

type sourceCatalog struct {
	Declarations map[string]*sourceDeclaration
	Assertions   map[string]assertionDescriptor
}
