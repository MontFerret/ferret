package debugger

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/MontFerret/ferret/v2/pkg/internal/operator"
	"github.com/MontFerret/ferret/v2/pkg/parser"
	"github.com/MontFerret/ferret/v2/pkg/parser/fql"
	"github.com/MontFerret/ferret/v2/pkg/runtime"
	"github.com/MontFerret/ferret/v2/pkg/vm"
)

func evaluateExpression(ctx context.Context, text string, scope evalScope) (runtime.Value, error) {
	if strings.TrimSpace(text) == "" {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug expression is empty")
	}

	if ctx == nil {
		ctx = context.Background()
	}

	p := parser.New(text)
	listener := newDebugEvalErrorListener()
	p.RemoveErrorListeners()
	p.AddErrorListener(listener)
	expr := p.Expression()

	if len(listener.errs) > 0 {
		return nil, errors.Join(listener.errs...)
	}

	if !p.AtEOF() {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debug expression contains unsupported trailing input")
	}

	return evalDebugExpression(ctx, expr, scope)
}

func evalDebugExpression(ctx context.Context, expr fql.IExpressionContext, scope evalScope) (runtime.Value, error) {
	node, ok := expr.(*fql.ExpressionContext)

	if !ok || node == nil {
		return nil, unsupportedDebugExpression(expr)
	}

	if node.UnaryOperator() != nil {
		value, err := evalDebugExpression(ctx, node.GetRight(), scope)
		if err != nil {
			return nil, err
		}

		return evalDebugUnary(node.UnaryOperator().GetText(), value, scope.values)
	}

	if node.GetLeft() != nil && node.GetRight() != nil {
		left, err := evalDebugExpression(ctx, node.GetLeft(), scope)
		if err != nil {
			return nil, err
		}

		op := ""
		if node.LogicalAndOperator() != nil {
			op = node.LogicalAndOperator().GetText()
		} else if node.LogicalOrOperator() != nil {
			op = node.LogicalOrOperator().GetText()
		}

		if strings.EqualFold(op, "AND") {
			leftBool, ok := left.(runtime.Boolean)
			if !ok {
				return nil, debugEvalTypeError(left, runtime.TypeBoolean.Name(), scope.values)
			}

			if !leftBool {
				return runtime.False, nil
			}
		} else if strings.EqualFold(op, "OR") {
			leftBool, ok := left.(runtime.Boolean)
			if !ok {
				return nil, debugEvalTypeError(left, runtime.TypeBoolean.Name(), scope.values)
			}

			if leftBool {
				return runtime.True, nil
			}
		}

		right, err := evalDebugExpression(ctx, node.GetRight(), scope)
		if err != nil {
			return nil, err
		}

		return evalDebugLogical(op, left, right)
	}
	if node.GetCondition() != nil {
		condition, err := evalDebugExpression(ctx, node.GetCondition(), scope)
		if err != nil {
			return nil, err
		}

		ok, isBool := condition.(runtime.Boolean)
		if !isBool {
			return nil, debugEvalTypeError(condition, runtime.TypeBoolean.Name(), scope.values)
		}

		if ok && node.GetOnTrue() != nil {
			return evalDebugExpression(ctx, node.GetOnTrue(), scope)
		}

		return evalDebugExpression(ctx, node.GetOnFalse(), scope)
	}

	if node.Predicate() != nil {
		return evalDebugPredicate(ctx, node.Predicate(), scope)
	}

	return nil, unsupportedDebugExpression(expr)
}

func evalDebugPredicate(ctx context.Context, predicate fql.IPredicateContext, scope evalScope) (runtime.Value, error) {
	node, ok := predicate.(*fql.PredicateContext)
	if !ok || node == nil {
		return nil, unsupportedDebugExpression(predicate)
	}

	if node.GetLeft() != nil && node.GetRight() != nil {
		if node.EqualityOperator() == nil {
			return nil, unsupportedDebugExpression(predicate)
		}
		left, err := evalDebugPredicate(ctx, node.GetLeft(), scope)

		if err != nil {
			return nil, err
		}

		right, err := evalDebugPredicate(ctx, node.GetRight(), scope)
		if err != nil {
			return nil, err
		}

		return evalDebugComparison(ctx, node.EqualityOperator().GetText(), left, right)
	}

	return evalDebugAtom(ctx, node.ExpressionAtom(), scope)
}

