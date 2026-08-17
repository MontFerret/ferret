package analyzer

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	apicatalog "github.com/MontFerret/specs/pkg/api/catalog"
)

const stdlibPackageMarker = "/pkg/stdlib/"

var (
	standardLibraryCategories = []apicatalog.Category{
		{ID: "arrays", Title: "Arrays", Description: "Global functions for working with arrays."},
		{ID: "collections", Title: "Collections", Description: "Global functions for working with collections."},
		{ID: "datetime", Title: "Date & Time", Description: "Global functions for working with dates and times."},
		{ID: "io", Title: "I/O", Description: "Functions for working with files, networks, and other input/output operations."},
		{ID: "math", Title: "Math", Description: "Mathematical and numeric global functions."},
		{ID: "objects", Title: "Objects", Description: "Global functions for working with objects."},
		{ID: "path", Title: "Path", Description: "Global functions for working with paths."},
		{ID: "strings", Title: "Strings", Description: "Global functions for working with strings."},
		{ID: "testing", Title: "Testing", Description: "Assertion functions for testing Ferret queries."},
		{ID: "types", Title: "Types", Description: "Global functions for type checks and conversions."},
		{ID: "utils", Title: "Utilities", Description: "General-purpose global utility functions."},
	}

	// categoryOverrides is intentionally keyed by canonical Ferret identity.
	// Add only source-layout exceptions; the ordinary path derives from pkg/stdlib.
	categoryOverrides = map[functionIdentity]string{}
)

type functionIdentity struct {
	Namespace string
	Name      string
}

func (identity functionIdentity) String() string {
	if identity.Namespace == "" {
		return identity.Name
	}

	return identity.Namespace + "::" + identity.Name
}

func buildStandardLibraryCatalog(
	version string,
	registered []registeredSignature,
	metadata []apicatalog.Category,
	overrides map[functionIdentity]string,
) (*apicatalog.Catalog, error) {
	categories := make([]apicatalog.Category, len(metadata))
	categoryFunctions := make(map[string]map[functionIdentity]struct{}, len(metadata))
	for index, category := range metadata {
		categories[index] = apicatalog.Category{
			ID:          category.ID,
			Title:       category.Title,
			Description: category.Description,
		}
		categoryFunctions[category.ID] = make(map[functionIdentity]struct{})
	}

	resolved := make(map[functionIdentity]string)
	problems := make([]error, 0)
	for _, entry := range registered {
		categoryID, err := resolveCategory(entry, overrides)
		if err != nil {
			problems = append(problems, err)

			continue
		}

		functions, exists := categoryFunctions[categoryID]
		if !exists {
			problems = append(problems, fmt.Errorf("function %s resolves to unknown category %q", entry.QualifiedName, categoryID))

			continue
		}

		identity := functionIdentity{Namespace: entry.Namespace, Name: entry.Name}
		if previous, exists := resolved[identity]; exists && previous != categoryID {
			problems = append(problems, fmt.Errorf("function %s overloads resolve to categories %q and %q", identity, previous, categoryID))

			continue
		}

		resolved[identity] = categoryID
		functions[identity] = struct{}{}
	}

	for index := range categories {
		functions := categoryFunctions[categories[index].ID]
		if len(functions) == 0 {
			problems = append(problems, fmt.Errorf("category %q has no functions", categories[index].ID))

			continue
		}

		categories[index].Functions = make([]apicatalog.FunctionRef, 0, len(functions))
		for function := range functions {
			categories[index].Functions = append(categories[index].Functions, apicatalog.FunctionRef{
				Namespace: function.Namespace,
				Name:      function.Name,
			})
		}

		sort.Slice(categories[index].Functions, func(left, right int) bool {
			leftFunction := categories[index].Functions[left]
			rightFunction := categories[index].Functions[right]
			if leftFunction.Namespace != rightFunction.Namespace {
				return leftFunction.Namespace < rightFunction.Namespace
			}

			return leftFunction.Name < rightFunction.Name
		})
	}

	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}

	return &apicatalog.Catalog{
		SchemaVersion: apicatalog.SchemaVersion,
		ID:            moduleID,
		Version:       version,
		Categories:    categories,
	}, nil
}

func resolveCategory(entry registeredSignature, overrides map[functionIdentity]string) (string, error) {
	identity := functionIdentity{Namespace: entry.Namespace, Name: entry.Name}
	if categoryID, exists := overrides[identity]; exists {
		return categoryID, nil
	}

	_, suffix, found := strings.Cut(entry.PackagePath, stdlibPackageMarker)
	if !found || suffix == "" {
		return "", fmt.Errorf("function %s source package %q is outside pkg/stdlib", identity, entry.PackagePath)
	}

	categoryID, _, _ := strings.Cut(suffix, "/")

	return categoryID, nil
}
