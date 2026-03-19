package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	"github.com/MontFerret/ferret/v2/pkg/file"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
)

func assertNoCellOps(t *testing.T, prog *bytecode.Program) {
	t.Helper()

	for _, op := range []bytecode.Opcode{bytecode.OpMakeCell, bytecode.OpLoadCell, bytecode.OpStoreCell} {
		if got := countOpcode(prog, op); got != 0 {
			t.Fatalf("expected no %s opcodes, got %d", op, got)
		}
	}
}

func TestVarSyntaxErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
			VAR
			RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "VAR missing variable name"),
		ErrorCase(
			`
			VAR x =
			RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'x'",
				Hint:    "Did you forget to provide a value?",
			}, "VAR missing assignment value"),
		ErrorCase(
			`
			VAR _ = 1
			RETURN 0
		`, E{
				Kind: parserd.SyntaxError,
			}, "VAR cannot use discard binding"),
		ErrorCase(
			`
				VAR x = 0
				RETURN (x = 1)
			`, E{
				Kind: parserd.SyntaxError,
			}, "Assignment is not allowed inside expressions"),
		ErrorCase(
			`
				VAR x = 0
				x +=
				RETURN x
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '+=' for variable 'x'",
				Hint:    "Did you forget to provide a value?",
			}, "Compound assignment missing assignment value"),
		ErrorCase(
			`
				VAR x = 0
				RETURN (x += 1)
			`, E{
				Kind: parserd.SyntaxError,
			}, "Compound assignment is not allowed inside expressions"),
	})
}

func TestVarCompoundAssignmentMissingValueDiagnosticSpan(t *testing.T) {
	src := "VAR x = 0\nx +=\nRETURN x"

	_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).Compile(file.NewSource("var_compound_span", src))
	if err == nil {
		t.Fatal("expected compilation error")
	}

	diag := firstCompilationError(err)
	if diag == nil {
		t.Fatal("expected diagnostic")
	}

	if diag.Kind != parserd.SyntaxError {
		t.Fatalf("expected syntax error, got %s", diag.Kind)
	}

	if diag.Message != "Expected expression after '+=' for variable 'x'" {
		t.Fatalf("unexpected diagnostic message: %q", diag.Message)
	}

	if len(diag.Spans) == 0 {
		t.Fatal("expected diagnostic spans")
	}

	wantStart := strings.Index(src, "+=") + len("+=")
	got := diag.Spans[0].Span
	if got.Start != wantStart || got.End != wantStart+1 {
		t.Fatalf("expected span [%d,%d), got [%d,%d)", wantStart, wantStart+1, got.Start, got.End)
	}
}

