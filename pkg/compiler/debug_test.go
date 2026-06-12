package compiler

import (
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode/artifact"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestWithDebugInfoEmitsLogicalPointsAndForcesO0(t *testing.T) {
	program, err := New(WithOptimizationLevel(O1), WithDebugInfo()).Compile(
		source.New("debug.fql", "LET x = 1\nVAR y = 2\ny = y + x\nRETURN y"),
	)
	if err != nil {
		t.Fatal(err)
	}

	if program.Metadata.OptimizationLevel != int(O0) {
		t.Fatalf("expected O0, got O%d", program.Metadata.OptimizationLevel)
	}

	if len(program.Metadata.DebugPoints) != 4 {
		t.Fatalf("expected four debug points, got %#v", program.Metadata.DebugPoints)
	}

	if len(program.Metadata.DebugPoints[0].Bindings) != 0 {
		t.Fatalf("first declaration must not be visible before execution: %#v", program.Metadata.DebugPoints[0])
	}

	if got := program.Metadata.DebugPoints[1].Bindings; len(got) != 1 || got[0].Name != "x" {
		t.Fatalf("unexpected second point bindings: %#v", got)
	}

	if got := program.Metadata.DebugPoints[2].Bindings; len(got) != 2 {
		t.Fatalf("unexpected assignment bindings: %#v", got)
	}
}

func TestDebugInfoArtifactRoundTrip(t *testing.T) {
	program, err := New(WithDebugInfo()).Compile(source.New("debug.fql", "FOR i IN 1..2\n  RETURN i"))
	if err != nil {
		t.Fatal(err)
	}

	data, err := artifact.Marshal(program, artifact.Options{Format: artifact.FormatMsgPack})
	if err != nil {
		t.Fatal(err)
	}

	decoded, err := artifact.Unmarshal(data)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(decoded.Metadata.DebugPoints, program.Metadata.DebugPoints) {
		t.Fatalf("debug points mismatch: got %#v, want %#v", decoded.Metadata.DebugPoints, program.Metadata.DebugPoints)
	}
}

func TestNormalCompilationDoesNotEmitDebugPoints(t *testing.T) {
	program, err := New(WithOptimizationLevel(O1)).Compile(source.NewAnonymous("LET x = 1\nRETURN x"))
	if err != nil {
		t.Fatal(err)
	}

	if len(program.Metadata.DebugPoints) != 0 {
		t.Fatalf("normal compilation emitted debug points: %#v", program.Metadata.DebugPoints)
	}

	if program.Metadata.OptimizationLevel != int(O1) {
		t.Fatalf("normal compilation optimization changed: O%d", program.Metadata.OptimizationLevel)
	}
}

func TestDebugInfoIncludesUDFArgumentsCapturesAndLoopVariables(t *testing.T) {
	program, err := New(WithDebugInfo()).Compile(source.New("bindings.fql", `
LET base = 10
FUNC add(value) => base + value
RETURN (
  FOR item IN 1..2
    RETURN add(item)
)`))
	if err != nil {
		t.Fatal(err)
	}

	var udfBindings, loopBindings map[string]bool
	for _, point := range program.Metadata.DebugPoints {
		names := make(map[string]bool, len(point.Bindings))

		for _, binding := range point.Bindings {
			names[binding.Name] = true
		}

		if point.FunctionID >= 0 {
			udfBindings = names
		}

		line, _ := program.Source.LocationAt(point.Span)
		if line == 6 {
			loopBindings = names
		}
	}

	if !udfBindings["base"] || !udfBindings["value"] {
		t.Fatalf("UDF point missing argument/capture bindings: %#v", udfBindings)
	}

	if !loopBindings["item"] || !loopBindings["base"] {
		t.Fatalf("loop point missing visible bindings: %#v", loopBindings)
	}
}

func TestDebugInfoPreservesLexicalShadowing(t *testing.T) {
	program, err := New(WithDebugInfo()).Compile(source.New("shadow.fql", `LET x = 1
RETURN (
  FOR x IN [2]
    RETURN x
)`))
	if err != nil {
		t.Fatal(err)
	}

	for _, point := range program.Metadata.DebugPoints {
		line, _ := program.Source.LocationAt(point.Span)

		if line != 4 {
			continue
		}

		count := 0
		for _, binding := range point.Bindings {
			if binding.Name == "x" {
				count++
			}
		}

		if count != 1 {
			t.Fatalf("expected one visible shadowed binding, got %#v", point.Bindings)
		}

		return
	}

	t.Fatal("expected loop return debug point")
}
