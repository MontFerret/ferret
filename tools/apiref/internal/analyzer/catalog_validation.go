package analyzer

import (
	"errors"
	"fmt"

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

	apiFunctions := make(map[functionIdentity]struct{})
	apiNamespaces := make(map[string]map[string]struct{}, len(reference.Namespaces))
	for _, namespace := range reference.Namespaces {
		functions := make(map[string]struct{}, len(namespace.Functions))
		apiNamespaces[namespace.Name] = functions
		for _, function := range namespace.Functions {
			functions[function.Name] = struct{}{}
			apiFunctions[functionIdentity{Namespace: namespace.Name, Name: function.Name}] = struct{}{}
		}
	}

	categorized := make(map[functionIdentity]string)
	for _, category := range catalog.Categories {
		for _, function := range category.Functions {
			identity := functionIdentity{Namespace: function.Namespace, Name: function.Name}
			namespace, exists := apiNamespaces[function.Namespace]

			if !exists {
				problems = append(problems, fmt.Errorf("catalog category %q references unknown API namespace %q", category.ID, function.Namespace))
			} else if _, exists := namespace[function.Name]; !exists {
				problems = append(problems, fmt.Errorf("catalog category %q references unknown function %q in API namespace %q", category.ID, function.Name, function.Namespace))
			}

			if previous, exists := categorized[identity]; exists {
				problems = append(problems, fmt.Errorf("function %q is assigned to categories %q and %q", identity, previous, category.ID))
			}

			categorized[identity] = category.ID
		}
	}

	for function := range apiFunctions {
		if _, exists := categorized[function]; !exists {
			problems = append(problems, fmt.Errorf("function %q is not assigned to a catalog category", function))
		}
	}

	return errors.Join(problems...)
}
