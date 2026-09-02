package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestDestructuringLoweringUsesAssertionsAndOptionalConstantLoads(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.OptimizationNone, compiler.OptimizationFull} {
		prog := compileWithLevel(t, level, `
LET { user: { name }, values: [first, _] } = @payload
RETURN [name, first]
`)

		counts := map[bytecode.Opcode]int{}
		for _, inst := range prog.Bytecode {
			counts[inst.Opcode]++
		}

		if got, want := counts[bytecode.OpAssertDestructure], 3; got != want {
			t.Fatalf("O%d assertion count = %d, want %d", level, got, want)
		}

		if got, want := counts[bytecode.OpLoadKeyOptionalConst], 3; got != want {
			t.Fatalf("O%d object load count = %d, want %d", level, got, want)
		}

		if got, want := counts[bytecode.OpLoadIndexOptionalConst], 1; got != want {
			t.Fatalf("O%d array load count = %d, want %d", level, got, want)
		}
	}
}

func TestDestructuringLoweringSkipsIgnoredStructuredChildren(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.OptimizationNone, compiler.OptimizationFull} {
		ignored := compileWithLevel(t, level, `
LET {
    kept,
    ignored: { nested: [_, _] },
    emptyObject: {},
    emptyArray: []
} = @payload
RETURN kept
`)
		direct := compileWithLevel(t, level, `
LET { kept, ignored: _, emptyObject: _, emptyArray: _ } = @payload
RETURN kept
`)

		counts := map[bytecode.Opcode]int{}
		for _, inst := range ignored.Bytecode {
			counts[inst.Opcode]++
		}

		if got, want := counts[bytecode.OpAssertDestructure], 1; got != want {
			t.Fatalf("O%d assertion count = %d, want %d", level, got, want)
		}

		if got, want := counts[bytecode.OpLoadKeyOptionalConst], 1; got != want {
			t.Fatalf("O%d object load count = %d, want %d", level, got, want)
		}

		if got := counts[bytecode.OpLoadIndexOptionalConst]; got != 0 {
			t.Fatalf("O%d array load count = %d, want 0", level, got)
		}

		if got, want := ignored.Registers, direct.Registers; got != want {
			t.Fatalf("O%d register count = %d, direct ignore uses %d", level, got, want)
		}

		if got, want := len(ignored.Constants), len(direct.Constants); got != want {
			t.Fatalf("O%d constant count = %d, direct ignore uses %d", level, got, want)
		}
	}
}

func TestDestructuringLoweringRetainsMixedStructuredChildren(t *testing.T) {
	for _, level := range []compiler.OptimizationLevel{compiler.OptimizationNone, compiler.OptimizationFull} {
		prog := compileWithLevel(t, level, `
LET { nested: [_, kept, _], ignored: { child: [_] } } = @payload
RETURN kept
`)

		counts := map[bytecode.Opcode]int{}
		for _, inst := range prog.Bytecode {
			counts[inst.Opcode]++
		}

		if got, want := counts[bytecode.OpAssertDestructure], 2; got != want {
			t.Fatalf("O%d assertion count = %d, want %d", level, got, want)
		}

		if got, want := counts[bytecode.OpLoadKeyOptionalConst], 1; got != want {
			t.Fatalf("O%d object load count = %d, want %d", level, got, want)
		}

		if got, want := counts[bytecode.OpLoadIndexOptionalConst], 1; got != want {
			t.Fatalf("O%d array load count = %d, want %d", level, got, want)
		}
	}
}

