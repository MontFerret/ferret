package compiler

import (
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	// SymbolID identifies a semantic symbol within one Analysis result.
	// Zero is reserved as an invalid ID. IDs are not stable across analyses.
	SymbolID uint32

	// SymbolKind classifies a source-visible semantic symbol.
	SymbolKind uint8

	// CallKind classifies a resolved function call.
	CallKind uint8

	// ValueType is a type fact established by the existing compiler pipeline.
	// Unknown means the compiler has not established a type; Any is a known
	// dynamic value.
	ValueType uint8

	// SyntaxTokenKind classifies a parser-independent source token.
	SyntaxTokenKind uint8

	// Symbol describes a source-visible declaration or analysis-local bind parameter.
	Symbol struct {
		Name            string
		DeclarationSpan source.Span
		SelectionSpan   source.Span
		ID              SymbolID
		Kind            SymbolKind
		Type            ValueType
		Mutable         bool
		HasDeclaration  bool
	}

	// Reference is a resolved source reference to a symbol.
	Reference struct {
		Symbol SymbolID
		Span   source.Span
	}

	// Call describes a compiler-resolved call site.
	Call struct {
		Name          string
		Identity      string
		ArgumentSpans []source.Span
		Span          source.Span
		CalleeSpan    source.Span
		Target        SymbolID
		Kind          CallKind
		ResultType    ValueType
	}

	// TypeFact describes the compiler-established type for an expression span.
	TypeFact struct {
		Span source.Span
		Type ValueType
	}

	// SyntaxToken describes one non-whitespace token from the analyzed source.
	// Span uses zero-based, half-open UTF-8 byte offsets.
	SyntaxToken struct {
		Kind SyntaxTokenKind
		// Word is the canonical language-word identity, or SyntaxWordUnknown.
		Word SyntaxWord
		Span source.Span
	}

	analysisScopeID uint32

	analysisScope struct {
		ID     analysisScopeID
		Parent analysisScopeID
		Span   source.Span
		Depth  int
	}

	analysisSymbolMetadata struct {
		Function   SymbolID
		Scope      analysisScopeID
		Activation int
	}

	analysisData struct {
		symbols        []Symbol
		references     []Reference
		calls          []Call
		diagnostics    []*diagnostics.Diagnostic
		typeFacts      []TypeFact
		syntaxTokens   []SyntaxToken
		scopes         []analysisScope
		symbolMetadata []analysisSymbolMetadata
		sourceLength   int
	}
)

const (
	SyntaxTokenKindUnknown SyntaxTokenKind = iota
	SyntaxTokenKindIdentifier
	SyntaxTokenKindNamespace
	SyntaxTokenKindKeyword
	SyntaxTokenKindString
	SyntaxTokenKindNumber
	SyntaxTokenKindDuration
	SyntaxTokenKindComment
	SyntaxTokenKindOperator
	SyntaxTokenKindPunctuation
)

// InvalidSymbolID represents the absence of a semantic symbol target.
const InvalidSymbolID SymbolID = 0

const (
	SymbolKindBinding SymbolKind = iota + 1
	SymbolKindFunctionParameter
	SymbolKindUDF
	SymbolKindBindParameter
	SymbolKindLoopBinding
	SymbolKindMatchBinding
	SymbolKindCollectBinding
	SymbolKindNamespaceAlias
)

const (
	CallKindUDF CallKind = iota + 1
	CallKindBuiltin
	CallKindHost
)

const (
	ValueTypeUnknown ValueType = iota
	ValueTypeAny
	ValueTypeNone
	ValueTypeInteger
	ValueTypeFloat
	ValueTypeDuration
	ValueTypeBoolean
	ValueTypeString
	ValueTypeArray
	ValueTypeObject
	ValueTypeList
	ValueTypeMap
)
