package compiler

import (
	"errors"
	"reflect"
	"strings"
	"sync"
	"testing"

	diagnosticspkg "github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestAnalyzeRecordsPublicSemanticsAndExactSpans(t *testing.T) {
	query := `USE WEB::HTML AS html
LET outer = 1
FUNC add(p) => outer + p
LET value = add(@input)
RETURN [LENGTH([value]), html::PARSE("<p/>")]`

	analysis, err := New().Analyze(source.New("analysis.fql", query))
	if err != nil {
		t.Fatal(err)
	}

	outer := requireAnalysisSymbol(t, analysis, "outer", SymbolKindBinding)
	assertAnalysisSpanText(t, query, outer.DeclarationSpan, "LET outer = 1")
	assertAnalysisSpanText(t, query, outer.SelectionSpan, "outer")
	if outer.Type != ValueTypeInteger {
		t.Fatalf("outer type = %v, want integer", outer.Type)
	}

	value := requireAnalysisSymbol(t, analysis, "value", SymbolKindBinding)
	assertAnalysisSpanText(t, query, value.DeclarationSpan, "LET value = add(@input)")
	assertAnalysisSpanText(t, query, value.SelectionSpan, "value")

	add := requireAnalysisSymbol(t, analysis, "add", SymbolKindUDF)
	assertAnalysisSpanText(t, query, add.DeclarationSpan, "FUNC add(p) => outer + p")
	assertAnalysisSpanText(t, query, add.SelectionSpan, "add")

	params := analysis.FunctionParameters(add.ID)
	if len(params) != 1 || params[0].Name != "p" || params[0].Kind != SymbolKindFunctionParameter {
		t.Fatalf("parameters = %+v, want p", params)
	}
	assertAnalysisSpanText(t, query, params[0].SelectionSpan, "p")

	bind := requireAnalysisSymbol(t, analysis, "input", SymbolKindBindParameter)
	if bind.HasDeclaration || bind.DeclarationSpan != (source.Span{}) || bind.SelectionSpan != (source.Span{}) {
		t.Fatalf("bind parameter has fabricated declaration: %+v", bind)
	}
	bindRefs := analysis.ReferencesTo(bind.ID)
	if len(bindRefs) != 1 {
		t.Fatalf("bind parameter references = %+v", bindRefs)
	}
	assertAnalysisSpanText(t, query, bindRefs[0].Span, "@input")

	alias := requireAnalysisSymbol(t, analysis, "html", SymbolKindNamespaceAlias)
	assertAnalysisSpanText(t, query, alias.SelectionSpan, "html")
	aliasRefs := analysis.ReferencesTo(alias.ID)
	if len(aliasRefs) != 1 {
		t.Fatalf("alias references = %+v", aliasRefs)
	}
	assertAnalysisSpanText(t, query, aliasRefs[0].Span, "html")

	calls := analysis.Calls()
	if len(calls) != 3 {
		t.Fatalf("calls = %+v, want 3", calls)
	}

	udfCall := requireAnalysisCall(t, calls, "add")
	if udfCall.Kind != CallKindUDF || udfCall.Target != add.ID || udfCall.Identity != "" || udfCall.ResultType != ValueTypeAny {
		t.Fatalf("UDF call = %+v", udfCall)
	}
	assertAnalysisSpanText(t, query, udfCall.Span, "add(@input)")
	assertAnalysisSpanText(t, query, udfCall.CalleeSpan, "add")
	if len(udfCall.ArgumentSpans) != 1 {
		t.Fatalf("UDF arguments = %+v", udfCall.ArgumentSpans)
	}
	assertAnalysisSpanText(t, query, udfCall.ArgumentSpans[0], "@input")

	builtinCall := requireAnalysisCall(t, calls, "LENGTH")
	if builtinCall.Kind != CallKindBuiltin || builtinCall.Target != 0 || builtinCall.Identity != "length" {
		t.Fatalf("builtin call = %+v", builtinCall)
	}
	assertAnalysisSpanText(t, query, builtinCall.CalleeSpan, "LENGTH")

	hostCall := requireAnalysisCall(t, calls, "WEB::HTML::PARSE")
	if hostCall.Kind != CallKindHost || hostCall.Target != 0 || hostCall.Identity != "web::html::parse" || hostCall.ResultType != ValueTypeAny {
		t.Fatalf("host call = %+v", hostCall)
	}
	assertAnalysisSpanText(t, query, hostCall.CalleeSpan, "html::PARSE")

	valueAt, ok := analysis.SymbolAt(strings.LastIndex(query, "value"))
	if !ok || valueAt.ID != value.ID {
		t.Fatalf("SymbolAt(value) = %+v, %t", valueAt, ok)
	}

	callAt, ok := analysis.CallAt(strings.Index(query, "<p/>"))
	if !ok || callAt.Identity != "web::html::parse" {
		t.Fatalf("CallAt(<p/>) = %+v, %t", callAt, ok)
	}

	fact, ok := analysis.TypeAt(strings.LastIndex(query, "[value]"))
	if !ok || fact.Type != ValueTypeArray {
		t.Fatalf("TypeAt([value]) = %+v, %t", fact, ok)
	}
}

