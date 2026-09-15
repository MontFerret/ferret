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
		{ID: "arrays", Title: "Arrays", Description: "Functions for creating, transforming, querying, and modifying arrays."},
		{ID: "collections", Title: "Collections", Description: "Functions for working with collections and collection values."},
		{ID: "crypto", Title: "Crypto", Description: "Functions for hashing, encoding, and generating secure values."},
		{ID: "datetime", Title: "Date & Time", Description: "Functions for creating, parsing, formatting, and manipulating dates and times."},
		{ID: "encoding", Title: "Encoding", Description: "Functions for encoding, decoding, serializing, and escaping values."},
		{ID: "io", Title: "I/O", Description: "Functions for working with files, networks, and other input and output operations."},
		{ID: "math", Title: "Math", Description: "Functions for mathematical operations and numeric calculations."},
		{ID: "objects", Title: "Objects", Description: "Functions for creating, transforming, querying, and modifying objects."},
		{ID: "path", Title: "Path", Description: "Functions for constructing, inspecting, and manipulating paths."},
		{ID: "random", Title: "Random", Description: "Pseudo-random value generation functions in the random namespace."},
		{ID: "strings", Title: "Strings", Description: "Functions for creating, transforming, searching, and inspecting strings."},
		{ID: "testing", Title: "Testing", Description: "Functions for assertions and testing Ferret queries."},
		{ID: "types", Title: "Types", Description: "Functions for inspecting, checking, and converting value types."},
		{ID: "utils", Title: "Utilities", Description: "General-purpose utility functions."},
	}

	// categoryOverrides is intentionally keyed by canonical Ferret identity.
	// Add only source-layout exceptions; the ordinary path derives from pkg/stdlib.
	categoryOverrides = map[functionIdentity]string{
		{Name: "range"}: "arrays",
	}
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
