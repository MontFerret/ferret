package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

func TestCoalesceLoweringUsesNoneBranch(t *testing.T) {
	program := compileWithLevel(t, compiler.None, "RETURN @value ?? 42")
	code := program.Bytecode

	jumpIndex := opcodeIndex(code, bytecode.OpJumpIfNone)
	if jumpIndex < 0 {
		t.Fatalf("expected bytecode to contain %s", bytecode.OpJumpIfNone)
	}

	if count := opcodeCount(code, bytecode.OpLoadParam); count != 1 {
		t.Fatalf("expected the left operand to be evaluated once, got %d parameter loads", count)
	}

	fallbackTarget := int(code[jumpIndex].Operands[0])
	if fallbackTarget <= jumpIndex+2 || fallbackTarget >= len(code) {
		t.Fatalf("invalid fallback target %d for jump at %d", fallbackTarget, jumpIndex)
	}

	presentMove := code[jumpIndex+1]
	if presentMove.Opcode != bytecode.OpMoveTracked {
		t.Fatalf("expected tracked present-value move after NONE check, got %s", presentMove.Opcode)
	}

	if got, want := presentMove.Operands[1], code[jumpIndex].Operands[1]; got != want {
		t.Fatalf("present branch moved %s, want checked operand %s", got, want)
	}

	skipFallback := code[jumpIndex+2]
	if skipFallback.Opcode != bytecode.OpJump {
		t.Fatalf("expected present branch to skip fallback, got %s", skipFallback.Opcode)
	}

	doneTarget := int(skipFallback.Operands[0])
	if doneTarget <= fallbackTarget || doneTarget > len(code) {
		t.Fatalf("invalid done target %d for fallback at %d", doneTarget, fallbackTarget)
	}

	fallbackMove := opcodeIndexBetween(code, bytecode.OpMoveTracked, fallbackTarget, doneTarget)
	if fallbackMove < 0 {
		t.Fatal("expected tracked fallback-value move inside fallback branch")
	}

	if got, want := code[fallbackMove].Operands[0], presentMove.Operands[0]; got != want {
		t.Fatalf("fallback branch wrote %s, want result register %s", got, want)
	}
}

func TestCoalesceMissingFallbackDiagnostic(t *testing.T) {
	input := "RETURN 1 ??"
	_, err := mustNewCompiler(t, compiler.WithOptimizationLevel(compiler.None)).Compile(source.NewAnonymous(input))
	if err == nil {
		t.Fatal("expected compilation to fail")
	}

	diagnostic := firstCompilationError(err)
	if diagnostic == nil {
		t.Fatalf("expected diagnostic error, got %T", err)
	}

	if diagnostic.Kind != parserd.SyntaxError {
		t.Fatalf("unexpected diagnostic kind: got %q, want %q", diagnostic.Kind, parserd.SyntaxError)
	}

	if got, want := diagnostic.Message, "Expected right-hand expression after '??'"; got != want {
		t.Fatalf("unexpected diagnostic message: got %q, want %q", got, want)
	}

	if got, want := diagnostic.Hint, "Provide a fallback expression, e.g. value ?? fallback."; got != want {
		t.Fatalf("unexpected diagnostic hint: got %q, want %q", got, want)
	}

	if len(diagnostic.Spans) != 1 {
		t.Fatalf("expected one diagnostic span, got %d", len(diagnostic.Spans))
	}

	operatorStart := strings.Index(input, "??")
	span := diagnostic.Spans[0]
	if !span.Main {
		t.Fatal("expected primary diagnostic span")
	}

	if got, want := span.Span.Start, operatorStart; got != want {
		t.Fatalf("unexpected span start: got %d, want %d", got, want)
	}

	if got, want := span.Span.End, operatorStart+2; got != want {
		t.Fatalf("unexpected span end: got %d, want %d", got, want)
	}
}

func TestCoalesceMissingFallbackDiagnosticIgnoresComments(t *testing.T) {
	_, err := mustNewCompiler(t, compiler.WithOptimizationLevel(compiler.None)).Compile(source.NewAnonymous("RETURN 1 + // ??"))
	if err == nil {
		t.Fatal("expected compilation to fail")
	}

	diagnostic := firstCompilationError(err)
	if diagnostic == nil {
		t.Fatalf("expected diagnostic error, got %T", err)
	}

	if got := diagnostic.Message; got == "Expected right-hand expression after '??'" {
		t.Fatalf("trailing comment was treated as a coalescing operator: %q", got)
	}
}

func opcodeIndex(code []bytecode.Instruction, opcode bytecode.Opcode) int {
	return opcodeIndexBetween(code, opcode, 0, len(code))
}

func opcodeIndexBetween(code []bytecode.Instruction, opcode bytecode.Opcode, start, end int) int {
	for i := start; i < end; i++ {
		if code[i].Opcode == opcode {
			return i
		}
	}

	return -1
}

func opcodeCount(code []bytecode.Instruction, opcode bytecode.Opcode) int {
	count := 0

	for _, instruction := range code {
		if instruction.Opcode == opcode {
			count++
		}
	}

	return count
}
