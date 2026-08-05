package runtime

import "github.com/MontFerret/ferret/v2/pkg/internal/operator"

func unaryOperatorTypeError(op operator.Unary, value Value) error {
	return Error(
		ErrInvalidOperation,
		operator.CannotApplyUnary(op, TypeName(TypeOf(value))),
	)
}
