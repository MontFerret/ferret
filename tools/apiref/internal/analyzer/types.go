package analyzer

import (
	"go/ast"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"

	"github.com/MontFerret/ferret/v2/pkg/runtime"
)

const moduleID = "montferret/core"

type (
	// Options identifies the source tree, API version, and authoritative runtime registry to analyze.
	Options struct {
		Functions *runtime.Functions
		Root      string
		Version   string
	}

	// Artifacts contains the API and presentation catalog generated from one analysis pass.
	Artifacts struct {
		Reference *api.Reference
		Catalog   *apicatalog.Catalog
	}

	registeredSignature struct {
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

	sourceDeclaration struct {
		Expression    ast.Expr
		Key           string
		Name          string
		PackagePath   string
		File          string
		Documentation string
		Factory       string
		Line          int
	}

	assertionDescriptor struct {
		Declaration *sourceDeclaration
		Min         int
		Max         int
	}

	sourceCatalog struct {
		Declarations map[string]*sourceDeclaration
		Assertions   map[string]assertionDescriptor
	}
)
