package internal

import (
	"testing"

	"github.com/antlr4-go/antlr/v4"

	"github.com/MontFerret/ferret/v2/pkg/bytecode"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/core"
	"github.com/MontFerret/ferret/v2/pkg/compiler/internal/optimization"
	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type bindingCompilerTestState struct {
	body    *fql.BodyContext
	errors  *parserd.ErrorHandler
	front   *CompilationFrontend
	program *fql.ProgramContext
	session *CompilationSession
	source  *source.Source
}

func TestBindingCompilerCapturedMutableDeclarationPromotesToCellStorage(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR base = 1
FUNC outer() (
  FUNC inner() (
    base = base + 1
    RETURN base
  )
  RETURN inner()
)
RETURN outer()
`)

	prepareBindingCompilerTestState(t, state)

	decl := findBindingCompilerTestDeclaration(t, state, "base")
	if !state.front.Bindings.IsPromotedDeclaration(decl) {
		t.Fatal("expected captured mutable declaration to be promoted before lowering")
	}

	state.front.Statements.Compile(state.body)
	state.front.UDFs.CompileAll()

	assertBindingCompilerTestNoErrors(t, state)

	binding, ok := state.session.Function.Symbols.ResolveBinding("base")
	if !ok {
		t.Fatal("expected top-level binding 'base' to exist")
	}

	if binding.Storage != core.BindingStorageCell {
		t.Fatalf("unexpected binding storage: got %v want %v", binding.Storage, core.BindingStorageCell)
	}

	instructions := state.session.Program.Emitter.Bytecode()
	for _, opcode := range []bytecode.Opcode{bytecode.OpMakeCell, bytecode.OpLoadCell, bytecode.OpStoreCell} {
		if got := countBindingCompilerTestOpcode(instructions, opcode); got == 0 {
			t.Fatalf("expected %s in emitted bytecode", opcode)
		}
	}
}

func TestBindingCompilerImmutableLetReassignmentReportsError(t *testing.T) {
	state := newBindingCompilerTestState(t, `
LET x = 1
x = 2
RETURN x
`)

	prepareBindingCompilerTestState(t, state)
	state.front.Statements.Compile(state.body)
	state.front.UDFs.CompileAll()

	if !state.errors.HasErrors() {
		t.Fatal("expected reassignment diagnostic")
	}

	diag := state.errors.Errors().First()
	if diag == nil {
		t.Fatal("expected first diagnostic")
	}

	if diag.Kind != parserd.SemanticError {
		t.Fatalf("unexpected diagnostic kind: got %v want %v", diag.Kind, parserd.SemanticError)
	}

	if diag.Message != "Variable 'x' cannot be reassigned" {
		t.Fatalf("unexpected diagnostic message: got %q", diag.Message)
	}

	if diag.Hint != "Declare it with VAR if you need to update it." {
		t.Fatalf("unexpected diagnostic hint: got %q", diag.Hint)
	}
}

func TestBindingCompilerStringCompoundAssignmentUsesConcatPath(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR text = ""
text += "a"
RETURN text
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAddConst); got != 1 {
		t.Fatalf("unexpected %s count: got %d want 1", bytecode.OpAddConst, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAdd); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpAdd, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpConcat); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpConcat, got)
	}
}

func TestBindingCompilerNumericCompoundAssignmentUsesArithmeticPath(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR total = 1
total += 2
RETURN total
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAdd); got != 1 {
		t.Fatalf("unexpected %s count: got %d want 1", bytecode.OpAdd, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAddConst); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpAddConst, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpConcat); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpConcat, got)
	}
}

func TestBindingCompilerCapturedReadOnlyMutableDoesNotPromote(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR x = 1
FUNC f() (
  RETURN x
)
RETURN f()
`)

	prepareBindingCompilerTestState(t, state)

	decl := findBindingCompilerTestDeclaration(t, state, "x")
	if state.front.Bindings.IsPromotedDeclaration(decl) {
		t.Fatal("expected read-only captured mutable declaration not to be promoted")
	}

	state.front.Statements.Compile(state.body)
	state.front.UDFs.CompileAll()

	assertBindingCompilerTestNoErrors(t, state)

	binding, ok := state.session.Function.Symbols.ResolveBinding("x")
	if !ok {
		t.Fatal("expected top-level binding 'x' to exist")
	}

	if binding.Storage != core.BindingStorageValue {
		t.Fatalf("unexpected binding storage: got %v want %v", binding.Storage, core.BindingStorageValue)
	}

	instructions := state.session.Program.Emitter.Bytecode()
	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpMakeCell); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpMakeCell, got)
	}
}