func TestDestructuringDuplicateBindingDiagnostic(t *testing.T) {
	src := `LET { name, nested: [name] } = @payload
RETURN name`
	_, err := mustNewCompiler(t).Compile(source.New("duplicate_destructure.fql", src))
	if err == nil {
		t.Fatal("expected duplicate binding error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatalf("expected diagnostic, got %T", err)
	}

	if got, want := diag.Kind, parserd.NameError; got != want {
		t.Fatalf("diagnostic kind = %s, want %s", got, want)
	}

	if got, want := diag.Message, `duplicate binding "name" in destructuring pattern`; got != want {
		t.Fatalf("diagnostic message = %q, want %q", got, want)
	}

	if len(diag.Spans) != 2 {
		t.Fatalf("diagnostic spans = %d, want 2", len(diag.Spans))
	}

	if got, want := diag.Spans[0].Label, "duplicate binding"; got != want {
		t.Fatalf("main label = %q, want %q", got, want)
	}

	if got, want := diag.Spans[1].Label, "first bound here"; got != want {
		t.Fatalf("secondary label = %q, want %q", got, want)
	}
}

func TestForDestructuringDuplicateBindingDiagnostic(t *testing.T) {
	_, err := mustNewCompiler(t).Compile(source.NewAnonymous(`
FOR { value, nested: [value] } IN @values
    RETURN value
`))
	if err == nil {
		t.Fatal("expected duplicate loop binding error")
	}

	diag := firstCompilationError(err)
	if diag == nil || diag.Message != `duplicate binding "value" in destructuring pattern` {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

func TestDestructuredLeavesParticipateInForwardDiagnostics(t *testing.T) {
	_, err := mustNewCompiler(t).Compile(source.NewAnonymous(`
LET before = value
LET { value } = @payload
RETURN before
`))
	if err == nil {
		t.Fatal("expected forward binding error")
	}

	diag := firstCompilationError(err)
	if diag == nil || !strings.Contains(diag.Message, "used before declaration") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

func TestDestructuredBindingsConflictWithFunctionNames(t *testing.T) {
	_, err := mustNewCompiler(t).Compile(source.NewAnonymous(`
FUNC value() => 1
LET { value } = { value: 2 }
RETURN value
`))
	if err == nil {
		t.Fatal("expected function and destructured binding conflict")
	}

	diag := firstCompilationError(err)
	if diag == nil || diag.Message != "Variable 'value' is already defined" {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}

func TestDestructuredBindingsAppearInDebugMetadata(t *testing.T) {
	prog, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(`
LET { first, second: alias } = { first: 1, second: 2 }
RETURN [first, alias]
`))
	if err != nil {
		t.Fatal(err)
	}

	found := map[string]bool{}
	for _, point := range prog.Metadata.DebugPoints {
		for _, binding := range point.Bindings {
			if binding.Name == "first" || binding.Name == "alias" {
				found[binding.Name] = true
			}
		}
	}

	for _, name := range []string{"first", "alias"} {
		if !found[name] {
			t.Fatalf("debug metadata does not expose %q", name)
		}
	}
}

func TestForDestructuredBindingsAppearInDebugMetadata(t *testing.T) {
	prog, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(`
RETURN (
    FOR { first, second: alias } IN [{ first: 1, second: 2 }]
        RETURN [first, alias]
)
`))
	if err != nil {
		t.Fatal(err)
	}

	found := map[string]bool{}
	for _, point := range prog.Metadata.DebugPoints {
		for _, binding := range point.Bindings {
			if binding.Name == "first" || binding.Name == "alias" {
				found[binding.Name] = true
			}
		}
	}

	for _, name := range []string{"first", "alias"} {
		if !found[name] {
			t.Fatalf("debug metadata does not expose FOR binding %q", name)
		}
	}
}

func TestMutableDestructuredBindingsRetainDebugStorageMetadata(t *testing.T) {
	prog, err := mustNewCompiler(t, compiler.WithDebugInfo()).Compile(source.NewAnonymous(`
VAR { captured, local } = { captured: 1, local: 2 }
FUNC increment() {
    captured += 1
    RETURN captured
}
RETURN [increment(), local]
`))
	if err != nil {
		t.Fatal(err)
	}

	capturedMutableCell := false
	localMutableValue := false
	for pointIndex := range prog.Metadata.DebugPoints {
		for bindingIndex := range prog.Metadata.DebugPoints[pointIndex].Bindings {
			binding := &prog.Metadata.DebugPoints[pointIndex].Bindings[bindingIndex]
			switch binding.Name {
			case "captured":
				capturedMutableCell = capturedMutableCell || binding.Mutable && binding.Cell
			case "local":
				localMutableValue = localMutableValue || binding.Mutable && !binding.Cell
			}
		}
	}

	if !capturedMutableCell {
		t.Fatal("debug metadata does not expose captured VAR as a mutable cell")
	}

	if !localMutableValue {
		t.Fatal("debug metadata does not expose uncaptured VAR as a mutable value binding")
	}
}

func TestForDestructuredBindingsRemainLexicallyScoped(t *testing.T) {
	_, err := mustNewCompiler(t).Compile(source.NewAnonymous(`
LET values = (
    FOR { value } IN [{ value: 1 }]
        RETURN value
)
RETURN value
`))
	if err == nil {
		t.Fatal("expected loop binding scope error")
	}

	diag := firstCompilationError(err)
	if diag == nil || !strings.Contains(diag.Message, "is not defined") {
		t.Fatalf("unexpected diagnostic: %v", err)
	}
}
