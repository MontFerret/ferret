package analyzer

import (
	"strings"
	"testing"
)

func TestParseSignatureDocumentationRequiresCompleteMetadata(t *testing.T) {
	registered := registeredSignature{QualifiedName: "TEST", Arity: 1}

	for _, test := range []struct {
		name          string
		documentation string
		reason        string
	}{
		{name: "missing", reason: "has no Ferret API documentation"},
		{name: "malformed", documentation: "Test.\n@param {Any} value - Value.\n@return {Any} Value.", reason: "malformed Ferret API documentation"},
		{name: "missing prose", documentation: "@param value {Any} Value.\n@return {Any} Value.", reason: "non-empty Ferret-facing prose"},
		{name: "missing return", documentation: "Test.\n@param value {Any} Value.", reason: "exactly one @return"},
		{name: "arity mismatch", documentation: "Test.\n@return {Any} Value.", reason: "fixed arity 1"},
	} {
		t.Run(test.name, func(t *testing.T) {
			declaration := testDeclaration(test.documentation)
			_, err := parseSignatureDocumentation(registered, declaration)
			if err == nil || !strings.Contains(err.Error(), test.reason) {
				t.Fatalf("error = %v, want reason containing %q", err, test.reason)
			}

			for _, expected := range []string{"sample.go:12", "package example.com/sample", "Ferret TEST", "Go declaration Sample"} {
				if !strings.Contains(err.Error(), expected) {
					t.Fatalf("error = %q, want diagnostic context %q", err, expected)
				}
			}
		})
	}
}

func TestParseSignatureDocumentationSupportsVariadicMetadata(t *testing.T) {
	registered := registeredSignature{QualifiedName: "JOIN", Variadic: true}
	declaration := testDeclaration("Joins values.\n@param value {Any, repeated} Value to join.\n@return {String} Joined value.")

	signature, err := parseSignatureDocumentation(registered, declaration)
	if err != nil {
		t.Fatal(err)
	}

	if !signature.Variadic || len(signature.Parameters) != 1 {
		t.Fatalf("signature = %#v, want one variadic logical parameter", signature)
	}
}

func TestParseAssertionDocumentationExpandsPrefixes(t *testing.T) {
	descriptor := assertionDescriptor{
		Declaration: testDeclaration("Tests values.\n@param actual {Any} Actual value.\n@param expected {Any} Expected value.\n@param message {String} Failure message.\n@return {Boolean} Assertion result."),
		Min:         2,
		Max:         3,
	}

	for arity := 2; arity <= 3; arity++ {
		registered := registeredSignature{QualifiedName: "t::eq", Arity: arity}
		signature, err := parseAssertionDocumentation(registered, descriptor)
		if err != nil {
			t.Fatalf("arity %d: %v", arity, err)
		}

		if len(signature.Parameters) != arity {
			t.Fatalf("arity %d generated %d parameters", arity, len(signature.Parameters))
		}
	}
}

func testDeclaration(documentation string) *sourceDeclaration {
	return &sourceDeclaration{
		Name:          "Sample",
		PackagePath:   "example.com/sample",
		File:          "sample.go",
		Line:          12,
		Documentation: documentation,
	}
}