func evalDebugAtom(ctx context.Context, atom fql.IExpressionAtomContext, scope evalScope) (runtime.Value, error) {
	node, ok := atom.(*fql.ExpressionAtomContext)
	if !ok || node == nil {
		return nil, unsupportedDebugExpression(atom)
	}

	if node.GetLeft() != nil && node.GetRight() != nil {
		left, err := evalDebugAtom(ctx, node.GetLeft(), scope)
		if err != nil {
			return nil, err
		}

		right, err := evalDebugAtom(ctx, node.GetRight(), scope)
		if err != nil {
			return nil, err
		}

		op := ""
		if node.MultiplicativeOperator() != nil {
			op = node.MultiplicativeOperator().GetText()
		} else if node.AdditiveOperator() != nil {
			op = node.AdditiveOperator().GetText()
		} else {
			return nil, unsupportedDebugExpression(atom)
		}

		return evalDebugArithmetic(ctx, op, left, right)
	}

	switch {
	case node.Literal() != nil:
		return evalDebugLiteral(node.Literal())
	case node.Variable() != nil:
		value, exists := scope.locals[node.Variable().GetText()]
		if !exists {
			return nil, runtime.Errorf(runtime.ErrNotFound, "debug variable %q", node.Variable().GetText())
		}

		return value, nil
	case node.Param() != nil:
		name := strings.TrimPrefix(node.Param().GetText(), "@")
		value, exists := scope.params.Get(name)

		if !exists {
			return nil, runtime.Errorf(runtime.ErrNotFound, "debug parameter %q", name)
		}

		return value, nil
	case node.MemberExpression() != nil:
		return evalDebugMember(ctx, node.MemberExpression(), scope)
	case node.Expression() != nil && node.ForExpression() == nil && node.WaitForExpression() == nil && node.ErrorOperator() == nil && node.RecoveryTails() == nil:
		return evalDebugExpression(ctx, node.Expression(), scope)
	default:
		return nil, unsupportedDebugExpression(atom)
	}
}