func TestVarErrors(t *testing.T) {
	RunUseCases(t, []UseCase{
		ErrorCase(
			`
				LET x = 1
				x = 2
				RETURN x
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "LET remains immutable"),
		ErrorCase(
			`
				LET x = 1
				x += 1
				RETURN x
			`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "LET remains immutable for compound assignment"),
		ErrorCase(
			`
				x = 1
				RETURN 0
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'x' is not defined",
			}, "Assignment target must already exist"),
		ErrorCase(
			`
			FUNC bump(x) (
			  x = x + 1
			  RETURN x
			)
			RETURN bump(1)
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "Parameters cannot be reassigned"),
		ErrorCase(
			`
			FOR i WHILE i < 2
			  i = i + 1
			  RETURN i
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'i' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "FOR WHILE variables cannot be reassigned"),
		ErrorCase(
			`
			LET obj = {}
			obj.x = 1
			RETURN obj
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Assignment target must be a local variable name",
				Hint:    "Property and index assignment are not supported. Use UPDATE for structural changes.",
			}, "Property assignment is rejected"),
		ErrorCase(
			`
				LET obj = {}
				obj.x += 1
				RETURN obj
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Assignment target must be a local variable name",
				Hint:    "Property and index assignment are not supported. Use UPDATE for structural changes.",
			}, "Compound property assignment is rejected"),
		ErrorCase(
			`
				LET arr = [0]
				arr[0] = 1
				RETURN arr
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Assignment target must be a local variable name",
				Hint:    "Property and index assignment are not supported. Use UPDATE for structural changes.",
			}, "Index assignment is rejected"),
		ErrorCase(
			`
				LET arr = [0]
				arr[0] += 1
				RETURN arr
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Assignment target must be a local variable name",
				Hint:    "Property and index assignment are not supported. Use UPDATE for structural changes.",
			}, "Compound index assignment is rejected"),
		ErrorCase(
			`
				VAR x = 1
				FUNC outer() (
			  LET x = 2
			  x = 3
			  RETURN x
			)
			RETURN outer()
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "Nearest shadowed binding controls reassignment"),
	})
}

func TestVarSupportedStatementPositionsCompile(t *testing.T) {
	expressions := []string{
		`
VAR counter = 1
counter = counter + 1
RETURN counter
`,
		`
FUNC run() (
  VAR total = 1
  total = total + 1
  RETURN total
)
RETURN run()
`,
		`
FOR item IN [1, 2]
  VAR current = item
  current = current + 1
  RETURN current
		`,
		`
	VAR i = 0
	FOR WHILE i < 2
	  LET current = i
	  i = i + 1
	  RETURN current
	`,
		`
	VAR total = 10
	total += 1
	total -= 2
	total *= 3
	total /= 3
	RETURN total
	`,
		`
	FUNC run() (
	  VAR total = 10
	  total += 1
	  total -= 2
	  total *= 3
	  total /= 3
	  RETURN total
	)
	RETURN run()
	`,
		`
	VAR i = 0
	FOR WHILE i < 3
	  LET current = i
	  i += 1
	  i -= 0
	  i *= 1
	  i /= 1
	  RETURN current
	`,
	}

	for _, expr := range expressions {
		_ = compileWithLevel(t, compiler.O0, expr)
	}
}

func TestVarRegisterBackedReassignmentAvoidsCellOps(t *testing.T) {
	expr := `
VAR x = 1
x = x + 1
RETURN x
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		prog := compileWithLevel(t, level, expr)
		assertNoCellOps(t, prog)
	}
}

func TestVarRegisterBackedCompoundAssignmentAvoidsCellOps(t *testing.T) {
	expr := `
	VAR x = 1
	x += 1
	RETURN x
	`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		prog := compileWithLevel(t, level, expr)
		assertNoCellOps(t, prog)
	}
}

func TestVarReadOnlyCaptureStaysByValueAcrossOptimizationLevels(t *testing.T) {
	expr := `
	VAR base = 1
	FUNC getBase() => base
RETURN getBase()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		prog := compileWithLevel(t, level, expr)
		assertNoCellOps(t, prog)
	}
}

func TestVarWriteCaptureUsesCellOpsAcrossOptimizationLevels(t *testing.T) {
	expr := `
VAR base = 1
FUNC setBase(v) (
  base = v
  RETURN base
)
RETURN setBase(2)
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		prog := compileWithLevel(t, level, expr)

		if got := countOpcode(prog, bytecode.OpMakeCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpMakeCell, level)
		}

		if got := countOpcode(prog, bytecode.OpLoadCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpLoadCell, level)
		}

		if got := countOpcode(prog, bytecode.OpStoreCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpStoreCell, level)
		}
	}
}

func TestVarWriteCaptureCompoundAssignmentUsesCellOpsAcrossOptimizationLevels(t *testing.T) {
	expr := `
	VAR base = 1
	FUNC addToBase(v) (
	  base += v
	  RETURN base
	)
	RETURN addToBase(2)
	`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		prog := compileWithLevel(t, level, expr)

		if got := countOpcode(prog, bytecode.OpMakeCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpMakeCell, level)
		}

		if got := countOpcode(prog, bytecode.OpLoadCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpLoadCell, level)
		}

		if got := countOpcode(prog, bytecode.OpStoreCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpStoreCell, level)
		}
	}
}

