package parser_test

import (
	"fmt"
	"testing"

	"github.com/antlr4-go/antlr/v4"

	parserd "github.com/MontFerret/ferret/v2/pkg/parser/diagnostics"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/source"
)

type udfMemberStatementKind string

const (
	udfMemberLet        udfMemberStatementKind = "LET"
	udfMemberVar        udfMemberStatementKind = "VAR"
	udfMemberAssignment udfMemberStatementKind = "assignment"
	udfMemberExpression udfMemberStatementKind = "expression"
)

type udfMemberPathShape struct {
	kind     string
	optional bool
}

func TestUDFMemberStatementsKeepCompleteMemberExpression(t *testing.T) {
	tests := []struct {
		name          string
		statement     string
		member        string
		source        string
		statementKind udfMemberStatementKind
		paths         []udfMemberPathShape
	}{
		{
			name:          "LET mixed member initializer",
			statement:     `LET price = product.attributes["data-price"]`,
			member:        `product.attributes["data-price"]`,
			source:        "variable",
			statementKind: udfMemberLet,
			paths: []udfMemberPathShape{
				{kind: "dot"},
				{kind: "computed"},
			},
		},
		{
			name:          "VAR computed member initializer",
			statement:     `VAR price = product["price"]`,
			member:        `product["price"]`,
			source:        "variable",
			statementKind: udfMemberVar,
			paths:         []udfMemberPathShape{{kind: "computed"}},
		},
		{
			name:          "dot member reassignment",
			statement:     `price = product.price`,
			member:        `product.price`,
			source:        "variable",
			statementKind: udfMemberAssignment,
			paths:         []udfMemberPathShape{{kind: "dot"}},
		},
		{
			name:          "optional member expression statement",
			statement:     `product?.details?.["price"]`,
			member:        `product?.details?.["price"]`,
			source:        "variable",
			statementKind: udfMemberExpression,
			paths: []udfMemberPathShape{
				{kind: "dot", optional: true},
				{kind: "computed", optional: true},
			},
		},
		{
			name:          "function result member initializer",
			statement:     `LET price = BUILD_PRODUCT().price`,
			member:        `BUILD_PRODUCT().price`,
			source:        "function call",
			statementKind: udfMemberLet,
			paths:         []udfMemberPathShape{{kind: "dot"}},
		},
		{
			name:          "object literal member initializer",
			statement:     `VAR price = { price: 42 }.price`,
			member:        `{price:42}.price`,
			source:        "object literal",
			statementKind: udfMemberVar,
			paths:         []udfMemberPathShape{{kind: "dot"}},
		},
		{
			name:          "array literal member reassignment",
			statement:     `price = [{ price: 42 }][0].price`,
			member:        `[{price:42}][0].price`,
			source:        "array literal",
			statementKind: udfMemberAssignment,
			paths: []udfMemberPathShape{
				{kind: "computed"},
				{kind: "dot"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := parseSingleUDFMemberStatement(t, tt.statement)
			statements := block.AllFunctionStatement()
			if len(statements) != 1 {
				t.Fatalf("expected one UDF statement, got %d", len(statements))
			}

			expression := udfMemberStatementExpression(t, statements[0], tt.statementKind)
			members := collectUDFMemberNodes[*fql.MemberExpressionContext](expression)
			if len(members) != 1 {
				t.Fatalf("expected one member expression, got %d", len(members))
			}

			member := members[0]
			if got := member.GetText(); got != tt.member {
				t.Fatalf("unexpected member expression: got %q, want %q", got, tt.member)
			}
			if got := udfMemberSourceKind(member.MemberExpressionSource()); got != tt.source {
				t.Fatalf("unexpected member source: got %q, want %q", got, tt.source)
			}

			paths := member.AllMemberExpressionPath()
			if len(paths) != len(tt.paths) {
				t.Fatalf("unexpected member path count: got %d, want %d", len(paths), len(tt.paths))
			}
			for i, want := range tt.paths {
				if got := udfMemberPathKind(paths[i]); got != want.kind {
					t.Errorf("unexpected member path %d: got %q, want %q", i, got, want.kind)
				}
				if got := paths[i].ErrorOperator() != nil; got != want.optional {
					t.Errorf("unexpected optional flag for member path %d: got %t, want %t", i, got, want.optional)
				}
			}
		})
	}
}

func TestUDFFunctionBlockPreservesArbitraryExpressionStatement(t *testing.T) {
	block := parseSingleUDFMemberStatement(t, `1 + 2`)
	statements := block.AllFunctionStatement()
	if len(statements) != 1 {
		t.Fatalf("expected one UDF statement, got %d", len(statements))
	}

	statement := statements[0].ExpressionStatement()
	if statement == nil {
		t.Fatal("expected an expression statement")
	}
	if got := statement.Expression().GetText(); got != "1+2" {
		t.Fatalf("unexpected expression statement: got %q, want %q", got, "1+2")
	}
	if members := collectUDFMemberNodes[*fql.MemberExpressionContext](statement.Expression()); len(members) != 0 {
		t.Fatalf("expected no member expressions, got %d", len(members))
	}
}

func parseSingleUDFMemberStatement(t *testing.T, statement string) *fql.FunctionBlockContext {
	t.Helper()

	input := fmt.Sprintf(`
FUNC EXTRACT(product) {
	%s
	RETURN NONE
}

RETURN NONE
`, statement)
	src := source.NewAnonymous(input)
	stream := antlr.NewInputStream(input)
	lexer := fql.NewFqlLexer(stream)
	tokens := antlr.NewCommonTokenStream(lexer, antlr.TokenDefaultChannel)
	history := parserd.NewTokenHistory(20)
	errors := parserd.NewErrorHandler(src, 5)
	parser := fql.NewFqlParser(parserd.NewTrackingTokenStream(tokens, history))
	parser.BuildParseTrees = true
	parser.RemoveErrorListeners()
	parser.AddErrorListener(parserd.NewErrorListener(src, errors, history))

	program := parser.Program().(*fql.ProgramContext)
	if errors.HasErrors() {
		t.Fatalf("unexpected parse errors:\n%s", errors.Errors().Format())
	}

	declarations := collectUDFMemberNodes[*fql.FunctionDeclarationContext](program)
	if len(declarations) != 1 {
		t.Fatalf("expected one function declaration, got %d", len(declarations))
	}

	block, ok := declarations[0].FunctionBody().FunctionBlock().(*fql.FunctionBlockContext)
	if !ok || block == nil {
		t.Fatal("expected a block-form function")
	}

	return block
}

func udfMemberStatementExpression(
	t *testing.T,
	statement fql.IFunctionStatementContext,
	kind udfMemberStatementKind,
) fql.IExpressionContext {
	t.Helper()

	switch kind {
	case udfMemberLet:
		declaration := statement.VariableDeclaration()
		if declaration == nil || declaration.Let() == nil {
			t.Fatal("expected a LET declaration")
		}
		return declaration.Expression()
	case udfMemberVar:
		declaration := statement.VariableDeclaration()
		if declaration == nil || declaration.Var() == nil {
			t.Fatal("expected a VAR declaration")
		}
		return declaration.Expression()
	case udfMemberAssignment:
		assignment := statement.AssignmentStatement()
		if assignment == nil {
			t.Fatal("expected an assignment statement")
		}
		return assignment.Expression()
	case udfMemberExpression:
		expression := statement.ExpressionStatement()
		if expression == nil {
			t.Fatal("expected an expression statement")
		}
		return expression.Expression()
	default:
		t.Fatalf("unsupported statement kind %q", kind)
		return nil
	}
}

func udfMemberSourceKind(source fql.IMemberExpressionSourceContext) string {
	switch {
	case source.Variable() != nil:
		return "variable"
	case source.FunctionCall() != nil:
		return "function call"
	case source.ObjectLiteral() != nil:
		return "object literal"
	case source.ArrayLiteral() != nil:
		return "array literal"
	default:
		return "unknown"
	}
}

func udfMemberPathKind(path fql.IMemberExpressionPathContext) string {
	switch {
	case path.PropertyName() != nil:
		return "dot"
	case path.ComputedPropertyName() != nil:
		return "computed"
	default:
		return "unknown"
	}
}

func collectUDFMemberNodes[T any](tree antlr.Tree) []T {
	if tree == nil {
		return nil
	}

	var nodes []T
	if node, ok := tree.(T); ok {
		nodes = append(nodes, node)
	}
	for i := 0; i < tree.GetChildCount(); i++ {
		nodes = append(nodes, collectUDFMemberNodes[T](tree.GetChild(i))...)
	}

	return nodes
}
