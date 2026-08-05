package runtime

import "github.com/MontFerret/ferret/v2/pkg/internal/operator"

func binaryOperatorTypeError(op operator.Binary, left, right Value) error {
	return Error(
		ErrInvalidOperation,
		operator.CannotApply(op, TypeName(TypeOf(left)), TypeName(TypeOf(right))),
	)
}

func incompatibleComparisonError(left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"comparison cannot be applied to %s and %s",
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
	)
}
