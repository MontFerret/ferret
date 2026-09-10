package uapi

import (
	"github.com/MontFerret/api"
	apidiagnostics "github.com/MontFerret/api/diagnostics"
	ferretdiagnostics "github.com/MontFerret/ferret/v2/pkg/diagnostics"
)

func wrapDiagnosticError(err error) error {
	if err == nil {
		return nil
	}

	var values apidiagnostics.Diagnostics
	var nativeFound bool
	seen := make(map[*ferretdiagnostics.Diagnostic]struct{})

	var visit func(error)
	visit = func(current error) {
		if current == nil {
			return
		}

		// Revisit the original tree rather than counting both an earlier
		// projection and the diagnostics from which it was constructed.
		if projected, ok := current.(*diagnosticError); ok {
			visit(projected.cause)

			return
		}

		if portable, ok := current.(apidiagnostics.Diagnostics); ok {
			values = append(values, portable...)

			return
		}

		// Inspect this node only: errors.As would select the first descendant
		// and hide later branches or change diagnostic ordering.
		if diagnostic, ok := current.(*ferretdiagnostics.Diagnostic); ok { //nolint:errorlint // Explicit error-tree traversal requires a node-local assertion.
			if diagnostic == nil {
				return
			}

			if _, exists := seen[diagnostic]; exists {
				return
			}

			seen[diagnostic] = struct{}{}
			nativeFound = true
			values = append(values, convertDiagnostic(diagnostic))
		}

		switch current := current.(type) { //nolint:errorlint // Traverse immediate children in their original order.
		case interface{ Unwrap() []error }:
			for _, child := range current.Unwrap() {
				visit(child)
			}
		case interface{ Unwrap() error }:
			visit(current.Unwrap())
		}
	}

	visit(err)

	if !nativeFound {
		return err
	}

	return newDiagnosticError(err, values)
}

func convertDiagnostic(value *ferretdiagnostics.Diagnostic) apidiagnostics.Diagnostic {
	result := apidiagnostics.Diagnostic{
		Source:  api.NewSource(value.Source.Name(), value.Source.Content()),
		Kind:    apidiagnostics.Kind(value.Kind.String()),
		Message: value.Message, Hint: value.Hint, Note: value.Note,
		Annotations: make([]apidiagnostics.Annotation, len(value.Spans)),
	}

	for index, span := range value.Spans {
		result.Annotations[index] = apidiagnostics.Annotation{
			Range: value.Source.RangeAt(span.Span), Message: span.Label, Primary: span.Main,
		}
	}

	return result
}
