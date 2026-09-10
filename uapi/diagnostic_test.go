package uapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/MontFerret/api"
	apidiagnostics "github.com/MontFerret/api/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func TestDiagnosticsPreserveErrorTreeAndOwnSources(t *testing.T) {
	firstCause := errors.New("first cause")
	secondCause := errors.New("second cause")
	first := diagnostics.NewUnexpectedErrorWith(source.New("buffer://first", "éx\nRETURN 1"), "first", firstCause)
	first.Spans = []diagnostics.ErrorSpan{
		diagnostics.NewMainErrorSpan(source.Span{Start: 4, End: 10}, "return"),
		diagnostics.NewSecondaryErrorSpan(source.Span{Start: 2, End: 3}, "after unicode"),
	}

	first.Hint, first.Note = "hint", "note"
	second := diagnostics.NewUnexpectedErrorWith(source.NewAnonymous("RETURN 2"), "second", secondCause)
	second.Spans = []diagnostics.ErrorSpan{diagnostics.NewMainErrorSpan(source.Span{Start: 0, End: 6}, "second")}
	set := diagnostics.NewDiagnosticsOf([]*diagnostics.Diagnostic{second})
	input := errors.Join(fmt.Errorf("wrapped: %w", &vm.RuntimeError{Diagnostic: first}), set, first)
	output := wrapDiagnosticError(input)

	if !errors.Is(output, firstCause) || !errors.Is(output, secondCause) {
		t.Fatalf("lost original causes: %v", output)
	}

	var native *vm.RuntimeError
	if !errors.As(output, &native) || native.Diagnostic != first {
		t.Fatal("lost native runtime error")
	}

	var portable apidiagnostics.Diagnostics
	if !errors.As(output, &portable) || len(portable) != 2 {
		t.Fatalf("diagnostics = %+v", portable)
	}

	if portable[0].Source != api.NewSource("buffer://first", "éx\nRETURN 1") || portable[1].Source != api.NewAnonymousSource("RETURN 2") {
		t.Fatalf("source identities = %+v", portable)
	}

	if portable[0].Hint != "hint" || portable[0].Note != "note" || len(portable[0].Annotations) != 2 {
		t.Fatalf("diagnostic metadata = %+v", portable[0])
	}

	annotation := portable[0].Annotations[0]
	if annotation.Range.Position != (api.Position{Line: 2, Column: 1}) || annotation.Range.Span != (api.Span{Start: 4, End: 10}) || !annotation.Primary || annotation.Message != "return" {
		t.Fatalf("annotation = %+v", annotation)
	}

	if position := portable[0].Annotations[1].Range.Position; position != (api.Position{Line: 1, Column: 3}) {
		t.Fatalf("native columns must count UTF-8 bytes: %+v", position)
	}

	first.Spans[0].Label = "changed"

	if portable[0].Annotations[0].Message != "return" {
		t.Fatal("portable annotation aliases native storage")
	}

	if output.Error() != input.Error() {
		t.Fatal("projection changed the native error message")
	}
}

func TestDiagnosticWithoutSourceRemainsWithoutSource(t *testing.T) {
	input := diagnostics.NewUnexpectedError(source.Source{}, "hook failed")
	output := wrapDiagnosticError(input)

	var portable apidiagnostics.Diagnostics
	if !errors.As(output, &portable) || len(portable) != 1 || portable[0].Source != (api.Source{}) {
		t.Fatalf("source-less diagnostic = %+v", portable)
	}

	plain := errors.New("plain")
	if !errors.Is(wrapDiagnosticError(plain), plain) || wrapDiagnosticError(nil) != nil {
		t.Fatal("plain errors should retain identity")
	}
}

func TestEqualDiagnosticsRemainDistinctAndOrdered(t *testing.T) {
	firstCause, secondCause := errors.New("first"), errors.New("second")
	src := source.New("buffer://same", "RETURN 1")
	first := diagnostics.NewUnexpectedErrorWith(src, "same diagnostic", firstCause)
	second := diagnostics.NewUnexpectedErrorWith(src, "same diagnostic", secondCause)
	first.Note, second.Note = "first note", "second note"
	output := wrapDiagnosticError(errors.Join(first, second, first))

	var portable apidiagnostics.Diagnostics
	if !errors.As(output, &portable) || len(portable) != 2 ||
		portable[0].Note != "first note" || portable[1].Note != "second note" {
		t.Fatalf("distinct diagnostic order=%+v", portable)
	}

	if !errors.Is(output, firstCause) || !errors.Is(output, secondCause) {
		t.Fatal("distinct native causes were lost")
	}
}

func TestNestedProjectionExposesAllDiagnostics(t *testing.T) {
	first := diagnostics.NewUnexpectedError(source.Source{}, "first")
	second := diagnostics.NewUnexpectedError(source.Source{}, "second")
	err := wrapDiagnosticError(errors.Join(wrapDiagnosticError(first), second))
	var portable apidiagnostics.Diagnostics
	if !errors.As(err, &portable) || len(portable) != 2 || portable[0].Message != "first" || portable[1].Message != "second" {
		t.Fatalf("nested projection=%+v err=%v", portable, err)
	}

	if !errors.Is(err, first) || !errors.Is(err, second) {
		t.Fatal("projection lost native causes")
	}
}

func TestProjectionPreservesPortableDiagnosticSiblings(t *testing.T) {
	portable := apidiagnostics.Diagnostics{{Message: "portable"}}
	native := diagnostics.NewUnexpectedError(source.Source{}, "native")
	err := wrapDiagnosticError(errors.Join(portable, wrapDiagnosticError(native)))
	var projected apidiagnostics.Diagnostics
	if !errors.As(err, &projected) || len(projected) != 2 || projected[0].Message != "portable" || projected[1].Message != "native" {
		t.Fatalf("mixed diagnostic projection=%+v err=%v", projected, err)
	}

	plain := errors.Join(portable, context.Canceled)
	if wrapDiagnosticError(plain) != plain {
		t.Fatal("portable-only error tree was needlessly wrapped")
	}
}

func TestCompilationProducesPortableByteCoordinates(t *testing.T) {
	runtime := newTestRuntime(t)
	for _, input := range []struct {
		source api.Source
		line   int
		column int
	}{
		{api.NewSource("buffer://unicode", `RETURN "é" + missing`), 1, 15},
		{api.NewAnonymousSource("LET prefix = \"é\"\nRETURN missing"), 2, 8},
	} {
		plan, err := runtime.Compile(t.Context(), input.source)
		if plan != nil {
			_ = plan.Close()
			t.Fatal("invalid source published a plan")
		}

		var portable apidiagnostics.Diagnostics
		if !errors.As(err, &portable) || len(portable) != 1 {
			t.Fatalf("compiler diagnostics=%+v error=%v", portable, err)
		}

		start := strings.Index(input.source.Content, "missing")
		found := false
		for _, annotation := range portable[0].Annotations {
			if !annotation.Primary {
				continue
			}

			found = true

			if annotation.Range != (api.Range{
				Location: api.Location{SourceName: input.source.Name, Position: api.Position{Line: input.line, Column: input.column}},
				Span:     api.Span{Start: start, End: start + len("missing")},
			}) {
				t.Fatalf("compiler range=%+v for %q", annotation.Range, input.source.Content)
			}
		}

		if !found || portable[0].Source != input.source {
			t.Fatalf("compiler source/primary annotation=%+v", portable[0])
		}
	}
}
