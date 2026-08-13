package compiler

import (
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func buildAnalysis(
	src *source.Source,
	snapshot internal.SemanticSnapshot,
	errors *parserd.ErrorHandler,
	syntaxTokens []SyntaxToken,
) *Analysis {
	data := analysisData{sourceLength: sourceLength(src)}
	data.syntaxTokens = append([]SyntaxToken(nil), syntaxTokens...)
	data.symbols = make([]Symbol, len(snapshot.Symbols))
	data.symbolMetadata = make([]analysisSymbolMetadata, len(snapshot.Symbols))

	for i, symbol := range snapshot.Symbols {
		data.symbols[i] = Symbol{
			ID:              SymbolID(symbol.ID),
			Name:            symbol.Name,
			Kind:            publicSymbolKind(symbol.Kind),
			Type:            publicValueType(symbol.Type),
			DeclarationSpan: symbol.Declaration,
			SelectionSpan:   symbol.Selection,
			Mutable:         symbol.Mutable,
			HasDeclaration:  symbol.HasDeclaration,
		}
		data.symbolMetadata[i] = analysisSymbolMetadata{
			Function:   SymbolID(symbol.Function),
			Scope:      analysisScopeID(symbol.Scope),
			Activation: symbol.Activation,
		}
	}

	data.references = make([]Reference, len(snapshot.References))
	for i, reference := range snapshot.References {
		data.references[i] = Reference{Symbol: SymbolID(reference.Symbol), Span: reference.Span}
	}

	data.calls = make([]Call, len(snapshot.Calls))
	for i, call := range snapshot.Calls {
		data.calls[i] = Call{
			Target:        SymbolID(call.Target),
			Name:          call.Name,
			Identity:      call.Identity,
			Kind:          publicCallKind(call.Kind),
			ResultType:    publicValueType(call.Type),
			Span:          call.Span,
			CalleeSpan:    call.CalleeSpan,
			ArgumentSpans: append([]source.Span(nil), call.ArgumentSpans...),
		}
	}

	data.typeFacts = make([]TypeFact, len(snapshot.TypeFacts))
	for i, fact := range snapshot.TypeFacts {
		data.typeFacts[i] = TypeFact{Span: fact.Span, Type: publicValueType(fact.Type)}
	}

	data.scopes = make([]analysisScope, len(snapshot.Scopes))
	for i, scope := range snapshot.Scopes {
		data.scopes[i] = analysisScope{
			ID:     analysisScopeID(scope.ID),
			Parent: analysisScopeID(scope.Parent),
			Span:   scope.Span,
			Depth:  scope.Depth,
		}
	}

	data.diagnostics = diagnosticsFromHandler(src, errors)

	return newAnalysis(data)
}

func diagnosticsFromHandler(src *source.Source, handler *parserd.ErrorHandler) []*diagnostics.Diagnostic {
	if handler == nil || handler.Errors() == nil {
		return nil
	}

	out := make([]*diagnostics.Diagnostic, 0, handler.Errors().Size())
	byteOffsets := analysisByteOffsets(src)

	for _, diagnostic := range handler.Errors().Errors() {
		cloned := cloneDiagnostic(diagnostic)
		if cloned != nil {
			for i := range cloned.Spans {
				cloned.Spans[i].Span = analysisByteSpan(byteOffsets, cloned.Spans[i].Span)
			}
		}

		out = append(out, cloned)
	}

	return out
}

func sourceLength(src *source.Source) int {
	if src == nil {
		return 0
	}

	return len(src.Content())
}

func publicSymbolKind(kind internal.SemanticSymbolKind) SymbolKind {
	switch kind {
	case internal.SemanticSymbolBinding:
		return SymbolKindBinding
	case internal.SemanticSymbolFunctionParameter:
		return SymbolKindFunctionParameter
	case internal.SemanticSymbolUserFunction:
		return SymbolKindUDF
	case internal.SemanticSymbolBindParameter:
		return SymbolKindBindParameter
	case internal.SemanticSymbolLoopBinding:
		return SymbolKindLoopBinding
	case internal.SemanticSymbolMatchBinding:
		return SymbolKindMatchBinding
	case internal.SemanticSymbolCollectBinding:
		return SymbolKindCollectBinding
	case internal.SemanticSymbolNamespaceAlias:
		return SymbolKindNamespaceAlias
	default:
		return 0
	}
}

func publicCallKind(kind internal.SemanticCallKind) CallKind {
	switch kind {
	case internal.SemanticCallUserFunction:
		return CallKindUDF
	case internal.SemanticCallBuiltin:
		return CallKindBuiltin
	case internal.SemanticCallHost:
		return CallKindHost
	default:
		return 0
	}
}

func publicValueType(typ core.ValueType) ValueType {
	switch typ {
	case core.TypeUnknown:
		return ValueTypeUnknown
	case core.TypeNone:
		return ValueTypeNone
	case core.TypeInt:
		return ValueTypeInteger
	case core.TypeFloat:
		return ValueTypeFloat
	case core.TypeDuration:
		return ValueTypeDuration
	case core.TypeString:
		return ValueTypeString
	case core.TypeBool:
		return ValueTypeBoolean
	case core.TypeArray:
		return ValueTypeArray
	case core.TypeObject:
		return ValueTypeObject
	case core.TypeList:
		return ValueTypeList
	case core.TypeMap:
		return ValueTypeMap
	case core.TypeAny:
		return ValueTypeAny
	default:
		return ValueTypeUnknown
	}
}