func TestAnalyzeVisibilityShadowingAndBindingKinds(t *testing.T) {
	query := `VAR x = 0
LET {left, nested: [right]} = {left: 1, nested: [2]}
FUNC inspect(x) {
  RETURN FOR item, index IN [x]
    COLLECT group = item WITH COUNT INTO count
    RETURN MATCH group { bound => [x, group, count, bound], _ => [] }
}
RETURN inspect(x)`

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	requireAnalysisSymbol(t, analysis, "left", SymbolKindBinding)
	requireAnalysisSymbol(t, analysis, "right", SymbolKindBinding)
	item := requireAnalysisSymbol(t, analysis, "item", SymbolKindLoopBinding)
	requireAnalysisSymbol(t, analysis, "index", SymbolKindLoopBinding)
	group := requireAnalysisSymbol(t, analysis, "group", SymbolKindCollectBinding)
	count := requireAnalysisSymbol(t, analysis, "count", SymbolKindCollectBinding)
	bound := requireAnalysisSymbol(t, analysis, "bound", SymbolKindMatchBinding)

	sourcePos := strings.Index(query, "[x]") + 1
	if hasAnalysisSymbol(analysis.VisibleSymbols(sourcePos), item.ID) {
		t.Fatalf("loop binding visible in its input expression: %+v", analysis.VisibleSymbols(sourcePos))
	}

	collectExprPos := strings.Index(query, "group = item") + len("group = ")
	visibleAtCollect := analysis.VisibleSymbols(collectExprPos)
	if !hasAnalysisSymbol(visibleAtCollect, item.ID) || hasAnalysisSymbol(visibleAtCollect, group.ID) || hasAnalysisSymbol(visibleAtCollect, count.ID) {
		t.Fatalf("unexpected COLLECT visibility: %+v", visibleAtCollect)
	}
	midCollectPos := strings.Index(query, "WITH COUNT")
	if visible := analysis.VisibleSymbols(midCollectPos); hasAnalysisSymbol(visible, group.ID) || hasAnalysisSymbol(visible, count.ID) {
		t.Fatalf("COLLECT bindings visible before clause end: %+v", visible)
	}

	matchPos := strings.Index(query, "MATCH group") + len("MATCH ")
	visibleAtMatch := analysis.VisibleSymbols(matchPos)
	if !hasAnalysisSymbol(visibleAtMatch, group.ID) || !hasAnalysisSymbol(visibleAtMatch, count.ID) || hasAnalysisSymbol(visibleAtMatch, bound.ID) {
		t.Fatalf("unexpected post-COLLECT visibility: %+v", visibleAtMatch)
	}

	boundResultPos := strings.Index(query, "[x, group, count, bound]")
	if !hasAnalysisSymbol(analysis.VisibleSymbols(boundResultPos), bound.ID) {
		t.Fatalf("MATCH binding is not visible in its arm result")
	}

	outerX := analysisSymbolsNamed(analysis.Symbols(), "x", SymbolKindBinding)
	parameterX := analysisSymbolsNamed(analysis.Symbols(), "x", SymbolKindFunctionParameter)
	if len(outerX) != 1 || len(parameterX) != 1 || outerX[0].ID == parameterX[0].ID {
		t.Fatalf("shadowed x symbols = outer %+v, parameter %+v", outerX, parameterX)
	}
	if !outerX[0].Mutable || parameterX[0].Mutable {
		t.Fatalf("mutability = outer %t, parameter %t", outerX[0].Mutable, parameterX[0].Mutable)
	}

	innerXPos := strings.Index(query, "[x]") + 1
	innerX, ok := analysis.SymbolAt(innerXPos)
	if !ok || innerX.ID != parameterX[0].ID {
		t.Fatalf("inner x resolved to %+v, want parameter %+v", innerX, parameterX[0])
	}

	outerXPos := strings.LastIndex(query, "inspect(x)") + len("inspect(")
	resolvedOuter, ok := analysis.SymbolAt(outerXPos)
	if !ok || resolvedOuter.ID != outerX[0].ID {
		t.Fatalf("outer x resolved to %+v, want %+v", resolvedOuter, outerX[0])
	}
}

