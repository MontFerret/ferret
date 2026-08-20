package analyzer

import (
	"reflect"
	"strings"
	"testing"

	"github.com/MontFerret/specs/pkg/api"
)

func TestBuildReferenceGroupsAndSortsRuntimeSignatures(t *testing.T) {
	catalog := &sourceCatalog{
		Declarations: map[string]*sourceDeclaration{
			"example.Zero": declarationWithDocs("Zero", "Returns zero.\n@return {Object?} Zero."),
			"example.One":  declarationWithDocs("One", "Returns one.\n@param value {String | Array | Object} Value.\n@return {[Object]} Value."),
			"example.Many": declarationWithDocs("Many", "Returns values.\n@param value {[Int | Float]} Values.\n@return {Any[]} Values."),
		},
		Assertions: map[string]assertionDescriptor{},
	}
	registered := []registeredSignature{
		{QualifiedName: "nested::item", Namespace: "nested", Name: "item", Symbol: "example.Many", Variadic: true},
		{QualifiedName: "root", Name: "root", Symbol: "example.One", Arity: 1},
		{QualifiedName: "nested::item", Namespace: "nested", Name: "item", Symbol: "example.Zero", Arity: 0},
	}

	reference, err := buildReference("1.2.3", registered, catalog)
	if err != nil {
		t.Fatal(err)
	}

	if reference.ID != moduleID || reference.Version != "1.2.3" || reference.SchemaVersion != api.SchemaVersion {
		t.Fatalf("identity = %#v", reference)
	}

	if got, want := namespaceNames(reference.Namespaces), []string{"", "nested"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("namespaces = %v, want %v", got, want)
	}

	signatures := reference.Namespaces[1].Functions[0].Signatures
	if len(signatures) != 2 || signatures[0].Variadic || !signatures[1].Variadic {
		t.Fatalf("signature ordering = %#v, want fixed before variadic", signatures)
	}

	rootSignature := reference.Namespaces[0].Functions[0].Signatures[0]
	if rootSignature.Parameters[0].Type.Kind != api.TypeKindUnion || len(rootSignature.Parameters[0].Type.Types) != 3 {
		t.Fatalf("root parameter type = %#v, want three-member union", rootSignature.Parameters[0].Type)
	}

	if rootSignature.Return.Type.Kind != api.TypeKindList || rootSignature.Return.Type.Element.Name != "Object" {
		t.Fatalf("root return type = %#v, want list of Object", rootSignature.Return.Type)
	}

	if signatures[1].Parameters[0].Type.Kind != api.TypeKindList || signatures[1].Parameters[0].Type.Element.Kind != api.TypeKindUnion {
		t.Fatalf("variadic parameter type = %#v, want list of numeric union", signatures[1].Parameters[0].Type)
	}
}

func TestValidateAssertionAritiesRequiresContiguousRange(t *testing.T) {
	descriptor := assertionDescriptor{
		Declaration: declarationWithDocs("Equal", "Tests equality.\n@param actual {Any} Actual.\n@param expected {Any} Expected.\n@param message {String} Message.\n@return {Boolean} Result."),
		Min:         2,
		Max:         3,
	}

	if err := validateAssertionArities("t::eq", descriptor, []int{3, 2}); err != nil {
		t.Fatalf("validate contiguous overloads: %v", err)
	}

	if err := validateAssertionArities("t::eq", descriptor, []int{2}); err == nil {
		t.Fatal("expected missing overload to fail")
	}
}

func TestBuildReferenceReportsUnresolvedDeclarationContext(t *testing.T) {
	registered := []registeredSignature{{
		QualifiedName: "PUBLIC_NAME",
		Symbol:        "example.com/stdlib.GoName",
		PackagePath:   "example.com/stdlib",
		File:          "stdlib/value.go",
		Line:          42,
	}}

	_, err := buildReference("1.2.3", registered, &sourceCatalog{
		Declarations: map[string]*sourceDeclaration{},
		Assertions:   map[string]assertionDescriptor{},
	})
	if err == nil {
		t.Fatal("expected unresolved declaration to fail")
	}

	for _, expected := range []string{"stdlib/value.go:42", "package example.com/stdlib", "Ferret PUBLIC_NAME", "Go declaration example.com/stdlib.GoName"} {
		if !strings.Contains(err.Error(), expected) {
			t.Fatalf("error = %q, want diagnostic context %q", err, expected)
		}
	}
}

func declarationWithDocs(name, documentation string) *sourceDeclaration {
	return &sourceDeclaration{
		Name:          name,
		PackagePath:   "example",
		File:          "example.go",
		Line:          1,
		Documentation: documentation,
	}
}

func namespaceNames(namespaces []api.Namespace) []string {
	result := make([]string, 0, len(namespaces))
	for _, namespace := range namespaces {
		result = append(result, namespace.Name)
	}

	return result
}
