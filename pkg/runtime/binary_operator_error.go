package runtime

func binaryOperatorTypeError(operator string, left, right Value) error {
	return Errorf(
		ErrInvalidOperation,
		"operator '%s' cannot be applied to %s and %s",
		operator,
		TypeName(TypeOf(left)),
		TypeName(TypeOf(right)),
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