func evalDebugMember(ctx context.Context, member fql.IMemberExpressionContext, scope evalScope) (runtime.Value, error) {
	node, ok := member.(*fql.MemberExpressionContext)
	if !ok || node == nil || node.RecoveryTails() != nil {
		return nil, unsupportedDebugExpression(member)
	}

	source := node.MemberExpressionSource()
	var value runtime.Value
	var err error

	switch {
	case source.Variable() != nil:
		var exists bool
		value, exists = scope.locals[source.Variable().GetText()]

		if !exists {
			return nil, runtime.Errorf(runtime.ErrNotFound, "debug variable %q", source.Variable().GetText())
		}
	case source.Param() != nil:
		name := strings.TrimPrefix(source.Param().GetText(), "@")

		var exists bool
		value, exists = scope.params.Get(name)

		if !exists {
			return nil, runtime.Errorf(runtime.ErrNotFound, "debug parameter %q", name)
		}
	case source.Expression() != nil:
		value, err = evalDebugExpression(ctx, source.Expression(), scope)
	default:
		return nil, unsupportedDebugExpression(member)
	}

	if err != nil {
		return nil, err
	}

	for _, path := range node.AllMemberExpressionPath() {
		if path.ErrorOperator() != nil || path.ArrayExpansion() != nil || path.ArrayContraction() != nil || path.ArrayQuestionMark() != nil || path.ArrayApply() != nil {
			return nil, unsupportedDebugExpression(path)
		}

		var key runtime.Value
		if property := path.PropertyName(); property != nil {
			text := property.GetText()

			if property.Param() != nil {
				name := strings.TrimPrefix(text, "@")

				var exists bool
				key, exists = scope.params.Get(name)

				if !exists {
					return nil, runtime.Errorf(runtime.ErrNotFound, "debug parameter %q", name)
				}
			} else if property.StringLiteral() != nil {
				text, err = unquoteDebugString(text)
				if err != nil {
					return nil, err
				}

				key = runtime.NewString(text)
			} else {
				key = runtime.NewString(text)
			}
		} else if computed := path.ComputedPropertyName(); computed != nil {
			key, err = evalDebugExpression(ctx, computed.Expression(), scope)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, unsupportedDebugExpression(path)
		}

		value, err = scope.values.Lookup(value, key)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}

func evalDebugLiteral(literal fql.ILiteralContext) (runtime.Value, error) {
	switch {
	case literal.NoneLiteral() != nil:
		return runtime.None, nil
	case literal.BooleanLiteral() != nil:
		return runtime.NewBoolean(strings.EqualFold(literal.BooleanLiteral().GetText(), "true")), nil
	case literal.IntegerLiteral() != nil:
		value, err := strconv.ParseInt(literal.IntegerLiteral().GetText(), 10, 64)
		return runtime.NewInt64(value), err
	case literal.FloatLiteral() != nil:
		value, err := strconv.ParseFloat(literal.FloatLiteral().GetText(), 64)
		return runtime.NewFloat(value), err
	case literal.DurationLiteral() != nil:
		return runtime.ParseDuration(literal.DurationLiteral().GetText())
	case literal.StringLiteral() != nil && literal.StringLiteral().TemplateLiteral() == nil:
		value, err := unquoteDebugString(literal.StringLiteral().GetText())
		return runtime.NewString(value), err
	default:
		return nil, unsupportedDebugExpression(literal)
	}
}

func evalDebugUnary(op string, value runtime.Value, _ vm.DebugValueAccess) (runtime.Value, error) {
	if strings.EqualFold(op, "NOT") {
		op = "!"
	}

	unary, ok := operator.ParseUnary(op)
	if !ok {
		return nil, runtime.Errorf(runtime.ErrInvalidArgument, "unsupported debugger unary operation %s", op)
	}

	switch unary {
	case operator.Not:
		return runtime.Not(value)
	case operator.Positive:
		return runtime.Positive(value)
	case operator.Negative:
		return runtime.Negative(value)
	}

	return nil, runtime.Errorf(runtime.ErrInvalidArgument, "unsupported debugger unary operation %s", op)
}

func evalDebugLogical(op string, left, right runtime.Value) (runtime.Value, error) {
	l, lok := left.(runtime.Boolean)
	r, rok := right.(runtime.Boolean)

	if !lok || !rok {
		return nil, runtime.Error(runtime.ErrInvalidArgument, "debugger logical operations require booleans")
	}

	if strings.EqualFold(op, "AND") {
		return runtime.NewBoolean(bool(l) && bool(r)), nil
	}

	if strings.EqualFold(op, "OR") {
		return runtime.NewBoolean(bool(l) || bool(r)), nil
	}

	return nil, unsupportedDebugExpression(nil)
}

func evalDebugArithmetic(ctx context.Context, op string, left, right runtime.Value) (runtime.Value, error) {
	binary, ok := operator.ParseBinary(op)
	if !ok {
		return nil, unsupportedDebugExpression(nil)
	}

	switch binary {
	case operator.Add:
		return runtime.Add(ctx, left, right)
	case operator.Subtract:
		return runtime.Subtract(ctx, left, right)
	case operator.Multiply:
		return runtime.Multiply(ctx, left, right)
	case operator.Divide:
		return runtime.Divide(ctx, left, right)
	case operator.Modulus:
		return runtime.Modulo(ctx, left, right)
	default:
		return nil, unsupportedDebugExpression(nil)
	}
}

func evalDebugComparison(ctx context.Context, operatorText string, left, right runtime.Value) (runtime.Value, error) {
	op, ok := operator.ParseBinary(operatorText)
	if !ok {
		return nil, unsupportedDebugExpression(nil)
	}

	switch op {
	case operator.Equal:
		return runtime.EqualValues(ctx, left, right)
	case operator.NotEqual:
		equal, err := runtime.EqualValues(ctx, left, right)
		return !equal, err
	}

	if !op.IsRelational() {
		return nil, unsupportedDebugExpression(nil)
	}

	cmp, err := compareDebugValues(ctx, op, left, right)
	switch op {
	case operator.Greater:
		return runtime.NewBoolean(cmp > 0), err
	case operator.Less:
		return runtime.NewBoolean(cmp < 0), err
	case operator.GreaterOrEqual:
		return runtime.NewBoolean(cmp >= 0), err
	case operator.LessOrEqual:
		return runtime.NewBoolean(cmp <= 0), err
	default:
		return nil, unsupportedDebugExpression(nil)
	}
}

func compareDebugValues(
	ctx context.Context,
	op operator.Binary,
	left, right runtime.Value,
) (runtime.Ordering, error) {
	result, err := runtime.CompareValues(ctx, left, right)
	if err == nil {
		return result, nil
	}

	if !errors.Is(err, runtime.ErrInvalidOperation) {
		return result, err
	}

	return runtime.Equal, runtime.Error(
		runtime.ErrInvalidOperation,
		operator.CannotApply(
			op,
			runtime.TypeName(runtime.TypeOf(left)),
			runtime.TypeName(runtime.TypeOf(right)),
		),
	)
}

func unquoteDebugString(text string) (string, error) {
	if len(text) < 2 || (text[0] != '\'' && text[0] != '"') || text[len(text)-1] != text[0] {
		return "", runtime.Error(runtime.ErrInvalidArgument, "invalid debugger string literal")
	}

	content := text[1 : len(text)-1]
	var b strings.Builder
	b.Grow(len(content))

	for i := 0; i < len(content); i++ {
		if content[i] != '\\' || i+1 >= len(content) {
			b.WriteByte(content[i])

			continue
		}

		i++
		switch content[i] {
		case 'n':
			b.WriteByte('\n')
		case 't':
			b.WriteByte('\t')
		default:
			b.WriteByte('\\')
			b.WriteByte(content[i])
		}
	}

	return b.String(), nil
}

func unsupportedDebugExpression(value any) error {
	base := runtime.Error(runtime.ErrInvalidOperation, "expression is not supported by the safe debugger evaluator")
	if value == nil {
		return base
	}

	if text, ok := value.(interface{ GetText() string }); ok {
		return fmt.Errorf("%w: %s", base, text.GetText())
	}

	return base
}

func debugEvalTypeError(value runtime.Value, expected string, access vm.DebugValueAccess) error {
	return runtime.Errorf(runtime.ErrInvalidArgument, "debugger expression requires %s, got %s", expected, access.TypeName(value))
}