func TestAnalyzeCapturesResolveOriginalSymbolsAndRejectFalseIdentities(t *testing.T) {
	query := `LET base = 1
FUNC outer(p) {
  LET local = p
  FUNC middle(q) {
    FUNC inner(r) => base + local + p + q + r
    RETURN inner(q)
	  }
	  RETURN middle(local)
	}
RETURN outer(base)`

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	for _, tc := range []struct {
		name string
		kind SymbolKind
	}{
		{name: "base", kind: SymbolKindBinding},
		{name: "local", kind: SymbolKindBinding},
		{name: "p", kind: SymbolKindFunctionParameter},
		{name: "q", kind: SymbolKindFunctionParameter},
		{name: "r", kind: SymbolKindFunctionParameter},
	} {
		symbols := analysisSymbolsNamed(analysis.Symbols(), tc.name, tc.kind)
		if len(symbols) != 1 {
			t.Fatalf("%s symbols = %+v, want one original declaration", tc.name, symbols)
		}

		if len(analysis.ReferencesTo(symbols[0].ID)) == 0 {
			t.Fatalf("%s has no resolved references", tc.name)
		}
	}

	invalid := `FUNC f() => later
LET later = 1
RETURN [f(), missing]`
	partial, err := New().Analyze(source.NewAnonymous(invalid))
	if err == nil || partial == nil || len(partial.Diagnostics()) == 0 {
		t.Fatalf("Analyze invalid = analysis %v, err %v", partial, err)
	}

	later := requireAnalysisSymbol(t, partial, "later", SymbolKindBinding)
	if refs := partial.ReferencesTo(later.ID); len(refs) != 0 {
		t.Fatalf("use-before-declaration received false identity: %+v", refs)
	}

	for _, reference := range partial.References() {
		assertValidAnalysisSymbolID(t, partial, reference.Symbol)
		if queryTextAt(invalid, reference.Span) == "missing" {
			t.Fatalf("unresolved reference received identity: %+v", reference)
		}
	}
}