func TestBindingCompilerSingleLevelCapturedReassignmentPromotesToCell(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR counter = 0
FUNC inc() (
  counter = counter + 1
  RETURN counter
)
RETURN inc()
`)

	prepareBindingCompilerTestState(t, state)

	decl := findBindingCompilerTestDeclaration(t, state, "counter")
	if !state.front.Bindings.IsPromotedDeclaration(decl) {
		t.Fatal("expected single-level captured mutable declaration to be promoted")
	}

	state.front.Statements.Compile(state.body)
	state.front.UDFs.CompileAll()

	assertBindingCompilerTestNoErrors(t, state)

	binding, ok := state.session.Function.Symbols.ResolveBinding("counter")
	if !ok {
		t.Fatal("expected top-level binding 'counter' to exist")
	}

	if binding.Storage != core.BindingStorageCell {
		t.Fatalf("unexpected binding storage: got %v want %v", binding.Storage, core.BindingStorageCell)
	}
}

func TestBindingCompilerImmutableLetCompoundReassignmentReportsError(t *testing.T) {
	state := newBindingCompilerTestState(t, `
LET x = 1
x += 2
RETURN x
`)

	compileBindingCompilerTestState(t, state)

	if !state.errors.HasErrors() {
		t.Fatal("expected reassignment diagnostic for compound assignment on immutable LET")
	}

	diag := state.errors.Errors().First()
	if diag == nil {
		t.Fatal("expected first diagnostic")
	}

	if diag.Kind != parserd.SemanticError {
		t.Fatalf("unexpected diagnostic kind: got %v want %v", diag.Kind, parserd.SemanticError)
	}

	if diag.Message != "Variable 'x' cannot be reassigned" {
		t.Fatalf("unexpected diagnostic message: got %q", diag.Message)
	}
}

func TestBindingCompilerSubtractionCompoundAssignmentUsesArithmeticPath(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR total = 10
total -= 3
RETURN total
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpSub); got != 1 {
		t.Fatalf("unexpected %s count: got %d want 1", bytecode.OpSub, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpConcat); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpConcat, got)
	}
}

func TestBindingCompilerMultiplicationCompoundAssignmentUsesArithmeticPath(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR total = 2
total *= 5
RETURN total
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpMul); got != 1 {
		t.Fatalf("unexpected %s count: got %d want 1", bytecode.OpMul, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpConcat); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpConcat, got)
	}
}

func TestBindingCompilerSequentialStringConcatAssignments(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR text = ""
text += "hello"
text += " world"
RETURN text
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAddConst); got != 2 {
		t.Fatalf("unexpected %s count: got %d want 2", bytecode.OpAddConst, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpAdd); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpAdd, got)
	}
}

func TestBindingCompilerDivisionCompoundAssignmentUsesArithmeticPath(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR total = 100
total /= 4
RETURN total
`)

	compileBindingCompilerTestState(t, state)
	assertBindingCompilerTestNoErrors(t, state)

	instructions := state.session.Program.Emitter.Bytecode()

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpDiv); got != 1 {
		t.Fatalf("unexpected %s count: got %d want 1", bytecode.OpDiv, got)
	}

	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpConcat); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpConcat, got)
	}
}

func TestBindingCompilerFailedInitializerDoesNotDeclareBinding(t *testing.T) {
	state := newBindingCompilerTestState(t, `
LET x = missing
RETURN 1
`)

	prepareBindingCompilerTestState(t, state)
	state.front.Statements.Compile(state.body)

	if !state.errors.HasErrors() {
		t.Fatal("expected initializer diagnostic")
	}

	if _, ok := state.session.Function.Symbols.ResolveBinding("x"); ok {
		t.Fatal("expected failed declaration not to bind x")
	}
}

func TestBindingCompilerFailedInitializerDoesNotPoisonLaterDeclarations(t *testing.T) {
	state := newBindingCompilerTestState(t, `
LET x = missing
LET y = 1
RETURN y
`)

	prepareBindingCompilerTestState(t, state)
	state.front.Statements.Compile(state.body)

	if !state.errors.HasErrors() {
		t.Fatal("expected initializer diagnostic")
	}

	if _, ok := state.session.Function.Symbols.ResolveBinding("x"); ok {
		t.Fatal("expected failed declaration not to bind x")
	}

	binding, ok := state.session.Function.Symbols.ResolveBinding("y")
	if !ok {
		t.Fatal("expected y binding to exist")
	}

	if binding.Register == bytecode.NoopOperand {
		t.Fatal("expected y binding to use a real register")
	}

	instructions := state.session.Program.Emitter.Bytecode()
	if got, want := len(instructions), 2; got != want {
		t.Fatalf("unexpected instruction count: got %d want %d", got, want)
	}

	if got, want := instructions[0].Opcode, bytecode.OpLoadConst; got != want {
		t.Fatalf("unexpected first opcode: got %s want %s", got, want)
	}

	if got := instructions[0].Operands[0]; got == bytecode.NoopOperand {
		t.Fatal("expected valid destination register for y declaration")
	}

	if got, want := instructions[1].Opcode, bytecode.OpReturn; got != want {
		t.Fatalf("unexpected second opcode: got %s want %s", got, want)
	}

	if got, want := instructions[1].Operands[0], instructions[0].Operands[0]; got != want {
		t.Fatalf("unexpected return register: got %s want %s", got, want)
	}
}