func TestVarReassignmentOutsideLoopKeepsExactType(t *testing.T) {
	expr := `
	VAR x = [1, 2]
	x = { value: 1 }
RETURN x[0]
`

	prog := compileWithLevel(t, compiler.O0, expr)
	assertNoCellOps(t, prog)

	if !hasOpcode(prog.Bytecode, bytecode.OpLoadKeyConst) {
		t.Fatalf("expected OpLoadKeyConst after straight-line reassignment")
	}

	if hasOpcode(prog.Bytecode, bytecode.OpLoadPropertyConst) {
		t.Fatalf("did not expect OpLoadPropertyConst after straight-line reassignment")
	}
}

func TestVarReassignmentInLoopWidenTypeForValueBindings(t *testing.T) {
	expr := `
VAR x = [1, 2]
LET ignored = (
  FOR item IN @items
    FILTER item
    x = { value: item }
    RETURN item
)
RETURN x[0]
`

	prog := compileWithLevel(t, compiler.O0, expr)
	assertNoCellOps(t, prog)

	if !hasOpcode(prog.Bytecode, bytecode.OpLoadPropertyConst) {
		t.Fatalf("expected OpLoadPropertyConst after loop-scoped conflicting reassignment")
	}

	if hasOpcode(prog.Bytecode, bytecode.OpLoadIndexConst) || hasOpcode(prog.Bytecode, bytecode.OpLoadKeyConst) {
		t.Fatalf("did not expect exact container load opcode after loop-scoped conflicting reassignment")
	}
}

func TestVarReassignmentInLoopWidenTypeForCellBindings(t *testing.T) {
	expr := `
VAR x = [1, 2]
FUNC touch(v) (
  x = v
  RETURN x
)
LET ignored = (
  FOR item IN @items
    FILTER item
    x = { value: item }
    RETURN item
)
RETURN x[0]
`

	prog := compileWithLevel(t, compiler.O0, expr)

	if got := countOpcode(prog, bytecode.OpMakeCell); got == 0 {
		t.Fatalf("expected OpMakeCell for captured mutable binding")
	}

	if got := countOpcode(prog, bytecode.OpLoadCell); got == 0 {
		t.Fatalf("expected OpLoadCell for captured mutable binding")
	}

	if got := countOpcode(prog, bytecode.OpStoreCell); got == 0 {
		t.Fatalf("expected OpStoreCell for captured mutable binding")
	}

	if !hasOpcode(prog.Bytecode, bytecode.OpLoadPropertyConst) {
		t.Fatalf("expected OpLoadPropertyConst after loop-scoped conflicting reassignment through cell binding")
	}

	if hasOpcode(prog.Bytecode, bytecode.OpLoadIndexConst) || hasOpcode(prog.Bytecode, bytecode.OpLoadKeyConst) {
		t.Fatalf("did not expect exact container load opcode after loop-scoped conflicting reassignment through cell binding")
	}
}

func TestVarReassignmentInLoopPreservesSameTypePrecision(t *testing.T) {
	expr := `
VAR x = [1, 2]
LET ignored = (
  FOR item IN @items
    FILTER item
    x = [item]
    RETURN item
)
RETURN x[0]
`

	prog := compileWithLevel(t, compiler.O0, expr)
	assertNoCellOps(t, prog)

	if !hasOpcode(prog.Bytecode, bytecode.OpLoadIndexConst) {
		t.Fatalf("expected OpLoadIndexConst after same-type loop reassignment")
	}

	if hasOpcode(prog.Bytecode, bytecode.OpLoadPropertyConst) {
		t.Fatalf("did not expect OpLoadPropertyConst after same-type loop reassignment")
	}
}
