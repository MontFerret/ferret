package compiler_test

import (
	"strings"
	"testing"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/source"
	"github.com/MontFerret/ferret/v2/test/spec"
	. "github.com/MontFerret/ferret/v2/test/spec/compile"
	"github.com/MontFerret/ferret/v2/test/spec/compile/inspect"
)

func assertNoCellOps(t *testing.T, prog *bytecode.Program) {
	t.Helper()

	for _, op := range []bytecode.Opcode{bytecode.OpMakeCell, bytecode.OpLoadCell, bytecode.OpStoreCell} {
		if got := inspect.CountOpcode(prog, op); got != 0 {
			t.Fatalf("expected no %s opcodes, got %d", op, got)
		}
	}
}

func TestVarSyntaxErrors(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		Failure(
			`
			VAR
			RETURN 1
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected variable name",
				Hint:    "Did you forget to provide a variable name?",
			}, "VAR missing variable name"),
		Failure(
			`
			VAR x =
			RETURN x
		`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '=' for variable 'x'",
				Hint:    "Did you forget to provide a value?",
			}, "VAR missing assignment value"),
		Failure(
			`
			VAR _ = 1
			RETURN 0
		`, E{
				Kind: parserd.SyntaxError,
			}, "VAR cannot use discard binding"),
		Failure(
			`
				VAR x = 0
				RETURN (x = 1)
			`, E{
				Kind: parserd.SyntaxError,
			}, "Assignment is not allowed inside expressions"),
		Failure(
			`
				VAR x = 0
				x +=
				RETURN x
			`, E{
				Kind:    parserd.SyntaxError,
				Message: "Expected expression after '+=' for variable 'x'",
				Hint:    "Did you forget to provide a value?",
			}, "Compound assignment missing assignment value"),
		Failure(
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

	_, err := compiler.New(compiler.WithOptimizationLevel(compiler.O0)).Compile(source.New("var_compound_span", src))
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
	RunSpecs(t, []spec.Spec{
		Failure(
			`
				LET x = 1
				x = 2
				RETURN x
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "LET remains immutable"),
		Failure(
			`
				LET x = 1
				x += 1
				RETURN x
			`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'x' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "LET remains immutable for compound assignment"),
		Failure(
			`
				x = 1
				RETURN 0
		`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'x' is not defined",
			}, "Assignment target must already exist"),
		Failure(
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
		Failure(
			`
			FOR i WHILE i < 2
			  i = i + 1
			  RETURN i
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Variable 'i' cannot be reassigned",
				Hint:    "Declare it with VAR if you need to update it.",
			}, "FOR WHILE variables cannot be reassigned"),
		Failure(
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

func TestVarErrorsFunctionAssignmentTargets(t *testing.T) {
	RunSpecsLevels(t, []spec.Spec{
		Failure(
			`
			FUNC test() => 1
			test.foo = 42
			RETURN NONE
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Function 'test' cannot be used as an assignment target",
				Hint:    "Call it as test(...), or assign to a declared VAR binding instead.",
			}, "UDF path assignment reports function-specific diagnostic"),
		Failure(
			`
			FUNC test() => 1
			test = 42
			RETURN NONE
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Function 'test' cannot be used as an assignment target",
				Hint:    "Call it as test(...), or assign to a declared VAR binding instead.",
			}, "UDF bare assignment reports function-specific diagnostic"),
		Failure(
			`
			FUNC test() => 1
			test += 1
			RETURN NONE
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Function 'test' cannot be used as an assignment target",
				Hint:    "Call it as test(...), or assign to a declared VAR binding instead.",
			}, "UDF compound assignment reports function-specific diagnostic"),
		Failure(
			`
			FUNC outer() (
			  FUNC inner() => 1
			  inner.foo = 42
			  RETURN NONE
			)
			RETURN outer()
		`, E{
				Kind:    parserd.SemanticError,
				Message: "Function 'inner' cannot be used as an assignment target",
				Hint:    "Call it as inner(...), or assign to a declared VAR binding instead.",
			}, "Nested UDF path assignment reports function-specific diagnostic"),
	}, compiler.O0, compiler.O1)
}