func TestAnalyzeNestedLoopsForwardUDFsAndByteOffsets(t *testing.T) {
	query := `LET label = "é"
FUNC first(value) => second(value)
FUNC second(value) => value
RETURN FOR outer IN [[1]]
  RETURN FOR inner IN outer
    RETURN [label, first(inner)]`

	analysis, err := New().Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	label := requireAnalysisSymbol(t, analysis, "label", SymbolKindBinding)
	if label.SelectionSpan.Start != strings.Index(query, "label") || label.SelectionSpan.End-label.SelectionSpan.Start != len("label") {
		t.Fatalf("label byte span = %+v", label.SelectionSpan)
	}
	labelReference := analysis.ReferencesTo(label.ID)
	if len(labelReference) != 1 || labelReference[0].Span.Start != strings.LastIndex(query, "label") {
		t.Fatalf("label reference byte span = %+v", labelReference)
	}

	outer := requireAnalysisSymbol(t, analysis, "outer", SymbolKindLoopBinding)
	inner := requireAnalysisSymbol(t, analysis, "inner", SymbolKindLoopBinding)
	innerSource := strings.Index(query, "IN outer") + len("IN ")
	visible := analysis.VisibleSymbols(innerSource)
	if !hasAnalysisSymbol(visible, outer.ID) || hasAnalysisSymbol(visible, inner.ID) {
		t.Fatalf("nested loop source visibility = %+v", visible)
	}

	second := requireAnalysisSymbol(t, analysis, "second", SymbolKindUDF)
	assertAnalysisSpanText(t, query, second.SelectionSpan, "second")
	secondParams := analysis.FunctionParameters(second.ID)
	if len(secondParams) != 1 {
		t.Fatalf("second parameters = %+v", secondParams)
	}
	assertAnalysisSpanText(t, query, secondParams[0].SelectionSpan, "value")
	forward := requireAnalysisCall(t, analysis.Calls(), "second")
	if forward.Kind != CallKindUDF || forward.Target != second.ID {
		t.Fatalf("forward call = %+v, want target %+v", forward, second)
	}
}

func TestAnalyzeVisitsEveryMatchArm(t *testing.T) {
	query := `RETURN MATCH 1 {
  1 => FIRST(),
  2 => SECOND(),
  _ => THIRD(),
}`

	analysis, err := New(WithOptimizationLevel(O1), WithDebugInfo()).Analyze(source.NewAnonymous(query))
	if err != nil {
		t.Fatal(err)
	}

	calls := analysis.Calls()
	for _, name := range []string{"FIRST", "SECOND", "THIRD"} {
		call := requireAnalysisCall(t, calls, name)
		if call.Kind != CallKindHost {
			t.Fatalf("%s call = %+v", name, call)
		}
	}
}

func TestAnalyzePartialResultsEmptySyntaxAndDefensiveCopies(t *testing.T) {
	compiler := New()

	empty, err := compiler.Analyze(source.NewAnonymous(""))
	if err == nil || empty == nil || len(empty.Diagnostics()) != 1 || len(empty.Symbols()) != 0 {
		t.Fatalf("empty analysis = %+v, err %v", empty, err)
	}

	syntax, err := compiler.Analyze(source.NewAnonymous("LET = RETURN"))
	if err == nil || syntax == nil || len(syntax.Diagnostics()) == 0 || len(syntax.Symbols()) != 0 {
		t.Fatalf("syntax analysis = %+v, err %v", syntax, err)
	}

	query := `LET value = @input RETURN HOST(value, 1)`
	analysis, err := compiler.Analyze(source.New("copy.fql", query))
	if err != nil {
		t.Fatal(err)
	}

	symbols := analysis.Symbols()
	symbols[0].Name = "changed"
	if analysis.Symbols()[0].Name == "changed" {
		t.Fatal("Symbols returned mutable backing storage")
	}

	references := analysis.References()
	references[0].Span = source.Span{}
	if analysis.References()[0].Span == (source.Span{}) {
		t.Fatal("References returned mutable backing storage")
	}

	calls := analysis.Calls()
	calls[0].ArgumentSpans[0] = source.Span{}
	if analysis.Calls()[0].ArgumentSpans[0] == (source.Span{}) {
		t.Fatal("Calls returned mutable nested storage")
	}

	facts := analysis.TypeFacts()
	facts[0].Span = source.Span{}
	if analysis.TypeFacts()[0].Span == (source.Span{}) {
		t.Fatal("TypeFacts returned mutable backing storage")
	}

	diagnosticQuery := "LET text = \"é\"\nRETURN missing"
	diagnosticAnalysis, err := compiler.Analyze(source.New("diagnostic.fql", diagnosticQuery))
	if err == nil {
		t.Fatal("expected diagnostic")
	}

	diagnostics := diagnosticAnalysis.Diagnostics()
	assertAnalysisSpanText(t, diagnosticQuery, diagnostics[0].Spans[0].Span, "missing")
	diagnosticErr, ok := err.(*diagnosticspkg.Diagnostic)
	if !ok {
		t.Fatalf("diagnostic error type = %T", err)
	}
	assertAnalysisSpanText(t, diagnosticQuery, diagnosticErr.Spans[0].Span, "missing")
	diagnostics[0].Message = "changed"
	diagnostics[0].Spans[0].Label = "changed"
	if diagnosticAnalysis.Diagnostics()[0].Message == "changed" || diagnosticAnalysis.Diagnostics()[0].Spans[0].Label == "changed" {
		t.Fatal("Diagnostics returned mutable backing storage")
	}
}