func TestBindingCompilerPromotedFailedInitializerSkipsCellBinding(t *testing.T) {
	state := newBindingCompilerTestState(t, `
VAR base = missing
FUNC outer() (
  FUNC inner() (
    base = 1
    RETURN base
  )
  RETURN inner()
)
RETURN 1
`)

	prepareBindingCompilerTestState(t, state)

	decl := findBindingCompilerTestDeclaration(t, state, "base")
	if !state.front.Bindings.IsPromotedDeclaration(decl) {
		t.Fatal("expected base declaration to be promoted before lowering")
	}

	state.front.Statements.Compile(state.body)

	if !state.errors.HasErrors() {
		t.Fatal("expected initializer diagnostic")
	}

	if _, ok := state.session.Function.Symbols.ResolveBinding("base"); ok {
		t.Fatal("expected failed promoted declaration not to bind base")
	}

	instructions := state.session.Program.Emitter.Bytecode()
	if got := countBindingCompilerTestOpcode(instructions, bytecode.OpMakeCell); got != 0 {
		t.Fatalf("unexpected %s count: got %d want 0", bytecode.OpMakeCell, got)
	}
}

func newBindingCompilerTestState(t *testing.T, query string) *bindingCompilerTestState {
	t.Helper()

	src := source.New("binding_compiler_test.fql", query)
	errors := parserd.NewErrorHandler(src, 10)
	session := NewCompilationSession(src, errors, optimization.LevelNone)
	front := NewCompilationFrontend(session)
	program := parseBindingCompilerTestProgram(t, src, errors)

	body, ok := program.Body().(*fql.BodyContext)
	if !ok || body == nil {
		t.Fatal("expected program body context")
	}

	return &bindingCompilerTestState{
		body:    body,
		errors:  errors,
		front:   front,
		program: program,
		session: session,
		source:  src,
	}
}

func parseBindingCompilerTestProgram(t *testing.T, src *source.Source, errors *parserd.ErrorHandler) *fql.ProgramContext {
	t.Helper()

	input := antlr.NewInputStream(src.Content())
	lexer := fql.NewFqlLexer(input)
	stream := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	history := parserd.NewTokenHistory(10)
	parser := fql.NewFqlParser(parserd.NewTrackingTokenStream(stream, history))
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserd.NewErrorListener(src, errors, history))

	program, ok := parser.Program().(*fql.ProgramContext)
	if !ok || program == nil {
		t.Fatal("expected program context")
	}

	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	return program
}

func prepareBindingCompilerTestState(t *testing.T, state *bindingCompilerTestState) {
	t.Helper()

	state.front.UDFCatalog.BuildCatalog(state.program)
	state.front.CaptureAnalyzer.AnalyzeProgram(state.body)
}

func compileBindingCompilerTestState(t *testing.T, state *bindingCompilerTestState) {
	t.Helper()

	prepareBindingCompilerTestState(t, state)
	state.front.Statements.Compile(state.body)
	state.front.UDFs.CompileAll()
}

func findBindingCompilerTestDeclaration(t *testing.T, state *bindingCompilerTestState, name string) antlr.ParserRuleContext {
	t.Helper()

	for _, stmt := range state.body.AllBodyStatement() {
		if stmt == nil || stmt.VariableDeclaration() == nil {
			continue
		}

		decl := stmt.VariableDeclaration()
		if state.front.Bindings.declarationName(decl) == name {
			return decl.(antlr.ParserRuleContext)
		}
	}

	t.Fatalf("expected declaration for %q", name)

	return nil
}

func assertBindingCompilerTestNoErrors(t *testing.T, state *bindingCompilerTestState) {
	t.Helper()

	if state.errors.HasErrors() {
		t.Fatalf("unexpected diagnostics:\n%s", state.errors.Errors().Format())
	}
}

func countBindingCompilerTestOpcode(instructions []bytecode.Instruction, opcode bytecode.Opcode) int {
	count := 0

	for _, inst := range instructions {
		if inst.Opcode == opcode {
			count++
		}
	}

	return count
}
