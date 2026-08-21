package analyzer

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
	"github.com/MontFerret/specs/pkg/validation"
)

// Generate builds and validates the montferret/core API Reference and catalog from one registry and source analysis.
func Generate(ctx context.Context, options Options) (*Artifacts, error) {
	if strings.TrimSpace(options.Root) == "" {
		return nil, fmt.Errorf("repository root is required")
	}

	registered, err := registeredSignatures(options.Functions)
	if err != nil {
		return nil, err
	}

	catalog, err := loadSourceCatalog(ctx, options.Root)
	if err != nil {
		return nil, err
	}

	reference, err := buildReference(options.Version, registered, catalog)
	if err != nil {
		return nil, err
	}

	if err := api.Validate(reference); err != nil {
		var validationErrors *validation.Errors
		if errors.As(err, &validationErrors) {
			problems := make([]error, 0, len(validationErrors.Violations))
			for _, violation := range validationErrors.Violations {
				problems = append(problems, fmt.Errorf("%s (%s): %s", violation.Path, violation.Rule, violation.Message))
			}

			return nil, fmt.Errorf("validate generated %s API Reference %s: %w", moduleID, options.Version, errors.Join(problems...))
		}

		return nil, fmt.Errorf("validate generated %s API Reference %s: %w", moduleID, options.Version, err)
	}

	stdlibCatalog, err := buildStandardLibraryCatalog(options.Version, registered, standardLibraryCategories, categoryOverrides)
	if err != nil {
		return nil, err
	}

	if err := apicatalog.Validate(stdlibCatalog); err != nil {
		return nil, fmt.Errorf("validate generated %s Standard Library catalog %s: %w", moduleID, options.Version, err)
	}

	if err := validateCatalogAgainstReference(reference, stdlibCatalog); err != nil {
		return nil, fmt.Errorf("validate generated %s API/catalog pair %s: %w", moduleID, options.Version, err)
	}

	return &Artifacts{Reference: reference, Catalog: stdlibCatalog}, nil
}

func buildReference(version string, registered []registeredSignature, catalog *sourceCatalog) (*api.Reference, error) {
	namespaces := make(map[string]map[string][]api.Signature)
	assertionArities := make(map[string][]int)
	problems := make([]error, 0)

	for _, entry := range registered {
		var signature api.Signature
		var err error

		if descriptor, isAssertion := catalog.Assertions[entry.QualifiedName]; isAssertion {
			signature, err = parseAssertionDocumentation(entry, descriptor)
			assertionArities[entry.QualifiedName] = append(assertionArities[entry.QualifiedName], entry.Arity)
		} else {
			declaration, exists := catalog.Declarations[entry.Symbol]
			if !exists {
				problems = append(problems, fmt.Errorf(
					"%s:%d: package %s; Ferret %s; Go declaration %s: source declaration was not found in pkg/stdlib",
					entry.File,
					entry.Line,
					entry.PackagePath,
					entry.QualifiedName,
					entry.Symbol,
				))

				continue
			}

			signature, err = parseSignatureDocumentation(entry, declaration)
		}

		if err != nil {
			problems = append(problems, err)

			continue
		}

		functions := namespaces[entry.Namespace]
		if functions == nil {
			functions = make(map[string][]api.Signature)
			namespaces[entry.Namespace] = functions
		}

		functions[entry.Name] = append(functions[entry.Name], signature)
	}

	for qualifiedName, arities := range assertionArities {
		descriptor := catalog.Assertions[qualifiedName]
		if err := validateAssertionArities(qualifiedName, descriptor, arities); err != nil {
			problems = append(problems, err)
		}
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	reference := &api.Reference{
		SchemaVersion: api.SchemaVersion,
		ID:            moduleID,
		Version:       version,
		Namespaces:    make([]api.Namespace, 0, len(namespaces)),
	}

	for namespaceName, functionMap := range namespaces {
		namespace := api.Namespace{
			Name:      namespaceName,
			Functions: make([]api.Function, 0, len(functionMap)),
		}

		for functionName, signatures := range functionMap {
			sort.Slice(signatures, func(i, j int) bool {
				if signatures[i].Variadic != signatures[j].Variadic {
					return !signatures[i].Variadic
				}

				return len(signatures[i].Parameters) < len(signatures[j].Parameters)
			})

			namespace.Functions = append(namespace.Functions, api.Function{
				Name:       functionName,
				Signatures: signatures,
			})
		}

		sort.Slice(namespace.Functions, func(i, j int) bool {
			return namespace.Functions[i].Name < namespace.Functions[j].Name
		})

		reference.Namespaces = append(reference.Namespaces, namespace)
	}

	sort.Slice(reference.Namespaces, func(i, j int) bool {
		return reference.Namespaces[i].Name < reference.Namespaces[j].Name
	})

	return reference, nil
}

func validateAssertionArities(qualifiedName string, descriptor assertionDescriptor, arities []int) error {
	sort.Ints(arities)
	if len(arities) != descriptor.Max-descriptor.Min+1 {
		return diagnostic(
			descriptor.Declaration,
			qualifiedName,
			"assertion descriptor",
			fmt.Sprintf("runtime overload count %d does not match assertion args range %d..%d", len(arities), descriptor.Min, descriptor.Max),
		)
	}

	for index, arity := range arities {
		expected := descriptor.Min + index
		if arity != expected {
			return diagnostic(
				descriptor.Declaration,
				qualifiedName,
				"assertion descriptor",
				fmt.Sprintf("runtime overload arity %d is not contiguous with assertion args range %d..%d", arity, descriptor.Min, descriptor.Max),
			)
		}
	}

	return nil
}