func TestAnalyzeIsConcurrentSafeAndOptionIndependent(t *testing.T) {
	query := `LET base = 1
FUNC add(value) => base + value
RETURN add(@value)`
	src := source.NewAnonymous(query)
	shared := New(WithOptimizationLevel(O1), WithDebugInfo())

	want, err := shared.Analyze(src)
	if err != nil {
		t.Fatal(err)
	}

	const workers = 24
	errs := make(chan error, workers)
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)

		go func() {
			defer wg.Done()

			got, analyzeErr := shared.Analyze(src)
			if analyzeErr != nil {
				errs <- analyzeErr

				return
			}

			if !reflect.DeepEqual(got.Symbols(), want.Symbols()) || !reflect.DeepEqual(got.References(), want.References()) || !reflect.DeepEqual(got.Calls(), want.Calls()) || !reflect.DeepEqual(got.TypeFacts(), want.TypeFacts()) {
				errs <- errors.New("concurrent analysis result mismatch")
			}
		}()
	}

	wg.Wait()
	close(errs)

	for concurrentErr := range errs {
		t.Fatal(concurrentErr)
	}

	o0, err := New(WithOptimizationLevel(O0)).Analyze(src)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(o0.Symbols(), want.Symbols()) || !reflect.DeepEqual(o0.References(), want.References()) || !reflect.DeepEqual(o0.Calls(), want.Calls()) || !reflect.DeepEqual(o0.TypeFacts(), want.TypeFacts()) {
		t.Fatal("Analyze depends on compiler optimization or debug options")
	}
}

func requireAnalysisSymbol(t *testing.T, analysis *Analysis, name string, kind SymbolKind) Symbol {
	t.Helper()

	symbols := analysisSymbolsNamed(analysis.Symbols(), name, kind)
	if len(symbols) != 1 {
		t.Fatalf("symbols named %q kind %v = %+v, want one", name, kind, symbols)
	}

	return symbols[0]
}

func analysisSymbolsNamed(symbols []Symbol, name string, kind SymbolKind) []Symbol {
	out := make([]Symbol, 0)
	for _, symbol := range symbols {
		if symbol.Name == name && symbol.Kind == kind {
			out = append(out, symbol)
		}
	}

	return out
}

func requireAnalysisCall(t *testing.T, calls []Call, name string) Call {
	t.Helper()

	for _, call := range calls {
		if call.Name == name {
			return call
		}
	}

	t.Fatalf("call %q not found in %+v", name, calls)

	return Call{}
}

func assertAnalysisSpanText(t *testing.T, query string, span source.Span, want string) {
	t.Helper()

	if got := queryTextAt(query, span); got != want {
		t.Fatalf("span %+v text = %q, want %q", span, got, want)
	}
}

func queryTextAt(query string, span source.Span) string {
	if span.Start < 0 || span.End < span.Start || span.End > len(query) {
		return ""
	}

	return query[span.Start:span.End]
}

func hasAnalysisSymbol(symbols []Symbol, id SymbolID) bool {
	for _, symbol := range symbols {
		if symbol.ID == id {
			return true
		}
	}

	return false
}

func assertValidAnalysisSymbolID(t *testing.T, analysis *Analysis, id SymbolID) {
	t.Helper()

	for _, symbol := range analysis.Symbols() {
		if symbol.ID == id {
			return
		}
	}

	t.Fatalf("reference has invalid symbol ID %d", id)
}
