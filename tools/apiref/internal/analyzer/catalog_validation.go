package analyzer

import (
	"errors"
	"fmt"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

func validateCatalogAgainstReference(reference *api.Reference, catalog *apicatalog.Catalog) error {
	problems := make([]error, 0)
	if reference.ID != catalog.ID {
		problems = append(problems, fmt.Errorf("catalog id %q does not match API id %q", catalog.ID, reference.ID))
	}

	if reference.Version != catalog.Version {
		problems = append(problems, fmt.Errorf("catalog version %q does not match API version %q", catalog.Version, reference.Version))
	}

	globalFunctions := make(map[string]struct{})
	namespaceRoots := make(map[string]struct{})
	for _, namespace := range reference.Namespaces {
		if namespace.Name == "" {
			for _, function := range namespace.Functions {
				globalFunctions[function.Name] = struct{}{}
			}

			continue
		}

		root, _, _ := strings.Cut(namespace.Name, "::")
		namespaceRoots[root] = struct{}{}
	}

	categorized := make(map[string]string)
	for _, category := range catalog.Categories {
		for _, function := range category.Functions {
			if _, exists := globalFunctions[function]; !exists {
				problems = append(problems, fmt.Errorf("catalog category %q references unknown global function %q", category.ID, function))
			}

			if previous, exists := categorized[function]; exists {
				problems = append(problems, fmt.Errorf("global function %q is assigned to categories %q and %q", function, previous, category.ID))
			}

			categorized[function] = category.ID
		}
	}

	for function := range globalFunctions {
		if _, exists := categorized[function]; !exists {
			problems = append(problems, fmt.Errorf("global function %q is not assigned to a catalog category", function))
		}
	}

	declaredRoots := make(map[string]struct{}, len(catalog.NamespaceRoots))
	for _, root := range catalog.NamespaceRoots {
		declaredRoots[root] = struct{}{}
		if _, exists := namespaceRoots[root]; !exists {
			problems = append(problems, fmt.Errorf("catalog namespace root %q does not cover an API namespace", root))
		}
	}

	for root := range namespaceRoots {
		if _, exists := declaredRoots[root]; !exists {
			problems = append(problems, fmt.Errorf("API namespace root %q is not declared by the catalog", root))
		}
	}

	return errors.Join(problems...)
}