func TestVarErrorsFunctionAssignmentTargetSpanLabel(t *testing.T) {
	src := `
FUNC test() => 1
test.foo = 42
RETURN NONE
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		_, err := compiler.New(compiler.WithOptimizationLevel(level)).Compile(source.NewAnonymous(src))
		if err == nil {
			t.Fatalf("expected compilation error at O%d", level)
		}

		diag := firstCompilationError(err)
		if diag == nil {
			t.Fatalf("expected diagnostic at O%d", level)
		}

		if diag.Kind != parserd.SemanticError {
			t.Fatalf("expected semantic error at O%d, got %s", level, diag.Kind)
		}

		if diag.Message != "Function 'test' cannot be used as an assignment target" {
			t.Fatalf("unexpected diagnostic message at O%d: %q", level, diag.Message)
		}

		if diag.Hint != "Call it as test(...), or assign to a declared VAR binding instead." {
			t.Fatalf("unexpected diagnostic hint at O%d: %q", level, diag.Hint)
		}

		if len(diag.Spans) == 0 {
			t.Fatalf("expected diagnostic span at O%d", level)
		}

		if diag.Spans[0].Label != "function is not a writable binding" {
			t.Fatalf("unexpected diagnostic span label at O%d: %q", level, diag.Spans[0].Label)
		}
	}
}

func TestVarFunctionNameShadowedByBindingUsesBinding(t *testing.T) {
	expr := `
FUNC target() => 1
FUNC outer() (
  LET target = {}
  target.foo = 42
  RETURN target
)
RETURN outer()
`

	for _, level := range []compiler.OptimizationLevel{compiler.O0, compiler.O1} {
		_ = compileWithLevel(t, level, expr)
	}
}

func TestDirectMutationCompile(t *testing.T) {
	RunSpecs(t, []spec.Spec{
		ProgramCheck(`
			LET obj = {}
			obj.x = 1
			RETURN obj
		`, func(program *bytecode.Program) error {
			if got := inspect.CountOpcode(program, bytecode.OpSetKeyConst); got == 0 {
				t.Fatalf("expected %s opcode", bytecode.OpSetKeyConst)
			}

			return nil
		}, "Property assignment compiles to key write"),
		ProgramCheck(`
			LET arr = [0]
			arr[0] = 1
			RETURN arr
		`, func(program *bytecode.Program) error {
			if got := inspect.CountOpcode(program, bytecode.OpSetIndexConst); got == 0 {
				t.Fatalf("expected %s opcode", bytecode.OpSetIndexConst)
			}

			return nil
		}, "Index assignment compiles to index write"),
		ProgramCheck(`
			LET obj = { count: 1 }
			obj?.count += 1
			RETURN obj
		`, func(program *bytecode.Program) error {
			if got := inspect.CountOpcode(program, bytecode.OpLoadKeyOptionalConst); got == 0 {
				t.Fatalf("expected %s opcode", bytecode.OpLoadKeyOptionalConst)
			}

			if got := inspect.CountOpcode(program, bytecode.OpSetKeyConst); got == 0 {
				t.Fatalf("expected %s opcode", bytecode.OpSetKeyConst)
			}

			return nil
		}, "Safe augmented assignment compiles to optional read and write"),
		Failure(
			`
				LOWER("x") = 1
				RETURN 0
			`, E{
				Kind: parserd.SyntaxError,
			}, "Function call assignment target is invalid"),
		Failure(
			`
				(1 + 2) = 3
				RETURN 0
			`, E{
				Kind: parserd.SyntaxError,
			}, "Expression assignment target is invalid"),
		Failure(
			`
				LET obj = []
				obj[*] = 1
				RETURN obj
			`, E{
				Kind: parserd.SyntaxError,
			}, "Array operator assignment target is invalid"),
		Failure(
			`
				LET obj = []
				obj?[0] = 1
				RETURN obj
			`, E{
				Kind: parserd.SyntaxError,
			}, "Malformed safe index assignment target is invalid"),
		Failure(
			`
				missing?.x = 1
				RETURN 0
			`, E{
				Kind:    parserd.NameError,
				Message: "Variable 'missing' is not defined",
			}, "Safe assignment still requires a declared root"),
		Failure(
			`
				VAR obj = {}
				obj += 1
				RETURN obj
			`, E{
				Kind:    parserd.SemanticError,
				Message: "Operator '+=' cannot be applied to this assignment target",
				Hint:    "Use a numeric binding for arithmetic assignment, or a string binding with +=.",
			}, "Invalid augmented assignment target types are rejected"),
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

		if got := inspect.CountOpcode(prog, bytecode.OpMakeCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpMakeCell, level)
		}

		if got := inspect.CountOpcode(prog, bytecode.OpLoadCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpLoadCell, level)
		}

		if got := inspect.CountOpcode(prog, bytecode.OpStoreCell); got == 0 {
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

		if got := inspect.CountOpcode(prog, bytecode.OpMakeCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpMakeCell, level)
		}

		if got := inspect.CountOpcode(prog, bytecode.OpLoadCell); got == 0 {
			t.Fatalf("expected %s in optimized level %v", bytecode.OpLoadCell, level)
		}

		if got := inspect.CountOpcode(prog, bytecode.OpStoreCell); got == 0 {
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

	if !inspect.HasOpcode(prog, bytecode.OpLoadKeyConst) {
		t.Fatalf("expected OpLoadKeyConst after straight-line reassignment")
	}

	if inspect.HasOpcode(prog, bytecode.OpLoadPropertyConst) {
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

	if !inspect.HasOpcode(prog, bytecode.OpLoadPropertyConst) {
		t.Fatalf("expected OpLoadPropertyConst after loop-scoped conflicting reassignment")
	}

	if inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) || inspect.HasOpcode(prog, bytecode.OpLoadKeyConst) {
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

	if got := inspect.CountOpcode(prog, bytecode.OpMakeCell); got == 0 {
		t.Fatalf("expected OpMakeCell for captured mutable binding")
	}

	if got := inspect.CountOpcode(prog, bytecode.OpLoadCell); got == 0 {
		t.Fatalf("expected OpLoadCell for captured mutable binding")
	}

	if got := inspect.CountOpcode(prog, bytecode.OpStoreCell); got == 0 {
		t.Fatalf("expected OpStoreCell for captured mutable binding")
	}

	if !inspect.HasOpcode(prog, bytecode.OpLoadPropertyConst) {
		t.Fatalf("expected OpLoadPropertyConst after loop-scoped conflicting reassignment through cell binding")
	}

	if inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) || inspect.HasOpcode(prog, bytecode.OpLoadKeyConst) {
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

	if !inspect.HasOpcode(prog, bytecode.OpLoadIndexConst) {
		t.Fatalf("expected OpLoadIndexConst after same-type loop reassignment")
	}

	if inspect.HasOpcode(prog, bytecode.OpLoadPropertyConst) {
		t.Fatalf("did not expect OpLoadPropertyConst after same-type loop reassignment")
	}
}
