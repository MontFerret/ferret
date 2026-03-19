package compiler_test

import (
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
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
	})
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
			FOR i = 0 WHILE i < 2 STEP i = i + 1
			  i = i + 1
			  RETURN i
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'i' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "STEP variables cannot be reassigned"),
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
