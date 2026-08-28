package internal

import (
	"sort"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type (
	SemanticSymbolID uint32

	SemanticScopeID uint32

	SemanticSymbolKind uint8

	SemanticCallKind uint8

	SemanticSymbol struct {
		Name           string
		Declaration    source.Span
		Selection      source.Span
		Type           core.ValueType
		Activation     int
		ID             SemanticSymbolID
		Function       SemanticSymbolID
		Scope          SemanticScopeID
		Kind           SemanticSymbolKind
		Mutable        bool
		HasDeclaration bool
	}

	SemanticReference struct {
		Symbol SemanticSymbolID
		Span   source.Span
	}

	SemanticCall struct {
		Name          string
		Identity      string
		ArgumentSpans []source.Span
		Span          source.Span
		CalleeSpan    source.Span
		Type          core.ValueType
		Target        SemanticSymbolID
		Kind          SemanticCallKind
	}

	SemanticTypeFact struct {
		Span source.Span
		Type core.ValueType
	}

	SemanticScope struct {
		ID     SemanticScopeID
		Parent SemanticScopeID
		Span   source.Span
		Depth  int
	}

	SemanticSnapshot struct {
		Symbols    []SemanticSymbol
		References []SemanticReference
		Calls      []SemanticCall
		TypeFacts  []SemanticTypeFact
		Scopes     []SemanticScope
	}

	// SemanticRecorder captures immutable-source semantic decisions while the
	// normal compiler frontend resolves and lowers the program.
	SemanticRecorder struct {
		symbols          []SemanticSymbol
		references       []SemanticReference
		calls            []SemanticCall
		typeFacts        []SemanticTypeFact
		scopes           []SemanticScope
		scopeStack       []SemanticScopeID
		symbolByBinding  map[core.BindingID]SemanticSymbolID
		symbolByFunction map[*core.UDFInfo]SemanticSymbolID
		scopeByUDF       map[*core.UDFScope]SemanticScopeID
		functionByScope  map[SemanticScopeID]SemanticSymbolID
		bindParameters   map[string]SemanticSymbolID
		aliases          map[string]SemanticSymbolID
		referenceKeys    map[SemanticReference]struct{}
		callBySpan       map[source.Span]int
		typeFactBySpan   map[source.Span]int
		byteOffsets      []int
	}
)

const (
	SemanticSymbolBinding SemanticSymbolKind = iota + 1
	SemanticSymbolFunctionParameter
	SemanticSymbolUserFunction
	SemanticSymbolBindParameter
	SemanticSymbolLoopBinding
	SemanticSymbolMatchBinding
	SemanticSymbolCollectBinding
	SemanticSymbolNamespaceAlias
)

const (
	SemanticCallUserFunction SemanticCallKind = iota + 1
	SemanticCallBuiltin
	SemanticCallHost
)

func NewSemanticRecorder(src source.Source) *SemanticRecorder {
	end := 0
	if !src.Empty() {
		end = src.Length()
	}

	r := &SemanticRecorder{
		symbolByBinding:  make(map[core.BindingID]SemanticSymbolID),
		symbolByFunction: make(map[*core.UDFInfo]SemanticSymbolID),
		scopeByUDF:       make(map[*core.UDFScope]SemanticScopeID),
		functionByScope:  make(map[SemanticScopeID]SemanticSymbolID),
		bindParameters:   make(map[string]SemanticSymbolID),
		aliases:          make(map[string]SemanticSymbolID),
		referenceKeys:    make(map[SemanticReference]struct{}),
		callBySpan:       make(map[source.Span]int),
		typeFactBySpan:   make(map[source.Span]int),
		byteOffsets:      semanticByteOffsets(src),
	}

	root := r.newScope(source.Span{Start: 0, End: end}, 0)
	r.scopeStack = append(r.scopeStack, root)

	return r
}

func (r *SemanticRecorder) Snapshot() SemanticSnapshot {
	if r == nil {
		return SemanticSnapshot{}
	}

	calls := make([]SemanticCall, len(r.calls))
	for i := range r.calls {
		calls[i] = r.calls[i]
		calls[i].ArgumentSpans = append([]source.Span(nil), r.calls[i].ArgumentSpans...)
	}

	return SemanticSnapshot{
		Symbols:    append([]SemanticSymbol(nil), r.symbols...),
		References: append([]SemanticReference(nil), r.references...),
		Calls:      calls,
		TypeFacts:  append([]SemanticTypeFact(nil), r.typeFacts...),
		Scopes:     append([]SemanticScope(nil), r.scopes...),
	}
}

func (r *SemanticRecorder) RegisterGlobalUDFScope(scope *core.UDFScope) {
	if r == nil || scope == nil || len(r.scopeStack) == 0 {
		return
	}

	r.scopeByUDF[scope] = r.scopeStack[0]
}

func (r *SemanticRecorder) RecordUserFunction(fn *core.UDFInfo) SemanticSymbolID {
	if r == nil || fn == nil || fn.Decl == nil || fn.Scope == nil || fn.BodyScope == nil {
		return 0
	}

	if id, ok := r.symbolByFunction[fn]; ok {
		return id
	}

	ownerScope, ok := r.scopeByUDF[fn.Scope]
	if !ok {
		return 0
	}

	decl := r.normalizeSpan(ruleSpan(fn.Decl))
	selection := decl

	if name := fn.Decl.FunctionName(); name != nil {
		selection = r.normalizeSpan(ruleSpan(name))
	}

	id := r.addSymbol(SemanticSymbol{
		Name:           fn.DisplayName,
		Kind:           SemanticSymbolUserFunction,
		Type:           core.TypeAny,
		Declaration:    decl,
		Selection:      selection,
		Scope:          ownerScope,
		Activation:     r.scope(ownerScope).Span.Start,
		HasDeclaration: true,
	})
	r.symbolByFunction[fn] = id

	bodyRuleSpan := ruleSpan(fn.Decl)

	if body := fn.Decl.FunctionBody(); body != nil {
		bodyRuleSpan = ruleSpan(body)
	}

	bodySpan := r.normalizeSpan(bodyRuleSpan)

	functionScope := r.newScope(bodySpan, ownerScope)
	r.scopeByUDF[fn.BodyScope] = functionScope
	r.functionByScope[functionScope] = id

	for _, param := range fn.Params {
		paramCtx := functionParameterByID(fn.Decl, param.ID)
		if paramCtx == nil {
			continue
		}

		span := ruleSpan(paramCtx)
		r.RecordBinding(param.ID, param.Name, SemanticSymbolFunctionParameter, span, span, false, core.TypeAny, bodyRuleSpan.Start, id, functionScope)
	}

	return id
}

func (r *SemanticRecorder) PushUDFScope(fn *core.UDFInfo) bool {
	if r == nil || fn == nil {
		return false
	}

	scope, ok := r.scopeByUDF[fn.BodyScope]
	if !ok {
		return false
	}

	r.scopeStack = append(r.scopeStack, scope)

	return true
}

func (r *SemanticRecorder) EnterScope(span source.Span) SemanticScopeID {
	if r == nil {
		return 0
	}

	id := r.newScope(r.normalizeSpan(span), r.CurrentScope())
	r.scopeStack = append(r.scopeStack, id)

	return id
}

func (r *SemanticRecorder) ExitScope() {
	if r == nil || len(r.scopeStack) <= 1 {
		return
	}

	r.scopeStack = r.scopeStack[:len(r.scopeStack)-1]
}

func (r *SemanticRecorder) CurrentScope() SemanticScopeID {
	if r == nil || len(r.scopeStack) == 0 {
		return 0
	}

	return r.scopeStack[len(r.scopeStack)-1]
}

func (r *SemanticRecorder) RecordBinding(
	id core.BindingID,
	name string,
	kind SemanticSymbolKind,
	declaration source.Span,
	selection source.Span,
	mutable bool,
	typ core.ValueType,
	activation int,
	function SemanticSymbolID,
	scope SemanticScopeID,
) SemanticSymbolID {
	if r == nil || id == core.InvalidBindingID || name == "" || name == core.IgnorePseudoVariable || name == core.PseudoVariable {
		return 0
	}

	if existing, ok := r.symbolByBinding[id]; ok {
		r.updateSymbolType(existing, typ)

		return existing
	}

	if scope == 0 {
		scope = r.CurrentScope()
	}

	declaration = r.normalizeSpan(declaration)
	selection = r.normalizeSpan(selection)
	activation = r.normalizeOffset(activation)

	symbol := SemanticSymbol{
		Function:       function,
		Name:           name,
		Kind:           kind,
		Type:           typ,
		Declaration:    declaration,
		Selection:      selection,
		Scope:          scope,
		Activation:     activation,
		Mutable:        mutable,
		HasDeclaration: true,
	}

	symbolID := r.addSymbol(symbol)
	r.symbolByBinding[id] = symbolID

	return symbolID
}

func (r *SemanticRecorder) UpdateBindingType(id core.BindingID, typ core.ValueType) {
	if r == nil {
		return
	}

	symbolID, ok := r.symbolByBinding[id]
	if !ok {
		return
	}

	r.updateSymbolType(symbolID, typ)
}

func (r *SemanticRecorder) RecordBindingReference(id core.BindingID, span source.Span) {
	if r == nil || id == core.InvalidBindingID {
		return
	}

	span = r.normalizeSpan(span)
	if !validSpan(span) {
		return
	}

	symbol, ok := r.symbolByBinding[id]
	if !ok {
		return
	}

	r.recordReference(symbol, span)
}

func (r *SemanticRecorder) RecordUserFunctionReference(fn *core.UDFInfo, span source.Span) {
	if r == nil || fn == nil {
		return
	}

	span = r.normalizeSpan(span)
	if !validSpan(span) {
		return
	}

	symbol, ok := r.symbolByFunction[fn]
	if !ok {
		return
	}

	r.recordReference(symbol, span)
}

func (r *SemanticRecorder) RecordBindParameter(name string, span source.Span) {
	if r == nil || name == "" {
		return
	}

	span = r.normalizeSpan(span)
	if !validSpan(span) {
		return
	}

	symbol, ok := r.bindParameters[name]
	if !ok {
		symbol = r.addSymbol(SemanticSymbol{
			Name:       name,
			Kind:       SemanticSymbolBindParameter,
			Type:       core.TypeAny,
			Scope:      r.scopeStack[0],
			Activation: r.scope(r.scopeStack[0]).Span.Start,
		})
		r.bindParameters[name] = symbol
	}

	r.recordReference(symbol, span)
}

func (r *SemanticRecorder) RecordNamespaceAlias(name string, declaration, selection source.Span) {
	if r == nil || name == "" {
		return
	}

	declaration = r.normalizeSpan(declaration)
	selection = r.normalizeSpan(selection)
	if !validSpan(selection) {
		return
	}

	symbol := r.addSymbol(SemanticSymbol{
		Name:           name,
		Kind:           SemanticSymbolNamespaceAlias,
		Type:           core.TypeAny,
		Declaration:    declaration,
		Selection:      selection,
		Scope:          r.scopeStack[0],
		Activation:     r.scope(r.scopeStack[0]).Span.Start,
		HasDeclaration: true,
	})
	r.aliases[name] = symbol
}

func (r *SemanticRecorder) RecordNamespaceAliasReference(name string, span source.Span) {
	if r == nil || name == "" {
		return
	}

	span = r.normalizeSpan(span)
	if !validSpan(span) {
		return
	}

	symbol, ok := r.aliases[name]
	if !ok {
		return
	}

	r.recordReference(symbol, span)
}

func (r *SemanticRecorder) RecordCall(call SemanticCall) {
	if r == nil {
		return
	}

	call.Span = r.normalizeSpan(call.Span)
	call.CalleeSpan = r.normalizeSpan(call.CalleeSpan)
	call.ArgumentSpans = append([]source.Span(nil), call.ArgumentSpans...)

	for i := range call.ArgumentSpans {
		call.ArgumentSpans[i] = r.normalizeSpan(call.ArgumentSpans[i])
	}

	if !validSpan(call.Span) || !validSpan(call.CalleeSpan) {
		return
	}

	if idx, ok := r.callBySpan[call.Span]; ok {
		r.calls[idx] = call

		return
	}

	r.callBySpan[call.Span] = len(r.calls)
	r.calls = append(r.calls, call)
}

func (r *SemanticRecorder) RecordExpressionType(ctx antlr.ParserRuleContext, typ core.ValueType) {
	if r == nil || ctx == nil {
		return
	}

	span := r.normalizeSpan(ruleSpan(ctx))
	if !validSpan(span) {
		return
	}

	if idx, ok := r.typeFactBySpan[span]; ok {
		r.typeFacts[idx].Type = core.JoinValueTypes(r.typeFacts[idx].Type, typ)

		return
	}

	r.typeFactBySpan[span] = len(r.typeFacts)
	r.typeFacts = append(r.typeFacts, SemanticTypeFact{Span: span, Type: typ})
}

func (r *SemanticRecorder) UserFunctionSymbol(fn *core.UDFInfo) SemanticSymbolID {
	if r == nil || fn == nil {
		return 0
	}

	return r.symbolByFunction[fn]
}

func (r *SemanticRecorder) CurrentFunctionSymbol() SemanticSymbolID {
	if r == nil {
		return 0
	}

	current := r.CurrentScope()
	for current != 0 {
		if function := r.functionByScope[current]; function != 0 {
			return function
		}

		current = r.scope(current).Parent
	}

	return 0
}

func (r *SemanticRecorder) Sort() {
	if r == nil {
		return
	}

	sort.SliceStable(r.references, func(i, j int) bool {
		if r.references[i].Span.Start != r.references[j].Span.Start {
			return r.references[i].Span.Start < r.references[j].Span.Start
		}

		return r.references[i].Span.End < r.references[j].Span.End
	})
	sort.SliceStable(r.calls, func(i, j int) bool {
		if r.calls[i].Span.Start != r.calls[j].Span.Start {
			return r.calls[i].Span.Start < r.calls[j].Span.Start
		}

		return r.calls[i].Span.End < r.calls[j].Span.End
	})
	sort.SliceStable(r.typeFacts, func(i, j int) bool {
		if r.typeFacts[i].Span.Start != r.typeFacts[j].Span.Start {
			return r.typeFacts[i].Span.Start < r.typeFacts[j].Span.Start
		}

		return r.typeFacts[i].Span.End < r.typeFacts[j].Span.End
	})
}

func (r *SemanticRecorder) addSymbol(symbol SemanticSymbol) SemanticSymbolID {
	symbol.ID = SemanticSymbolID(len(r.symbols) + 1)
	r.symbols = append(r.symbols, symbol)

	return symbol.ID
}

func (r *SemanticRecorder) recordReference(symbol SemanticSymbolID, span source.Span) {
	reference := SemanticReference{Symbol: symbol, Span: span}
	if _, exists := r.referenceKeys[reference]; exists {
		return
	}

	r.referenceKeys[reference] = struct{}{}
	r.references = append(r.references, reference)
}

func (r *SemanticRecorder) updateSymbolType(id SemanticSymbolID, typ core.ValueType) {
	idx := int(id) - 1
	if idx < 0 || idx >= len(r.symbols) {
		return
	}

	r.symbols[idx].Type = core.JoinValueTypes(r.symbols[idx].Type, typ)
}

func (r *SemanticRecorder) newScope(span source.Span, parent SemanticScopeID) SemanticScopeID {
	depth := 0
	if parent != 0 {
		depth = r.scope(parent).Depth + 1
	}

	id := SemanticScopeID(len(r.scopes) + 1)
	r.scopes = append(r.scopes, SemanticScope{ID: id, Parent: parent, Span: span, Depth: depth})

	return id
}

func (r *SemanticRecorder) scope(id SemanticScopeID) SemanticScope {
	idx := int(id) - 1
	if idx < 0 || idx >= len(r.scopes) {
		return SemanticScope{}
	}

	return r.scopes[idx]
}

func (r *SemanticRecorder) normalizeOffset(offset int) int {
	if r == nil || offset < 0 || offset >= len(r.byteOffsets) {
		return offset
	}

	return r.byteOffsets[offset]
}

func (r *SemanticRecorder) normalizeSpan(span source.Span) source.Span {
	if r == nil || span.Start < 0 {
		return span
	}

	return source.Span{
		Start: r.normalizeOffset(span.Start),
		End:   r.normalizeOffset(span.End),
	}
}
