package analyzer

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/MontFerret/specs/pkg/api"
	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

const stdlibPackageMarker = "/pkg/stdlib/"

var (
	standardLibraryCategories = []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Global functions for working with arrays."},
		{ID: "collections", Title: "Collections", Description: "Global functions for working with collections."},
		{ID: "datetime", Title: "Date & Time", Description: "Global functions for working with dates and times."},
		{ID: "math", Title: "Math", Description: "Mathematical and numeric global functions."},
		{ID: "objects", Title: "Objects", Description: "Global functions for working with objects."},
		{ID: "path", Title: "Path", Description: "Global functions for working with paths."},
		{ID: "strings", Title: "Strings", Description: "Global functions for working with strings."},
		{ID: "types", Title: "Types", Description: "Global functions for type checks and conversions."},
		{ID: "utils", Title: "Utilities", Description: "General-purpose global utility functions."},
	}

	// categoryOverrides is intentionally keyed by the canonical global Ferret name.
	// Add only source-layout exceptions; the ordinary path derives from pkg/stdlib.
	categoryOverrides = map[string]string{}
)

func buildStandardLibraryCatalog(
	version string,
	registered []registeredSignature,
	reference *api.Reference,
	metadata []apicatalog.Category,
	overrides map[string]string,
) (*apicatalog.Catalog, error) {
	categories := make([]apicatalog.Category, len(metadata))
	categoryFunctions := make(map[string]map[string]struct{}, len(metadata))
	for index, category := range metadata {
		categories[index] = apicatalog.Category{
			ID:          category.ID,
			Title:       category.Title,
			Description: category.Description,
		}
		categoryFunctions[category.ID] = make(map[string]struct{})
	}

	resolved := make(map[string]string)
	problems := make([]error, 0)
	for _, entry := range registered {
		if entry.Namespace != "" {
			continue
		}

		categoryID, err := resolveCategory(entry, overrides)
		if err != nil {
			problems = append(problems, err)

			continue
		}

		functions, exists := categoryFunctions[categoryID]
		if !exists {
			problems = append(problems, fmt.Errorf("global function %s resolves to unknown category %q", entry.Name, categoryID))

			continue
		}

		if previous, exists := resolved[entry.Name]; exists && previous != categoryID {
			problems = append(problems, fmt.Errorf("global function %s overloads resolve to categories %q and %q", entry.Name, previous, categoryID))

			continue
		}

		resolved[entry.Name] = categoryID
		functions[entry.Name] = struct{}{}
	}

	for index := range categories {
		functions := categoryFunctions[categories[index].ID]
		if len(functions) == 0 {
			problems = append(problems, fmt.Errorf("category %q has no global functions", categories[index].ID))

			continue
		}

		categories[index].Functions = make([]string, 0, len(functions))
		for function := range functions {
			categories[index].Functions = append(categories[index].Functions, function)
		}

		sort.Strings(categories[index].Functions)
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	return &apicatalog.Catalog{
		SchemaVersion:  apicatalog.SchemaVersion,
		ID:             moduleID,
		Version:        version,
		Categories:     categories,
		NamespaceRoots: discoverNamespaceRoots(reference),
	}, nil
}

func resolveCategory(entry registeredSignature, overrides map[string]string) (string, error) {
	if categoryID, exists := overrides[entry.Name]; exists {
		return categoryID, nil
	}

	_, suffix, found := strings.Cut(entry.PackagePath, stdlibPackageMarker)
	if !found || suffix == "" {
		return "", fmt.Errorf("global function %s source package %q is outside pkg/stdlib", entry.Name, entry.PackagePath)
	}

	categoryID, _, _ := strings.Cut(suffix, "/")

	return categoryID, nil
}

func discoverNamespaceRoots(reference *api.Reference) []string {
	roots := make(map[string]struct{})
	for _, namespace := range reference.Namespaces {
		if namespace.Name == "" {
			continue
		}

		root, _, _ := strings.Cut(namespace.Name, "::")
		roots[root] = struct{}{}
	}

	result := make([]string, 0, len(roots))
	for root := range roots {
		result = append(result, root)
	}

	sort.Strings(result)

	return result
}
