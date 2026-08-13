package compiler

import (
	"sort"

	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

// Analysis is an immutable semantic snapshot produced by Compiler.Analyze.
// Its accessors return defensive copies.
type Analysis struct {
	data analysisData
}

func newAnalysis(data analysisData) *Analysis {
	return &Analysis{data: data}
}

// Symbols returns all source-visible symbols in deterministic analysis-local ID order.
func (a *Analysis) Symbols() []Symbol {
	if a == nil {
		return nil
	}

	return append([]Symbol(nil), a.data.symbols...)
}

// Symbol returns the symbol with the analysis-local ID.
func (a *Analysis) Symbol(id SymbolID) (Symbol, bool) {
	if a == nil || id == 0 || int(id) > len(a.data.symbols) {
		return Symbol{}, false
	}

	return a.data.symbols[int(id)-1], true
}

// References returns all resolved references in source order.
func (a *Analysis) References() []Reference {
	if a == nil {
		return nil
	}

	return append([]Reference(nil), a.data.references...)
}

// ReferencesTo returns references resolved to symbol in source order.
func (a *Analysis) ReferencesTo(symbol SymbolID) []Reference {
	if a == nil || symbol == 0 {
		return nil
	}

	out := make([]Reference, 0)
	for _, reference := range a.data.references {
		if reference.Symbol == symbol {
			out = append(out, reference)
		}
	}

	return out
}

// Calls returns all resolved calls in source order.
func (a *Analysis) Calls() []Call {
	if a == nil {
		return nil
	}

	out := make([]Call, len(a.data.calls))
	for i := range a.data.calls {
		out[i] = cloneCall(a.data.calls[i])
	}

	return out
}

// Diagnostics returns deep copies of all source diagnostics in compiler order.
func (a *Analysis) Diagnostics() []*diagnostics.Diagnostic {
	if a == nil {
		return nil
	}

	return cloneDiagnostics(a.data.diagnostics)
}

// TypeFacts returns all compiler-established expression facts in source order.
func (a *Analysis) TypeFacts() []TypeFact {
	if a == nil {
		return nil
	}

	return append([]TypeFact(nil), a.data.typeFacts...)
}

// FunctionParameters returns the declared parameters for a UDF symbol.
func (a *Analysis) FunctionParameters(function SymbolID) []Symbol {
	if a == nil || function == 0 {
		return nil
	}

	out := make([]Symbol, 0)
	for i, symbol := range a.data.symbols {
		if symbol.Kind == SymbolKindFunctionParameter && a.data.symbolMetadata[i].Function == function {
			out = append(out, symbol)
		}
	}

	return out
}

// SymbolAt returns the symbol whose declaration selection or reference contains offset.
// The offset is a zero-based UTF-8 byte offset into the analyzed source.
func (a *Analysis) SymbolAt(offset int) (Symbol, bool) {
	if !a.validOffset(offset) {
		return Symbol{}, false
	}

	best := source.Span{}
	var found SymbolID

	for _, symbol := range a.data.symbols {
		if spanContains(symbol.SelectionSpan, offset) && narrowerSpan(symbol.SelectionSpan, best) {
			found = symbol.ID
			best = symbol.SelectionSpan
		}
	}

	for _, reference := range a.data.references {
		if spanContains(reference.Span, offset) && narrowerSpan(reference.Span, best) {
			found = reference.Symbol
			best = reference.Span
		}
	}

	return a.Symbol(found)
}

// ReferenceAt returns the narrowest semantic reference containing offset.
// Declaration selection spans are not references. The offset is a zero-based
// UTF-8 byte offset into the analyzed source.
func (a *Analysis) ReferenceAt(offset int) (Reference, bool) {
	if !a.validOffset(offset) {
		return Reference{}, false
	}

	best := source.Span{}
	var found Reference
	ok := false

	for _, reference := range a.data.references {
		if spanContains(reference.Span, offset) && narrowerSpan(reference.Span, best) {
			found = reference
			best = reference.Span
			ok = true
		}
	}

	return found, ok
}

// CallAt returns the narrowest call whose complete Call.Span contains offset.
// The offset may be inside the callee or any argument. Consumers that need to
// detect the callee specifically should test the offset against Call.CalleeSpan.
// The offset is a zero-based UTF-8 byte offset into the analyzed source.
func (a *Analysis) CallAt(offset int) (Call, bool) {
	if !a.validOffset(offset) {
		return Call{}, false
	}

	best := source.Span{}
	var found *Call

	for i := range a.data.calls {
		call := &a.data.calls[i]
		if spanContains(call.Span, offset) && narrowerSpan(call.Span, best) {
			found = call
			best = call.Span
		}
	}

	if found == nil {
		return Call{}, false
	}

	return cloneCall(*found), true
}

// TypeAt returns the narrowest expression type fact containing offset.
// The offset is a zero-based UTF-8 byte offset into the analyzed source.
func (a *Analysis) TypeAt(offset int) (TypeFact, bool) {
	if !a.validOffset(offset) {
		return TypeFact{}, false
	}

	best := source.Span{}
	var found TypeFact
	ok := false

	for _, fact := range a.data.typeFacts {
		if spanContains(fact.Span, offset) && narrowerSpan(fact.Span, best) {
			found = fact
			best = fact.Span
			ok = true
		}
	}

	return found, ok
}

// VisibleSymbols returns source-visible symbols active at offset after lexical shadowing.
// The offset is a zero-based UTF-8 byte offset into the analyzed source.
func (a *Analysis) VisibleSymbols(offset int) []Symbol {
	if !a.validOffset(offset) {
		return nil
	}

	type candidate struct {
		symbol Symbol
		depth  int
	}

	candidates := make([]candidate, 0)
	for i, symbol := range a.data.symbols {
		metadata := a.data.symbolMetadata[i]
		if offset < metadata.Activation || !a.scopeContains(metadata.Scope, offset) {
			continue
		}

		candidates = append(candidates, candidate{symbol: symbol, depth: a.scope(metadata.Scope).Depth})
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].depth != candidates[j].depth {
			return candidates[i].depth > candidates[j].depth
		}

		return candidates[i].symbol.ID > candidates[j].symbol.ID
	})

	seen := make(map[string]struct{}, len(candidates))
	out := make([]Symbol, 0, len(candidates))
	for _, candidate := range candidates {
		key := semanticNamespace(candidate.symbol.Kind) + "\x00" + candidate.symbol.Name
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		out = append(out, candidate.symbol)
	}

	sort.Slice(out, func(i, j int) bool {
		return out[i].ID < out[j].ID
	})

	return out
}

func (a *Analysis) validOffset(offset int) bool {
	return a != nil && offset >= 0 && offset <= a.data.sourceLength
}

func (a *Analysis) scope(id analysisScopeID) analysisScope {
	if a == nil || id == 0 || int(id) > len(a.data.scopes) {
		return analysisScope{}
	}

	return a.data.scopes[int(id)-1]
}

func (a *Analysis) scopeContains(id analysisScopeID, offset int) bool {
	scope := a.scope(id)
	if scope.ID == 0 {
		return false
	}

	if scope.Parent == 0 && offset == a.data.sourceLength {
		return true
	}

	return spanContains(scope.Span, offset)
}
