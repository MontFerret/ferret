package compiler

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
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
	assertSourcePointsMatchDebugPoints(t, program)

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

	for _, format := range []artifact.FormatID{artifact.FormatJSON, artifact.FormatMsgPack} {
		t.Run(fmt.Sprintf("format_%d", format), func(t *testing.T) {
			data, err := artifact.Marshal(program, artifact.Options{Format: format})
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
			if !reflect.DeepEqual(decoded.Bytecode, program.Bytecode) {
				t.Fatalf("bytecode mismatch: got %#v, want %#v", decoded.Bytecode, program.Bytecode)
			}
			assertSourcePointsMatchDebugPoints(t, decoded)
		})
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

	for pc, inst := range program.Bytecode {
		if inst.Opcode == bytecode.OpSourcePoint {
			t.Fatalf("normal compilation emitted source point at pc %d", pc)
		}
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

func assertSourcePointsMatchDebugPoints(t *testing.T, program *bytecode.Program) {
	t.Helper()

	ids := make(map[bytecode.DebugPointID]struct{}, len(program.Metadata.DebugPoints))

	for pointIndex, point := range program.Metadata.DebugPoints {
		if point.ID != bytecode.DebugPointID(pointIndex) {
			t.Fatalf("debug point %d has non-monotonic id %d", pointIndex, point.ID)
		}
		if _, exists := ids[point.ID]; exists {
			t.Fatalf("duplicate debug point id %d", point.ID)
		}
		ids[point.ID] = struct{}{}

		if point.PC < 0 || point.PC >= len(program.Bytecode) {
			t.Fatalf("debug point %d pc %d out of range", point.ID, point.PC)
		}

		inst := program.Bytecode[point.PC]
		if inst.Opcode != bytecode.OpSourcePoint || bytecode.DebugPointID(inst.Operands[0]) != point.ID {
			t.Fatalf("debug point %d does not match source point at pc %d: %#v", point.ID, point.PC, inst)
		}
	}
}
